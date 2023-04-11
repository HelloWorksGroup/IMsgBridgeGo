package node

import (
	"github.com/phuslu/log"

	"github.com/HelloWorksGroup/IMSuperGroup/imnode"
)

type node struct {
	logger       *log.Logger
	stdioChan    string
	validGroupID []string
	msgHandler   func(gid string, msg *imnode.IMMsg)
}

func Setup(setting map[string]any, logger *log.Logger) imnode.IMNode {
	n := new(node)
	n.logger = logger
	n.validGroupID = imnode.ConvertSettingGroups2StrSlice(setting["groups"])
	return n
}

func (n *node) GroupIDValid(gid string) bool {
	for _, id := range n.validGroupID {
		if id == gid {
			return true
		}
	}
	return false
}
func (n *node) Start() error {
	// node start
	return nil
}
func (n *node) BeforeStop() {
}
func (n *node) Stop() {
	// node stop
}
func (n *node) Stdio(markdown string) {
	if n.GroupIDValid(n.stdioChan) {
		// send mark down to stdio Channel
	}
}
func (n *node) SendStdioLog(markdown string) {
}

func (n *node) RouteMsg2Group(gid string, msg *imnode.IMMsg) {
	// route msg to gid
}
func (n *node) RouteImg2GroupByBytes(gid string, img []byte) {
}
func (n *node) RouteImg2GroupByUrl(gid string, url string) {
}

func (n *node) SendMsg2Group(gid string, msg string) {
	// send msg to gid
}
func (n *node) SendImg2GroupByBytes(gid string, img []byte) {
}

func (n *node) SendImg2GroupByUrl(gid string, url string) {
}

func (n *node) SetMsgHandler(handler func(gid string, msg *imnode.IMMsg)) {
	n.msgHandler = handler
}

func (n *node) Name() string {
	return "Node-" + "blahblahblah"
}
