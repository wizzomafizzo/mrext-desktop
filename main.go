package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
	"net/http"
)

//go:embed all:frontend/dist
var assets embed.FS

func takeScreenshot() error {
	config, _, err := LoadIni()
	if err != nil {
		return err
	}

	_, err = http.Post(
		"http://"+config.Host+"/api/screenshots",
		"application/json",
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func hotkeys() {
	hk := hotkey.New(nil, hotkey.KeyF11)
	err := hk.Register()
	if err != nil {
		fmt.Print("Error registering hotkey: ", err)
		panic(err)
	}
	for {
		<-hk.Keydown()
		err = takeScreenshot()
		if err != nil {
			fmt.Print("Error taking screenshot: ", err)
			panic(err)
		}
	}
}

func main() {
	go func() {
		mainthread.Init(hotkeys)
	}()

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Remote",
		Width:  600,
		Height: 200,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 20, G: 20, B: 20, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		fmt.Println(err)
		println("Error:", err.Error())
	}
}
