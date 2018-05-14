package Nano

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func listen(addr string, opts ...Option) {
	for _, opt := range opts {
		opt(handler.options)
	}

	hbdEncode()
	startupComponents()

	// create global ticker instance, timer precision could be customized
	// by SetTimerPrecision
	globalTicker = time.NewTicker(timerPrecision)

	// startup logic dispatcher
	go handler.dispatch()

	go func() {
		// if isWs {
		// 	listenAndServeWS(addr)
		// } else {
		listenAndServe(addr)
		// }
	}()

	logger.Println(fmt.Sprintf("starting application %s, listen at %s", app.name, addr))
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

	// stop server
	select {
	case <-env.die:
		logger.Println("The app will shutdown in a few seconds")
	case s := <-sg:
		logger.Println("got signal", s)
	}

	logger.Println("server is stopping...")

	// shutdown all components registered by application, that
	// call by reverse order against register
	shutdownComponents()
}

// Enable current server accept connection
func listenAndServe(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println(err.Error())
			continue
		}

		go handler.handle(conn)
	}
}
