package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"time"
)

type Keyword struct {
	ID        uint      `gorm:"AUTO_INCREMENT" json:"id"`
	Keyword   string    `json:"keyword"`
	CreatedAt time.Time `json:"created_at"`
}

func GetKeywords(db *gorm.DB) (keywords []Keyword, err error) {
	err = db.Find(&keywords).Error
	if err != nil {
		err = errors.Wrap(err, "GetKeywords")
		return
	}
	return
}

func SaveKeyword(db *gorm.DB, kw string) (keyword Keyword, err error) {
	var count int
	err = db.Model(&Keyword{}).Where("keyword = ?", kw).Count(&count).Error
	if count == 0 {
		keyword.Keyword = kw
		err = CreateKeyword(db, &keyword)
	}
	return
}

func CreateKeyword(db *gorm.DB, keyword *Keyword) (err error) {
	err = db.Create(keyword).Error
	if err != nil {
		err = errors.Wrap(err, "CreateKeyword")
		return
	}
	return
}

func RemoveKeyword(db *gorm.DB, kw string) (err error) {
	err = db.Where("keyword = ?", kw).Delete(Keyword{}).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveKeyword")
		return
	}
	return
}
