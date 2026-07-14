package chat

import (
	"errors"
	"g_kk_ch/internal/model"
	"g_kk_ch/pkg/apperror"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService interface {
	// Chat
	TalkChat(msg *model.WSMessage) (*model.ChatMessage, error)
	GetChatMessages(roomID uuid.UUID, cursor *time.Time, limit int) ([]model.ChatMessage, error)

	// Room
	FindRoomByID(roomID uuid.UUID) (*model.Room, error)
	CreateRoom(room *model.Room, participant *model.RoomParticipant) error
	ReadRoom(roomID uuid.UUID, userID uuid.UUID) error
	IsRoomParticipant(roomID, userID uuid.UUID) (bool, error)
	GetUserRooms(userID uuid.UUID) ([]RoomListResponse, error)

	// Room Participant
	AddParticipants(participants []*model.RoomParticipant) error
}

type chatService struct {
	chatRepository ChatRepository
}

func NewChatService(chatRepository ChatRepository) ChatService {
	return &chatService{
		chatRepository,
	}
}

func (s *chatService) TalkChat(msg *model.WSMessage) (*model.ChatMessage, error) {
	senderID := msg.SenderID
	var targetRoomID uuid.UUID

	// [CASE 1] 대화방 ID(RoomID)가 없는 상태로 메시지가 왔을 때 (첫 메시지 시작)
	if msg.RoomID == nil || *msg.RoomID == uuid.Nil {
		if msg.ReceiverID == nil || *msg.ReceiverID == uuid.Nil {
			return nil, apperror.ErrNeedReceiverID
		}

		// 1. 이미 두 사람이 속한 1:1 방이 DB에 존재하는지 레포지토리 조회
		existingRoomID, err := s.chatRepository.ExistsSingleChat(senderID, *msg.ReceiverID)
		if err == nil && existingRoomID != uuid.Nil {
			// 기존에 쓰던 1:1 방이 있다면 해당 방 ID를 재사용
			targetRoomID = existingRoomID
		} else {
			// 진짜 처음 톡을 나누는 사이라면 트랜잭션을 열고 방과 참여자를 새로 개설
			targetRoomID = uuid.New()

			txErr := s.chatRepository.WithTx(func(tx *gorm.DB) error {
				// 트랜잭션 전용 레포지토리 주입
				txRepo := s.chatRepository.NewWithTx(tx)

				// 텔레그램형 1:1 방 생성 (이름은 비워둠)
				room := &model.Room{
					ID:        targetRoomID,
					IsGroup:   false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := txRepo.CreateRoom(room); err != nil {
					return err
				}

				// 참여자 등록 (나)
				me := &model.RoomParticipant{
					ID:       uuid.New(),
					RoomID:   targetRoomID,
					UserID:   senderID,
					JoinedAt: time.Now(),
				}
				if err := txRepo.CreateParticipant(me); err != nil {
					return err
				}

				// 참여자 등록 (상대방)
				you := &model.RoomParticipant{
					ID:       uuid.New(),
					RoomID:   targetRoomID,
					UserID:   *msg.ReceiverID,
					JoinedAt: time.Now(),
				}
				if err := txRepo.CreateParticipant(you); err != nil {
					return err
				}

				return nil
			})

			if txErr != nil {
				return nil, txErr
			}
		}
	} else {
		// [CASE 2] 이미 기존 대화방 안에서 정상적으로 톡을 주고받는 상태
		targetRoomID = *msg.RoomID
	}

	// 2. 최종 결정된 RoomID로 채팅 메시지 모델 빌드 및 DB 저장
	chatMessage := &model.ChatMessage{
		ID:        uuid.New(),
		RoomID:    targetRoomID,
		SenderID:  senderID,
		Content:   msg.Content,
		CreatedAt: time.Now(), // 방 생성 이후 0.00000s 정도의 차이가 나야 에러가 뜨지 않음
	}

	if err := s.chatRepository.SaveChatMessage(chatMessage); err != nil {
		return nil, err
	}

	return chatMessage, nil
}

func (s *chatService) FindRoomByID(roomID uuid.UUID) (*model.Room, error) {
	return s.chatRepository.FindRoomByID(roomID)
}

func (s *chatService) CreateRoom(room *model.Room, participant *model.RoomParticipant) error {
	return s.chatRepository.WithTx(func(tx *gorm.DB) error {
		txRepo := s.chatRepository.NewWithTx(tx)

		if err := txRepo.CreateRoom(room); err != nil {
			return err
		}

		participant.RoomID = room.ID
		if err := txRepo.CreateParticipant(participant); err != nil {
			return err
		}

		return nil
	})
}

func (s *chatService) AddParticipants(participants []*model.RoomParticipant) error {
	if len(participants) == 0 {
		return nil
	}

	// 1. 방 존재 여부 검증
	room, err := s.chatRepository.FindRoomByID(participants[0].RoomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrInvalidRoomID
		}
		return err
	}
	if room == nil {
		return apperror.ErrInvalidRoomID
	}

	return s.chatRepository.WithTx(func(tx *gorm.DB) error {
		txRepo := s.chatRepository.NewWithTx(tx)

		for _, participant := range participants {
			if err := txRepo.CreateParticipant(participant); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *chatService) ReadRoom(roomID uuid.UUID, userID uuid.UUID) error {
	return s.chatRepository.ReadRoom(roomID, userID)
}

func (s *chatService) IsRoomParticipant(roomID, userID uuid.UUID) (bool, error) {
	return s.chatRepository.IsRoomParticipant(roomID, userID)
}

func (s *chatService) GetChatMessages(roomID uuid.UUID, cursor *time.Time, limit int) ([]model.ChatMessage, error) {
	// 기본값 방어 코드: limit이 없거나 음수면 기본 30개로 지정
	if limit <= 0 {
		limit = 30
	}
	// 과도한 데이터를 한 번에 요청하는 것을 방지 (Max 100개 제한)
	if limit > 100 {
		limit = 100
	}

	return s.chatRepository.GetChatMessages(roomID, cursor, limit)
}

func (s *chatService) GetUserRooms(userID uuid.UUID) ([]RoomListResponse, error) {
	return s.chatRepository.GetUserRooms(userID)
}
