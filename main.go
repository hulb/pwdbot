package main

import (
	"log"
	"os"
	"pwdbot/handlers"
	"pwdbot/structs"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	BOTTOKEN string
	FILENAME string
)

func init() {
	BOTTOKEN = os.Getenv("BOTTOKEN")
	if BOTTOKEN == "" {
		log.Fatalln("bot token required.exit")
	}
}

func main() {
	// DEBUG set http proxy
	// proxy, _ := url.Parse("http://127.0.0.1:7890")
	// tr := &http.Transport{
	// 	Proxy:           http.ProxyURL(proxy),
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{
	// 	Transport: tr,
	// }

	b, err := tb.NewBot(tb.Settings{
		Token:  BOTTOKEN,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	structs.RegisterBot(b)

	for cmd, handlerFunc := range handlers.CmdHandler {
		b.Handle(cmd, handlerFunc)
	}

	b.Start()
}
