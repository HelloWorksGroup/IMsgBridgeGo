package main

import (
	"github.com/Nigh/MiraiGo-Template-Mod/bot"
	"github.com/Nigh/MiraiGo-Template-Mod/config"

	qq "local/rt"
)

func qqbotInit() {
	// utils.WriteLogToFS(utils.LogInfoLevel, utils.LogWithStack)
	qq.SetGroupID(qqGroupCode)
	qq.OnMsg(MsgRouteQQ2KOOK)
}

func qqGetKOOK(content string) {
	qq.MsgRouteKOOK2QQ(content)
}

func qqbotStart() {
	config.Init()
	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	err := bot.Login()
	if err == nil {
		// 登录成功，保存 token 信息
		bot.SaveToken()
	}

	// 刷新好友列表，群列表
	bot.RefreshList()
}
func qqbotSayhello() {
}

func qqbotSaybye() {
}

func qqbotStop() {
}
