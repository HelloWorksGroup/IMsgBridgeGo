package main

import (
	"fmt"

	kook "kookNode"
	vc "vcNode"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/spf13/viper"
)

var db *scribble.Driver

func convertMap2StrStr(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		switch val := v.(type) {
		case string:
			result[k] = val
		default:
			result[k] = fmt.Sprintf("%v", val)
		}
	}
	return result
}

func convertAny2StrSlice(v interface{}) []string {
	var s []string
	switch val := v.(type) {
	case []string:
		s = val
	case []interface{}:
		for _, v := range val {
			s = append(s, fmt.Sprintf("%v", v))
		}
	default:
		s = append(s, fmt.Sprintf("%v", val))
	}
	return s
}

func GetConfig() {
	db, _ = scribble.New("./database", nil)
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	s := viper.Get("nodes").([]any)
	for _, v := range s {
		node := v.(map[string]any)
		switch node["type"] {
		case "kook":
			n := kook.Setup(node, &gLog)
			n.SetMsgHandler(globalMsgRouter)
			nodes = append(nodes, n)
		case "qq":
		case "vc":
			n := vc.Setup(node, &gLog)
			n.SetMsgHandler(globalMsgRouter)
			nodes = append(nodes, n)
		case "webhook":
		}
		fmt.Println(node["type"].(string) + ":")
		fmt.Print("\t")
		fmt.Println(node)
	}

	s1 := viper.Get("supergroups").([]any)
	for _, groups := range s1 {
		superGroups = append(superGroups, convertAny2StrSlice(groups))
	}

	// msgCacheSetup()
}
