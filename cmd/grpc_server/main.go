package main

import (
	"log"

	"github.com/vbulash/chat-server/internal/api"
)

func main() {
	err := api.AppRun()
	if err != nil {
		log.Fatal(err)
	}
}
