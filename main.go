package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed bin/RePKG.exe
var repkgBin []byte

func main() {

	// 判断bin/RePKG.exe是否存在，不存在则提取
	info, err := os.Stat("RePKG.exe")
	if os.IsNotExist(err) || info.Size() == 0 {
		// 提取RePKG.exe
		os.WriteFile("RePKG.exe", repkgBin, 0777)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "repkg-gui",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 248, G: 248, B: 248, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
