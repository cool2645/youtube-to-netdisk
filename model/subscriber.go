package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/yanzay/log"
	"strconv"
	"time"
)

type Subscriber struct {
	ID          uint `gorm:"AUTO_INCREMENT"`
	User        string
	Platform    string
	Level       int
	MessageType string
	CreatedAt   time.Time
}

type TGSubscriber struct {
	ChatID int64
	Level  int
}

type QQSubscriber struct {
	ChatID      int64
	MessageType string
	Level       int
}

func ListTGSubscribers(db *gorm.DB) (tgSubscribers []TGSubscriber, err error) {
	subscribers, err := GetSubscribers(db)
	if err != nil {
		return
	}
	for _, v := range subscribers {
		if v.Platform != "Telegram" {
			continue
		}
		chatID, err := strconv.ParseInt(v.User, 10, 64)
		if err != nil {
			log.Error(err)
		}
		tgSubscribers = append(tgSubscribers, TGSubscriber{ChatID: chatID, Level: v.Level})
	}
	return
}

func ListQQSubscribers(db *gorm.DB) (qqSubscribers []QQSubscriber, err error) {
	subscribers, err := GetSubscribers(db)
	if err != nil {
		return
	}
	for _, v := range subscribers {
		if v.Platform != "QQ" {
			continue
		}
		chatID, err := strconv.ParseInt(v.User, 10, 64)
		if err != nil {
			log.Error(err)
		}
		qqSubscribers = append(qqSubscribers, QQSubscriber{ChatID: chatID, Level: v.Level, MessageType: v.MessageType})
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

func SaveQQSubscriber(db *gorm.DB, chatID int64, messageType string, level int) (newSubscriber Subscriber, err error) {
	var subscriber Subscriber
	err = db.Where("platform = ?", "QQ").
		Where("user = ?", strconv.FormatInt(chatID, 10)).
		Where("message_type = ?", messageType).First(&subscriber).Error
	if err == nil {
		subscriber.Level = level
		newSubscriber, err = UpdateSubscriber(db, subscriber)
	} else if err.Error() == "record not found" {
		subscriber.User = strconv.FormatInt(chatID, 10)
		subscriber.MessageType = messageType
		subscriber.Platform = "QQ"
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

func RemoveTGSubscriber(db *gorm.DB, chatID int64) (err error) {
	err = db.Where("platform = ?", "Telegram").
		Where("user = ?", strconv.FormatInt(chatID, 10)).
		Delete(Subscriber{}).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveTGSubscriber")
		return
	}
	return
}

func RemoveQQSubscriber(db *gorm.DB, chatID int64, messageType string) (err error) {
	err = db.Where("platform = ?", "QQ").
		Where("user = ?", strconv.FormatInt(chatID, 10)).
		Where("message_type = ?", messageType).
		Delete(Subscriber{}).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveQQSubscriber")
		return
	}
	return
}
