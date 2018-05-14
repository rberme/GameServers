package Nano

import (
	"Nano/Component"
	"Nano/Internal/Codec"
	"Nano/Internal/Message"
	"Nano/Internal/Packet"
	"Nano/Session"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"time"
)

// Unhandled message buffer size
const packetBacklog = 1024
const funcBacklog = 1 << 8

var (
	// handler service singleton
	handler = newHandlerService()
	// 序列化的数据
	hrd []byte // 握手回应数据
	hbd []byte // 心跳包数据
)

func hbdEncode() {
	data, err := json.Marshal(map[string]interface{}{
		"code": 200,
		"sys":  map[string]float64{"heartbeat": env.heartbeat.Seconds()},
	})
	if err != nil {
		panic(err)
	}

	hrd, err = Codec.Encode(Packet.Handshake, data)
	if err != nil {
		panic(err)
	}

	hbd, err = Codec.Encode(Packet.Heartbeat, nil)
	if err != nil {
		panic(err)
	}
}

type (
	handlerService struct {
		services       map[string]*Component.Service // 所有注册的服务
		handlers       map[string]*Component.Handler // 所有的处理方法
		chLocalProcess chan unhandledMessage         // packets that process locally
		chCloseSession chan *Session.Session         // closed session
		chFunction     chan func()                   // function that called in logic gorontine
		options        *options
	}

	unhandledMessage struct {
		agent   *agent
		lastMid uint
		handler reflect.Method
		args    []reflect.Value
	}
)

func newHandlerService() *handlerService {
	h := &handlerService{
		services:       make(map[string]*Component.Service),
		handlers:       make(map[string]*Component.Handler),
		chLocalProcess: make(chan unhandledMessage, packetBacklog),
		chCloseSession: make(chan *Session.Session, packetBacklog),
		chFunction:     make(chan func(), funcBacklog),
		options:        &options{},
	}
	return h
}

func (me *handlerService) handle(conn net.Conn) {
	// 创建一个客户端代理并开始 "写go程"
	agent := newAgent(conn, me.options)

	// startup write goroutine
	go agent.write()

	if env.debug {
		logger.Println(fmt.Sprintf("New session established: %s", agent.String()))
	}

	// guarantee agent related resource be destroyed
	defer func() {
		agent.Close()
		if env.debug {
			logger.Println(fmt.Sprintf("Session read goroutine exit, SessionID=%d, UID=%d", agent.session.ID(), agent.session.UID()))
		}
	}()

	// read loop
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			logger.Println(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
		}

		// TODO(warning): decoder use slice for performance, packet data should be copy before next Decode
		packets, err := agent.decoder.Decode(buf[:n])
		if err != nil {
			logger.Println(err.Error())
			return
		}

		if len(packets) < 1 {
			continue
		}

		// process all packet
		for i := range packets {
			if err := me.processPacket(agent, packets[i]); err != nil {
				logger.Println(err.Error())
				return
			}
		}
	}
}

func (me *handlerService) processPacket(agent *agent, p *Packet.Packet) error {
	switch p.Type {
	case Packet.Handshake:
		if _, err := agent.conn.Write(hrd); err != nil {
			return err
		}

		agent.setStatus(statusHandshake)
		if env.debug {
			logger.Println(fmt.Sprintf("Session handshake Id=%d, Remote=%s", agent.session.ID(), agent.conn.RemoteAddr()))
		}

	case Packet.HandshakeAck:
		agent.setStatus(statusWorking)
		if env.debug {
			logger.Println(fmt.Sprintf("Receive handshake ACK Id=%d, Remote=%s", agent.session.ID(), agent.conn.RemoteAddr()))
		}

	case Packet.Data:
		if agent.status() < statusWorking {
			return fmt.Errorf("receive data on socket which not yet ACK, session will be closed immediately, remote=%s",
				agent.conn.RemoteAddr().String())
		}

		msg, err := Message.Decode(p.Data)
		if err != nil {
			return err
		}
		me.processMessage(agent, msg)

	case Packet.Heartbeat:
		// expected
	}

	agent.lastAt = time.Now().Unix()
	return nil
}

func (me *handlerService) processMessage(agent *agent, msg *Message.Message) {
	var lastMid uint
	switch msg.Type {
	case Message.Request: //客户端向服务端请求,并回应给客户端
		lastMid = msg.ID
	case Message.Notify: //客户端通知给服务端,不需要回复
		lastMid = 0
	}

	handler, ok := me.handlers[msg.Route]
	if !ok {
		logger.Println(fmt.Sprintf("nano/handler: %s not found(forgot registered?)", msg.Route))
		return
	}

	if pipe := me.options.pipeline; pipe != nil {
		pipe.Inbound().Process(agent.session, *msg)
	}

	var payload = msg.Data
	var data interface{}
	if handler.IsRawArg {
		data = payload
	} else {
		data = reflect.New(handler.Type.Elem()).Interface()
		err := serializer.Unmarshal(payload, data)
		if err != nil {
			logger.Println("deserialize error", err.Error())
			return
		}
	}

	if env.debug {
		logger.Println(fmt.Sprintf("UID=%d, Message={%s}, Data=%+v", agent.session.UID(), msg.String(), data))
	}

	args := []reflect.Value{handler.Receiver, agent.srv, reflect.ValueOf(data)}
	me.chLocalProcess <- unhandledMessage{agent, lastMid, handler.Method, args}
}

// call handler with protected
func pcall(method reflect.Method, args []reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(fmt.Sprintf("nano/dispatch: %v", err))
			println(stack())
		}
	}()

	if r := method.Func.Call(args); len(r) > 0 {
		if err := r[0].Interface(); err != nil {
			logger.Println(err.(error).Error())
		}
	}
}

// dispatch message to corresponding logic handler
func (me *handlerService) dispatch() {
	// close chLocalProcess & chCloseSession when application quit
	defer func() {
		close(me.chLocalProcess)
		close(me.chCloseSession)
		globalTicker.Stop()
	}()

	// handle packet that sent to chLocalProcess
	for {
		select {
		case m := <-me.chLocalProcess: // logic dispatch
			m.agent.lastMid = m.lastMid
			pcall(m.handler, m.args)

		case s := <-me.chCloseSession: // session closed callback
			onSessionClosed(s)

		case fn := <-me.chFunction:
			pinvoke(fn)

		case <-globalTicker.C: // execute cron task
			cron()

		case t := <-timerManager.chCreatedTimer: // new timers
			timerManager.timers[t.id] = t

		case id := <-timerManager.chClosingTimer: // closing timers
			delete(timerManager.timers, id)

		case <-env.die: // application quit signal
			return
		}
	}
}

// call handler with protected
func pinvoke(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(fmt.Sprintf("nano/invoke: %v", err))
			println(stack())
		}
	}()
	fn()
}

func (me *handlerService) register(comp Component.IComponent, opts []Component.Option) error {
	s := Component.NewService(comp, opts)

	if _, ok := me.services[s.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", s.Name)
	}

	if err := s.ExtractHandler(); err != nil {
		return err
	}

	// register all handlers
	me.services[s.Name] = s
	for name, handler := range s.Handlers {
		me.handlers[fmt.Sprintf("%s.%s", s.Name, name)] = handler
	}
	return nil
}

func onSessionClosed(s *Session.Session) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(fmt.Sprintf("nano/onSessionClosed: %v", err))
			println(stack())
		}
	}()

	env.muCallbacks.RLock()
	defer env.muCallbacks.RUnlock()

	if len(env.callbacks) < 1 {
		return
	}

	for _, fn := range env.callbacks {
		fn(s)
	}
}

// DumpServices outputs all registered services
func (me *handlerService) DumpServices() {
	for name := range me.handlers {
		logger.Println("registered service", name)
	}
}
