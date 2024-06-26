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

var appName string = "QQ Hime"

var buildVersion string = appName + " 0050"

func buildUpdateLog() string {
	updateLog := ""
	updateLog += "1. 更新协议库等\n"
	updateLog += "\n\nHelloWorks-QQ Hime@[GitHub](https://github.com/HelloWorksGroup/Route2QQ-bot)"
	return updateLog
}

type IMNode interface {
	init()
	start() error
	stop()
	registMsgHandler()
	routeMsg2Group(gid string, uid string, msg string)
	sendMsg2Group(gid string, msg string)
	// sendMsg2Person(uid string, msg string)
	sendImg2GroupByBytes(gid string, img []byte)
	sendImg2GroupByUrl(gid string, url string)
	name() string
}

type handlerRule struct {
	matcher string
	getter  func(ctxCommon *kook.EventDataGeneral, matchs []string, reply func(string) string)
}

func kookLog(markdown string) {
	gLog.Info().Msgf("kq-log:%s", markdown)
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

var gLog log.Logger

func prog(state overseer.State) {
	fmt.Printf("App#[%s] start ...\n", state.ID)
	GetConfig()

	gLog = log.Logger{
		Level: log.InfoLevel,
		Writer: &log.MultiEntryWriter{
			&log.ConsoleWriter{ColorOutput: true},
			&log.FileWriter{
				Filename:   "kq.log",
				MaxSize:    512 << 10,
				MaxBackups: 16,
				LocalTime:  true},
		},
	}

	s := kook.New(token, plog.NewLogger(&gLog))
	me, _ := s.UserMe()
	fmt.Println("ID=" + me.ID)
	botID = me.ID
	s.AddHandler(markdownMessageHandler)
	s.AddHandler(imageMessageHandler)
	s.Open()
	localSession = s

	fmt.Println("KOOK node online.")

	qqbotInit()
	qqbotStart()

	fmt.Println("QQ node online.")

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
	msgCache.gc()

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
		Fetcher:  &fetcher.File{Path: "Route2QQ"},
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
		for k, v := range kook2qqRouteMap {
			if ctx.Common.TargetID == k {
				go kookMsgToQQGroup(ctx, k, v)
			}
		}
		for kookGid, v := range kook2vcRouteMap {
			if ctx.Common.TargetID == kookGid {
				go kookMsgToVC(ctx, kookGid, v.Url, v.Gid, v.Secret)
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
