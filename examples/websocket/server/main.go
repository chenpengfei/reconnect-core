package main

import (
	"context"
	cw "github.com/chenpengfei/context-wrapper"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("upgrade")
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.WithError(err).Error("read message")
			break
		}

		log.WithField("message_type", mt).
			WithField("message", string(message)).
			Info("received new data from client, write back...")

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.WithError(err).Error("write message")
			break
		}
	}
}

func main() {
	ctx := cw.WithSignal(context.Background())

	address := "localhost:8080"

	go func() {
		http.HandleFunc("/echo", echo)
		log.Fatal(http.ListenAndServe(address, nil))
	}()

	<-ctx.Done()

	log.Info("I have to go...")
}
