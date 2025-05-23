package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"repkg-gui/util"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed bin/RePKG.exe
var repkgBin []byte

type FileLoader struct {
	http.Handler
	AppBasePath *string
}

func CoverFileLoader(appBasePath *string) *FileLoader {
	return &FileLoader{
		AppBasePath: appBasePath,
	}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := filepath.Clean(*h.AppBasePath + "/" + req.URL.Path)
	fileData, err := util.GetImgByte(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
		return
	}
	res.Write(fileData)
}

func (a *App) clearRePkg(ctx context.Context) {
	// 清除RePKG.exe
	err := os.Remove("RePKG.exe")
	if err != nil {
		println("Error:", err.Error())
	}
}



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
		Title:  "RePKG-GUI v1.2.1",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: CoverFileLoader(&app.AppBasePath),
		},
		BackgroundColour: &options.RGBA{R: 248, G: 248, B: 248, A: 1},
		OnStartup:        app.startup,
		OnShutdown:		  app.clearRePkg,
		Bind: []interface{}{
			app,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
