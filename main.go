package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"strconv"
	"time"

	"math/rand"

	kcard "local/khlcard"

	qq "local/rt"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/lonelyevil/khl"
	"github.com/lonelyevil/khl/log_adapter/plog"
	"github.com/phuslu/log"
	"github.com/spf13/viper"

	"github.com/jinzhu/copier"
)

// TODO:
// qq 二维码发送至 kook 登录
// qq 接口由 kook stdio 频道控制
// 实现 kook 图片转发至 qq

var appName string = "QQ Hime"

func buildUpdateLog() string {
	updateLog := ""
	updateLog += "1. 为SendKCard失败添加报错\n"
	updateLog += "\n\nHelloWorks-QQ Hime@[GitHub](https://github.com/HelloWorksGroup/KOOK2QQ-bot)"
	return updateLog
}

var buildVersion string = appName + " 0020"

// stdout频道
var stdoutChannel string

// 转发map
var routeMap map[string]string

// 邀请map
var kookInviteUrl map[string]string

type handlerRule struct {
	matcher string
	getter  func(ctxCommon *khl.EventDataGeneral, matchs []string, reply func(string) string)
}

var masterID string
var botID string

var localSession *khl.Session

func updateKMsg(msgId string, content string) error {
	return localSession.MessageUpdate((&khl.MessageUpdate{
		MessageUpdateBase: khl.MessageUpdateBase{
			MsgID:   msgId,
			Content: content,
		},
	}))
}

func sendKCard(target string, content string) (resp *khl.MessageResp, err error) {

	return localSession.MessageCreate((&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeCard,
			TargetID: target,
			Content:  content,
		},
	}))
}
func sendMarkdown(target string, content string) (resp *khl.MessageResp, err error) {
	return localSession.MessageCreate((&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: target,
			Content:  content,
		},
	}))
}

func sendMarkdownDirect(target string, content string) (mr *khl.MessageResp, err error) {
	return localSession.DirectMessageCreate(&khl.DirectMessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: target,
			Content:  content,
		},
	})
}

func kookLog(markdown string) {
	localTime := time.Now().Local()
	strconv.Itoa(localTime.Hour())
	strconv.Itoa(localTime.Minute())
	strconv.Itoa(localTime.Second())
	tstr := fmt.Sprintf("%02d:%02d:%02d", localTime.Hour(), localTime.Minute(), localTime.Second())
	fmt.Println("["+tstr+" KOOK LOG]:", markdown)
	if stdoutChannel != "0" {
		sendMarkdown(stdoutChannel, "`"+tstr+"` "+markdown)
	}
}

var token string

func routeMapSetup() {
	routeMap = make(map[string]string, 0)
	s := viper.Get("kook2qq").(map[string]any)
	for k, v := range s {
		vs := v.(string)
		if k != v {
			if _, ok := routeMap[k]; !ok {
				routeMap[k] = vs
			}
			if _, ok := routeMap[vs]; !ok {
				routeMap[vs] = k
			}
		}
	}
}
func kookInviteUrlSetup() {
	kookInviteUrl = make(map[string]string, 0)
	s := viper.Get("kookinvite").(map[string]any)
	for k, v := range s {
		vs := v.(string)
		if _, ok := kookInviteUrl[k]; !ok {
			kookInviteUrl[k] = vs
		}
	}
}
func kookMergeMapSetup() {
	kookMergeMap = make(map[string]KookLastMsg, 0)
}
func getConfig() {
	rand.Seed(time.Now().UnixNano())
	viper.SetDefault("token", "0")
	viper.SetDefault("stdoutChannel", "0")
	viper.SetDefault("masterID", "")
	viper.SetDefault("oldversion", "0.0.0")
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	masterID = viper.Get("masterID").(string)
	stdoutChannel = viper.Get("stdoutChannel").(string)
	fmt.Println("stdoutChannel=" + stdoutChannel)
	token = viper.Get("token").(string)
	fmt.Println("token=" + token)
	routeMapSetup()
	kookInviteUrlSetup()
	kookMergeMapSetup()
}

func prog(state overseer.State) {
	fmt.Printf("App#[%s] start ...\n", state.ID)
	getConfig()

	l := log.Logger{
		Level:  log.InfoLevel,
		Writer: &log.ConsoleWriter{},
	}

	s := khl.New(token, plog.NewLogger(&l))
	me, _ := s.UserMe()
	fmt.Println("ID=" + me.ID)
	botID = me.ID
	s.AddHandler(markdownMessageHandler)
	s.AddHandler(imageMessageHandler)
	s.Open()
	localSession = s

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println(appName + " is now running.")

	qqbotInit()
	qqbotStart()

	if viper.Get("oldversion").(string) != buildVersion {
		go func() {
			<-time.After(time.Second * time.Duration(3))
			card := kcard.KHLCard{}
			card.Init()
			card.Card.Theme = "success"
			card.AddModule_header(appName + " 热更新完成")
			card.AddModule_divider()
			card.AddModule_markdown("当前版本号：`" + buildVersion + "`")
			card.AddModule_markdown("**更新内容：**\n" + buildUpdateLog())
			sendKCard(stdoutChannel, card.String())
		}()
	}

	viper.Set("oldversion", buildVersion)
	viper.WriteConfig()

	kookLog("系统已启动")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, overseer.SIGUSR2)
	<-sc

	kookLog("系统即将关闭")

	fmt.Println("Bot will shutdown after 1 second.")

	<-time.After(time.Second * time.Duration(1))
	qqbotStop()
	// Cleanly close down the KHL session.
	s.Close()
}

func main() {
	overseer.Run(overseer.Config{
		Required: true,
		Program:  prog,
		Fetcher:  &fetcher.File{Path: "KOOK2QQ"},
		Debug:    false,
	})
}

func kookMsgToQQGroup(ctx *khl.KmarkdownMessageContext, groupId string) {
	if _, ok := kookMergeMap[ctx.Common.TargetID]; ok {
		kookMergeMap[ctx.Common.TargetID] = KookLastMsg{}
	}
	channel := ctx.Common.TargetID
	name := ctx.Extra.Author.Nickname
	content := ctx.Common.Content

	fmt.Println("[KOOK Markdown]:", channel, name, content)
	id, _ := strconv.ParseInt(groupId, 10, 64)
	qq.SendToQQGroup(name+" 转发自 KOOK:\n"+content, id)
}

func imageHandler(ctx *khl.ImageMessageContext) {
	if _, ok := kookMergeMap[ctx.Common.TargetID]; ok {
		kookMergeMap[ctx.Common.TargetID] = KookLastMsg{}
	}
	fmt.Println("[KOOK Image]:", ctx.Extra.Author.Nickname, ctx.Extra.Attachments.URL)
	var title string
	for k, v := range routeMap {
		if ctx.Common.TargetID == k {
			gid, _ := strconv.ParseInt(v, 10, 64)
			// TODO: more cases
			casen := rand.Intn(100)
			if casen <= 20 {
				title = "[图片未通过QQ审查]"
			} else if casen <= 40 {
				title = "[当前版本QQ不支持的消息]"
			} else if casen <= 60 {
				title = "[图片转发至QQ失败]"
			} else if casen <= 80 {
				title = "[未能成功转发图片]"
			} else if casen <= 100 {
				title = "[请进入KOOK端查看图片]"
			}
			var inviteStr string = ""
			if _, ok := kookInviteUrl[k]; ok {
				inviteStr = "\n邀请链接：" + kookInviteUrl[k]
			}
			qq.SendToQQGroup(ctx.Extra.Author.Nickname+" 转发自 KOOK:\n"+title+"\n"+path.Base(ctx.Extra.Attachments.URL)+"\n请使用KOOK查看。"+inviteStr, gid)
		}
	}
}

func qqMsgHandler(msg *message.GroupMessage) {
	for k, v := range routeMap {
		gid := strconv.FormatInt(msg.GroupCode, 10)
		if gid == k {
			name := msg.Sender.CardName
			if name == "" {
				name = msg.Sender.Nickname
			}
			qqMsgToKook(msg.Sender.Uin, v, name, qq.GroupMsgParse(msg))
		}
	}
}

type KookLastMsg struct {
	lastCard      kcard.KHLCard
	lastUid       int64
	lastMsgTime   int64
	lastMsgId     string
	lastCardStack int
}

var kookMergeMap map[string]KookLastMsg

// DONE: 相同用户短时间连续发言自动合并
func qqMsgToKook(uid int64, channel string, name string, msgs []qq.QQMsg) {
	var card kcard.KHLCard
	// 是否合并消息
	var merge bool = false
	var entry KookLastMsg
	fmt.Println("[MergeRoutine]:")
	fmt.Println("\tchannel=", channel)
	if kmm, ok := kookMergeMap[channel]; ok {
		entry = kmm
		fmt.Println("\tuid=", uid)
		fmt.Println("\tlastUid=", entry.lastUid)
		fmt.Println("\tlastMsgTimeDiff=", time.Now().Unix()-entry.lastMsgTime)
		fmt.Println("\tlastCardStack=", entry.lastCardStack)
		if uid == entry.lastUid && time.Now().Unix()-entry.lastMsgTime < 300 && entry.lastCardStack < 10 {
			entry.lastCardStack += 1
			card = entry.lastCard
			merge = true
		}
	}
	if !merge {
		if _, ok := kookMergeMap[channel]; !ok {
			kookMergeMap[channel] = KookLastMsg{}
			entry = kookMergeMap[channel]
		}
		card = kcard.KHLCard{}
		card.Init()
		card.Card.Theme = "success"
		card.AddModule_markdown("**`" + name + "`** 转发自 QQ:\n---")
	}
	for _, v := range msgs {
		switch v.Type {
		case 0:
			card.AddModule_markdown(v.Content)
		case 1:
			card.AddModule_image(v.Content)
		}
	}
	if !merge {
		resp, err := sendKCard(channel, card.String())
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("消息转发失败", channel, card.String())
			sendMarkdown(channel, "消息转发失败")
			entry.lastMsgId = ""
		} else {
			entry.lastMsgId = resp.MsgID
		}
	} else {
		updateKMsg(entry.lastMsgId, card.String())
	}

	entry.lastMsgTime = time.Now().Unix()
	entry.lastUid = uid
	copier.Copy(&entry.lastCard, &card)
	kookMergeMap[channel] = entry
}

func markdownMessageHandler(ctx *khl.KmarkdownMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	switch ctx.Common.TargetID {
	case botID:
		directMessageHandler(ctx.Common)
	case stdoutChannel:
		stdinHandler(ctx)
	default:
		for k, v := range routeMap {
			if ctx.Common.TargetID == k {
				kookMsgToQQGroup(ctx, v)
			}
		}
	}
}

func imageMessageHandler(ctx *khl.ImageMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	imageHandler(ctx)
}
