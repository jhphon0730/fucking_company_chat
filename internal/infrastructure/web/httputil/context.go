package httputil

import (
	"g_kk_ch/pkg/apperror"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserIDByContext(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, apperror.ErrUnauthorized
	}
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperror.ErrUnauthorized
	}
	return userID, nil
}
