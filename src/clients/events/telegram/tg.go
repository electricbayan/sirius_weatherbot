package telegram

import (
	"github.com/electric_bayan/weatherbot/clients/telegram"
)

type Proccessor struct {
	tg     *telegram.Client
	offset int
}

func New(client *telegram.Client) {

}
