package main

import (
	"context"
	"embed"
	"ui/services"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	httpClient := services.NewHTTPClientService()

	err := wails.Run(&options.App{
		Title:  "JHKIM_KAKAO",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			httpClient.Startup(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			_ = httpClient.DisconnectWebSocket()
			return false
		},
		OnShutdown: func(ctx context.Context) {
			_ = httpClient.DisconnectWebSocket()
		},
		Bind: []interface{}{
			app,
			httpClient,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
