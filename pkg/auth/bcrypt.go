package auth

import (
	"g_kk_ch/internal/config"
	"g_kk_ch/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), utils.InterfaceToInt(cfg.BCRYPT_COST))
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
