package websocket

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID uuid.UUID       // 유저 식별자
	Conn   *websocket.Conn // Gorilla 웹소켓 커넥션 객체
	Send   chan []byte     // 이 유저에게 보낼 메시지를 쌓아두는 버퍼 채널
}

// 클라이언트가 서버로 전송한 메시지를 받는 함수
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c // 에러 나거나 연결 끊기면 퇴장 창구로 던짐
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			// 💡 유저가 정상적으로 나갔을 때(1000, 1001)는 에러 로그를 찍지 않도록 예외 처리합니다.
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,     // 1001 (새로고침, 페이지 이동 등)
				websocket.CloseNormalClosure, // 1000 (Disconnect 버튼 클릭 등 정상 종료)
			) {
				// 💡 여기에 걸린 것만 진짜 통신 장애나 서버 다운 등의 '비정상 에러'입니다.
				log.Printf("[Client] 비정상 Read 에러: %v", err)
			}
			break
		}

		hub.Broadcast <- message
	}
}

// 클라이언트 ?? 에게 메시지를 전송해주는 함수
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		// Hub가 client.Send 채널에 메시지를 밀어 넣어줄 때까지 대기합니다.
		message, ok := <-c.Send
		if !ok {
			// Hub가 채널을 닫았으면(Close) 클라이언트에게 종료 메시지를 보냅니다.
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		// 실제 웹소켓 텍스트 메시지 형태로 브라우저에 전송합니다.
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)
		w.Close()
	}
}
