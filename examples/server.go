package examples

import (
	ws "github.com/shenweijiekdel/light-websocket-server"
	"log"
	"net/http"
)

func startServer() {
	var ops ws.Option
	server := ws.NewServer(ops)
	server.EndpointHandler(handleEndpoint)
	err := server.Listen("/endpoint", "0.0.0.0:8988")
	if err != nil {
		log.Printf("server listen error: %v\n", err)
	}
}

func handleEndpoint(r *http.Request, accept func(endpointId string) ws.IEndpoint, reject func(httpStatus int)) {
	log.Printf("handleEndpoint: \n")

	for k, v := range r.Header {
		log.Printf("%s: %v\n", k, v)
	}

	endpoint := accept("endpointId")
	endpoint.SetCloseHandler(func() {
		log.Printf("handleEndpointClosed: %s\n", endpoint.Id())
	})
	endpoint.SetMessageHandler(func(bytes []byte) {
		log.Printf("handleEndpointMessage: %s\n", bytes)
	})
	endpoint.StartLoopSync()
}
