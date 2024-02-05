package main

import (
	"log"

	"github.com/tigorlazuardi/redmage/app"
)

func main() {
	if err := app.Start(nil); err != nil {
		log.Fatal(err)
	}
}
