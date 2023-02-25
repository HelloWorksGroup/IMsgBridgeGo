package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/lonelyevil/kook"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/spf13/viper"
)

// stdout频道
var stdoutChannel string

// 转发map
var kook2qqRouteMap map[string]string

// [qqgroup] [vcsetting]
type vocechatInstance struct {
	Url    string
	Secret string
	Gid    string
}

var vc2qqRouteMap map[string]vocechatInstance

// 邀请map
var kookInviteUrl map[string]string

var masterID string
var botID string

var localSession *kook.Session

var token string

var db *scribble.Driver

func routeMapInit() {
	vc2qqRouteMap = make(map[string]vocechatInstance, 0)
	kook2qqRouteMap = make(map[string]string, 0)
	kookInviteUrl = make(map[string]string, 0)
}
func routeMapSetupOld() {
	s := viper.Get("kook2qq").(map[string]any)
	for k, v := range s {
		vs := v.(string)
		if k != v {
			if _, ok := kook2qqRouteMap[k]; !ok {
				kook2qqRouteMap[k] = vs
			}
			if _, ok := kook2qqRouteMap[vs]; !ok {
				kook2qqRouteMap[vs] = k
			}
		}
	}
}
func kookInviteUrlSetup() {
	s := viper.Get("kookinvite").(map[string]any)
	for k, v := range s {
		vs := v.(string)
		if _, ok := kookInviteUrl[k]; !ok {
			kookInviteUrl[k] = vs
		}
	}
}

var vcCount int

func RouteMapSetup() {
	vcCount = 0
	s := viper.Get("routes").([]any)
	for _, newmap := range s {
		fmt.Println(newmap)
		route := newmap.(map[string]any)
		fmt.Println(route["type"])
		if route["type"] == "kook2qq" {
			if route["host"] != nil && route["qqgroup"] != nil {
				kook2qqRouteMap[route["host"].(string)] = route["qqgroup"].(string)
				kook2qqRouteMap[route["qqgroup"].(string)] = route["host"].(string)
				if route["hostinvite"] != nil {
					kookInviteUrl[route["host"].(string)] = route["hostinvite"].(string)
				}
			}
		}
		if route["type"] == "vc2qq" {
			if route["gid"] != nil && route["vcapi"] != nil && route["secret"] != nil && route["qqgroup"] != nil {
				vc2qqRouteMap[route["qqgroup"].(string)] = vocechatInstance{Url: route["vcapi"].(string), Secret: route["secret"].(string), Gid: route["gid"].(string)}
				vcCount += 1
			}
		}
	}
}

func GetConfig() {
	rand.Seed(time.Now().UnixNano())
	db, _ = scribble.New("./database", nil)
	viper.SetDefault("token", "0")
	viper.SetDefault("stdoutChannel", "0")
	viper.SetDefault("masterID", "")
	viper.SetDefault("oldversion", "0.0.0")
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	masterID = viper.Get("masterID").(string)
	stdoutChannel = viper.Get("stdoutChannel").(string)
	vcport, _ := strconv.Atoi(viper.Get("vcport").(string))
	fmt.Println("stdoutChannel=" + stdoutChannel)
	token = viper.Get("token").(string)

	routeMapInit()
	routeMapSetupOld()
	RouteMapSetup()

	kookInviteUrlSetup()
	kookLastCacheSetup()
	msgCacheSetup()

	if vcCount > 0 && vcport > 0 {
		go voceChatBot(vcport)
	}
}

func beforeShutdown() {
	msgCache.Backup()
}
