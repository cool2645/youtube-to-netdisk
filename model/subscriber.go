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
	ID        uint `gorm:"AUTO_INCREMENT"`
	User      string
	Platform  string
	Level     int
	CreatedAt time.Time
}

type TGSubscriber struct {
	ChatID int64
	Level  int
}

func ListTGSubscribers(db *gorm.DB) (tgSubscribers []TGSubscriber, err error) {
	subscribers, err := GetSubscribers(db)
	if err != nil {
		return
	}
	for _, v := range subscribers {
		chatID, err := strconv.ParseInt(v.User, 10, 64)
		if err != nil {
			log.Error(err)
		}
		tgSubscribers = append(tgSubscribers, TGSubscriber{ChatID: chatID, Level: v.Level})
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

func SaveTelegramSubscriber(db *gorm.DB, chatID int64, level int) (newSubscriber Subscriber, err error) {
	var subscriber Subscriber
	err = db.Where("platform = ?", "Telegram").
		Where("user = ?", strconv.FormatInt(chatID, 10)).First(&subscriber).Error
	if err == nil {
		subscriber.Level = level
		newSubscriber, err = UpdateSubscriber(db, subscriber)
	} else if err.Error() == "record not found" {
		subscriber.User = strconv.FormatInt(chatID, 10)
		subscriber.Platform = "Telegram"
		subscriber.Level = level
		newSubscriber, err = CreateSubscriber(db, subscriber)
	}
	return
}

func UpdateSubscriber(db *gorm.DB, subscriber Subscriber) (newSubscriber Subscriber, err error) {
	db.Save(&subscriber)
	if err != nil {
		err = errors.Wrap(err, "UpdateSubscriber")
		return
	}
	newSubscriber = subscriber
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
