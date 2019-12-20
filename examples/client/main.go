package main

import (
	"bufio"
	"context"
	"encoding/binary"
	cw "github.com/chenpengfei/context-wrapper"
	rc "github.com/chenpengfei/reconnect-core"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	ctx := cw.WithSignal(context.Background())

	address := "localhost:8080"

	re := rc.NewReconnection(ctx, "tcp", address)
	re.OnConnect(func(conn *rc.Reconnection) {
		log.WithField("address", address).Info("connected to server")
		go func() {
			head := make([]byte, 4)
			writer := bufio.NewWriter(conn)

			for {
				data := time.Now().String()
				binary.LittleEndian.PutUint32(head, uint32(len(data)))
				n, err := writer.Write(head)
				if err != nil {
					conn.Close()
					return
				}
				if n != len(head) {
					panic("header bytes length not enough")
				}
				n, err = writer.Write([]byte(data))
				if err != nil {
					conn.Close()
					return
				}
				if n != len(data) {
					panic("data bytes length not enough")
				}

				log.WithField("data", data).Info("send data to server")

				writer.Flush()

				time.Sleep(time.Second)
			}
		}()
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
