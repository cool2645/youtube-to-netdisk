package broadcaster

import (
	"fmt"
	"github.com/catsworld/qq-bot-api/cqcode"
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rikakomoe/ritorudemonriri/bot"
	qq "github.com/rikakomoe/ritorudemonriri/qq-bot-api"
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
var bot botapi.Bot

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
	bot = botapi.NewBot()
	go pushMessage(ch)
	cqcode.StrictCommand = true
	updates := bot.GetUpdatesChan(botapi.UpdateConfig{})
	for update := range updates {
		if !update.IsCommand() {
			continue
		}
		cmd, args := update.Command()
		dispatchCmd(db, cmd, update, args)
	}
}

var cmdMap = map[string]func(db *gorm.DB, update botapi.Update, args []string) string{
	"carrier_subscribe#yt2nd":   carrierSubscribe,
	"carrier_unsubscribe#yt2nd": carrierUnsubscribe,
	"ping#yt2nd":                ping,
	"status#yt2nd":              ping,
	"start#yt2nd":               help,
	"help#yt2nd":                help,
}

func carrierUnsubscribe(db *gorm.DB, update botapi.Update, args []string) string {
	if update.Messenger == ririsdk.Telegram {
		return tgStop(db, update.TGUpdate.Message)
	} else {
		return qqStop(db, update.QQUpdate.Message)
	}
}

func carrierSubscribe(db *gorm.DB, update botapi.Update, args []string) string {
	level := Detailed
	if len(args) > 0 && (args[0] == "--condense" || args[0] == "—condense") {
		level = Condensed
	}
	if update.Messenger == ririsdk.Telegram {
		return tgStart(db, update.TGUpdate.Message, level)
	} else {
		return qqStart(db, update.QQUpdate.Message, level)
	}
}

func dispatchCmd(db *gorm.DB, cmd string, update botapi.Update, args []string) {
	if len(cmd) < len("#yt2nd") || cmd[len(cmd)-len("#yt2nd"):] != "#yt2nd" {
		cmd = cmd + "#yt2nd"
	}
	c, ok := cmdMap[cmd]
	if !ok {
		return
	}
	text := c(db, update, args)
	mc := botapi.MessageConfig{}
	mc.DisableWebPreview = true
	update.Reply(mc, text)
}

func tgReplyMarkdownMessage(text string, reqChatID int64) {
	sc := botapi.NewTGSendConfig(reqChatID)
	sc.ParseMode = botapi.FormatMarkdown
	sc.DisableWebPreview = true
	bot.Send(sc, text)
}

func qqSendMessage(text string, reqChatID int64, reqChatType string) {
	sc := botapi.NewQQSendConfig(reqChatID, reqChatType)
	bot.Send(sc, text)
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
	keyStr := m.Chat.Type + strconv.FormatInt(m.Chat.ID, 10)
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

func help(db *gorm.DB, update botapi.Update, args []string) string {
	return "虹咲搬运机器人：\n/carrier_subscribe#yt2nd - 订阅搬运机器人的通知（详细）\n/carrier_subscribe#yt2nd --condense - 订阅搬运机器人的通知（精简）\n" +
		"/carrier_unsubscribe#yt2nd - 退订搬运机器人的通知\n（#号及后面的部分可以省略）"
}

func ping(db *gorm.DB, update botapi.Update, args []string) string {
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
