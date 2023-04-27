package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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
	Listen(path string, addr string) error
	EndpointHandler(handler EndpointHandler)
}

type ServerImpl struct {
	option          Option
	upgrade         *websocket.Upgrader
	endpointHandler EndpointHandler
}

func (s *ServerImpl) Listen(path string, addr string) error {
	if s.endpointHandler == nil {
		return errors.New("please set endpoint handler")
	}

	http.HandleFunc(path, s.handleWebsocket)
	server := http.Server{
		Addr:        addr,
		ReadTimeout: 2 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *ServerImpl) EndpointHandler(handler EndpointHandler) {
	s.endpointHandler = handler
}

func (s *ServerImpl) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	s.endpointHandler(r, func(endpointId string) IEndpoint {
		responseHeader := http.Header{
			"Sec-WebSocket-Protocol": r.Header.Values("Sec-WebSocket-Protocol"),
		}

		conn, err := s.upgrade.Upgrade(w, r, responseHeader)
		if err != nil {
			w.WriteHeader(400)
			_, _ = w.Write([]byte("error"))
		}

		return NewEndpoint(endpointId, conn)
	}, func(httpStatus int) {
		w.WriteHeader(httpStatus)
	})
}
