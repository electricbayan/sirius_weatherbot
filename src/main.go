package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/electric_bayan/weather_bot/config"
	"github.com/electric_bayan/weather_bot/db"
	"github.com/electric_bayan/weather_bot/fsm"
	"github.com/electric_bayan/weather_bot/weatherapi"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	// coords := weatherapi.SendGeocoderRequest("Yekaterinburg")
	// weatherapi.SendWeatherRequest(coords)
	conf := config.New()

	ctx := context.Background()

	redis_client := fsm.New()

	bot, err := telego.NewBot(conf.TgAPIkey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		fmt.Println("Error with polling", err)
	}
	bh, err := telegohandler.NewBotHandler(bot, updates)
	if err != nil {
		fmt.Println("Error during creatin handler", err)
	}
	// main handler
	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		chatID := message.Chat.ID
		strid := strconv.Itoa(int(chatID))
		current_state, err := redis_client.Get(ctx, strid).Result()

		if err != nil {
			fmt.Println("Error during getting state")
		}
		if message.Text == "/start" {
			_, err := bot.SendMessage(&telego.SendMessageParams{
				ChatID: telego.ChatID{ID: chatID},
				Text:   "Enter your city.",
			})
			if err != nil {
				fmt.Println("Error during requesting msg", err)
			}

			err = redis_client.Set(ctx, strid, "CityWaiting", 0).Err()
			if err != nil {
				fmt.Println("Error with redis", err)
			}
		}
		if current_state == "CityWaiting" {

			coords, err := weatherapi.SendGeocoderRequest(message.Text)
			fmt.Println(coords)
			if err != nil {
				fmt.Println("error with geocoder", err)
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatID},
					Text:   "Wrong City..",
				})
				if err != nil {
					fmt.Println("Error during requesting msg", err)
				}
			} else {
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatID},
					Text:   "How frequently would you like to get updates?.",
				})
				if err != nil {
					fmt.Println("Error during requesting msg", err)
				}
				db.InsertUser(int(chatID), message.Text)
			}
			err = redis_client.Set(ctx, strid, "", 0).Err()
			if err != nil {
				fmt.Println("Error with redis", err)
			}
			// fmt.Println(chatID)
		}
	})

	bh.Start()

	bh.Stop()

}
