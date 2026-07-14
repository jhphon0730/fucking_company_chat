package user

import (
	"errors"
	"g_kk_ch/internal/model"
	"g_kk_ch/pkg/apperror"
	"g_kk_ch/pkg/auth"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Login(loginID, password string) (*model.User, string, error)
	Create(user *model.User) error

	ValidateToken(token string) (uuid.UUID, error)

	FindAll() ([]*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	FindByLoginID(loginID string) (*model.User, error)
	ExistsUser(userID uuid.UUID) (bool, error)
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Login(loginID, password string) (*model.User, string, error) {
	user, token, err := s.userRepository.Login(loginID, password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperror.ErrUserNotFound
		}

		if errors.Is(err, apperror.ErrUserInvalidPassword) {
			return nil, "", err
		}

		return nil, "", err
	}
	return user, token, nil
}

func (s *userService) Create(user *model.User) error {
	if err := s.userRepository.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperror.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (s *userService) ValidateToken(token string) (uuid.UUID, error) {
	claims, err := auth.ValidateJWTToken(token)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserID, nil
}

func (s *userService) FindByID(id uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *userService) FindByLoginID(loginID string) (*model.User, error) {
	user, err := s.userRepository.FindByLoginID(loginID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *userService) ExistsUser(userID uuid.UUID) (bool, error) {
	return s.userRepository.ExistsUser(userID)
}

func (s *userService) FindAll() ([]*model.User, error) {
	return s.userRepository.FindAll()
}
