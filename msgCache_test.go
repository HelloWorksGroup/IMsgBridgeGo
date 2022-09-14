package main_test

import (
	"fmt"
	kcard "local/khlcard"
	"testing"
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
	s.instances[found].MsgCache = append(s.instances[found].MsgCache, msgDetail{MsgId: mid, Uid: uid})
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

func TestMsgCache(t *testing.T) {
	var msgCache AllChannelInstances
	msgCache.GetMsg("0001", "msg001", "qq123456")
	msgCache.GetMsg("0001", "msg002", "qq223456")
	msgCache.GetMsg("0001", "msg003", "qq323456")
	msgCache.GetMsg("0001", "msg004", "qq123456")
	msgCache.GetMsg("0002", "msg001", "qq123457")
	msgCache.GetMsg("0002", "msg002", "qq223457")
	msgCache.GetMsg("0002", "msg003", "qq323457")
	msgCache.GetMsg("0002", "msg004", "qq123457")

	fmt.Println(msgCache.WhomReply("0002", "msg002"))
	fmt.Println(msgCache.WhomReply("0001", "msg004"))
	fmt.Println(msgCache.WhomReply("0001", "msg006"))
}
