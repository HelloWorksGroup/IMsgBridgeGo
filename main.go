package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	kcard "local/khlcard"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/lonelyevil/kook"
	"github.com/lonelyevil/kook/log_adapter/plog"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
)

// TODO:
// qq 二维码发送至 kook 登录
// qq 接口由 kook stdio 频道控制
// 实现 kook 图片转发至 qq

var appName string = "QQ Hime"

var buildVersion string = appName + " 0026"

func buildUpdateLog() string {
	updateLog := ""
	updateLog += "1. 增加消息缓存，为消息回复的支持做准备"
	updateLog += "\n\nHelloWorks-QQ Hime@[GitHub](https://github.com/HelloWorksGroup/KOOK2QQ-bot)"
	return updateLog
}

type handlerRule struct {
	matcher string
	getter  func(ctxCommon *kook.EventDataGeneral, matchs []string, reply func(string) string)
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

func prog(state overseer.State) {
	fmt.Printf("App#[%s] start ...\n", state.ID)
	getConfig()

	l := log.Logger{
		Level:  log.InfoLevel,
		Writer: &log.ConsoleWriter{},
	}

	s := kook.New(token, plog.NewLogger(&l))
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

	kookLog("系统已完全启动")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, overseer.SIGUSR2)
	sig := <-sc
	if sig == overseer.SIGUSR2 {
		kookLog("检测到二进制变更，系统即将进行快速重启")
	} else {
		kookLog("接收到外部指令，系统即将关闭")
	}
	beforeShutdown()

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

func markdownMessageHandler(ctx *kook.KmarkdownMessageContext) {
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
				go kookMsgToQQGroup(ctx, k, v)
			}
		}
	}
}

func imageMessageHandler(ctx *kook.ImageMessageContext) {
	if ctx.Extra.Author.Bot {
		return
	}
	imageHandler(ctx)
}
