package main

import "github.com/lonelyevil/kook"

func updateKMsg(msgId string, content string) error {
	return localSession.MessageUpdate((&kook.MessageUpdate{
		MessageUpdateBase: kook.MessageUpdateBase{
			MsgID:   msgId,
			Content: content,
		},
	}))
}

func sendKCard(target string, content string) (resp *kook.MessageResp, err error) {
	return localSession.MessageCreate((&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			Type:     kook.MessageTypeCard,
			TargetID: target,
			Content:  content,
		},
	}))
}
func sendMarkdown(target string, content string) (resp *kook.MessageResp, err error) {
	return localSession.MessageCreate((&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			Type:     kook.MessageTypeKMarkdown,
			TargetID: target,
			Content:  content,
		},
	}))
}

func sendMarkdownDirect(target string, content string) (mr *kook.MessageResp, err error) {
	return localSession.DirectMessageCreate(&kook.DirectMessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			Type:     kook.MessageTypeKMarkdown,
			TargetID: target,
			Content:  content,
		},
	})
}
