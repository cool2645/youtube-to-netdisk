package broadcaster

import (
	"time"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yanzay/log"
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

var subscribedChats = make(map[int64]int64)
var mux sync.RWMutex
var ch = make(chan string)
var bot *tg.BotAPI

func ServeTelegram(db *gorm.DB, apiKey string) {
	log.Infof("Reading subscribed chats from database %s", time.Now())
	subscribers, err := model.ListSubscribers(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range subscribers {
		subscribedChats[v] = v
	}
	log.Warningf("%v %s", subscribedChats, time.Now())
	log.Infof("Started serve telegram %s", time.Now())
	bot, err = tg.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	go pushMessage(ch)
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		m := update.Message
		if m.IsCommand() {
			switch m.Command() {
			case "carrier_subscribe":
				ReplyMessage(start(db, m), m.Chat.ID)
			case "carrier_unsubscribe":
				ReplyMessage(stop(db, m), m.Chat.ID)
			}
		}
	}
}

func pushMessage(c chan string) {
	var m string
	for {
		m = <-c
		mux.RLock()
		for _, v := range subscribedChats {
			msg := tg.NewMessage(v, m)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
		mux.RUnlock()
	}
}

func start(db *gorm.DB, m *tg.Message) string {
	mux.Lock()
	defer mux.Unlock()
	subscribedChats[m.Chat.ID] = m.Chat.ID
	_, err := model.SaveSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "You have set up subscription of yt2nd for this chat, pwp"
}

func stop(db *gorm.DB, m *tg.Message) string {
	mux.Lock()
	defer mux.Unlock()
	delete(subscribedChats, m.Chat.ID)
	err := model.RemoveSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "Your subscription of yt2nd is suspended, qaq"
}

func ReplyMessage(text string, reqChatID int64) {
	msg := tg.NewMessage(reqChatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func Broadcast(msg string) {
	ch <- msg
	return
}
