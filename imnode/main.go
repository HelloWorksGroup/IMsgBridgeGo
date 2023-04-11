package imnode

import "fmt"

type IMMsg struct {
	Type string
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
	GroupIDValid(gid string) bool

	RouteMsg2Group(gid string, msg *IMMsg)
	RouteImg2GroupByBytes(gid string, img []byte)
	RouteImg2GroupByUrl(gid string, url string)

	SendMsg2Group(gid string, msg string)
	SendImg2GroupByBytes(gid string, img []byte)
	SendImg2GroupByUrl(gid string, url string)

	SetMsgHandler(func(gid string, msg *IMMsg))

	Name() string
}

func ConvertAnySlice2StrSlice(slice []interface{}) []string {
	var s []string

	for _, v := range slice {
		switch val := v.(type) {
		case string:
			s = append(s, val)
		default:
			s = append(s, fmt.Sprintf("%v", val))
		}
	}

	return s
}

func ConvertSettingGroups2StrSlice(groups any) []string {
	return ConvertAnySlice2StrSlice(groups.([]any))
}
