package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"ui/services/model"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const websocketPath = "/ws"

func (s *HTTPClientService) ConnectWebSocket() error {
	s.chatMu.Lock()
	if s.chatConn != nil {
		s.chatMu.Unlock()
		return nil
	}
	s.chatMu.Unlock()

	endpoint, err := s.websocketEndpoint()
	if err != nil {
		return err
	}

	dialer := *websocket.DefaultDialer
	if s.client != nil {
		dialer.Jar = s.client.Jar
	}

	conn, _, err := dialer.Dial(endpoint, nil)
	if err != nil {
		return err
	}

	s.chatMu.Lock()
	s.chatConn = conn
	s.chatMu.Unlock()

	go s.readPump()
	return nil
}

func (s *HTTPClientService) DisconnectWebSocket() error {
	s.chatMu.Lock()
	conn := s.chatConn
	s.chatConn = nil
	s.chatMu.Unlock()

	if conn == nil {
		return nil
	}

	err := conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client closing connection"),
	)
	if err != nil {
		return conn.Close()
	}

	return conn.Close()
}

func (s *HTTPClientService) readPump() {
	s.chatMu.Lock()
	conn := s.chatConn
	s.chatMu.Unlock()

	if conn == nil {
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.cleanupWebSocketConnection(conn)
			return
		}

		var wsMsg model.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		runtime.EventsEmit(s.ctx, string(wsMsg.Type), wsMsg)
	}
}

func (s *HTTPClientService) cleanupWebSocketConnection(conn *websocket.Conn) {
	s.chatMu.Lock()
	defer s.chatMu.Unlock()

	if s.chatConn == conn {
		s.chatConn = nil
	}
}

func (s *HTTPClientService) websocketEndpoint() (string, error) {
	parsedURL, err := url.Parse(s.baseURL)
	if err != nil {
		return "", err
	}

	switch parsedURL.Scheme {
	case "http":
		parsedURL.Scheme = "ws"
	case "https":
		parsedURL.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported websocket scheme: %s", parsedURL.Scheme)
	}

	parsedURL.Path = websocketPath
	parsedURL.RawPath = websocketPath
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""

	return parsedURL.String(), nil
}
