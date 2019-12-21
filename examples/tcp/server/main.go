package main

import (
	"bufio"
	"context"
	"encoding/binary"
	cw "github.com/chenpengfei/context-wrapper"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

func main() {
	ctx := cw.WithSignal(context.Background())

	go func() {
		address := "localhost:8080"

		listener, err := net.Listen("tcp", address)
		if err != nil {
			log.WithError(err).Error("start server failed")
			return
		}
		defer listener.Close()

		log.WithField("port", address).Info("server is listening")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.WithError(err).Error("accept connection failed")
				return
			}

			log.Info("one of client has connected")

			go listen(ctx, conn)
		}
	}()

	<-ctx.Done()

	log.Info("I have to go...")
}

func listen(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	head := make([]byte, 4)
	reader := bufio.NewReader(conn)

	for {
		n, err := io.ReadFull(reader, head)
		if err != nil {
			conn.Close()
			return
		}
		if n != len(head) {
			panic("header bytes length not enough")
		}

		size := binary.LittleEndian.Uint32(head)
		buffer := make([]byte, size)

		n, err = io.ReadFull(reader, buffer)
		if err != nil {
			conn.Close()
			return
		}
		if uint32(n) != size {
			panic("data bytes length not enough")
		}

		log.WithField("data", string(buffer[:size])).Info("received new data from client")
	}
}
