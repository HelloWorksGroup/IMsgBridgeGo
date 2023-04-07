package imnode

type IMMsg struct {
	// 消息发送者的显示名称
	ShowName string
	// 消息发送者在宿主端的UID
	UID string
	// At的消息的发送者在宿主端的UID
	AtUID string
	// 消息内容
	Content string
}

type IMNode interface {
	Start() error
	BeforeStop()
	Stop()

	Stdio(markdown string)
	SendStdioLog(markdown string)

	RouteMsg2Group(gid string, uid string, msg IMMsg)
	RouteImg2GroupByBytes(gid string, img []byte)
	RouteImg2GroupByUrl(gid string, url string)

	SendMsg2Group(gid string, msg string)
	SendImg2GroupByBytes(gid string, img []byte)
	SendImg2GroupByUrl(gid string, url string)

	Name() string
}
