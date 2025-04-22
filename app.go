package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"repkg-gui/util"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	BasePath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	BasePath, err := util.FindWallpaperEnginePath()
	for err != nil {
		BasePath, err = runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
	}
	a.BasePath = BasePath
	fmt.Println("BasePath:", BasePath)
}

// Greet returns a greeting for the given name


func (a *App) GetWallpapers() string {
	// 获取壁纸信息
	wallpapers, err := util.GetWallpapers(a.BasePath)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return ""
	}
	return string(wallpapers)
}


func (a *App) GetImgBase64(path string) string {
	// 获取图标文件并转换为base64返回
	data, err := util.GetImgBase64(path)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return ""
	}
	return data
}


func (a *App) OpenDirInExploer(path string) {
	// 打开文件夹
	path = filepath.Clean(path)
	err := exec.Command("explorer.exe", path).Start()
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}


func (a *App) GetRepkgVersion() {
	options := runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "RePKG",
		Message:       "RePKG版本：" + util.RePkgVersion(),
		Buttons:       []string{"确定"},
		DefaultButton: "确定",
		CancelButton:  "",
	}
	runtime.MessageDialog(a.ctx, options)
}

func (a *App) ExtractPKG(source string, target string) {
	if source == ""{
		options := runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "RePKG",
			Message:       "源文件不存在 该壁纸可能没有PKG文件",
			Buttons:       []string{"确定"},
			DefaultButton: "确定",
			CancelButton:  "",
		}
		runtime.MessageDialog(a.ctx, options)
		return
	}
	target = filepath.Clean(target)
	err := util.RePkgExtract(source, target)
	if err != nil {
		options := runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "RePKG",
			Message:       "提取失败：" + err.Error(),
			Buttons:       []string{"确定"},
			DefaultButton: "确定",
			CancelButton:  "",
		}
		runtime.MessageDialog(a.ctx, options)
		return
	}
	err = exec.Command("explorer.exe", target).Start()
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}