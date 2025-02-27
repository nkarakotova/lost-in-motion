package main

import (
	"fmt"
	"lim/cmd/registry"
	"os"

	"github.com/charmbracelet/log"

	tech_registry "lim/internal/lim-console/registry"
	menu "lim/internal/lim-console/techUI"
)

func main() {
	app := registry.App{}

	errConfig := app.Config.ParseConfig("config.json", "./")
	if errConfig != nil {
		f, err := os.OpenFile(registry.DefaultLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %s", fmt.Errorf("%s", err.Error()))
		}
		defer f.Close()
		log.SetOutput(f)
		log.SetReportCaller(true)
		log.Fatal(errConfig)
	}

	err := app.Run()
	if err != nil {
		log.Fatal(fmt.Errorf("%s", err.Error()))
	}

	services := tech_registry.AppServiceFields(*app.Services)

	if app.Config.Mode == "tech" {
		app.Logger.Info("Start with tech ui!")
		menu.RunMenu(&services, app.Config.AdminLogin, app.Config.AdminPassword)
	} else {
		app.Logger.Error("Wrong app mode", "mode", app.Config.Mode)
	}
}
