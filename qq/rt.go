package route

// 本模块用于将QQ消息转发至KOOK，并将KOOK消息转发至QQ

import (
	"bytes"
	"strconv"
	"sync"

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

var externMsgHandler func(msg *message.GroupMessage)

func OnMsg(handler func(msg *message.GroupMessage)) {
	externMsgHandler = handler
}

func RouteKOOK2QQText(content string) {
	go func() {
		m := message.NewSendingMessage().Append(message.NewText(content))
		bot.Instance.SendGroupMessage(validGroupId, m)
	}()
}

func SendToQQGroupEx(e []message.IMessageElement, groupId int64) int32 {
	m := message.NewSendingMessage()
	for _, v := range e {
		m.Append(v)
	}
	ret := bot.Instance.SendGroupMessage(groupId, m)
	return ret.Id
}
func SendToQQGroup(content string, groupId int64) int32 {
	m := message.NewSendingMessage().Append(message.NewText(content))
	ret := bot.Instance.SendGroupMessage(groupId, m)
	return ret.Id
}

func UploadImgToQQGroup(img []byte, groupId int64) (msg message.IMessageElement, err error) {
	return bot.Instance.UploadImage(message.Source{
		PrimaryID:  groupId,
		SourceType: message.SourceGroup,
	}, bytes.NewReader(img))
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
		externMsgHandler(msg)
	})
}

func (a *rt) Start(bot *bot.Bot) {
}

func (a *rt) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func GroupMsgParse(msg *message.GroupMessage) (qqmsg []QQMsg) {
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
			if e.Target != bot.Instance.Uin {
				qqmsg = append(qqmsg, QQMsg{2, "[" + e.Display + "]"})
			}
		case *message.RedBagElement:
			qqmsg = append(qqmsg, QQMsg{0, "[红包:" + e.Title + "]"})
		case *message.ReplyElement:
			qqmsg = append(qqmsg, QQMsg{3, "[回复:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"})
		default:
			qqmsg = append(qqmsg, QQMsg{4, "[无法转发的消息类型]"})
		}
	}
	return
}

func GroupMsg2Markdown(msg *message.GroupMessage) (qqmsg string) {
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			qqmsg += e.Content + "\n"
		case *message.GroupImageElement:
			qqmsg += "![](" + e.Url + ")\n"
		case *message.FaceElement:
			qqmsg += "[表情:" + e.Name + "]\n"
		case *message.MarketFaceElement:
			qqmsg += "[商店表情:" + e.Name + "]\n"
		case *message.AtElement:
			if e.Target != bot.Instance.Uin {
				qqmsg += "[" + e.Display + "]"
			}
		case *message.RedBagElement:
			qqmsg += "[红包:" + e.Title + "]\n"
		case *message.ReplyElement:
			qqmsg += "[回复:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"
		default:
			qqmsg += "[无法转发的消息类型]\n"
		}
	}
	return
}
