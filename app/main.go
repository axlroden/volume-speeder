//go:build windows || darwin

package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	if err := wails.Run(getAppOptions(app)); err != nil {
		log.Fatal(err)
	}
}
