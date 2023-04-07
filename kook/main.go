package kookNode

import (
	"github.com/phuslu/log"

	"github.com/HelloWorksGroup/IMSuperGroup/imnode"
	"github.com/lonelyevil/kook"
	"github.com/lonelyevil/kook/log_adapter/plog"
)

type node struct {
	botID     string
	token     string
	logger    *log.Logger
	session   *kook.Session
	stdioChan string
}

func Setup(setting map[string]string, logger *log.Logger) imnode.IMNode {
	n := new(node)
	n.logger = logger
	n.token = setting["token"]
	return n
}

func (n *node) Start() error {
	n.session = kook.New(n.token, plog.NewLogger(n.logger))
	me, _ := n.session.UserMe()
	n.botID = me.ID
	n.session.AddHandler(n.markdownMessageHandler)
	n.session.AddHandler(n.imageMessageHandler)
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

func (n *node) RouteMsg2Group(gid string, uid string, msg imnode.IMMsg) {
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

func (n *node) Name() string {
	return "KOOK-" + n.token[:3]
}

func (n *node) markdownMessageHandler(ctx *kook.KmarkdownMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	switch ctx.Common.TargetID {
	default:
		// for k, v := range kook2qqRouteMap {
		// 	if ctx.Common.TargetID == k {
		// 		go kookMsgToQQGroup(ctx, k, v)
		// 	}
		// }
		// for kookGid, v := range kook2vcRouteMap {
		// 	if ctx.Common.TargetID == kookGid {
		// 		go kookMsgToVC(ctx, kookGid, v.Url, v.Gid, v.Secret)
		// 	}
		// }
	}
}

func (n *node) imageMessageHandler(ctx *kook.ImageMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	// imageHandler(ctx)
}
