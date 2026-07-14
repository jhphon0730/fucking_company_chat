package infrastructure

import (
	"errors"
	"g_kk_ch/internal/config"
	"g_kk_ch/internal/model"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db_instance *gorm.DB
	once        sync.Once
)

func InitDatabase() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	dsn := "host=" + cfg.POSTGRES.DB_HOST + " user=" + cfg.POSTGRES.DB_USER + " password=" + cfg.POSTGRES.DB_PASSWORD + " dbname=" + cfg.POSTGRES.DB_NAME + " port=" + cfg.POSTGRES.DB_PORT + " sslmode=" + cfg.POSTGRES.SSL_MODE + " TimeZone=" + cfg.POSTGRES.TIMEZONE
	db_instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return err
	}

	db, err := db_instance.DB()
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	// 풀 설정 추가
	db.SetMaxOpenConns(10)   // 최대 10개의 연결을 허용
	db.SetMaxIdleConns(10)   // 최대 10개의 유휴 연결을 허용
	db.SetConnMaxLifetime(0) // 연결의 최대 수명을 무제한으로 설정

	return nil
}

func GetDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		err = InitDatabase()
	})

	if err != nil {
		return nil, err
	}

	return db_instance, nil
}

func CloseDB() error {
	if db_instance == nil {
		return errors.New("database not initialized")
	}

	db, err := db_instance.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func Migration() error {
	if db_instance == nil {
		return errors.New("database not initialized")
	}

	log.Println("Starting database migration...")
	return db_instance.AutoMigrate(
		model.User{},
		model.ChatMessage{},
		model.Room{},
		model.RoomParticipant{},
	)
}
