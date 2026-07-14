package chat

import (
	"errors"
	"g_kk_ch/internal/infrastructure/web/httputil"
	"g_kk_ch/internal/infrastructure/web/websocket"
	"g_kk_ch/internal/model"
	"g_kk_ch/internal/user"
	"g_kk_ch/pkg/apperror"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatHandler interface {
	GetChatMessages(c *gin.Context)

	GetUserRooms(c *gin.Context)
	ReadRoom(c *gin.Context)
	CreateRoom(c *gin.Context)

	AddParticipants(c *gin.Context)
}

type chatHandler struct {
	hub *websocket.Hub

	chatService ChatService
	userService user.UserService
}

type RoomListResponse struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	LastMessage   string     `json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"` // 메시지가 없을 수도 있으므로 포인터
	UnreadCount   int        `json:"unread_count"`
}

func NewChatHandler(hub *websocket.Hub, chatService ChatService, userService user.UserService) ChatHandler {
	return &chatHandler{
		hub,

		chatService,
		userService,
	}
}

func (h *chatHandler) CreateRoom(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), err.Error())
		return
	}

	userID, err := httputil.GetUserIDByContext(c)
	if err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	user, err := h.userService.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.Fail(c, http.StatusNotFound, apperror.ErrUserNotFound.Error(), err.Error())
			return
		}
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	/* 방 생성 초기엔 맴버는 방 생성자 1명 */
	room := &model.Room{
		Name:    req.Name,
		IsGroup: true,
	}

	participant := &model.RoomParticipant{
		UserID:   user.ID,
		JoinedAt: time.Now(),
	}

	if err := h.chatService.CreateRoom(room, participant); err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	httputil.OKMessage(c, http.StatusCreated, "그룹방 생성 성공", nil)
}

func (h *chatHandler) AddParticipants(c *gin.Context) {
	// 1. URL에서 방 ID 추출 및 유효성 검증
	roomIDStr := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), apperror.ErrInvalidRoomID)
		return
	}

	// 2. 요청 바디 데이터(초대할 유저 ID 목록) 바인딩
	var req struct {
		UserIDs []uuid.UUID `json:"user_ids" binding:"required,min=1"` // 💡 최소 1명 이상 초대하도록 제한
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), "초대할 유저 ID 목록이 필요합니다.")
		return
	}

	// 3. 서비스 레이어에 넘겨줄 RoomParticipant 슬라이스 조립
	participants := make([]*model.RoomParticipant, len(req.UserIDs))
	for i, userID := range req.UserIDs {
		participants[i] = &model.RoomParticipant{
			ID:       uuid.New(),
			RoomID:   roomID, // URL에서 추출한 방 ID 주입
			UserID:   userID, // 초대받은 유저 ID 주입
			JoinedAt: time.Now(),
		}
	}

	// 4. 서비스 로직 호출 (방 존재 여부 검증 및 트랜잭션 처리 포함됨)
	if err := h.chatService.AddParticipants(participants); err != nil {
		// 서비스 레이어에서 정의한 커스텀 에러(예: apperror.ErrInvalidRoomID) 분기 처리 가능
		if errors.Is(err, apperror.ErrInvalidRoomID) {
			httputil.Fail(c, http.StatusNotFound, apperror.ErrInvalidRoomID.Error(), err.Error())
			return
		}

		// 그 외 서버 내부 에러
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	// 5. 성공 응답
	httputil.OKMessage(c, http.StatusOK, "멤버 초대 성공", nil)
}

func (h *chatHandler) ReadRoom(c *gin.Context) {
	roomIDStr := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), "올바르지 않은 방 ID 형식입니다.")
		return
	}

	userID, err := httputil.GetUserIDByContext(c)
	if err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	if _, err := h.userService.FindByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.Fail(c, http.StatusNotFound, apperror.ErrUserNotFound.Error(), err.Error())
			return
		}
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	if err := h.chatService.ReadRoom(roomID, userID); err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	readNotification := &model.WSMessage{
		Type:     model.TypeReadMessage,
		RoomID:   &roomID,
		SenderID: userID,
	}

	h.hub.BroadcastToRoom(roomID, readNotification)

	httputil.OKMessage(c, http.StatusOK, "읽음 처리 완료", nil)
}

func (h *chatHandler) GetChatMessages(c *gin.Context) {
	// 1. URL 파라미터에서 Room ID 추출 및 검증
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), "올바르지 않은 방 ID 형식입니다.")
		return
	}

	// 2. 쿼리 파라미터 파싱 (limit, cursor)
	limitStr := c.DefaultQuery("limit", "30")
	limit, _ := strconv.Atoi(limitStr)

	var cursorTime *time.Time
	cursorStr := c.Query("cursor")
	if cursorStr != "" {
		// 예: "2026-03-31T15:04:05Z" 포맷의 문자열을 time.Time으로 변환
		parsedTime, err := time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), "올바르지 않은 커서(시간) 형식입니다. RFC3339 규격을 사용하세요.")
			return
		}
		cursorTime = &parsedTime
	}

	// 3. 비즈니스 로직 호출
	messages, err := h.chatService.GetChatMessages(roomID, cursorTime, limit)
	if err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	// 4. 응답 데이터 구성 (다음 스크롤 시 사용할 next_cursor를 함께 내려주면 프론트엔드가 아주 편해합니다)
	var nextCursor string
	if len(messages) > 0 {
		// 가져온 메시지 중 가장 마지막(가장 과거) 메시지의 시간을 다음 커서로 지정
		nextCursor = messages[len(messages)-1].CreatedAt.Format(time.RFC3339)
	}

	response := gin.H{
		"messages":    messages,
		"next_cursor": nextCursor,
		"count":       len(messages),
	}

	httputil.OKMessage(c, http.StatusOK, "채팅 기록 조회 성공", response)
}

func (h *chatHandler) GetUserRooms(c *gin.Context) {
	userID, err := httputil.GetUserIDByContext(c)
	if err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	type RoomListResponse struct {
		ID            uuid.UUID  `json:"id"`
		Name          string     `json:"name"`
		Type          string     `json:"type"` // 1:1, 그룹 등
		LastMessage   string     `json:"last_message"`
		LastMessageAt *time.Time `json:"last_message_at"` // 메시지가 없을 수도 있으므로 포인터
		UnreadCount   int        `json:"unread_count"`
	}

	rooms, err := h.chatService.GetUserRooms(userID)
	if err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	httputil.OKMessage(c, http.StatusOK, "채팅방 목록 조회 성공", gin.H{
		"rooms": rooms,
	})
}
