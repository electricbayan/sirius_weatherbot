package main

import (
	"fmt"
	"log"

	"github.com/electric_bayan/weatherbot/config"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New()
	fmt.Println(conf.Postgres.DbHost)
}
