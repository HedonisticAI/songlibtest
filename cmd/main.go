package main

import (
	"log"
	"songlibtest/config"
	"songlibtest/internal/app"
	"songlibtest/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	Cfg, err := config.NewConfig()
	if err != nil {
		log.Print(err.Error())
	}
	Logger := logger.New(Cfg)
	App, err := app.New(*Cfg, Logger)
	if err != nil {
		Logger.Info("Couldn't build App")
		Logger.Debug(err)
	} else {
		Logger.Info("App running")
		App.Run()
	}

}
