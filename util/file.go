package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Wallpaper struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	PkgPath   string `json:"pkgPath"`
	CoverPath string `json:"coverPath"`
	Contentrating string `json:"contentrating"`
	Description string `json:"description"`

}

func GetWallpapers(path string) (result []byte, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	dirs, _ := os.ReadDir(path)
	Wallpapers := []Wallpaper{}
	for _, dir := range dirs {
		Wallpaper, err := GetWallpaperInfo(path + "/" +dir.Name())
		if err != nil {
			continue
		}
		Wallpapers = append(Wallpapers, Wallpaper)
	}
	result, err = json.MarshalIndent(Wallpapers, "", "  ")
	
	if err != nil {
		return nil, err
	}
	return result, nil
}


func GetWallpaperInfo (path string) (wallpaper Wallpaper, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return Wallpaper{}, err
	}
	if !fileInfo.IsDir() {
		return Wallpaper{}, fmt.Errorf("%s is not a directory", path)
	}

	// 读取project.json文件
	projectJsonFile, err := os.ReadFile(path + "/project.json")
	if err != nil {
		return Wallpaper{}, err
	}
	var projectJson map[string]interface{}
	err = json.Unmarshal(projectJsonFile, &projectJson)
	if err != nil {
		return Wallpaper{}, err
	}

	wallpaper.Name = projectJson["title"].(string)
	wallpaper.Path = path
	wallpaper.CoverPath = path + "/" + projectJson["preview"].(string)
	wallpaper.Contentrating = projectJson["contentrating"].(string)
	if projectJson["description"] == nil { 
		wallpaper.Description = "无描述"
	} else {
		wallpaper.Description = projectJson["description"].(string)
	}
	
	
	files, _ := os.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "scene.pkg" {
			wallpaper.PkgPath = path + "/scene.pkg"
		}
	}
	return wallpaper, nil
}

func FindWallpaperEnginePath() (string, error) {
	defaultPaths := []string{
		"C:/Program Files (x86)/Steam/steamapps/workshop/content/431960",
		"D:/Steam/steamapps/workshop/content/431960",
		"D:/SteamLibrary/steamapps/workshop/content/431960",
		"E:/SteamLibrary/steamapps/workshop/content/431960",
		"F:/SteamLibrary/steamapps/workshop/content/431960",
	}
	for _, path := range defaultPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("未找到Wallpaper Engine目录")
}

func GetImgBase64(path string) (base64str string, err error) {
	etxs := map[string]bool{
		"png": true,
		"jpg": true,
		"jpeg": true,
		"gif": true,
	}

	//获取文件后缀
	path1 := strings.Split(path, ".")
	ext := strings.ToLower(path1[len(path1)-1]) 
	if !etxs[ext] {
		println("Error:", "不支持的文件类型")
		return "", err
	}
	//读取文件
	raw, err := os.ReadFile(path)
	if err != nil {
		println("Error:", err.Error())
	}
	
	//将文件转换为base64
	dataType := "data:image/" + ext + ";base64,"
	base64str = dataType + base64.StdEncoding.EncodeToString(raw)
	return base64str, nil
}


func GetWallpaperProjectInfo(path string) (info string) {
	fmt.Println(path+"/project.json")
	_, err := os.Stat(path+"/project.json")
	if err != nil {
		return "获取失败 可能是project.json 文件不存在"
	}
	// 读取project.json文件
	projectJsonFile, err := os.ReadFile(path + "/project.json")
	if err != nil {
		return "获取失败 可能是project.json 文件存在错误"
	}
	return string(projectJsonFile)
}