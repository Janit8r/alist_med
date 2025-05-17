package model

import "gorm.io/gorm"

type MediaItem struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Title      string         `json:"title" gorm:"index"`
	Cover      string         `json:"cover"`         // 封面图片路径
	Fanart     string         `json:"fanart"`        // 海报图片路径
	StrmPath   string         `json:"strm_path"`     // 视频播放文件路径
	NfoPath    string         `json:"nfo_path"`      // nfo元数据文件路径
	Tags       string         `json:"tags"`          // 逗号分隔标签
	Subtitles  string         `json:"subtitles"`     // 逗号分隔字幕文件路径
	Poster     string         `json:"poster"`        // 备用海报
	FanartList string         `json:"fanart_list"`   // 备用fanart，逗号分隔
	CreatedAt  int64          `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt  int64          `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
} 