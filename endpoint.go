package ws

import (
	"github.com/gorilla/websocket"
	"github.com/shenweijiekdel/light-websocket-server/stanza"
	"log"
	"net"
)

type (
	DisconnectHandler func()
	MessageHandler    func([]byte)
)

type Endpoint struct {
	id                string
	stanzaId          uint32
	conn              connection
	pingHandler       func()
	disconnectHandler DisconnectHandler
	messageHandler    MessageHandler
}

func NewEndpoint(id string, conn *websocket.Conn) IEndpoint {
	return &Endpoint{
		id:   id,
		conn: newConnection(conn),
	}
}

func (e *Endpoint) SendMessage(message []byte) error {
	return e.Send(stanza.Message(message))
}

func (e *Endpoint) StartLoopSync() {
	e.initialize()
	e.conn.startReadLoopAsync()
}

func (e *Endpoint) Id() string {
	return e.id
}

func (e *Endpoint) RemoteAddress() net.Addr {
	return e.conn.remoteAddr()
}

func (e *Endpoint) SetPingHandler(handler func()) {
	e.pingHandler = handler
}

func (e *Endpoint) SetMessageHandler(handler MessageHandler) {
	e.messageHandler = handler
}

func (e *Endpoint) SetCloseHandler(f func()) {
	e.disconnectHandler = f
}

func (e *Endpoint) Close() {
	e.conn.close()
}

func (e *Endpoint) handleStanza(s stanza.Stanza) {
	if _, ok := s.(stanza.Ping); ok {
		e.handlePing()
	}

	if message, ok := s.(stanza.Message); ok {
		e.handleMessage(message)
	}
}

func (e *Endpoint) onStanza(b []byte) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("Endpoint [%s] recv recover panic: %v", e.id, err)
		}
	}()

	s, err := stanza.Decode(b)
	if err != nil {
		log.Printf("Endpoint [%s] stanza decode error: %v\n", e.id, err)
		return
	}

	e.handleStanza(s)
}

func (e *Endpoint) Kickoff() error {
	return e.Send(stanza.Kickoff{})
}

func (e *Endpoint) Send(s stanza.Stanza) error {
	b, err := stanza.Encode(s)
	if err != nil {
		return err
	}

	return e.conn.write(b)
}

func (e *Endpoint) handleClosed() {
	log.Printf("Endpoint [%s] close\n", e.id)
	if e.disconnectHandler != nil {
		e.disconnectHandler()
	}
}

func (e *Endpoint) handleFrame(t int, m []byte) {
	if t == websocket.BinaryMessage {
		e.onStanza(m)
	}
}

func (e *Endpoint) handlePing() {
	if e.pingHandler == nil {
		e.pingHandler = func() {
			err := e.Send(&stanza.Pong{})
			if err != nil {
				log.Printf("send pong error: %v", err)
			}
		}
	}

	e.pingHandler()
}

func (e *Endpoint) handleMessage(message stanza.Message) {
	if e.messageHandler != nil {
		e.messageHandler(message)
	}
}

func (e *Endpoint) initialize() {
	e.conn.setCloseHandler(e.handleClosed)
	e.conn.setFrameHandler(e.handleFrame)
}
