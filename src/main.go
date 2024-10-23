package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/electric_bayan/weather_bot/config"
	"github.com/electric_bayan/weather_bot/db"
	"github.com/electric_bayan/weather_bot/fsm"
	"github.com/electric_bayan/weather_bot/weatherapi"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	// coords := weatherapi.SendGeocoderRequest("Yekaterinburg")
	// weatherapi.SendWeatherRequest(coords)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			telego.InlineKeyboardButton{
				Text:         "Per minute",
				CallbackData: "minute",
			},
		),
		tu.InlineKeyboardRow(
			telego.InlineKeyboardButton{
				Text:         "Daily",
				CallbackData: "daily",
			},
			telego.InlineKeyboardButton{
				Text:         "Weekly",
				CallbackData: "weekly",
			},
		),
	)
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

			lat, lon, err := weatherapi.SendGeocoderRequest(message.Text)
			if err != nil {
				fmt.Println("error with geocoder", err)
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatID},
					Text:   "Wrong City.",
				})
				if err != nil {
					fmt.Println("Error during requesting msg", err)
				}
			} else {
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID:      telego.ChatID{ID: chatID},
					Text:        "How frequently would you like to get updates?",
					ReplyMarkup: keyboard,
				})
				if err != nil {
					fmt.Println("Error during requesting msg", err)
				}
				db.InsertUser(int(chatID), lat, lon)

				forecast := weatherapi.SendWeatherRequest(lat, lon)
				fmt.Println(forecast.AverageTemperature, forecast.RainStart)
			}
			err = redis_client.Set(ctx, strid, "", 0).Err()
			if err != nil {
				fmt.Println("Error with redis", err)
			}
		}
		if message.Text == "/frequency" {
			_, err := bot.SendMessage(&telego.SendMessageParams{
				ChatID:      telego.ChatID{ID: chatID},
				Text:        "How frequently would you like to get updates?",
				ReplyMarkup: keyboard,
			})
			if err != nil {
				fmt.Println("Error during requesting msg", err)
			}
		}
	})
	bh.HandleCallbackQuery(func(bot *telego.Bot, callback telego.CallbackQuery) {
		chatID := callback.Message.GetChat().ID
		frequency := 86400
		if callback.Data == "minute" {
			frequency = 60
		} else if callback.Data == "daily" {
			frequency = 86400
		} else if callback.Data == "weekly" {
			frequency = 604800
		}
		db.UpdateFrequency(int(chatID), frequency)

		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(callback.ID).WithText(""))
		_, err := bot.SendMessage(&telego.SendMessageParams{
			ChatID: telego.ChatID{ID: chatID},
			Text:   "You have subscribed.",
		})
		if err != nil {
			fmt.Println("Error during requesting msg", err)
		}
	})
	go weatherSender(*bot)
	bh.Start()

	bh.Stop()

}

func weatherSender(bot telego.Bot) {
	for {
		results := db.SelectNewMessages()
		for results.Next() {
			var id int
			var lat, lon float64
			results.Scan(&id, &lat, &lon)

			forecast := weatherapi.SendWeatherRequest(lat, lon)
			if forecast.IsRain {
				text := fmt.Sprintf("Current day forecast:\nAverage temperature: %.2f.\nCurrent temperature: %.2f\nThere will be no rain.", forecast.AverageTemperature, forecast.CurrentTemp)
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: int64(id)},
					Text:   text,
				})
				if err != nil {
					fmt.Println("Error during requesting msg", err)
				}
			} else {
				text := fmt.Sprintf("Current day forecast:\nAverage temperature: %.2f. \nCurrent temperature: %.2f\nThere is rain from %d to %d", forecast.AverageTemperature, forecast.CurrentTemp, forecast.RainStart, forecast.RainStop)
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: int64(id)},
					Text:   text,
				})
				if err != nil {
					fmt.Println("Error during requesting msg", err)
				}
			}

		}

		time.Sleep(30 * time.Second)
	}
}
