package examples

import (
	ws "github.com/shenweijiekdel/light-websocket-server"
	"log"
	"net/http"
	"time"
)

func startServer() {
	var ops ws.Option
	server := ws.NewServer(ops)
	http.HandleFunc("/endpoint", func(writer http.ResponseWriter, request *http.Request) {
		endpoint, err := server.UpgradeEndpoint("1", writer, request)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		endpoint.SetCloseHandler(func() {
			log.Printf("handleEndpointClosed: %s\n", endpoint.Id())
		})
		endpoint.SetMessageHandler(func(bytes []byte) {
			log.Printf("handleEndpointMessage: %s\n", bytes)
		})
		endpoint.StartLoopSync()

	})

	httpServer := http.Server{
		Addr:        "0.0.0.0:8988",
		ReadTimeout: 2 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return
	}
}
