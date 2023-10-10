package main

import (
	"context"
	_ "embed"
	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
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

func iniFilename() string {
	path, err := os.Executable()
	if err != nil {
		path = ""
	}

	return filepath.Join(filepath.Dir(path), "remote.ini")
}

type Config struct {
	Host           string `ini:"host"`
	SystemFilename string `ini:"system_filename"`
	GameFilename   string `ini:"game_filename"`
}

const SectionRemote = "remote"

func LoadIni() (Config, *ini.File, error) {
	cfg, err := ini.Load(iniFilename())
	if err != nil {
		cfg = ini.Empty()
	}

	section, err := cfg.GetSection(SectionRemote)
	if err != nil {
		section, err = cfg.NewSection(SectionRemote)
		if err != nil {
			return Config{}, nil, err
		}

		data := Config{
			Host:           "mister:8182",
			SystemFilename: "active_system.txt",
			GameFilename:   "active_game.txt",
		}

		err = section.ReflectFrom(&data)
		if err != nil {
			return Config{}, nil, err
		}

		err = cfg.SaveTo(iniFilename())
		if err != nil {
			return Config{}, nil, err
		}
	}

	config := Config{}
	err = section.MapTo(&config)
	if err != nil {
		return Config{}, nil, err
	}

	return config, cfg, nil
}

func SaveIni(config Config) error {
	_, cfg, err := LoadIni()
	if err != nil {
		return err
	}

	section, err := cfg.GetSection(SectionRemote)
	if err != nil {
		return err
	}

	err = section.ReflectFrom(&config)
	if err != nil {
		return err
	}

	err = cfg.SaveTo(iniFilename())
	if err != nil {
		return err
	}

	return nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	systray.Run(systemTray(a), nil)
}

func (a *App) SetIconOffline() {
	setIconOffline()
}

func (a *App) SetIconOnline() {
	setIconOnline()
}

func (a *App) GetHost() (string, error) {
	config, _, err := LoadIni()
	if err != nil {
		return "", err
	}
	return config.Host, nil
}

func (a *App) SetHost(host string) error {
	config, _, err := LoadIni()
	if err != nil {
		return err
	}

	config.Host = host

	err = SaveIni(config)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) WriteSystem(name string) error {
	path, err := os.Executable()
	if err != nil {
		path = ""
	}

	config, _, err := LoadIni()
	if err != nil {
		return err
	}

	path = filepath.Join(filepath.Dir(path), config.SystemFilename)
	err = os.WriteFile(path, []byte(name), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) WriteGame(name string) error {
	path, err := os.Executable()
	if err != nil {
		path = ""
	}

	config, _, err := LoadIni()
	if err != nil {
		return err
	}

	path = filepath.Join(filepath.Dir(path), config.GameFilename)
	err = os.WriteFile(path, []byte(name), 0644)
	if err != nil {
		return err
	}

	return nil
}
