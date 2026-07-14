package apperror

import "errors"

var (
	ErrInternalServerError = errors.New("내부 서버 오류입니다.")
	ErrInvalidRequest      = errors.New("잘못된 요청입니다.")
	ErrInvalidUUIDParam    = errors.New("올바르지 않은 형식의 URL입니다.")
	ErrCannotFindCookie    = errors.New("인증 쿠키가 존재하지 않습니다.")
	ErrUnauthorized        = errors.New("인증이 필요합니다.")
	ErrInvalidUser         = errors.New("인증된 사용자가 아닙니다.")
	ErrForbidden           = errors.New("권한이 없습니다.")

	// auth
	ErrInvalidToken = errors.New("invalid token")

	// user
	ErrUserInvalidPassword = errors.New("잘못된 비밀번호입니다.")
	ErrUserNotFound        = errors.New("사용자를 찾을 수 없습니다.")
	ErrUserAlreadyExists   = errors.New("사용자가 이미 존재합니다.")

	// chat
	ErrNeedReceiverID = errors.New("새로운 대화를 시작하려면 상대방(receiver_id) 정보가 필요합니다")
	ErrInvalidRoomID  = errors.New("방을 찾을 수 없습니다.")
)
