package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	qq "local/rt"
	"log"
	"net/http"
	"strconv"
)

type vocechatUser struct {
	Name   string `json:"name"`
	Uid    int    `json:"uid"`
	Gender int    `json:"gender"`
	IsBot  bool   `json:"is_bot"`
}
type vocechatDetail struct {
	Content     string `json:"content,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Type        string `json:"type"`
}
type vocechatTarget struct {
	Gid int `json:"gid,omitempty"`
	Uid int `json:"uid,omitempty"`
}
type vocechatJSON struct {
	Detail vocechatDetail `json:"detail"`
	From   int            `json:"from_uid"`
	Target vocechatTarget `json:"target"`
}

func voceChatBot(port int) {
	listen(port)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("headers: %v\n", r.Header)
	// fmt.Printf("body: %v\n", r.Body)

	// bytes, _ := io.ReadAll(r.Header)
	// fmt.Println(string(bytes))
	// bytes, _ = io.ReadAll(r.Body)
	// fmt.Println(string(bytes))

	if r.Header.Get("Content-Type") != "application/json" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var vcJSON vocechatJSON
	decoder.Decode(&vcJSON)
	// 获取用户信息
	getUser := func(v vocechatInstance) vocechatUser {
		r, err := http.Get(v.Url + "/user/" + strconv.Itoa(vcJSON.From))
		if err != nil {
			// Error log
		}
		decoder := json.NewDecoder(r.Body)
		var user vocechatUser
		decoder.Decode(&user)
		return user
	}
	// 判断是否是新消息
	if vcJSON.Detail.Type == "normal" {
		// 判断消息是否属于转发列表
		for qqgroup, v := range qq2vcRouteMap {
			if v.Gid == strconv.Itoa(vcJSON.Target.Gid) {
				user := getUser(v)
				// 判断消息是否是机器人发送
				if user.IsBot {
					return
				}
				// 转发消息至目的地
				id, _ := strconv.ParseInt(qqgroup, 10, 64)
				go qq.SendToQQGroup(user.Name+" 转发自 vocechat:\n"+vcJSON.Detail.Content, id)
			}
		}
		for kookGid, v := range kook2vcRouteMap {
			if v.Gid == strconv.Itoa(vcJSON.Target.Gid) {
				user := getUser(v)
				// 判断消息是否是机器人发送
				if user.IsBot {
					return
				}
				go vcMsgToKook(kookGid, user.Name, vcJSON.Detail.Content)
			}
		}
	}
}

func vocechatSend(url string, gid int, secret string, content string) {
	markdown := []byte(content)
	r, err := http.NewRequest("POST", url+"/bot/send_to_group/"+strconv.Itoa(gid), bytes.NewBuffer(markdown))
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "text/markdown")
	r.Header.Add("x-api-key", secret)
	client := &http.Client{}
	_, err = client.Do(r)
	if err != nil {
		// 发送失败
		return
	}
}

func listen(port int) {
	fmt.Println("Listen on port:", port)
	http.HandleFunc("/", handleWebhook)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
