package broadcaster

import (
	"time"
	"github.com/yanzay/log"
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/juzi5201314/cqhttp-go-sdk/cqcode"
	"fmt"
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
var tgMux sync.RWMutex
var ch = make(chan BroadcastMessage)
var bot *tg.BotAPI

func ServeTelegram(db *gorm.DB, apiKey string) {
	log.Infof("Reading subscribed chats from database %s", time.Now())
	tgSubscribers, err := model.ListTGSubscribers(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range tgSubscribers {
		tgSubscribedChats[v.ChatID] = v
	}
	log.Warningf("%v %s", tgSubscribedChats, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bot, err = tg.NewBotAPI(apiKey)
	log.Infof("Started serve telegram %s", time.Now())
	go pushMessage(ch)
	cqcode.StrictCommand = true
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
				if m.CommandArguments() == "--condense" || m.CommandArguments() == "—condense" {
					tgReplyMessage(tgStart(db, m, Condensed), m.Chat.ID)
				} else {
					tgReplyMessage(tgStart(db, m, Detailed), m.Chat.ID)
				}
			case "carrier_unsubscribe":
				tgReplyMessage(tgStop(db, m), m.Chat.ID)
			case "help":
				tgReplyMessage(help(), m.Chat.ID)
			case "start":
				tgReplyMessage(help(), m.Chat.ID)
			case "ping":
				tgReplyMessage(ping(db), m.Chat.ID)
			}
		}
	}
}

func tgReplyMarkdownMessage(text string, reqChatID int64) {
	msg := tg.NewMessage(reqChatID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}

func tgReplyMessage(text string, reqChatID int64) {
	msg := tg.NewMessage(reqChatID, text)
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}

func pushMessage(c chan BroadcastMessage) {
	var m BroadcastMessage
	for {
		m = <-c
		tgMux.RLock()
		for _, v := range tgSubscribedChats {
			if m.Level >= v.Level {
				tgReplyMarkdownMessage(m.Message, v.ChatID)
			}
		}
		tgMux.RUnlock()
	}
}

func tgStart(db *gorm.DB, m *tg.Message, level int) string {
	tgMux.Lock()
	defer tgMux.Unlock()
	tgSubscribedChats[m.Chat.ID] = model.TGSubscriber{ChatID: m.Chat.ID, Level: level}
	_, err := model.SaveTelegramSubscriber(db, m.Chat.ID, level)
	if err != nil {
		log.Fatal(err)
	}
	if level == Condensed {
		return "您已在此会话中订阅精简通知，pwp"
	}
	return "您已在此会话中订阅详细通知，pwp"
}

func tgStop(db *gorm.DB, m *tg.Message) string {
	tgMux.Lock()
	defer tgMux.Unlock()
	delete(tgSubscribedChats, m.Chat.ID)
	err := model.RemoveTGSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "您的订阅已取消，qaq"
}

func help() string {
	return "/carrier_subscribe - 订阅搬运机器人的通知（详细）\n/carrier_subscribe --condense - 订阅搬运机器人的通知（精简）\n" +
		"/carrier_unsubscribe - 退订搬运机器人的通知\n/help - 显示此帮助\n/ping - 测试是否在线"
}

func ping(db *gorm.DB) string {
	keywords, err := model.GetKeywords(db)
	if err != nil {
		log.Error(err)
		return "搬运机器人监听关键字读取失败"
	}
	kwds := make([]string, 0)
	for _, keyword := range keywords {
		kwds = append(kwds, keyword.Keyword)
	}
	return fmt.Sprintf("搬运机器人工作正常，监听关键字：%v", kwds)
}

func Broadcast(msg BroadcastMessage) {
	ch <- msg
	return
}
