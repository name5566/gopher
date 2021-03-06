/*
下载
*/
package gopher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Filename    string `json:"filename"`
	Summary     string `json:"summary"`
	Size        string `json:"size"`
	SHA1        string `json:"sha1"`
	URL         string `json:"url"`
	Recommended bool   `json:"recommended,omitempty"`
}

type Version struct {
	Version string `json:"version"`
	Files   []File `json:"files`
	Date    string `json:"date"`
}

type FileInfo struct {
	Filename string
	Size     int64 // bytes
}

func (info FileInfo) HumanSize() string {
	if info.Size < 1024 {
		return fmt.Sprintf("%d B", info.Size)
	} else if info.Size < 1024*1024 {
		return fmt.Sprintf("%d K", info.Size/1024)
	} else {
		return fmt.Sprintf("%d M", info.Size/1024/1024)
	}
}

type VersionInfo struct {
	Name  string
	Files []FileInfo
}

// 获取版本信息
// downloadPaht: 下载路径
// categoryLength: 分类路径
func getVersions(downloadPath string, categoryLength int) []VersionInfo {
	fileLength := categoryLength + 1
	versions := []VersionInfo{}

	var version VersionInfo

	first := true
	filepath.Walk(downloadPath, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		if path == downloadPath {
			return nil
		}

		temp := strings.Split(path, "/")
		if len(temp) == categoryLength {
			// 版本文件夹
			if !first {
				versions = append(versions, version)
			} else {
				first = false
			}

			version = VersionInfo{
				Name:  info.Name(),
				Files: []FileInfo{},
			}
		} else if len(temp) == fileLength {
			// 文件
			version.Files = append(version.Files, FileInfo{
				Filename: info.Name(),
				Size:     info.Size(),
			})
		}

		return nil
	})

	versions = append(versions, version)

	// 倒序排列
	count := len(versions)
	for i := 0; i < count/2; i++ {
		versions[i], versions[count-i-1] = versions[count-i-1], versions[i]
	}

	return versions
}

func downloadGoHandler(handler *Handler) {
	handler.renderTemplate("download.html", BASE, map[string]interface{}{
		"versions": getVersions("/data/gopher/static/go", 6),
		"active":   "download",
	})
}

func downloadLiteIDEHandler(handler *Handler) {
	handler.renderTemplate("download/liteide.html", BASE, map[string]interface{}{
		"versions": getVersions("./static/liteide", 3),
		"active":   "download",
	})
}
