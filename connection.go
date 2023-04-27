package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"sync"
)

type connection interface {
	startReadLoopAsync()

	setFrameHandler(func(t int, m []byte))

	setCloseHandler(func())

	write(b []byte) error

	close()

	remoteAddr() net.Addr
}

type connectionImpl struct {
	conn *websocket.Conn

	closeHandler func()
	frameHandler func(t int, m []byte)
	mutex        sync.RWMutex
}

func newConnection(conn *websocket.Conn) *connectionImpl {
	wrap := &connectionImpl{conn: conn}
	return wrap
}

func (c *connectionImpl) startReadLoopAsync() {
	go func() {
		defer func() {
			_ = c.conn.Close()
			c.handleClose()
		}()

		for {
			t, m, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("ReadMessage error: %v", err)
				return
			}

			c.onFrame(t, m)
		}
	}()
}

func (c *connectionImpl) remoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *connectionImpl) setFrameHandler(handler func(t int, m []byte)) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.frameHandler = handler
}

func (c *connectionImpl) setCloseHandler(handler func()) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.closeHandler = handler
}

func (c *connectionImpl) write(b []byte) error {
	return c.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (c *connectionImpl) onFrame(t int, m []byte) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()

	if c.frameHandler != nil {
		c.frameHandler(t, m)
	}
}

func (c *connectionImpl) close() {
	_ = c.conn.Close()
}

func (c *connectionImpl) handleClose() {
	c.mutex.Lock()
	closeHandler := c.closeHandler
	c.closeHandler = nil
	c.mutex.Unlock()

	if closeHandler != nil {
		closeHandler()
	}
}
