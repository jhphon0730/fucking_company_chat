package user

import (
	"errors"
	"g_kk_ch/internal/infrastructure/web/httputil"
	"g_kk_ch/internal/model"
	"g_kk_ch/pkg/apperror"
	"g_kk_ch/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	FindAll(c *gin.Context)
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) FindAll(c *gin.Context) {
	users, err := h.userService.FindAll()
	if err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	type userResponse struct {
		ID      uuid.UUID `json:"id"`
		LoginID string    `json:"login_id"`
		Name    string    `json:"name"`

		CreatedAt   time.Time  `json:"created_at"`
		LastLoginAt *time.Time `json:"last_login_at"`
	}

	resUsers := make([]*userResponse, len(users))
	for i, user := range users {
		u := &userResponse{
			user.ID,
			user.LoginID,
			user.Name,
			user.CreatedAt,
			user.LastLoginAt,
		}

		resUsers[i] = u
	}

	httputil.OKMessage(c, http.StatusOK, "사용자 조회 성공", gin.H{
		"users": resUsers,
	})
}

func (h *userHandler) Register(c *gin.Context) {
	type registerRequest struct {
		LoginID  string `json:"login_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), err.Error())
		return
	}

	hashedPassword, err := auth.GenerateHashPassword(req.Password)
	if err != nil {
		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	user := &model.User{
		LoginID:      req.LoginID,
		Name:         req.Name,
		PasswordHash: hashedPassword,
	}

	if err := h.userService.Create(user); err != nil {
		if errors.Is(err, apperror.ErrUserAlreadyExists) {
			httputil.Fail(c, http.StatusConflict, apperror.ErrUserAlreadyExists.Error(), err.Error())
		} else {
			httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		}
		return
	}

	httputil.OKMessage(c, http.StatusCreated, "회원가입 성공", nil)
}

func (h *userHandler) Login(c *gin.Context) {
	type loginRequest struct {
		LoginID  string `json:"login_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Fail(c, http.StatusBadRequest, apperror.ErrInvalidRequest.Error(), err.Error())
		return
	}

	user, token, err := h.userService.Login(req.LoginID, req.Password)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) || errors.Is(err, apperror.ErrUserInvalidPassword) {
			httputil.Fail(c, http.StatusUnauthorized, err.Error(), err.Error())
			return
		}

		httputil.Fail(c, http.StatusInternalServerError, apperror.ErrInternalServerError.Error(), err.Error())
		return
	}

	// HttpOnly Cookie 세팅 (24시간)
	c.SetCookie(
		auth.JWT_COOKIE_NAME,
		token,
		int(24*time.Hour/time.Second), // maxAge (초)
		"/",
		"",    // domain: 빈 문자열 = 현재 도메인
		false, // secure: HTTPS 배포 시 true로 변경
		true,  // httpOnly
	)

	type loginResponse struct {
		ID      uuid.UUID `json:"id"`
		LoginID string    `json:"login_id"`
		Name    string    `json:"name"`

		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		LastLoginAt *time.Time `json:"last_login_at"`
	}

	httputil.OKMessage(c, http.StatusOK, "로그인 성공", gin.H{
		"user": &loginResponse{
			user.ID,
			user.LoginID,
			user.Name,
			user.CreatedAt,
			user.UpdatedAt,
			user.LastLoginAt,
		},
		"token": token,
	})
}
