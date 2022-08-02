package main

import (
	"math/rand"
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

func commonChanHandler(ctxCommon *khl.EventDataGeneral) {
	if ctxCommon.Type != khl.MessageTypeKMarkdown {
		return
	}

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
