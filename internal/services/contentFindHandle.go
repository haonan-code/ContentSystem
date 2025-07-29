package services

import (
	"contentsystem/internal/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Content struct {
	ID             int
	Title          string
	Description    string
	Author         string
	VideoURL       string
	Thumbnail      string
	Category       string
	Duration       time.Duration
	Resolution     string
	FileSize       int64
	Format         string
	Quality        int
	ApprovalStatus int
}

type ContentFindReq struct {
	ID       int    `json:"id"`        // 内容ID
	Title    string `json:"title"`     // 标题
	Author   string `json:"author"`    // 作者
	Page     int    `json:"page"`      // 页
	PageSize int    `json:"page_size"` // 页大小
}

type ContentFindRsp struct {
	Message  string    `json:"message"`
	Contents []Content `json:"contents"`
	Total    int64     `json:"total"`
}

func (c *CmsApp) ContentFind(ctx *gin.Context) {
	var req ContentFindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentDao := dao.NewContentDao(c.db)
	contentList, total, err := contentDao.Find(&dao.FindParams{
		ID:       req.ID,
		Title:    req.Title,
		Author:   req.Author,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contents := make([]Content, 0, len(contentList))
	log.Println("len(contentList)为：", len(contentList))
	for _, content := range contentList {
		contents = append(contents, Content{
			ID:             content.ID,
			Title:          content.Title,
			Description:    content.Description,
			Author:         content.Author,
			VideoURL:       content.VideoURL,
			Thumbnail:      content.Thumbnail,
			Category:       content.Category,
			Duration:       content.Duration,
			Resolution:     content.Resolution,
			FileSize:       content.FileSize,
			Format:         content.Format,
			Quality:        content.Quality,
			ApprovalStatus: content.ApprovalStatus,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &ContentFindRsp{
			Message:  fmt.Sprintf("ok"),
			Contents: contents,
			Total:    total,
		},
	})
}
