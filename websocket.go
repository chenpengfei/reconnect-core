package reconnect_core

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebsocketReconnection struct {
	*websocket.Conn
	*Reconnection
}

func NewWebsocketReconnection(ctx context.Context, urlStr string, requestHeader http.Header, opts ...Option) *WebsocketReconnection {
	wrc := &WebsocketReconnection{
		Reconnection: NewReconnection(ctx, opts...),
	}

	wrc.operation = func() error {
		conn, _, err := websocket.DefaultDialer.DialContext(ctx, urlStr, requestHeader)
		if err == nil {
			wrc.Conn = conn
		}
		return err
	}

	wrc.retry()

	return wrc
}

func (wrc *WebsocketReconnection) Close() error {
	err := wrc.Conn.Close()
	wrc.retry()
	return err
}
