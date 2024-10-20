package main

import (
	"log"

	"github.com/electric_bayan/weatherbot/clients/telegram"
	"github.com/electric_bayan/weatherbot/config"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New()
	tg_client := telegram.New(conf.TgAPIkey, conf.Host)
}
