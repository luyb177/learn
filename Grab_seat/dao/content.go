package dao

import (
	"gorm.io/gorm"
	"learn/Grab_seat/model"
	"log"
)

type ContentDAO interface {
	AddContent(content *model.Content) error
	FindContent(pn int) ([]model.Content, error)
}

type ContentDAOImpl struct {
	db *gorm.DB
}

func NewContentDAOImpl(db *gorm.DB) *ContentDAOImpl {
	return &ContentDAOImpl{
		db: db,
	}
}

// AddContent 添加信息
func (dao *ContentDAOImpl) AddContent(content *model.Content) error {
	result := dao.db.Create(content)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}

// FindContent 查询消息分页
func (dao *ContentDAOImpl) FindContent(pn int) ([]model.Content, error) {
	var contents []model.Content
	result := dao.db.Limit(10).Offset((pn - 1) * 10).Order("id desc").Find(&contents)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	return contents, nil
}
