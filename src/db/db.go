package db

import (
	"database/sql"
	"fmt"

	"github.com/electric_bayan/weather_bot/config"
)

func New() *sql.DB {
	conf := config.New()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgres.DbHost, conf.Postgres.DbPort, conf.Postgres.DbUser, conf.Postgres.DbPass, conf.Postgres.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error with postgres", err)
	}
	// defer db.Close()

	fmt.Println("Successfully connected to postgres.")
	return db
}
