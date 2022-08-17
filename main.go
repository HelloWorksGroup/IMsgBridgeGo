package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"math/rand"

	kcard "local/khlcard"

	qq "local/rt"

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
	updateLog += "1. 捕获未识别消息类型\n"
	updateLog += "\n\nHelloWorks-QQ Hime@[GitHub](https://github.com/HelloWorksGroup/KOOK2QQ-bot)"
	return updateLog
}

var buildVersion string = appName + " 0007"

// kook邀请链接
var kookUrl string

// kook频道
var kookChannel string

// stdout频道
var stdoutChannel string

// QQ群号
var qqGroup string
var qqGroupCode int64

type handlerRule struct {
	matcher string
	getter  func(ctxCommon *khl.EventDataGeneral, matchs []string, reply func(string) string)
}

var masterID string
var botID string

var localSession *khl.Session

// DONE: 相同用户短时间连续发言自动合并

// TODO: 待优化垃圾代码
var lastSpeakerId int64
var lastSpeakTime int64
var lastMsgId string
var lastCard kcard.KHLCard
var lastCardStack int = 0

func MsgRouteQQ2KOOK(id int64, name string, qqmsg []qq.QQMsg) {
	var card kcard.KHLCard
	// 是否合并消息
	var merge bool = false
	if id == lastSpeakerId && time.Now().Unix()-lastSpeakTime < 300 && lastCardStack < 10 {
		lastCardStack += 1
		card = lastCard
		merge = true
	} else {
		lastCardStack = 0
		card = kcard.KHLCard{}
		card.Init()
		card.Card.Theme = "success"
		card.AddModule_markdown("**`" + name + "`** from QQ:\n---")
	}

	for _, v := range qqmsg {
		switch v.Type {
		case 0:
			card.AddModule_markdown(v.Content)
		case 1:
			card.AddModule_image(v.Content)
		}
	}
	if !merge {
		resp, _ := sendKCard(kookChannel, card.String())
		lastMsgId = resp.MsgID
	} else {
		updateKMsg(lastMsgId, card.String())
	}

	lastSpeakTime = time.Now().Unix()
	lastSpeakerId = id
	copier.Copy(&lastCard, &card)
}

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

func getConfig() {
	rand.Seed(time.Now().UnixNano())
	viper.SetDefault("token", "0")
	viper.SetDefault("kookChannel", "0")
	viper.SetDefault("koolUrl", "")
	viper.SetDefault("stdoutChannel", "0")
	viper.SetDefault("qqGroup", "0")
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
	kookChannel = viper.Get("kookChannel").(string)
	fmt.Println("kookChannel=" + kookChannel)
	kookUrl = viper.Get("koolUrl").(string)
	fmt.Println("koolUrl=" + kookUrl)
	stdoutChannel = viper.Get("stdoutChannel").(string)
	fmt.Println("stdoutChannel=" + stdoutChannel)
	qqGroup = viper.Get("qqGroup").(string)
	qqGroupCode, _ = strconv.ParseInt(qqGroup, 10, 64)
	fmt.Println("qqGroupCode=", qqGroupCode)

	token = viper.Get("token").(string)
	fmt.Println("token=" + token)
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
	s.AddHandler(fileMessageHandler)
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

func markdownMessageHandler(ctx *khl.KmarkdownMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	switch ctx.Common.TargetID {
	case botID:
		directMessageHandler(ctx.Common)
	case kookChannel:
		markdownHandler(ctx)
	case stdoutChannel:
		stdinHandler(ctx)
	}
}

func imageMessageHandler(ctx *khl.ImageMessageContext) {
	if ctx.Extra.Author.Bot || ctx.Common.TargetID != kookChannel {
		return
	}
	imageHandler(ctx)
}

func fileMessageHandler(ctx *khl.FileMessageContext) {
	if ctx.Extra.Author.Bot || ctx.Common.TargetID != kookChannel {
		return
	}
	fileHandler(ctx)
}
