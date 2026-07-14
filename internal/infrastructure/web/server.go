package web

import (
	"context"
	"fmt"
	"g_kk_ch/internal/config"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type server struct {
	cfg     *config.Config
	startAt time.Time

	router  *gin.Engine
	httpSrv *http.Server
}

func ensureDataDirectories(cfg *config.Config) error {
	if cfg.UPLOAD_DIR != "" {
		if err := os.MkdirAll(cfg.UPLOAD_DIR, 0o755); err != nil {
			return fmt.Errorf("failed to create upload dir: %w", err)
		}
	}
	if cfg.BACKUP_DIR != "" {
		if err := os.MkdirAll(cfg.BACKUP_DIR, 0o755); err != nil {
			return fmt.Errorf("failed to create backup dir: %w", err)
		}
	}
	return nil
}

func NewServer(cfg *config.Config) (Server, error) {
	if err := ensureDataDirectories(cfg); err != nil {
		return nil, err
	}

	if cfg.GIN_MODE != "" {
		gin.SetMode(cfg.GIN_MODE)
	}

	router := gin.Default()

	// 배포 시에는 AllowOrigins, AllowHeaders 변경
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:5001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// OPTIONS
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	httpSrv := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: router,
	}

	s := &server{
		cfg:     cfg,
		startAt: time.Now(),

		router:  router,
		httpSrv: httpSrv,
	}

	return s, nil
}

func (s *server) Start() error {
	if err := s.registerRoutes(); err != nil {
		return fmt.Errorf("failed to register routes: %w", err)
	}

	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := s.httpSrv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
