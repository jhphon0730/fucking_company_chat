package chat

import (
	"g_kk_ch/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	WithTx(fn func(tx *gorm.DB) error) error
	NewWithTx(tx *gorm.DB) ChatRepository

	// Chat
	SaveChatMessage(chatMessage *model.ChatMessage) error
	ExistsSingleChat(senderID, receiverID uuid.UUID) (uuid.UUID, error)
	GetChatMessages(roomID uuid.UUID, cursor *time.Time, limit int) ([]model.ChatMessage, error)

	// Room
	FindRoomByID(roomID uuid.UUID) (*model.Room, error)
	CreateRoom(room *model.Room) error
	ReadRoom(roomID uuid.UUID, userID uuid.UUID) error
	IsRoomParticipant(roomID, userID uuid.UUID) (bool, error)
	GetUserRooms(userID uuid.UUID) ([]RoomListResponse, error)

	// Room Participant
	CreateParticipant(participant *model.RoomParticipant) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db,
	}
}

func (r *chatRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *chatRepository) NewWithTx(tx *gorm.DB) ChatRepository {
	return &chatRepository{db: tx}
}

func (r *chatRepository) SaveChatMessage(chatMessage *model.ChatMessage) error {
	return r.db.Create(chatMessage).Error
}

// 1:1 채팅 기록이 있는지 확인 -> 있다면 방 ID 반환, 없으면 빈 ID 반환
func (r *chatRepository) ExistsSingleChat(senderID, receiverID uuid.UUID) (uuid.UUID, error) {
	var result struct {
		ID uuid.UUID `gorm:"column:id"`
	}

	err := r.db.Table("rooms").
		Select("rooms.id").
		Joins("JOIN room_participants me ON me.room_id = rooms.id AND me.user_id = ?", senderID).
		Joins("JOIN room_participants you ON you.room_id = rooms.id AND you.user_id = ?", receiverID).
		Where("rooms.is_group = ?", false).
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return uuid.Nil, err
	}

	return result.ID, nil
}

func (r *chatRepository) FindRoomByID(roomID uuid.UUID) (*model.Room, error) {
	var room model.Room
	if err := r.db.First(&room, "id = ?", roomID).Error; err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *chatRepository) CreateRoom(room *model.Room) error {
	if room.ID == uuid.Nil {
		room.ID = uuid.New()
	}

	return r.db.Create(room).Error
}

func (r *chatRepository) CreateParticipant(participant *model.RoomParticipant) error {
	if participant.ID == uuid.Nil {
		participant.ID = uuid.New()
	}

	return r.db.Create(participant).Error
}

func (r *chatRepository) ReadRoom(roomID uuid.UUID, userID uuid.UUID) error {
	return r.db.Table("room_participants").
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("last_read_at", time.Now()).Error
}

func (r *chatRepository) IsRoomParticipant(roomID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Table("room_participants").
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *chatRepository) GetChatMessages(roomID uuid.UUID, cursor *time.Time, limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage

	// 기본 쿼리 빌더 시작 (특정 방의 메시지)
	query := r.db.Table("chat_messages").
		Where("room_id = ?", roomID).
		Order("created_at DESC"). // 최신 메시지부터 거꾸로 정렬
		Limit(limit)

	// 💡 커서(특정 시점)가 있다면, 그 시점보다 더 과거(이전)의 메시지만 조회
	if cursor != nil {
		query = query.Where("created_at < ?", cursor)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *chatRepository) GetUserRooms(userID uuid.UUID) ([]RoomListResponse, error) {
	var rooms []RoomListResponse

	query := `
		SELECT 
			r.id, 
			CASE 
				WHEN r.is_group = false THEN (
					-- 1:1 방이면: 이 방에 참여한 사람 중 '나'가 아닌 사람의 닉네임을 가져옴
					SELECT u.name
					FROM room_participants rp2
					INNER JOIN users u ON rp2.user_id = u.id
					WHERE rp2.room_id = r.id AND rp2.user_id != ?
					LIMIT 1
				)
				ELSE r.name 
			END AS name,
			r.is_group,
			(SELECT content FROM chat_messages cm WHERE cm.room_id = r.id ORDER BY created_at DESC LIMIT 1) AS last_message,
			(SELECT created_at FROM chat_messages cm WHERE cm.room_id = r.id ORDER BY created_at DESC LIMIT 1) AS last_message_at,
			(SELECT COUNT(*) FROM chat_messages cm WHERE cm.room_id = r.id AND cm.created_at > rp.last_read_at) AS unread_count
		FROM rooms r
		INNER JOIN room_participants rp ON r.id = rp.room_id
		WHERE rp.user_id = ?
		ORDER BY last_message_at DESC NULLS LAST
	`

	// 첫 번째 ? : rp2.user_id != ? (상대방 찾기 용도)
	// 두 번째 ? : rp.user_id = ? (내 방 목록 필터링 용도)
	if err := r.db.Raw(query, userID, userID).Scan(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
}
