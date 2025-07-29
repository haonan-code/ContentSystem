package dao

import (
	"contentsystem/internal/model"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type ContentDao struct {
	db *gorm.DB
}

func NewContentDao(db *gorm.DB) *ContentDao {
	return &ContentDao{db: db}
}

func (c *ContentDao) Create(detail model.ContentDetail) error {
	if err := c.db.Create(&detail).Error; err != nil {
		log.Printf("content create error =  %v", err)
		return err
	}
	return nil
}

func (c *ContentDao) IsExist(contentID int) (bool, error) {
	var detail model.ContentDetail
	err := c.db.Where("id = ?", contentID).First(&detail).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		fmt.Printf("ContentDao isExist = [%v]", err)
		return false, err
	}
	return true, nil
}

func (c *ContentDao) Update(id int, detail model.ContentDetail) error {
	if err := c.db.Where("id = ?", id).Updates(&detail).Error; err != nil {
		log.Printf("content update error =  %v", err)
		return err
	}
	return nil
}

func (c *ContentDao) Delete(id int) error {
	if err := c.db.Where("id = ?", id).Delete(&model.ContentDetail{}).Error; err != nil {
		log.Printf("content delete error =  %v", err)
		return err
	}
	return nil
}

type FindParams struct {
	ID       int
	Title    string
	Author   string
	Page     int
	PageSize int
}

func (c *ContentDao) Find(params *FindParams) ([]*model.ContentDetail, int64, error) {
	// 构造查询条件
	query := c.db.Model(&model.ContentDetail{})
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}
	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}
	// 总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var page, pageSize = 1, 10
	if params.Page > 0 {
		page = params.Page
	}
	if params.PageSize > 0 {
		pageSize = params.PageSize
	}

	offset := (page - 1) * pageSize
	var data []*model.ContentDetail
	if err := query.Offset(offset).Limit(pageSize).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
