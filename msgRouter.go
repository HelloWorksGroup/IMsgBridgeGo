package main

import (
	"context"
	"io/ioutil"
	kcard "local/khlcard"
	qq "local/rt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/jinzhu/copier"
	"github.com/lonelyevil/kook"
)

// TODO: 大量辣鸡代码待优化

// : 支持消息回复
// : 将kook回复标记转换成 qq @ uid

// TODO: 更多 kmarkdown tag 处理
// TODO: 处理KOOK服务器表情
// TODO: 移除对于机器人自身的at转发，QQ & KOOK
// 将(met)bot id(met)\s+ 变为 @ name
func kookMet2At(content string, channelId string) string {
	replaceMap := make(map[string]string)
	var c *kook.Channel = nil
	r := regexp.MustCompile(`\(met\)(\d+)\(met\)`)
	content = strings.ReplaceAll(content, "(met)"+botID+"(met)", "")
	for {
		matchs := r.FindStringSubmatch(content)
		if len(matchs) > 0 {
			if c == nil {
				c, _ = localSession.ChannelView(channelId)
			}
			if _, ok := replaceMap["(met)"+matchs[1]+"(met)"]; !ok {
				u, err := localSession.UserView(matchs[1], kook.UserViewWithGuildID(c.GuildID))
				if err == nil {
					content = strings.ReplaceAll(content, "(met)"+matchs[1]+"(met)", "@"+u.Nickname)
					replaceMap["(met)"+matchs[1]+"(met)"] = "@" + u.Nickname
				} else {
					content = strings.ReplaceAll(content, "(met)"+matchs[1]+"(met)", "@某人")
					replaceMap["(met)"+matchs[1]+"(met)"] = "@某人"
				}
			}
		} else {
			break
		}
	}

	return content
}

// 将 [foo](bar) 变为 bar
func kookLink2Link(content string) string {
	r := regexp.MustCompile(`\[.+?\]\((.+?)\)`)
	content = r.ReplaceAllString(content, " $1 ")
	return content
}

func kookMsgToVC(ctx *kook.KmarkdownMessageContext, guildId string, api string, gid string, secret string) {
	name := ctx.Extra.Author.Nickname
	content := ctx.Common.Content
	gidInt, _ := strconv.Atoi(gid)
	vocechatSend(api, gidInt, secret, "**`"+name+"`** 转发自 KOOK:\n"+content)

}

func kookMsgToQQGroup(ctx *kook.KmarkdownMessageContext, guildId string, groupId string) {
	gLog.Info().Msgf("kmsg-log:%v", ctx)
	if _, ok := kookLastCache[ctx.Common.TargetID]; ok {
		kookLastCache[ctx.Common.TargetID] = kookLastMsgs{}
	}
	channel := ctx.Common.TargetID
	name := ctx.Extra.Author.Nickname
	content := ctx.Common.Content

	content = kookMet2At(content, guildId)
	content = kookLink2Link(content)

	gLog.Info().Msgf("[KOOK Markdown]:[channel=%v][name=%v][content=%v]", channel, name, content)
	id, _ := strconv.ParseInt(groupId, 10, 64)

	var mid int32
	var replyUid, replyName string
	if ctx.Extra.Quote != nil {
		replyUid, replyName = msgCache.WhomReply(guildId, ctx.Extra.Quote.RongID)
	}
	if len(replyUid) > 0 {
		msgs := make([]message.IMessageElement, 0)
		quoteUid, _ := strconv.ParseInt(replyUid, 10, 64)
		msgs = append(msgs, message.NewText(name+" 转发自 KOOK:\n"))
		if len(replyName) > 0 {
			msgs = append(msgs, message.NewAt(quoteUid, "@"+replyName))
		} else {
			msgs = append(msgs, message.NewAt(quoteUid))
		}
		msgs = append(msgs, message.NewText(content))
		mid = qq.SendToQQGroupEx(msgs, id)
	} else {
		mid = qq.SendToQQGroup(name+" 转发自 KOOK:\n"+content, id)
	}
	msgCache.GetMsg(groupId, strconv.FormatInt(int64(mid), 10), ctx.Extra.Author.ID, name)
	gLog.Info().Msgf("[SEND QQ msg]:[ID=%d]", mid)
}

func imageHandler(ctx *kook.ImageMessageContext) {
	if _, ok := kookLastCache[ctx.Common.TargetID]; ok {
		kookLastCache[ctx.Common.TargetID] = kookLastMsgs{}
	}
	gLog.Info().Msgf("[KOOK Image]:[name=%v][url=%v]", ctx.Extra.Author.Nickname, ctx.Extra.Attachments.URL)
	for k, v := range kook2qqRouteMap {
		if ctx.Common.TargetID == k {
			gidStr := v
			go func() {
				gid, _ := strconv.ParseInt(gidStr, 10, 64)
				// 超时退出
				context, cancel := context.WithCancel(context.Background())
				defer cancel()
				ch := make(chan bool, 1)
				// 下载图片后发送至QQ群
				go func() {
					// mid = qq.SendToQQGroup(ctx.Extra.Author.Nickname+" 尝试转发图片自 KOOK:\n"+ctx.Extra.Attachments.URL, gid)
					// msgCache.GetMsg(strconv.FormatInt(gid, 10), strconv.FormatInt(int64(mid), 10), ctx.Extra.Author.ID, ctx.Extra.Author.Nickname)
					resp, err := http.Get(ctx.Extra.Attachments.URL)
					if err != nil {
						return
					}
					defer resp.Body.Close()
					data, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						panic(err)
					}

					msgs := make([]message.IMessageElement, 0)
					msgs = append(msgs, message.NewText(ctx.Extra.Author.Nickname+" 转发图片自 KOOK:"))
					imgMsg, err := qq.UploadImgToQQGroup(data, gid)
					if err == nil {
						msgs = append(msgs, imgMsg)
					} else {
						msgs = append(msgs, message.NewText("Failed..."))
					}
					mid := qq.SendToQQGroupEx(msgs, gid)
					msgCache.GetMsg(strconv.FormatInt(gid, 10), strconv.FormatInt(int64(mid), 10), ctx.Extra.Author.ID, ctx.Extra.Author.Nickname)
					gLog.Info().Msgf("[SEND QQ msg]:[ID=%d]", mid)

					select {
					default:
						ch <- true
					case <-context.Done():
						return
					}
				}()
				select {
				case <-ch:
					return
				case <-time.After(20 * time.Second):
					gLog.Info().Msgf("[SEND QQ IMG]:TIMEOUT")
				}
			}()
		}
	}
}

func qqMsgHandler(msg *message.GroupMessage) {
	gid := strconv.FormatInt(msg.GroupCode, 10)
	name := msg.Sender.DisplayName()
	if name == "" {
		name = msg.Sender.Nickname
	}
	for k, v := range kook2qqRouteMap {
		if gid == k {
			go qqMsgToKook(gid, msg.Sender.Uin, v, name, qq.GroupMsgParse(msg))
		}
	}

	for k, v := range qq2vcRouteMap {
		if gid == k {
			intgid, _ := strconv.Atoi(v.Gid)
			go vocechatSend(v.Url, intgid, v.Secret, "**`"+name+"`** 转发自 QQ:\n"+qq.GroupMsg2Markdown(msg))
		}
	}
}

func escapeToCleanUnicode(raw string) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, str)
	return clean, nil
}

func vcMsgToKook(channel string, name string, content string) {
	card := kcard.KHLCard{}
	card.Init()
	card.Card.Theme = "success"
	card.AddModule_markdown("**`" + name + "`** 转发自 vocechat:\n---")
	card.AddModule_markdown(content)
	sendKCard(channel, card.String())
}

// DONE: 相同用户短时间连续发言自动合并
func qqMsgToKook(gid string, uid int64, channel string, name string, msgs []qq.QQMsg) {
	var card kcard.KHLCard
	// 是否合并消息
	var merge bool = false
	var entry kookLastMsgs
	var cleanName string
	var err error
	gLog.Info().Msgf("qmsg-log:%v", msgs)
	if kmm, ok := kookLastCache[channel]; ok {
		entry = kmm
		if uid == entry.Uid && time.Now().Unix()-entry.MsgTime < 300 && entry.CardStack < 10 {
			entry.CardStack += 1
			card = entry.Card
			merge = true
		}
	}
	if !merge {
		if _, ok := kookLastCache[channel]; !ok {
			kookLastCache[channel] = kookLastMsgs{}
			entry = kookLastCache[channel]
		}
		card = kcard.KHLCard{}
		card.Init()
		card.Card.Theme = "success"
		cleanName, err = escapeToCleanUnicode(name)
		if err != nil {
			cleanName = "某姓名无法打印人士"
		}
		card.AddModule_markdown("**`" + cleanName + "`** 转发自 QQ:\n---")
	}
	var atCount int = 0
	var cachedStr string = ""
	cachedStrRelease := func() {
		if len(cachedStr) > 0 {
			card.AddModule_markdown(cachedStr)
			cachedStr = ""
		}
	}
	for _, v := range msgs {
		switch v.Type {
		case 0: // 可合并消息
			if len(v.Content) > 0 && v.Content != " " {
				// 忽略QQ回复消息自带的空白消息(一个0x20字符)
				cachedStr += v.Content
			}
		case 1:
			cachedStrRelease()
			card.AddModule_image(v.Content)
		case 2: // At
			atCount += 1
			if atCount == 1 {
				cachedStr += v.Content
			}
		case 3: // Reply
			var mid string
			r := regexp.MustCompile(`\d+`)
			matchs := r.FindStringSubmatch(v.Content)
			if len(matchs) > 0 {
				mid = matchs[0]
			}
			replyUid, _ := msgCache.WhomReply(gid, mid)
			if len(replyUid) > 0 {
				cachedStr += "(met)" + replyUid + "(met) "
			} else {
				cachedStr += v.Content
			}
		case 4: // Unknown
			cachedStrRelease()
			card.AddModule_markdown(v.Content)
		}
	}
	cachedStrRelease()
	if !merge {
		resp, err := sendKCard(channel, card.String())
		if err != nil {
			gLog.Error().Msgf("[SEND KOOK MSG]:[json=%s]", card.String())
			kookLog("消息转发失败")
			entry.MsgId = ""
		} else {
			entry.CardStack = 1
			entry.MsgId = resp.MsgID
			msgCache.GetMsg(channel, entry.MsgId, strconv.FormatInt(uid, 10), cleanName)
			gLog.Info().Msgf("[SEND KOOK MSG]:[ID=%s][json=%s]", entry.MsgId, card.String())
		}
	} else {
		updateKMsg(entry.MsgId, card.String())
	}

	entry.MsgTime = time.Now().Unix()
	entry.Uid = uid
	copier.Copy(&entry.Card, &card)
	kookLastCache[channel] = entry
}
