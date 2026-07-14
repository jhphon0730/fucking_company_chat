package services

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/gorilla/websocket"
)

type HTTPClientService struct {
	baseURL string

	ctx      context.Context
	client   *http.Client
	chatMu   sync.Mutex
	chatConn *websocket.Conn
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type APIResponse[T any] struct {
	Data    T         `json:"data"`
	Message string    `json:"message,omitempty"`
	Error   *APIError `json:"error"`
}

func NewHTTPClientService() *HTTPClientService {
	jar, _ := cookiejar.New(nil)
	return &HTTPClientService{
		baseURL: "http://192.168.0.85:5001/api/v1",

		client: &http.Client{
			Jar: jar,
		},
	}
}

func (s *HTTPClientService) Startup(ctx context.Context) {
	s.ctx = ctx
}
