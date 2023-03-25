package main

import (
	qq "qqNode"
)

func qqbotInit() {
	// utils.WriteLogToFS(utils.LogInfoLevel, utils.LogWithStack)
	qq.OnMsg(qqMsgHandler)
}

func qqbotStart() {
	qq.ConfigInit()
	// 快速初始化
	qq.Init()

	// 初始化 Modules
	qq.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	qq.UseProtocol(qq.AndroidWatch)

	// 登录
	err := qq.Login()
	if err == nil {
		// 登录成功，保存 token 信息
		qq.SaveToken()
	}

	// 刷新好友列表，群列表
	qq.RefreshList()
}
func qqbotSayhello() {
}

func qqbotSaybye() {
}

func qqbotStop() {
}
