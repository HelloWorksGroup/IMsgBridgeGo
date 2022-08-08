package main

import (
	"fmt"
	"math/rand"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/lonelyevil/khl"
)

var commOnce sync.Once
var commRules []handlerRule = []handlerRule{
	{`^qqping`, func(ctxCommon *khl.EventDataGeneral, s []string, f func(string) string) {
		delay := rand.Intn(500) + rand.Intn(100) + rand.Intn(50) + rand.Intn(14)
		<-time.After(time.Millisecond * time.Duration(delay+2000))
		f("来自 `QQ` 的回复: 字节=`256` 时间=`" + strconv.Itoa(delay) + "ms` TTL=`" + strconv.Itoa(61-rand.Intn(7)) + "`")
	}},
}

func stdinHandler(ctx *khl.KmarkdownMessageContext) {
	ctxCommon := ctx.Common

	reply := func(words string) string {
		resp, _ := sendMarkdown(kookChannel, words)
		return resp.MsgID
	}

	for n := range commRules {
		v := &commRules[n]
		r := regexp.MustCompile(v.matcher)
		matchs := r.FindStringSubmatch(ctxCommon.Content)
		if len(matchs) > 0 {
			go v.getter(ctxCommon, matchs, reply)
			return
		}
	}
}

func markdownHandler(ctx *khl.KmarkdownMessageContext) {
	// TODO: 待优化垃圾代码
	lastSpeakerId = 0
	fmt.Println("[KOOK Markdown]:", ctx.Extra.Author.Nickname, ctx.Common.Content)
	qqGetKookMarkdown(ctx.Extra.Author.Nickname + " from KOOK:\n" + ctx.Common.Content)
}

func imageHandler(ctx *khl.ImageMessageContext) {
	// TODO: 待优化垃圾代码
	lastSpeakerId = 0
	fmt.Println("[KOOK Image]:", ctx.Extra.Author.Nickname, ctx.Extra.Attachments.URL)
	if rand.Intn(100) <= 200 {
		var title string
		if rand.Intn(100) <= 50 {
			title = "[图片未通过审查]"
		} else {
			title = "[当前版本QQ不支持的消息]"
		}
		qqGetKookMarkdown(ctx.Extra.Author.Nickname + ":" + title + "\n" + path.Base(ctx.Extra.Attachments.URL) + "\n请访问 " + kookUrl + " 查看")
	} else {
		qqGetKookImage(ctx.Extra.Author.Nickname, ctx.Extra.Attachments.URL)
	}
}

func fileHandler(ctx *khl.FileMessageContext) {
	fmt.Println("[KOOK File]:", ctx.Extra.Author.Nickname, ctx.Extra.Attachments.URL)
	qqGetKookMarkdown(ctx.Extra.Author.Nickname + ":\n[当前QQ版本不支持的消息]\n请访问 " + kookUrl + " 查看")
}

func portMarkdown(ctxCommon *khl.EventDataGeneral, s []string, f func(string) string) {
	sendMarkdown(s[1], s[2])
	return
}

var directRules []handlerRule = []handlerRule{
	{`^\s*send\s*(\d+),(.*)$`, portMarkdown},
}

func directMessageHandler(ctxCommon *khl.EventDataGeneral) {
	if ctxCommon.AuthorID != masterID {
		sendMarkdownDirect(ctxCommon.AuthorID, "未授权的通信...访问拒绝 [-2]")
		return
	}
	reply := func(words string) string {
		resp, _ := sendMarkdownDirect(masterID, words)
		return resp.MsgID
	}
	fmt.Println("Master said: " + ctxCommon.Content)

	for n := range directRules {
		v := &directRules[n]
		r := regexp.MustCompile(v.matcher)
		matchs := r.FindStringSubmatch(ctxCommon.Content)
		if len(matchs) > 0 {
			go v.getter(ctxCommon, matchs, reply)
			return
		}
	}
}

func otherChanHandler(ctxCommon *khl.EventDataGeneral) {
	if ctxCommon.Type != khl.MessageTypeKMarkdown {
		return
	}
	reply := func(words string) string {
		resp, _ := sendMarkdown(ctxCommon.TargetID, words)
		return resp.MsgID
	}

	for n := range commRules {
		v := &commRules[n]
		r := regexp.MustCompile(v.matcher)
		matchs := r.FindStringSubmatch(ctxCommon.Content)
		if len(matchs) > 0 {
			go v.getter(ctxCommon, matchs, reply)
			return
		}
	}
}
