package services

import (
	"contentsystem/internal/dao"
	"contentsystem/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ContentUpdateReq struct {
	ID             int           `json:"id" binding:"required"`
	Title          string        `json:"title"`
	VideoURL       string        `json:"video_url"`
	Author         string        `json:"author"`
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

type ContentUpdateRsp struct {
	Message string `json:"message" binding:"required"`
}

func (c *CmsApp) ContentUpdate(ctx *gin.Context) {
	var req ContentUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentDao := dao.NewContentDao(c.db)
	ok, err := contentDao.IsExist(req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "内容不存在"})
		return
	}
	if err := contentDao.Update(req.ID, model.ContentDetail{
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
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &ContentUpdateRsp{
			Message: fmt.Sprintf("ok"),
		},
	})
}
