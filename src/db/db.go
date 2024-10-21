package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/electric_bayan/weather_bot/config"
)

func InsertUser(user_id int, city string) sql.Result {
	conf := config.New()
	psqlInfo := fmt.Sprintf("host=localhost port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgres.DbPort, conf.Postgres.DbUser, conf.Postgres.DbPass, conf.Postgres.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error with postgres", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to postgres.")

	currentTime := time.Now()
	stmt := fmt.Sprintf("INSERT INTO users(id, city, frequency, last_update) VALUES(%d, '%s', null, '%s')", user_id, city, currentTime.Format("2006-01-02 15:04:05"))
	result, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("error during insert user", err)
	}
	return result

}

// func insert
