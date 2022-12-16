package main

import (
	"flag"
	"log"
)

func main() {
	t := mustToken()
	//tgClient = telegram.New(token)
	//fetcher = fetcher.New()
	//processor = processor.New()
	//consumer.Start(fetcher,processor)
}

func mustToken() string {
	token := flag.String(name:"token-bot-token",value:"",usage:"token for access to telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is empty")
	}
	return *token
}