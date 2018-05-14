package Nano

import (
	"Nano/Session"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// VERSION returns current nano version
var VERSION = "0.0.1"

var (
	// app represents the current server process
	app = &struct {
		name    string    // current application name
		startAt time.Time // startup time
	}{}

	// env represents the environment of the current process, includes
	// work path and config path etc.
	env = &struct {
		wd          string                   // working path
		die         chan bool                // wait for end application
		heartbeat   time.Duration            // heartbeat internal
		checkOrigin func(*http.Request) bool // check origin when websocket enabled
		debug       bool                     // enable debug
		wsPath      string                   // WebSocket path(eg: ws://127.0.0.1/wsPath)

		// session closed handlers
		muCallbacks sync.RWMutex           // protect callbacks
		callbacks   []SessionClosedHandler // session关闭的时候触发的调用
	}{}
)

type (
	// SessionClosedHandler 一个回调方法,在一个session关闭或者中断后被调用
	SessionClosedHandler func(session *Session.Session)
)

// init default configs
func init() {
	// application initialize
	app.name = strings.TrimLeft(filepath.Base(os.Args[0]), "/")
	app.startAt = time.Now()

	// environment initialize
	if wd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		env.wd, _ = filepath.Abs(wd)
	}

	env.die = make(chan bool)
	env.heartbeat = 30 * time.Second
	env.debug = false
	env.muCallbacks = sync.RWMutex{}
	env.checkOrigin = func(_ *http.Request) bool { return true }
}
