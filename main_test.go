package main_test

import (
	"fmt"
	qq "qqNode"
	"testing"

	"github.com/spf13/viper"
)

func TestGenDevice(t *testing.T) {
	qq.GenRandomDevice()
}

var testmap map[string]string

func TestViper(t *testing.T) {
	testmap = make(map[string]string, 0)
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	s := viper.Get("kook2qq").(map[string]any)
	for k, v := range s {
		vs := v.(string)
		if k != v {
			if _, ok := testmap[k]; !ok {
				testmap[k] = vs
			}
			if _, ok := testmap[vs]; !ok {
				testmap[vs] = k
			}
		}
	}

	s1 := viper.Get("routes").([]any)
	for _, newmap := range s1 {
		// fmt.Println(newmap)
		route := newmap.(map[string]any)
		// fmt.Println(route["type"])
		if route["type"] == "kook2qq" {
			if route["host"] != nil && route["qqgroup"] != nil {
				testmap[route["host"].(string)] = route["qqgroup"].(string)
				testmap[route["qqgroup"].(string)] = route["host"].(string)
			}
		}
	}
	fmt.Println(testmap)
}
