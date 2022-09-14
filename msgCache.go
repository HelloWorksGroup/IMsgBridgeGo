package main

import kcard "local/khlcard"

type msgDetail struct {
	Uid       string
	Timestamp int64
}

type kookLastMsgs struct {
	Card      kcard.KHLCard
	Uid       int64
	MsgTime   int64
	MsgId     string
	CardStack int
}

type channelInstance struct {
	Id        string // 群号
	Target    string // 转发目标
	InviteUrl string // 邀请链接
	LastMsg   kookLastMsgs
	MsgCache  map[string]msgDetail
}

type AllChannelInstances struct {
	instances map[string]channelInstance
}

var kookLastCache map[string]kookLastMsgs

func kookLastCacheSetup() {
	kookLastCache = make(map[string]kookLastMsgs, 0)
}
