package websocket

import (
	"encoding/json"
	"g_kk_ch/internal/model"
	"log"
	"sync"

	"github.com/google/uuid"
)

type RoomValidator interface {
	TalkChat(msg *model.WSMessage) (*model.ChatMessage, error)
	IsRoomParticipant(roomID, userID uuid.UUID) (bool, error)
}

type Hub struct {
	mu sync.RWMutex

	Clients map[uuid.UUID]*Client

	Broadcast chan []byte

	Register   chan *Client
	Unregister chan *Client

	roomValidator RoomValidator
}

func NewHub(roomValidator RoomValidator) *Hub {
	return &Hub{
		mu:            sync.RWMutex{},
		Clients:       make(map[uuid.UUID]*Client),
		Broadcast:     make(chan []byte),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		roomValidator: roomValidator,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("[Hub] 유저 %s 접속완료. 현재 동시접속자: %d명", client.UserID, len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("[Hub] 유저 %s 접속종료. 현재 동시접속자: %d명", client.UserID, len(h.Clients))
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			var wsMessage model.WSMessage
			if err := json.Unmarshal(message, &wsMessage); err != nil {
				log.Printf("[Hub] 메시지 파싱 에러: %v", err)
				continue
			}

			// 💡 메시지 타입에 따른 핵심 비즈니스 로직 분기
			switch wsMessage.Type {
			case model.TypeTalkMessage:
				// 언제든 메시지를 보내면 방 검증 후 저장 및 상대방에게 전달하는 함수
				h.handleUserTalk(&wsMessage)

			case model.TypeLeaveMessage:
				// 그룹방 나가기 버튼을 눌렀을 때 대화방 멤버에서 제외하는 함수
				// h.handleRoomLeave(&wsMessage)

			default:
				log.Println(wsMessage)
			}
		}
	}
}

func (h *Hub) handleUserTalk(wsMsg *model.WSMessage) {
	saveMsg, err := h.roomValidator.TalkChat(wsMsg)
	if err != nil {
		log.Printf("[HUB] 메시지 전송 실패, %v", err.Error())
		return
	}

	broadcastMessage := model.WSMessage{
		Type:       model.TypeTalkMessage,
		SenderID:   saveMsg.SenderID,
		RoomID:     &saveMsg.RoomID,
		ReceiverID: wsMsg.ReceiverID,
		Content:    saveMsg.Content,
	}
	payload, err := json.Marshal(broadcastMessage)
	if err != nil {
		log.Printf("[HUB] 패킷 직렬화 실패, %v\n", err.Error())
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if senderClient, connected := h.Clients[wsMsg.SenderID]; connected {
		h.trySend(senderClient, payload)
	}

	if wsMsg.ReceiverID != nil {
		if receiverClient, connected := h.Clients[*wsMsg.ReceiverID]; connected {
			h.trySend(receiverClient, payload)
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID uuid.UUID, msg *model.WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal room read message: %v", err)
		return
	}

	for _, client := range h.Clients {
		if client.UserID == msg.SenderID {
			continue
		}

		isParticipant, err := h.roomValidator.IsRoomParticipant(roomID, client.UserID)
		if err != nil || !isParticipant {
			continue
		}

		h.trySend(client, payload)
	}
}

func (h *Hub) trySend(client *Client, payload []byte) {
	if client == nil {
		return
	}

	select {
	case client.Send <- payload:
	default:
		log.Printf("[Hub] 클라이언트 %s에게 전송이 큐 가득참, 메시지 드롭", client.UserID)
	}
}
