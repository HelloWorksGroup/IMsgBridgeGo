package route

// 本模块用于将QQ消息转发至KOOK，并将KOOK消息转发至QQ

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Nigh/MiraiGo-Template-Mod/bot"
)

type rt struct {
}

var instance *rt

type QQMsg struct {
	Type    int
	Content string
}

func init() {
	instance = &rt{}
	bot.RegisterModule(instance)
}

var validGroupId int64 = 0

func SetGroupID(n int64) {
	validGroupId = n
}

var msgRouteQQ2KOOK func(name string, msg []QQMsg)

func OnMsg(handler func(name string, msg []QQMsg)) {
	msgRouteQQ2KOOK = handler
}

func RouteKOOK2QQText(content string) {
	go func() {
		m := message.NewSendingMessage().Append(message.NewText(content))
		bot.Instance.SendGroupMessage(validGroupId, m)
	}()
}

func NewImageShare(url, title, image string) *(message.ServiceElement) {
	template := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?><msg flag="3" templateID="12345" action="web" brief="[KOOK图片] %s" serviceID="1" url="%s"><item layout="1"><title>%v</title><picture cover="%v"/></item><source/></msg>`,
		title, url, image, title)
	return &message.ServiceElement{
		Id:      1,
		Content: template,
		ResId:   url,
		SubType: "UrlShare",
	}
}

func RouteKOOK2QQImage(displayName string, imageUrl string, linkUrl string) {
	go func() {
		m := message.NewSendingMessage().Append(NewImageShare(linkUrl, "由 "+displayName+" 发送", imageUrl))
		bot.Instance.SendGroupMessage(validGroupId, m)
	}()
}

func (a *rt) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "kook.route",
		Instance: instance,
	}
}

func (a *rt) Init() {

}

func (a *rt) PostInit() {
}

func (a *rt) Serve(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		if msg.GroupCode != validGroupId {
			fmt.Println("[QQ]:", msg.GroupName, msg.Sender.Nickname+": "+msg.ToString())
			return
		}
		fmt.Println("[QQQQ]:", msg.Sender.Nickname+": "+msg.ToString())
		for _, elem := range msg.Elements {
			switch e := elem.(type) {
			case *message.GroupImageElement:
				fmt.Println("ImageURL=", e.Url)
			}
		}
		if msg.ToString() == "ping" {
			go func() {
				delay := rand.Intn(500) + rand.Intn(100) + rand.Intn(50) + rand.Intn(14)
				<-time.After(time.Millisecond * time.Duration(delay+2000))
				m := message.NewSendingMessage().Append(message.NewText("来自 \"QQ\" 的回复: 字节=256 时间=" + strconv.Itoa(delay) + "ms TTL=" + strconv.Itoa(61-rand.Intn(7))))
				c.SendGroupMessage(msg.GroupCode, m)
			}()
		} else {
			// DONE: 转发
			// fmt.Println("msgRouteQQ2KOOK", msg.Sender.Nickname, msg.ToString())
			go msgRouteQQ2KOOK(msg.Sender.Nickname, qqGroupMsgParse(msg))
		}
	})
}

func (a *rt) Start(bot *bot.Bot) {
}

func (a *rt) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func qqGroupMsgParse(msg *message.GroupMessage) (qqmsg []QQMsg) {
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			qqmsg = append(qqmsg, QQMsg{0, e.Content})
		case *message.GroupImageElement:
			qqmsg = append(qqmsg, QQMsg{1, e.Url})
		case *message.FaceElement:
			qqmsg = append(qqmsg, QQMsg{0, "[表情:" + e.Name + "]"})
		case *message.MarketFaceElement:
			qqmsg = append(qqmsg, QQMsg{0, "[商店表情:" + e.Name + "]"})
		case *message.AtElement:
			qqmsg = append(qqmsg, QQMsg{0, "[" + e.Display + "]"})
		case *message.RedBagElement:
			qqmsg = append(qqmsg, QQMsg{0, "[红包:" + e.Title + "]"})
		case *message.ReplyElement:
			qqmsg = append(qqmsg, QQMsg{0, "[回复:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"})
		}
	}
	return
}
