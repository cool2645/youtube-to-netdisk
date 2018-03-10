package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"strconv"
	"github.com/yanzay/log"
)

type Subscriber struct {
	ID        uint   `gorm:"AUTO_INCREMENT"`
	User      string
	Platform  string
	CreatedAt time.Time
}

func ListSubscribers(db *gorm.DB) (chats []int64, err error) {
	subscribers, err := GetSubscribers(db)
	if err != nil {
		return
	}
	for _, v := range subscribers {
		chatID, err := strconv.ParseInt(v.User, 10, 64)
		if err != nil {
			log.Error(err)
		}
		chats = append(chats, chatID)
	}
	return
}

func GetSubscribers(db *gorm.DB) (subscribers []Subscriber, err error) {
	err = db.Find(&subscribers).Error
	if err != nil {
		err = errors.Wrap(err, "GetSubscribers")
		return
	}
	return
}

func SaveSubscriber(db *gorm.DB, chatID int64) (newSubscriber Subscriber, err error) {
	var count int
	err = db.Model(&Subscriber{}).Where("platform = ?", "Telegram").
		Where("user = ?", strconv.FormatInt(chatID, 10)).Count(&count).Error
	if count == 0 {
		var subscriber Subscriber
		subscriber.User = strconv.FormatInt(chatID, 10)
		subscriber.Platform = "Telegram"
		newSubscriber, err = CreateSubscriber(db, subscriber)
	}
	return
}

func CreateSubscriber(db *gorm.DB, subscriber Subscriber) (newSubscriber Subscriber, err error) {
	err = db.Create(&subscriber).Error
	if err != nil {
		err = errors.Wrap(err, "SaveSubscriber")
		return
	}
	newSubscriber = subscriber
	return
}

func RemoveSubscriber(db *gorm.DB, chatID int64) (err error) {
	err = db.Where("platform = ?", "Telegram").
		Where("user = ?", strconv.FormatInt(chatID, 10)).
		Delete(Subscriber{}).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveSubscriber")
		return
	}
	return
}
