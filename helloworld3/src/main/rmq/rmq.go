package rmq

import (
	"llog"
	"log"
	"sync"
	"utils/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	//exchanges string
	//queues    string
	exchanges = sync.Map{}
	queues    = sync.Map{}

	hasMQ = false
)

// // Reader .
// type Reader interface {
// 	Read(msg *string) (err error)
// }

// func initRabbitMQCtrl(kind string,)

// SetupRMQ 初始化 参数格式：amqp://用户名:密码@地址:端口号/host
func SetupRMQ(rmqAddr string) (err error) {
	if channel == nil {
		conn, err = amqp.Dial(rmqAddr)
		if err != nil {
			return err
		}

		channel, err = conn.Channel()
		if err != nil {
			return err
		}
		hasMQ = true
	}
	return nil
}

// BindQueue 往已有交换器中加入队列
func BindQueue(exchange, queue string) {
	// _, err := channel.QueueInspect(queue)
	// if err == nil {
	// 	return
	// }
	// channel.QueueDeclare(queue, false, true, false, false, nil)
	// channel.QueueBind(queue, queue, exchange, false, nil)
}

// HasMQ 是否已经初始化
func HasMQ() bool {
	return hasMQ
}

// PublishMQ 发布消息
func PublishMQ(exchange, queue string, msg []byte) (err error) {
	//if exchanges == "" || !strings.Contains(exchanges, exchange) {
	if _, ok := exchanges.Load(exchange); ok == false {
		err = channel.ExchangeDeclare(exchange, "direct", false, false, false, false, nil)
		if err != nil {
			llog.Error(err.Error())
			return err
		}
		exchanges.Store(exchange, true)
		//exchanges += "  " + exchange + "  "
	}
	//if queues == "" || !strings.Contains(queues, queue) {
	if _, ok := queues.Load(queue); ok == false {
		_, err = channel.QueueDeclare(queue, false, true, false, false, nil)
		if err != nil {
			llog.Error(err.Error())
			return err
		}
		err = channel.QueueBind(queue, queue, exchange, false, nil)
		if err != nil {
			llog.Error(err.Error())
			return err
		}
		//queues += "  " + queue + "  "
		queues.Store(queue, true)
	}

	err = channel.Publish(exchange, queue, false, false, amqp.Publishing{
		ContentType: "text/plain", //"application/octet-stream", //
		Body:        msg,
	})
	if err != nil {
		exchanges.Delete(exchange)
		queues.Delete(queue)
		log.Printf(err.Error())
	}
	return nil
}

// ReceiveMQ 监听接收到的消息
func ReceiveMQ(exchange, queue string, reader func(msg []byte)) (err error) {
	// channel.ExchangeDelete(exchange, false, true)
	// channel.QueueDelete(queue, false, false, true)
	//if exchanges == "" || !strings.Contains(exchanges, exchange) {
	if _, ok := exchanges.Load(exchange); ok == false {
		err = channel.ExchangeDeclare(exchange, "direct", false, false, false, false, nil)
		if err != nil {
			llog.Error(err.Error())
			return err
		}
		exchanges.Store(exchange, true)
		//exchanges += "  " + exchange + "  "
	}
	//if queues == "" || !strings.Contains(queues, queue) {
	if _, ok := queues.Load(queue); ok == false {
		_, err = channel.QueueDeclare(queue, false, true, false, false, nil)
		if err != nil {
			llog.Error(err.Error())
			return err
		}
		err = channel.QueueBind(queue, queue, exchange, false, nil)
		if err != nil {
			llog.Error(err.Error())
			return err
		}
		queues.Store(queue, true)
		//queues += "  " + queue + "  "
	}

	msgs, err := channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		//fmt.Println(*msgs)
		for d := range msgs {
			//s := bytesToString(&(d.Body))
			reader(d.Body)
		}
	}()
	return nil
}

// Close 关闭连接
func Close() {
	channel.Close()
	conn.Close()
	hasMQ = false
}

// // Ping 测试连接是否正常
// func Ping() (err error) {
// 	if !hasMQ || channel == nil {
// 		return errors.New("RabbitMQ is not initialize")
// 	}
// 	err = channel.ExchangeDeclare("ping.ping", "topic", false, true, false, true, nil)
// 	if err != nil {
// 		return err
// 	}
// 	msgContent := "ping.ping"
// 	err = channel.Publish("ping.ping", "ping.ping", false, false, amqp.Publishing{
// 		ContentType: "text/plain",
// 		Body:        []byte(msgContent),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	err = channel.ExchangeDelete("ping.ping", false, false)
// 	return err
// }
