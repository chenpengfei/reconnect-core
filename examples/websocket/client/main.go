package main

import (
	"context"
	cw "github.com/chenpengfei/context-wrapper"
	rc "github.com/chenpengfei/reconnect-core"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"math"
	"net/url"
	"time"
)

func main() {
	ctx := cw.WithSignal(context.Background())

	address := "localhost:8080"
	u := url.URL{Scheme: "ws", Host: address, Path: "/echo"}

	re := rc.NewWebsocketReconnection(ctx, u.String(), nil,
		rc.WithMaxElapsedTime(time.Duration(math.MaxInt64)))
	re.OnConnect(func() {
		for {
			data := time.Now().String()
			log.WithField("data", data).Info("write message")
			err := re.WriteMessage(websocket.TextMessage, []byte(time.Now().String()))
			if err != nil {
				_ = re.Close()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	})
	re.OnNotify(func(err error, duration time.Duration) {
		log.WithError(err).WithField("next", duration).Error("retry...")
	})
	re.OnError(func(err error) {
		log.WithError(err).Error("connection has broken")
	})

	<-ctx.Done()

	log.Info("I have to go...")

	// send some raw data to server
	// echo -n "test out the server" | nc localhost 8080
}
