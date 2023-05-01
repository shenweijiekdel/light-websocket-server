package ws

import (
	"github.com/shenweijiekdel/light-websocket-server/stanza"
	"net"
)

type IEndpoint interface {
	Id() string
	StartLoopSync()
	RemoteAddress() net.Addr
	SendMessage(content []byte) error
	Send(stanza stanza.Stanza) error
	Kickoff() error
	SetPingHandler(handler func())
	SetCloseHandler(func())
	SetMessageHandler(handler MessageHandler)
	Close()
}
