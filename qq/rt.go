package route

// 本模块用于将QQ消息转发至KOOK，并将KOOK消息转发至QQ

import (
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &rt{}
var logger = utils.GetModuleLogger("kook.route")

type rt struct {
}

func (a *rt) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "kook.route",
		Instance: instance,
	}
}

func (a *rt) Init() {

}

func (a *rt) PostInit() {
}

var stop bool = true

func (a *rt) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		if stop {
			return
		}
		m := message.NewSendingMessage().Append(message.NewText("hello"))
		c.SendGroupMessage(msg.GroupCode, m)
	})
}

func (a *rt) Start(bot *bot.Bot) {
}

func (a *rt) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}
