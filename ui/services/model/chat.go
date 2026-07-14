package model

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:text" json:"name"`                               // 그룹방일 경우 방 이름, 1:1방이면 비워둠
	IsGroup   bool      `gorm:"type:boolean;not null;default:false" json:"is_group"` // 1:1방과 그룹방 구분 플래그
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/* 방 참여자 매핑 */
type RoomParticipant struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID   uuid.UUID `gorm:"type:uuid;index;not null" json:"room_id"`
	UserID   uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	JoinedAt time.Time `gorm:"not null" json:"joined_at"`

	LastReadAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"last_read_at"`

	// 관계 설정 (GORM 연관관계)
	Room Room `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;" json:"room,omitempty"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
}

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID    uuid.UUID `gorm:"type:uuid;index;not null" json:"room_id"`
	SenderID  uuid.UUID `gorm:"type:uuid;index;not null" json:"sender_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"index;not null" json:"created_at"` // 채팅 순서 정렬을 위한 인덱스

	// 관계 설정 (GORM 연관관계)
	Room   Room `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;" json:"room,omitempty"`
	Sender User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;" json:"sender,omitempty"`
}

/* DTO */
type RoomListResponse struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	IsGroup       bool       `json:"is_group"`
	LastMessage   string     `json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"`
	UnreadCount   int        `json:"unread_count"`
}
