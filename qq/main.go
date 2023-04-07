package qq

import (
	"github.com/phuslu/log"
)

type node struct {
	logger    *log.Logger
	stdioChan string
	// 用于对device描述文件和登陆session缓存进行唯一命名
	deviceName string
}

func QqbotInit() {
	// OnMsg(qqMsgHandler)
}

func QqbotStart() {
	ConfigInit()
	// 快速初始化
	Init()

	// 初始化 Modules
	StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	UseProtocol(AndroidWatch)

	// 登录
	err := Login()
	if err == nil {
		// 登录成功，保存 token 信息
		SaveToken()
	}

	// 刷新好友列表，群列表
	RefreshList()
}
