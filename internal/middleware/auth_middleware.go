package middleware

import (
	"g_kk_ch/internal/infrastructure/web/httputil"
	"g_kk_ch/pkg/apperror"
	"g_kk_ch/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(auth.JWT_COOKIE_NAME)
		if err != nil {
			httputil.Fail(c, http.StatusUnauthorized, apperror.ErrUnauthorized.Error(), "token cookie not found")
			c.Abort()
			return
		}

		claims, err := auth.ValidateJWTToken(tokenString)
		if err != nil {
			httputil.Fail(c, http.StatusUnauthorized, apperror.ErrInvalidUser.Error(), err.Error())
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
