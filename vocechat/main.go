package vcNode

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/phuslu/log"

	"github.com/HelloWorksGroup/IMSuperGroup/imnode"
)

type vchat struct {
	logger       *log.Logger
	stdioChan    string
	validGroupID []string
	msgHandler   func(gid string, msg *imnode.IMMsg)

	port   string
	secret string
	apiurl string
}

func Setup(setting map[string]any, logger *log.Logger) imnode.IMNode {
	n := new(vchat)
	n.logger = logger
	n.validGroupID = imnode.ConvertSettingGroups2StrSlice(setting["groups"])
	n.port = setting["port"].(string)
	n.secret = setting["secret"].(string)
	n.apiurl = setting["apiurl"].(string)
	return n
}

func (n *vchat) GroupIDValid(gid string) bool {
	for _, id := range n.validGroupID {
		if id == gid {
			return true
		}
	}
	return false
}
func (n *vchat) Start() error {
	http.HandleFunc("/", n.handleWebhook)
	go http.ListenAndServe(":"+n.port, nil)
	return nil
}
func (n *vchat) BeforeStop() {
}
func (n *vchat) Stop() {
	// vchat stop
}
func (n *vchat) Stdio(markdown string) {
	if n.GroupIDValid(n.stdioChan) {
		// send mark down to stdio Channel
	}
}
func (n *vchat) SendStdioLog(markdown string) {
	if n.stdioChan != "" {
		n.send(n.stdioChan, n.secret, markdown)
	}

}

func (n *vchat) RouteMsg2Group(gid string, msg *imnode.IMMsg) {
	md := "#### `" + msg.ShowName + "` 转发自 **" + msg.Type + "** :\n" + msg.Content
	n.send(gid, n.secret, md)
}
func (n *vchat) RouteImg2GroupByBytes(gid string, img []byte) {
}
func (n *vchat) RouteImg2GroupByUrl(gid string, url string) {
}

func (n *vchat) SendMsg2Group(gid string, msg string) {
	n.send(gid, n.secret, msg)
}
func (n *vchat) SendImg2GroupByBytes(gid string, img []byte) {
}

func (n *vchat) SendImg2GroupByUrl(gid string, url string) {
}

func (n *vchat) SetMsgHandler(handler func(gid string, msg *imnode.IMMsg)) {
	n.msgHandler = handler
}

func (n *vchat) Name() string {
	return "vocechat[" + n.secret[:3] + "]"
}

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

func (n *vchat) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var vcJSON vocechatJSON
	decoder.Decode(&vcJSON)
	// 获取用户信息
	// TODO: user info cache
	getUser := func() vocechatUser {
		r, err := http.Get(n.apiurl + "/user/" + strconv.Itoa(vcJSON.From))
		if err != nil {
			n.logger.Error().Msgf("getUser ERROR")
		}
		decoder := json.NewDecoder(r.Body)
		var user vocechatUser
		decoder.Decode(&user)
		return user
	}
	// 判断是否是新消息
	if vcJSON.Detail.Type == "normal" {
		user := getUser()
		// 判断消息是否是机器人发送
		if user.IsBot {
			return
		}
		msg := &imnode.IMMsg{}
		msg.Type = n.Name()
		msg.Content = vcJSON.Detail.Content
		msg.UID = strconv.Itoa(user.Uid)
		msg.ShowName = user.Name
		n.msgHandler(strconv.Itoa(vcJSON.Target.Gid), msg)
	}
}

func (n *vchat) send(gid string, secret string, content string) {
	markdown := []byte(content)
	r, err := http.NewRequest("POST", n.apiurl+"/bot/send_to_group/"+gid, bytes.NewBuffer(markdown))
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "text/markdown")
	r.Header.Add("x-api-key", secret)
	client := &http.Client{}
	_, err = client.Do(r)
	if err != nil {
		n.logger.Error().Msgf("send ERROR")
		return
	}
}
