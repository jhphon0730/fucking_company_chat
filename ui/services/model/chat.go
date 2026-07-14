package model

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:text"`                           // 그룹방일 경우 방 이름, 1:1방이면 비워둠
	IsGroup   bool      `gorm:"type:boolean;not null;default:false"` // 1:1방과 그룹방 구분 플래그
	CreatedAt time.Time
	UpdatedAt time.Time
}

/* 방 참여자 매핑 */
type RoomParticipant struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoomID   uuid.UUID `gorm:"type:uuid;index;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;index;not null"`
	JoinedAt time.Time `gorm:"not null"`

	LastReadAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// 관계 설정 (GORM 연관관계)
	Room Room `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoomID    uuid.UUID `gorm:"type:uuid;index;not null"`
	SenderID  uuid.UUID `gorm:"type:uuid;index;not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"index;not null"` // 채팅 순서 정렬을 위한 인덱스

	// 관계 설정 (GORM 연관관계)
	Room   Room `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE;"`
	Sender User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;"`
}
