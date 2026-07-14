package model

import "github.com/google/uuid"

type MessageType string

const (
	// TypeReadMessage는 특정 대화방을 클릭해 들어갔을 때, "이 방 메시지 다 읽었음"을 알리는 신호
	TypeReadMessage MessageType = "READ"

	// TypeTalkMessage는 일반적인 메시지 전송 (첫 메시지면 내부적으로 방과 연동)
	TypeTalkMessage MessageType = "TALK"

	// TypeLeaveMessage는 단체 그룹방에서 내가 '그룹 나가기'를 눌렀을 때 (방 참여자 테이블에서 삭제)
	TypeLeaveMessage MessageType = "LEAVE"

	TypeTalkGroupMessage MessageType = "TALK_GROUP"
)

type WSMessage struct {
	Type       MessageType `json:"type"`                  // READ, TALK, LEAVE 등...
	SenderID   uuid.UUID   `json:"sender_id"`             // 보낸 사람 (나)
	RoomID     *uuid.UUID  `json:"room_id,omitempty"`     // 대화 중인 방 ID (첫 채팅이라면 NULL, 첫 채팅 메시지 보낸 이후엔 RoomID가 생긴 이후로 웹소켓에 응답이 옴 / 굳이 넣지 않아도 됨)
	ReceiverID *uuid.UUID  `json:"receiver_id,omitempty"` // 1:1 클릭 후 첫 메시지 보낼 때 상대방 ID
	Content    string      `json:"content"`               // 메시지 내용
}
