package main

import (
	"log"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
