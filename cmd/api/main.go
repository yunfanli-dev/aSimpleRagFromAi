package main

import (
	"log"

	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/app"
)

// main starts the API process and exits on bootstrap failure.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
