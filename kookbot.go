package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/lonelyevil/kook"
)

var commOnce sync.Once
var commRules []handlerRule = []handlerRule{
	{`^qqping`, func(ctxCommon *kook.EventDataGeneral, s []string, f func(string) string) {
		delay := rand.Intn(500) + rand.Intn(100) + rand.Intn(50) + rand.Intn(14)
		<-time.After(time.Millisecond * time.Duration(delay+2000))
		f("来自 `QQ` 的回复: 字节=`256` 时间=`" + strconv.Itoa(delay) + "ms` TTL=`" + strconv.Itoa(61-rand.Intn(7)) + "`")
	}},
}

func stdinHandler(ctx *kook.KmarkdownMessageContext) {
	ctxCommon := ctx.Common

	reply := func(words string) string {
		resp, _ := sendMarkdown(stdoutChannel, words)
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

func portMarkdown(ctxCommon *kook.EventDataGeneral, s []string, f func(string) string) {
	sendMarkdown(s[1], s[2])
	return
}

var directRules []handlerRule = []handlerRule{
	{`^\s*send\s*(\d+),(.*)$`, portMarkdown},
}

func directMessageHandler(ctxCommon *kook.EventDataGeneral) {
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

func otherChanHandler(ctxCommon *kook.EventDataGeneral) {
	if ctxCommon.Type != kook.MessageTypeKMarkdown {
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
