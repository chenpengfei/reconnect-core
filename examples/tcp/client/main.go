package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"github.com/chenpengfei/backoff"
	cw "github.com/chenpengfei/context-wrapper"
	reconnect "github.com/chenpengfei/reconnect-core"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	ctx := cw.WithSignal(context.Background())

	address := "localhost:8080"

	rc := reconnect.NewReconnection(ctx)
	rc.BackOff = backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx)
	rc.OnConnect(func(raw net.Conn) {
		log.WithField("address", address).Info("connected to server")
		run(raw)
		rc.Dial("tcp", address)
	})
	rc.OnNotify(func(err error, duration time.Duration) {
		log.WithError(err).WithField("next", duration).Error("retry...")
	})
	rc.OnError(func(err error) {
		log.WithError(err).Error("connection has broken")
	})
	rc.Dial("tcp", address)

	<-ctx.Done()

	log.Info("I have to go...")

	// send some raw data to server
	// echo -n "test out the server" | nc localhost 8080
}

func run(conn net.Conn) {
	head := make([]byte, 4)
	writer := bufio.NewWriter(conn)

	for {
		data := time.Now().String()
		binary.LittleEndian.PutUint32(head, uint32(len(data)))
		n, err := writer.Write(head)
		if err != nil {
			_ = conn.Close()
			return
		}
		if n != len(head) {
			panic("header bytes length not enough")
		}
		n, err = writer.Write([]byte(data))
		if err != nil {
			_ = conn.Close()
			return
		}
		if n != len(data) {
			panic("data bytes length not enough")
		}

		log.WithField("data", data).Info("send data to server")

		_ = writer.Flush()

		time.Sleep(time.Second)
	}
}
