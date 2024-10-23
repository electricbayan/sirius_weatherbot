package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/electric_bayan/weather_bot/config"
)

func InsertUser(user_id int, lat float64, lon float64) sql.Result {
	conf := config.New()
	psqlInfo := fmt.Sprintf("host=localhost port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgres.DbPort, conf.Postgres.DbUser, conf.Postgres.DbPass, conf.Postgres.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error with postgres", err)
	}
	defer db.Close()

	currentTime := time.Now()
	currentTime = currentTime.Add(-time.Hour * 5)
	stmt := fmt.Sprintf("INSERT INTO users(id, frequency, last_update, lat, lon) VALUES(%d, null, '%s', %.4f, %.4f)", user_id, currentTime.Format("2006-01-02 15:04:05"), lat, lon)
	result, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("error during insert user", err)
	}
	return result

}

func UpdateFrequency(user_id int, frequency int) sql.Result {
	conf := config.New()
	psqlInfo := fmt.Sprintf("host=localhost port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgres.DbPort, conf.Postgres.DbUser, conf.Postgres.DbPass, conf.Postgres.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error with postgres", err)
	}
	defer db.Close()
	stmt := fmt.Sprintf("UPDATE USERS SET frequency=%d WHERE id=%d;", frequency, user_id)
	result, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("error during insert user", err)
	}
	return result
}

func SelectNewMessages() *sql.Rows {
	conf := config.New()
	psqlInfo := fmt.Sprintf("host=localhost port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgres.DbPort, conf.Postgres.DbUser, conf.Postgres.DbPass, conf.Postgres.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error with postgres", err)
	}
	defer db.Close()

	currentTime := time.Now()
	currentTime = currentTime.Add(-time.Hour * 5).Add(time.Second * 2)

	res, err := db.Query("SELECT id, lat, lon FROM users WHERE last_update < to_timestamp('" + currentTime.Format("2006-01-02 15:04:05") + "', 'YYYY-MM-DD hh24:mi:ss') - (frequency || ' seconds')::interval")
	if err != nil {
		fmt.Println("Error during database request", err)
	}
	upd, err := db.Query("update users set last_update=current_timestamp where id in (SELECT id FROM users WHERE last_update < to_timestamp('" + currentTime.Format("2006-01-02 15:04:05") + "', 'YYYY-MM-DD hh24:mi:ss') - (frequency || ' seconds')::interval)")
	if err != nil {
		fmt.Println(upd, err)
	}
	return res
}
