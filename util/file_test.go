package util

import (
	"testing"
)

func TestGetWallpapers(t *testing.T) {
	path, _ := FindWallpaperEnginePath()
	wallpapers, err := GetWallpapers(path)
	t.Log(string(wallpapers), err)
}