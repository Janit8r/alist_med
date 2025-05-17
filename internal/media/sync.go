package media

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/internal/model"
)

type NfoMeta struct {
	XMLName xml.Name `xml:"movie"`
	Title   string   `xml:"title"`
	Tag     string   `xml:"tag>name"`
	// 可根据实际nfo结构扩展
}

// 支持的文件后缀
var (
	strmExts   = []string{".strm"}
	nfoExts    = []string{".nfo"}
	subExts    = []string{".ass", ".ssa", ".sub", ".srt"}
	coverNames = []string{"cover.jpg", "cover.png", "poster.jpg", "poster.png"}
	fanartNames= []string{"fanart.jpg", "fanart.png"}
)

// SyncMediaLibrary 递归扫描并写入数据库
func SyncMediaLibrary(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		var item model.MediaItem
		item.StrmPath = ""
		item.NfoPath = ""
		item.Cover = ""
		item.Fanart = ""
		item.Poster = ""
		item.FanartList = ""
		item.Subtitles = ""
		item.Tags = ""
		item.Title = ""

		// 识别strm
		if contains(strmExts, ext) {
			item.StrmPath = path
		}
		// 识别nfo并解析
		if contains(nfoExts, ext) {
			item.NfoPath = path
			if meta, err := parseNfo(path); err == nil {
				item.Title = meta.Title
				item.Tags = meta.Tag
			}
		}
		// 识别字幕
		if contains(subExts, ext) {
			item.Subtitles = appendPath(item.Subtitles, path)
		}
		// 识别封面
		if contains(coverNames, d.Name()) {
			item.Cover = path
		}
		// 识别fanart
		if contains(fanartNames, d.Name()) {
			item.Fanart = path
		}
		// 识别poster
		if strings.Contains(strings.ToLower(d.Name()), "poster") {
			item.Poster = path
		}
		// 识别fanart列表
		if strings.Contains(strings.ToLower(d.Name()), "fanart") {
			item.FanartList = appendPath(item.FanartList, path)
		}
		// 仅当有strm或nfo时写入
		if item.StrmPath != "" || item.NfoPath != "" {
			db.GetDb().Where(model.MediaItem{StrmPath: item.StrmPath, NfoPath: item.NfoPath}).Assign(item).FirstOrCreate(&item)
		}
		return nil
	})
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func appendPath(orig, add string) string {
	if orig == "" {
		return add
	}
	return orig + "," + add
}

func parseNfo(path string) (NfoMeta, error) {
	var meta NfoMeta
	f, err := os.Open(path)
	if err != nil {
		return meta, err
	}
	defer f.Close()
	if err := xml.NewDecoder(f).Decode(&meta); err != nil {
		return meta, err
	}
	return meta, nil
} 