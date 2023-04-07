package kookNode

import "github.com/lonelyevil/kook"

func (n *node) updateKMsg(msgId string, content string) error {
	return n.session.MessageUpdate((&kook.MessageUpdate{
		MessageUpdateBase: kook.MessageUpdateBase{
			MsgID:   msgId,
			Content: content,
		},
	}))
}

func (n *node) sendKCard(target string, content string) (resp *kook.MessageResp, err error) {
	return n.session.MessageCreate((&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			Type:     kook.MessageTypeCard,
			TargetID: target,
			Content:  content,
		},
	}))
}
func (n *node) sendMarkdown(target string, content string) (resp *kook.MessageResp, err error) {
	return n.session.MessageCreate((&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			Type:     kook.MessageTypeKMarkdown,
			TargetID: target,
			Content:  content,
		},
	}))
}

func (n *node) sendMarkdownDirect(target string, content string) (mr *kook.MessageResp, err error) {
	return n.session.DirectMessageCreate(&kook.DirectMessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			Type:     kook.MessageTypeKMarkdown,
			TargetID: target,
			Content:  content,
		},
	})
}
