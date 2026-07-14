package user

import (
	"g_kk_ch/internal/model"
	"g_kk_ch/pkg/auth"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(fn func(tx *gorm.DB) error) error
	NewWithTx(tx *gorm.DB) UserRepository

	Login(loginID, password string) (*model.User, string, error)
	Create(user *model.User) error

	FindAll() ([]*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	FindByLoginID(loginID string) (*model.User, error)
	ExistsUser(userID uuid.UUID) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

/* 트랜잭션을 사용하여 userRepository 생성 */
func (r *userRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

/* 새로운 userRepository 생성 ( 트랜잭션을 사용하여 ) */
func (r *userRepository) NewWithTx(tx *gorm.DB) UserRepository {
	return &userRepository{
		db: tx,
	}
}

/* 로그인 */
func (r *userRepository) Login(loginID, password string) (*model.User, string, error) {
	user, err := r.FindByLoginID(loginID)
	if err != nil {
		return nil, "", err
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	if err := r.db.Model(user).UpdateColumn("last_login_at", &now).Error; err != nil {
		return nil, "", err
	}

	return user, token, nil
}

/* 새로운 user 생성 */
func (r *userRepository) Create(user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.db.Create(user).Error
}

/* user 조회 */
func (r *userRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

/* user 조회 ( loginID ) */
func (r *userRepository) FindByLoginID(loginID string) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "login_id = ?", loginID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

/* 사용자 존재 확인 */
func (r *userRepository) ExistsUser(userID uuid.UUID) (bool, error) {
	var user model.User
	// Limit(1)을 적용하여 1개만 찾으면 즉시 종료
	err := r.db.Select("id").Where("id = ?", userID).Limit(1).Find(&user).Error
	if err != nil {
		return false, err
	}

	return r.db.RowsAffected > 0, nil
}

func (r *userRepository) FindAll() ([]*model.User, error) {
	users := make([]*model.User, 0)

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
