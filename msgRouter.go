package main

import "github.com/HelloWorksGroup/IMSuperGroup/imnode"

var group2Node map[string]imnode.IMNode

// 建立群组id与Node的对应
func bindGroupToNode() {

}

func globalMsgRouter(from string, uid string, msg imnode.IMMsg) {
	for _, group := range superGroups {
		var groupHit bool = false
		for _, v := range group {
			if v == from {
				groupHit = true
				break
			}
		}
		if groupHit {
			for _, v := range group {
				if v != from {
					group2Node[v].RouteMsg2Group(v, uid, msg)
				}
			}
		}
	}
}
