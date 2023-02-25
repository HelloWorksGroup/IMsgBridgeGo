package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func voceChatBot(port int) {
	listen(port)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("headers: %v\n", r.Header)
	fmt.Printf("body: %v\n", r.Body)

	if r.Header.Get("Content-Type") != "application/json" {
		return
	}
	// decoder := json.NewDecoder(r.Body)
	// switch r.Header.Get("User-Agent") {
	// case "git-oschina-hook":
	// 	var t webhookGiteeJSON
	// 	decoder.Decode(&t)
	// 	go giteeHandler(t)
	// }
}

func listen(port int) {
	fmt.Println("Listen on port:", port)
	http.HandleFunc("/", handleWebhook)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
