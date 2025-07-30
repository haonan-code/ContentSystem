package services

import (
	"contentsystem/internal/dao"
	"contentsystem/internal/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	goflow "github.com/s8sg/goflow/v1"
	"net/http"
	"time"
)

type ContentCreateReq struct {
	Title          string        `json:"title" binding:"required"`
	VideoURL       string        `json:"video_url" binding:"required"`
	Author         string        `json:"author" binding:"required"`
	Description    string        `json:"description"`
	Thumbnail      string        `json:"thumbnail"`
	Category       string        `json:"category"`
	Duration       time.Duration `json:"duration"`
	Resolution     string        `json:"resolution"`
	FileSize       int64         `json:"fileSize"`
	Format         string        `json:"format"`
	Quality        int           `json:"quality"`
	ApprovalStatus int           `json:"approval_status"`
	UpdatedAt      time.Time     `json:"updated_at"`
	CreatedAt      time.Time     `json:"created_at"`
}

type ContentCreateRsp struct {
	Message string `json:"message" binding:"required"`
}

func (c *CmsApp) ContentCreate(ctx *gin.Context) {
	var req ContentCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentDao := dao.NewContentDao(c.db)
	id, err := contentDao.Create(model.ContentDetail{
		Title:          req.Title,
		Description:    req.Description,
		Author:         req.Author,
		VideoURL:       req.VideoURL,
		Thumbnail:      req.Thumbnail,
		Category:       req.Category,
		Duration:       req.Duration,
		Resolution:     req.Resolution,
		FileSize:       req.FileSize,
		Format:         req.Format,
		Quality:        req.Quality,
		ApprovalStatus: req.ApprovalStatus,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flowData := map[string]interface{}{
		"content_id": id,
	}
	data, _ := json.Marshal(flowData)
	if err := c.flowService.Execute("content-flow", &goflow.Request{
		Body: data,
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &ContentCreateRsp{
			Message: fmt.Sprintf("ok"),
		},
	})
}
