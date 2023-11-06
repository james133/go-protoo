package client

import (
	"crypto/tls"
	"net/http"
	"time"
	"github.com/james133/go-protoo/event_emitter"
	"github.com/james133/go-protoo/logger"
	"github.com/james133/go-protoo/transport"
	"github.com/gorilla/websocket"
)

const pingPeriod = 5 * time.Second

type WebSocketClient struct {
	IEventEmitter
	socket          *websocket.Conn
	transport Transport
	handleWebSocket func(transport *Transport)
}

func NewClient(url string, handleWebSocket func(transport *Transport)) *WebSocketClient {
	var client WebSocketClient
	client.IEventEmitter = NewEventEmitter()
	logger.Infof("Connecting to %s", url)

	responseHeader := http.Header{}
	responseHeader.Add("Sec-WebSocket-Protocol", "protoo")

	// only for testing
	tls_cfg := &tls.Config{
		InsecureSkipVerify: true,
	}

	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  tls_cfg,
	}

	socket, _, err := dialer.Dial(url, responseHeader)
	if err != nil {
		logger.Errorf("Dial failed: %v", err)
		return nil
	}
	client.socket = socket
	client.handleWebSocket = handleWebSocket
	client.transport = transport.NewWebsocketTransport(conn)
	if err2 := client.transport.Run(); err2 != nil {
		logger.Err(err2).Msg("transport.run")
	}
	client.handleWebSocket(client.transport)
	return &client
}

func (client *WebSocketClient) GetTransport() *transport.WebSocketTransport {
	return client.transport
}

func (client *WebSocketClient) Close() {
	client.transport.Close()
}
