package kookNode

import (
	"github.com/phuslu/log"

	"github.com/HelloWorksGroup/IMsgBridgeGo/imnode"
	"github.com/lonelyevil/kook"
	"github.com/lonelyevil/kook/log_adapter/plog"
)

type node struct {
	logger       *log.Logger
	stdioChan    string
	validGroupID []string
	msgHandler   func(gid string, msg *imnode.IMMsg)

	botID   string
	token   string
	session *kook.Session
}

func Setup(setting map[string]any, logger *log.Logger) imnode.IMNode {
	n := new(node)
	n.logger = logger
	n.token = setting["token"].(string)
	n.validGroupID = imnode.ConvertSettingGroups2StrSlice(setting["groups"])
	return n
}

func (n *node) GroupIDValid(gid string) bool {
	for _, id := range n.validGroupID {
		if id == gid {
			return true
		}
	}
	return false
}
func (n *node) Start() error {
	n.session = kook.New(n.token, plog.NewLogger(n.logger))
	me, _ := n.session.UserMe()
	n.botID = me.ID
	n.session.AddHandler(n.markdownMessageHandler)
	n.session.AddHandler(n.imageMessageHandler)
	// TODO:
	// if validGroupID is null, get all valid group.
	return n.session.Open()
}
func (n *node) BeforeStop() {
}
func (n *node) Stop() {
	n.session.Close()
}
func (n *node) Stdio(markdown string) {
	if n.stdioChan != "" {
		n.SendMsg2Group(n.stdioChan, markdown)
	}
}
func (n *node) SendStdioLog(markdown string) {
}

func (n *node) RouteMsg2Group(gid string, msg *imnode.IMMsg) {
	card := KCard{}
	card.Init()
	card.Card.Theme = "success"
	card.AddModule_markdown("**`" + msg.ShowName + "`** 转发自 " + msg.Type + " :\n---")
	card.AddModule_markdown(msg.Content)
	n.sendKCard(gid, card.String())
}
func (n *node) RouteImg2GroupByBytes(gid string, img []byte) {
}
func (n *node) RouteImg2GroupByUrl(gid string, url string) {
}

func (n *node) SendMsg2Group(gid string, msg string) {
	n.sendMarkdown(gid, msg)
}
func (n *node) SendImg2GroupByBytes(gid string, img []byte) {
}

func (n *node) SendImg2GroupByUrl(gid string, url string) {
}

func (n *node) SetMsgHandler(handler func(gid string, msg *imnode.IMMsg)) {
	n.msgHandler = handler
}

func (n *node) Name() string {
	return "KOOK[" + n.token[12:15] + "]"
}

func (n *node) markdownMessageHandler(ctx *kook.KmarkdownMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	msg := &imnode.IMMsg{}
	msg.Type = n.Name()
	msg.Content = ctx.Common.Content
	msg.UID = ctx.Common.AuthorID
	msg.ShowName = ctx.Extra.Author.Nickname
	n.msgHandler(ctx.Common.TargetID, msg)
}

func (n *node) imageMessageHandler(ctx *kook.ImageMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	msg := &imnode.IMMsg{}
	msg.Type = n.Name()
	// TODO:
	msg.Content = "[图片]"
	msg.UID = ctx.Common.AuthorID
	msg.ShowName = ctx.Extra.Author.Nickname
	n.msgHandler(ctx.Common.TargetID, msg)
}
