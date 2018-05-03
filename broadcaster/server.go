package broadcaster

import (
	"time"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yanzay/log"
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"github.com/rikakomoe/ritorudemonriri/ririsdk"
	"encoding/json"
)

type BroadcastMessage struct {
	Message string
	Level   int
}

const (
	Detailed  = iota
	Condensed
)

var tgSubscribedChats = make(map[int64]model.TGSubscriber)
var mux sync.RWMutex
var ch = make(chan BroadcastMessage)

func ServeTelegram(db *gorm.DB, addr string, key string) {
	log.Infof("Reading subscribed chats from database %s", time.Now())
	tgSubscribers, err := model.ListTGSubscribers(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range tgSubscribers {
		tgSubscribedChats[v.ChatID] = v
	}
	log.Warningf("%v %s", tgSubscribedChats, time.Now())
	log.Infof("Started serve telegram %s", time.Now())
	ririsdk.Init(addr, key, true)
	go pushMessage(ch)
	messages, _ := ririsdk.GetUpdatesChan()
	for message := range messages {
		b, err := json.Marshal(message)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Info(string(b))
		if message.Direction != ririsdk.IN {
			continue
		}
		switch message.Messenger {
		case ririsdk.Telegram:
			m := message.TelegramMessage
			if m.IsCommand() {
				switch m.Command() {
				case "carrier_subscribe":
					if m.CommandArguments() == "--condense" || m.CommandArguments() == "—condense" {
						replyMarkdownMessage(start(db, m, Condensed), m.Chat.ID)
					} else {
						replyMarkdownMessage(start(db, m, Detailed), m.Chat.ID)
					}
				case "carrier_unsubscribe":
					replyMarkdownMessage(stop(db, m), m.Chat.ID)
				}
			}
		}
	}
}

func replyMessage(text string, parseMode string, reqChatID int64) {
	msg := tg.NewMessage(reqChatID, text)
	msg.ParseMode = parseMode
	msg.DisableWebPagePreview = true
	ririsdk.PushMessage(ririsdk.Message{
		Direction:             ririsdk.OUT,
		Messenger:             ririsdk.Telegram,
		TelegramMessageConfig: &msg,
	})
}

func pushMessage(c chan BroadcastMessage) {
	var m BroadcastMessage
	for {
		m = <-c
		mux.RLock()
		for _, v := range tgSubscribedChats {
			if m.Level >= v.Level {
				replyMarkdownMessage(m.Message, v.ChatID)
			}
		}
		mux.RUnlock()
	}
}

func start(db *gorm.DB, m *tg.Message, level int) string {
	mux.Lock()
	defer mux.Unlock()
	tgSubscribedChats[m.Chat.ID] = model.TGSubscriber{ChatID: m.Chat.ID, Level: level}
	_, err := model.SaveTelegramSubscriber(db, m.Chat.ID, level)
	if err != nil {
		log.Fatal(err)
	}
	if level == Condensed {
		return "You have set up condensed subscription of yt2nd for this chat, pwp"
	}
	return "You have set up detailed subscription of yt2nd for this chat, pwp"
}

func stop(db *gorm.DB, m *tg.Message) string {
	mux.Lock()
	defer mux.Unlock()
	delete(tgSubscribedChats, m.Chat.ID)
	err := model.RemoveSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "Your subscription of yt2nd is suspended, qaq"
}

func replyMarkdownMessage(text string, reqChatID int64) {
	replyMessage(text, "Markdown", reqChatID)
}

func Broadcast(msg BroadcastMessage) {
	ch <- msg
	return
}
