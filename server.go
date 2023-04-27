package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type (
	EndpointHandler func(r *http.Request, accept func(endpointId string) IEndpoint, reject func(httpStatus int))
	CloseHandler    func(endpointId string, code int, text string)
)

func NewServer(option Option) Server {
	return &ServerImpl{
		option: parseOption(option),
		upgrade: &websocket.Upgrader{
			ReadBufferSize:   option.ReadBufferSize,
			WriteBufferSize:  option.WriteBufferSize,
			HandshakeTimeout: option.HandshakeTimeout,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type Server interface {
	UpgradeEndpoint(endpointId string, w http.ResponseWriter, r *http.Request) (IEndpoint, error)
}

type ServerImpl struct {
	option          Option
	upgrade         *websocket.Upgrader
	endpointHandler EndpointHandler
}

func (s *ServerImpl) UpgradeEndpoint(endpointId string, w http.ResponseWriter, r *http.Request) (IEndpoint, error) {
	responseHeader := http.Header{
		"Sec-WebSocket-Protocol": r.Header.Values("Sec-WebSocket-Protocol"),
	}

	conn, err := s.upgrade.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return NewEndpoint(endpointId, conn), nil
}

func (s *ServerImpl) EndpointHandler(handler EndpointHandler) {
	s.endpointHandler = handler
}
