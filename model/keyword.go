package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type Keyword struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Keyword   string
	CreatedAt time.Time
}


func GetKeywords(db *gorm.DB) (keywords []Keyword, err error) {
	err = db.Find(&keywords).Error
	if err != nil {
		err = errors.Wrap(err, "GetKeywords")
		return
	}
	return
}

func SaveKeyword(db *gorm.DB, kw string) (newKeyword Keyword, err error) {
	var count int
	err = db.Model(&Keyword{}).Where("keyword = ?", kw).Count(&count).Error
	if count == 0 {
		var keyword Keyword
		keyword.Keyword = kw
		newKeyword, err = CreateKeyword(db, keyword)
	}
	return
}

func CreateKeyword(db *gorm.DB, keyword Keyword) (newKeyword Keyword, err error) {
	err = db.Create(&keyword).Error
	if err != nil {
		err = errors.Wrap(err, "CreateKeyword")
		return
	}
	newKeyword = keyword
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
