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

func init() {
	instance = &rt{}
	bot.RegisterModule(instance)
}

var validGroupId int64 = 0

func SetGroupID(n int64) {
	validGroupId = n
}

var msgRouteQQ2KOOK func(name string, msg string)

func OnMsg(handler func(name string, msg string)) {
	msgRouteQQ2KOOK = handler
}

func MsgRouteKOOK2QQ(content string) {
	go func() {
		m := message.NewSendingMessage().Append(message.NewText(content))
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

func qqGroupMsgParse(msg *message.GroupMessage) (qqMsgStr string) {
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			qqMsgStr += e.Content
		case *message.GroupImageElement:
			qqMsgStr += "\n" + e.Url + "\n"
		case *message.FaceElement:
			qqMsgStr += "[表情:" + e.Name + "]"
		case *message.MarketFaceElement:
			qqMsgStr += "[商店表情:" + e.Name + "]"
		case *message.AtElement:
			qqMsgStr += "[" + e.Display + "]"
		case *message.RedBagElement:
			qqMsgStr += "[红包:" + e.Title + "]"
		case *message.ReplyElement:
			qqMsgStr += "[回复:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"
		}
	}
	return
}
