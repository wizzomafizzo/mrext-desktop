package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	goruntime "runtime"
)

//go:embed frontend/src/assets/images/misterkun.ico
var misterkunIco []byte

//go:embed frontend/src/assets/images/misterkun-online.ico
var misterkunIcoOnline []byte

//go:embed frontend/src/assets/images/misterkun.png
var misterkunPng []byte

//go:embed frontend/src/assets/images/misterkun-online.png
var misterkunPngOnline []byte

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func setIconOffline() {
	//goland:noinspection GoBoolExpressions
	if goruntime.GOOS == "windows" {
		systray.SetIcon(misterkunIco)
	} else {
		systray.SetIcon(misterkunPng)
	}
}

func setIconOnline() {
	//goland:noinspection GoBoolExpressions
	if goruntime.GOOS == "windows" {
		systray.SetIcon(misterkunIcoOnline)
	} else {
		systray.SetIcon(misterkunPngOnline)
	}
}

func systemTray(app *App) func() {
	return func() {
		setIconOffline()

		show := systray.AddMenuItem("Show", "Show The Window")
		systray.AddSeparator()
		exit := systray.AddMenuItem("Exit", "Quit The Program")

		show.Click(func() {
			runtime.WindowShow(app.ctx)
		})
		exit.Click(func() {
			os.Exit(0)
		})

		systray.SetOnClick(func() {
			runtime.WindowShow(app.ctx)
		})
		systray.SetOnRClick(func(menu systray.IMenu) {
			_ = menu.ShowMenu()
		})
	}
}

// startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	systray.Run(systemTray(a), nil)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SetIconOffline() {
	setIconOffline()
}

func (a *App) SetIconOnline() {
	setIconOnline()
}
