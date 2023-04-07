package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/HelloWorksGroup/IMSuperGroup/imnode"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
)

var appName string = "IMSuperGroup"

var buildVersion string = appName + " 0100"

func buildUpdateLog() string {
	updateLog := ""
	updateLog += "1. 重大更新，详见GitHub\n"
	updateLog += "\n\nHelloWorks-IMSuperGroup@[GitHub](https://github.com/HelloWorksGroup/IMSuperGroup)"
	return updateLog
}

var gLog log.Logger
var nodes []imnode.IMNode
var superGroups [][]string

// 在绑定了stdio通道的IM节点发送LOG
func remoteIMLog(markdown string) {
	gLog.Info().Msgf("REMOTE-LOG:%s", markdown)
	localTime := time.Now().Local()
	strconv.Itoa(localTime.Hour())
	strconv.Itoa(localTime.Minute())
	strconv.Itoa(localTime.Second())
	tstr := fmt.Sprintf("%02d:%02d:%02d", localTime.Hour(), localTime.Minute(), localTime.Second())
	// fmt.Println("["+tstr+" REMOTE LOG]:", markdown)
	// if stdoutChannel != "0" {
	// TODO:
	// kook.SendMarkdown(stdoutChannel, "`"+tstr+"` "+markdown)
	// }
	md := "`" + tstr + "` " + markdown
	for _, v := range nodes {
		v.SendStdioLog(md)
	}
}

func allNodeBeforeShutdown() {
	for _, v := range nodes {
		v.BeforeStop()
	}
	// msgCache.Backup()
}

func allNodeShutdown() {
	for _, v := range nodes {
		v.Stop()
	}
}

func prog(state overseer.State) {
	fmt.Printf("App#[%s] start ...\n", state.ID)
	nodes = make([]imnode.IMNode, 0)
	superGroups = make([][]string, 0)
	GetConfig()

	gLog = log.Logger{
		Level: log.InfoLevel,
		Writer: &log.MultiEntryWriter{
			&log.ConsoleWriter{ColorOutput: true},
			&log.FileWriter{
				Filename:   "IMR.log",
				MaxSize:    512 << 10,
				MaxBackups: 16,
				LocalTime:  true},
		},
	}
	// TODO: Nodes init & start
	for _, v := range nodes {
		gLog.Info().Msgf("Node [" + v.Name() + "] Starting")
		if v.Start() != nil {
			gLog.Error().Msgf("Node [" + v.Name() + "] Start FAILED!!!")
		} else {
			gLog.Info().Msgf("Node [" + v.Name() + "] ONLINE")
		}
	}
	if viper.Get("oldversion").(string) != buildVersion {
		go func() {
			<-time.After(time.Second * time.Duration(3))
			updateLog := appName + " 热更新完成"
			updateLog += "\n当前版本号：`" + buildVersion + "`\n"
			updateLog += "**更新内容：**\n" + buildUpdateLog()
			remoteIMLog(updateLog)
		}()
	}
	viper.Set("oldversion", buildVersion)
	viper.WriteConfig()

	remoteIMLog("系统已完全启动")
	// msgCache.gc()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, overseer.SIGUSR2)
	sig := <-sc
	if sig == overseer.SIGUSR2 {
		remoteIMLog("检测到二进制变更，系统即将进行快速重启")
	} else {
		remoteIMLog("接收到外部指令，系统即将关闭")
	}
	allNodeBeforeShutdown()

	gLog.Info().Msgf("Bot will shutdown after 1 second.")
	<-time.After(time.Second * time.Duration(1))
	allNodeShutdown()
	gLog.Info().Msgf("[SHUTDOWN]")
}

func main() {
	overseer.Run(overseer.Config{
		Required: true,
		Program:  prog,
		Fetcher:  &fetcher.File{Path: "IMSuperGroup"},
		Debug:    false,
	})
}
