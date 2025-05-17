package handles

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/alist-org/alist/v3/internal/media"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/db"
	"github.com/gin-gonic/gin"
)

type SyncMediaReq struct {
	Path string `json:"path" binding:"required"`
}

// POST /api/media/sync
func SyncMediaLibraryHandler(c *gin.Context) {
	var req SyncMediaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := media.SyncMediaLibrary(req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "同步完成"})
}

// GET /api/media/list?page=1&page_size=20
func ListMediaHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	dbConn := db.GetDb()
	var total int64
	var items []model.MediaItem
	dbConn.Model(&model.MediaItem{}).Count(&total)
	dbConn.Order("created_at desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&items)
	// 仅返回前4个标签
	for i := range items {
		tags := items[i].Tags
		if tags != "" {
			tagArr := splitTags(tags)
			if len(tagArr) > 4 {
				tagArr = tagArr[:4]
			}
			items[i].Tags = joinTags(tagArr)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"list": items,
	})
}

// GET /api/media/detail/:id
func MediaDetailHandler(c *gin.Context) {
	id := c.Param("id")
	var item model.MediaItem
	dbConn := db.GetDb()
	if err := dbConn.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该媒体"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func splitTags(tags string) []string {
	var res []string
	for _, t := range strings.Split(tags, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}

func joinTags(tags []string) string {
	return strings.Join(tags, ",")
} 