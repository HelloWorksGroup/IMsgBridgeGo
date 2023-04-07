package main_test

import (
	"fmt"
	qq "qqNode"
	"testing"

	"github.com/HelloWorksGroup/IMSuperGroup/imnode"
	"github.com/spf13/viper"
)

// 生成QQ随机设备描述文件
func TestGenDevice(t *testing.T) {
	qq.GenRandomDevice()
}

// 本地登陆QQ，以取得登陆session
func TestQQLogin(t *testing.T) {

}

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

func TestViper(t *testing.T) {
	nodes := make([]imnode.IMNode, 0)
	superGroups := make([][]string, 0)
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
			nodes = append(nodes, kook.Setup(convertMap2StrStr(node), nil))
		case "qq":
		case "vc":
		case "webhook":
		}
		fmt.Println(node["type"].(string) + ":")
		fmt.Print("\t")
		fmt.Println(node)
	}

	s1 := viper.Get("supergroups").([]any)
	for _, groups := range s1 {
		superGroups = append(superGroups, convertAny2StrSlice(groups))
		for _, gid := range groups.([]any) {
			fmt.Println(gid.(string))
		}
	}
	fmt.Println(nodes[0].Name())
	fmt.Println(superGroups)
}
