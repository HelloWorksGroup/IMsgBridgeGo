package main

import (
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
	// 判断是否是新消息
	if vcJSON.Detail.Type == "normal" {
		// 判断消息是否属于转发列表
		for qqgroup, v := range vc2qqRouteMap {
			if v.Gid == strconv.Itoa(vcJSON.Target.Gid) {
				// 获取用户名称
				r, err := http.Get(v.Url + "/user/" + strconv.Itoa(vcJSON.From))
				var name string = "未知姓名"
				if err == nil {
					decoder := json.NewDecoder(r.Body)
					var user vocechatUser
					decoder.Decode(&user)
					name = user.Name
				}
				// 转发消息至目的地
				id, _ := strconv.ParseInt(qqgroup, 10, 64)
				qq.SendToQQGroup(name+" 转发自 vocechat:\n"+vcJSON.Detail.Content, id)
				return
			}
		}
	}
}

func listen(port int) {
	fmt.Println("Listen on port:", port)
	http.HandleFunc("/", handleWebhook)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
