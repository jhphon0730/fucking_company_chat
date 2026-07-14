package websocket

import (
	"g_kk_ch/internal/infrastructure/web/httputil"
	"g_kk_ch/internal/user"
	"g_kk_ch/pkg/apperror"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler interface {
	ConnectWebSocket(c *gin.Context)
}

type webSocketHandler struct {
	Hub         *Hub
	userService user.UserService
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 테스트 편의 및 CORS 허용을 위해 true 리턴
		return true
	},
}

func NewWebSocketHandler(hub *Hub, userService user.UserService) WebSocketHandler {
	return &webSocketHandler{
		hub,
		userService,
	}
}

// ConnectWebSocket은 r.GET("/ws") 엔드포인트의 처리기입니다.
func (h *webSocketHandler) ConnectWebSocket(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil || tokenString == "" {
		// 쿠키가 없거나 읽을 수 없는 경우
		httputil.Fail(c, http.StatusUnauthorized, apperror.ErrInvalidRequest.Error(), "인증 쿠키가 존재하지 않습니다.")
		return
	}

	/* 요청을 보낸 사용자가 누군 지 검증 */
	userID, err := h.userService.ValidateToken(tokenString)
	if err != nil {
		httputil.Fail(c, http.StatusUnauthorized, apperror.ErrUserNotFound.Error(), err.Error())
		return
	}

	/* 요청 업글 */
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return // 에러시 upgrader가 내부적으로 응답하므로 return만 수행
	}

	// 이 유저만을 위한 통신 인스턴스(Client) 생성
	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256), // 버퍼 크기 256
	}

	// 주입하고 있던 Hub의 Register 채널로 유저 전달-> Hub가 명단에 올림
	h.Hub.Register <- client

	// 독립된 고루틴으로 이 유저의 읽기/쓰기 시작
	go client.WritePump()
	go client.ReadPump(h.Hub)
}
