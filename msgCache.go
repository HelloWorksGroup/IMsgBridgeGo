package main

import (
	kcard "local/khlcard"
	"strconv"
	"time"
)

type kookLastMsgs struct {
	Card      kcard.KHLCard
	Uid       int64
	MsgTime   int64
	MsgId     string
	CardStack int
}

type msgDetail struct {
	MsgId     string
	Uid       string
	Timestamp int64
}

type channelInstance struct {
	Id        string // 群号
	Target    string // 转发目标
	InviteUrl string // 邀请链接
	LastMsg   kookLastMsgs
	MsgCache  []msgDetail
}

type AllChannelInstances struct {
	instances []channelInstance
}

func (s *AllChannelInstances) GetMsg(gid string, mid string, uid string) {
	var found int = -1
	for k, v := range s.instances {
		if v.Id == gid {
			found = k
			break
		}
	}
	if found < 0 {
		s.instances = append(s.instances, channelInstance{Id: gid, LastMsg: kookLastMsgs{}, MsgCache: make([]msgDetail, 0)})
		for k, v := range s.instances {
			if v.Id == gid {
				found = k
				break
			}
		}
	}
	s.instances[found].MsgCache = append(s.instances[found].MsgCache, msgDetail{MsgId: mid, Uid: uid, Timestamp: time.Now().Unix()})
}
func (s *AllChannelInstances) WhomReply(gid string, mid string) (uid string) {
	var foundGid int = -1
	for k, v := range s.instances {
		if v.Id == gid {
			foundGid = k
			break
		}
	}
	if foundGid >= 0 {
		for _, v := range s.instances[foundGid].MsgCache {
			if v.MsgId == mid {
				return v.Uid
			}
		}
	}
	return ""
}

func (s *AllChannelInstances) gc() {
	now := time.Now().Unix()
	var startCacheDepth, endCacheDepth int = 0, 0
	for _, v := range s.instances {
		startCacheDepth += len(v.MsgCache)
	}
	for k := range s.instances {
		temp := s.instances[k].MsgCache[:0]
		for km := range s.instances[k].MsgCache {
			if now-s.instances[k].MsgCache[km].Timestamp < 86400 {
				temp = append(temp, s.instances[k].MsgCache[km])
			}
		}
		s.instances[k].MsgCache = temp
	}
	for _, v := range s.instances {
		endCacheDepth += len(v.MsgCache)
	}
	kookLog("执行GC......\nGC前消息缓存深度为`" + strconv.Itoa(startCacheDepth) + "`\nGC后消息缓存深度为`" + strconv.Itoa(endCacheDepth) + "`")
}

func (s *AllChannelInstances) init() {
	s.instances = make([]channelInstance, 0)
	go func() {
		ticker := time.NewTicker(60 * time.Minute)
		for range ticker.C {
			s.gc()
		}
		ticker.Stop()
	}()
}

var kookLastCache map[string]kookLastMsgs

func kookLastCacheSetup() {
	kookLastCache = make(map[string]kookLastMsgs, 0)
}

var msgCache AllChannelInstances

func msgCacheSetup() {
	msgCache.init()
}
