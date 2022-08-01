package main

import (
	"fmt"
	"regexp"

	"github.com/lonelyevil/khl"
)

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
