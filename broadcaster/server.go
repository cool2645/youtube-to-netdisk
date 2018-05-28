package broadcaster

import (
	"fmt"
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	qq "github.com/rikakomoe/ritorudemonriri/qq-bot-api"
	"github.com/catsworld/qq-bot-api/cqcode"
	"github.com/rikakomoe/ritorudemonriri/ririsdk"
	tg "github.com/rikakomoe/ritorudemonriri/telegram-bot-api"
	"github.com/yanzay/log"
	"strconv"
	"sync"
	"time"
)

type BroadcastMessage struct {
	Message string
	Level   int
}

const (
	Detailed = iota
	Condensed
)

var tgSubscribedChats = make(map[int64]model.TGSubscriber)
var qqSubscribedChats = make(map[string]model.QQSubscriber)
var tgMux sync.RWMutex
var qqMux sync.RWMutex
var ch = make(chan BroadcastMessage)
var tgBot, _ = tg.NewBotAPI()
var qqBot, _ = qq.NewBotAPI()

func ServeRiri(db *gorm.DB, addr string, key string) {
	log.Infof("Reading subscribed chats from database %s", time.Now())
	tgSubscribers, err := model.ListTGSubscribers(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range tgSubscribers {
		tgSubscribedChats[v.ChatID] = v
	}
	log.Warningf("%v %s", tgSubscribedChats, time.Now())
	qqSubscribers, err := model.ListQQSubscribers(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range qqSubscribers {
		keyStr := v.MessageType + strconv.FormatInt(v.ChatID, 10)
		qqSubscribedChats[keyStr] = v
	}
	log.Warningf("%v %s", qqSubscribedChats, time.Now())
	log.Infof("Started serve telegram %s", time.Now())
	ririsdk.Init(addr, key, true)
	go pushMessage(ch)
	cqcode.StrictCommand = true
	u := tg.NewUpdate(0)
	u2 := qq.NewWebhook("")
	updates, _ := ririsdk.GetUpdatesChan(0)
	for update := range updates {
		switch update.Messenger {
		case ririsdk.Telegram:
			tgUpdates, err := tgBot.GetUpdates(&u, update)
			if err != nil {
				continue
			}
			for _, tgUpdate := range tgUpdates {
				if tgUpdate.Message == nil {
					continue
				}
				m := tgUpdate.Message
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
		case ririsdk.CQHttp:
			cqUpdate, err := qqBot.GetUpdate(u2, update)
			if err != nil {
				continue
			}
			switch cqUpdate.PostType {
			case "message":
				m := cqUpdate.Message
				if err != nil {
					continue
				}
				if m.IsCommand() {
					cmd, args := m.Command()
					switch cmd {
					case "carrier_subscribe":
						if len(args) > 0 && (args[0] == "--condense" || args[0] == "—condense") {
							qqSendMessage(qqStart(db, m, Condensed), m.Chat.ID, m.Chat.Type)
						} else {
							qqSendMessage(qqStart(db, m, Detailed), m.Chat.ID, m.Chat.Type)
						}
					case "carrier_unsubscribe":
						qqSendMessage(qqStop(db, m), m.Chat.ID, m.Chat.Type)
					case "help":
						qqSendMessage(help(), m.Chat.ID, m.Chat.Type)
					case "start":
						qqSendMessage(help(), m.Chat.ID, m.Chat.Type)
					case "ping":
						qqSendMessage(ping(db), m.Chat.ID, m.Chat.Type)
					}
				}
			}
		}
	}
}

func tgReplyMarkdownMessage(text string, reqChatID int64) {
	msg := tg.NewMessage(reqChatID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	tgBot.Send(msg)
}

func tgReplyMessage(text string, reqChatID int64) {
	msg := tg.NewMessage(reqChatID, text)
	msg.DisableWebPagePreview = true
	tgBot.Send(msg)
}

func qqSendMessage(text string, reqChatID int64, reqChatType string) {
	qqBot.SendMessage(reqChatID, reqChatType, text)
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
		qqMux.RLock()
		for _, v := range qqSubscribedChats {
			if m.Level >= v.Level {
				qqSendMessage(m.Message, v.ChatID, v.MessageType)
			}
		}
		qqMux.RUnlock()
	}
}

func qqStart(db *gorm.DB, m *qq.Message, level int) string {
	qqMux.Lock()
	defer qqMux.Unlock()
	keyStr := m.Chat.Type + strconv.FormatInt(m.Chat.ID, 10)
	qqSubscribedChats[keyStr] = model.QQSubscriber{ChatID: m.Chat.ID, Level: level, MessageType: m.Chat.Type}
	_, err := model.SaveQQSubscriber(db, m.Chat.ID, m.Chat.Type, level)
	if err != nil {
		log.Fatal(err)
	}
	if level == Condensed {
		return "您已在此会话中订阅精简通知，pwp"
	}
	return "您已在此会话中订阅详细通知，pwp"
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

func qqStop(db *gorm.DB, m *qq.Message) string {
	qqMux.Lock()
	defer qqMux.Unlock()
	keyStr := m.Chat.Type + strconv.FormatInt(m.Chat.ID, 64)
	delete(qqSubscribedChats, keyStr)
	err := model.RemoveQQSubscriber(db, m.Chat.ID, m.Chat.Type)
	if err != nil {
		log.Fatal(err)
	}
	return "您的订阅已取消，qaq"
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
