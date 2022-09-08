package main_test

import (
	"fmt"
	"testing"

	"github.com/Nigh/MiraiGo-Template-Mod/bot"
	"github.com/spf13/viper"
)

func TestGenDevice(t *testing.T) {
	bot.GenRandomDevice()
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
	fmt.Println(testmap)
}
