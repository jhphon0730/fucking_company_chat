package main

import (
	"context"
	"g_kk_ch/internal/config"
	"g_kk_ch/internal/infrastructure"
	"g_kk_ch/internal/infrastructure/web"
	"g_kk_ch/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalln("failed load config", err.Error())
	}

	if err := infrastructure.InitDatabase(); err != nil {
		log.Fatalln("failed load database", err.Error())
	}

	if utils.InterfaceToBool(cfg.DB_AUTO_MIGRATE) {
		if err := infrastructure.Migration(); err != nil {
			log.Fatalln("failed migration database", err.Error())
		}
	}

	srv, err := web.NewServer(cfg)
	if err != nil {
		log.Fatalln("fialed load server", err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalln("failed start server", err.Error())
		}
	}()

	<-c

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("failed shutdown server", err.Error())
	}

	cancel()

	log.Println("Server gracefully stopped")
}
