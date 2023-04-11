package main

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/HelloWorksGroup/IMSuperGroup/imnode"
)

func escapeToCleanUnicode(raw string) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, str)
	return clean, nil
}

func globalMsgRouter(gid string, msg *imnode.IMMsg) {
	cleanName, err := escapeToCleanUnicode(msg.ShowName)
	if err == nil {
		msg.ShowName = cleanName
	}
	for _, group := range superGroups {
		var groupHit bool = false
		// if this group is contained in a supergroup
		for _, v := range group {
			if v == gid {
				groupHit = true
				break
			}
		}
		// send msg to other group in this supergroup
		if groupHit {
			for _, grp := range group {
				if grp != gid {
					for _, node := range nodes {
						if node.GroupIDValid(grp) {
							node.RouteMsg2Group(grp, msg)
						}
					}
				}
			}
		}
	}
}
