package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/electric_bayan/weather_bot/config"
)

type Coordinates struct {
	lat float64
	lon float64
}

type WeatherForecast struct {
	AverageTemperature float32
	Rain               float32
	Snow               float32
}

func SendGeocoderRequest(city string) (Coordinates, error) {
	conf := config.New()

	client := &http.Client{}
	// get city coords

	addr := "https://geocode-maps.yandex.ru/1.x/?apikey=" + conf.GeocoderAPIkey + "&geocode=" + city + "&format=json"
	// fmt.Println(addr)
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		fmt.Println("Error with geocode api", err)
	}
	req.Header.Add("Accept-Charset", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error with geocode api", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading json", err)
	}
	// fmt.Println(string(body))

	res_bytes := []byte(body)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(res_bytes, &jsonRes)
	response := jsonRes["response"].(map[string]interface{})
	geoobjectcol := response["GeoObjectCollection"].(map[string]interface{})
	featuremember := geoobjectcol["featureMember"].([]interface{})
	if len(featuremember) > 0 {
		ft := featuremember[0].(map[string]interface{})
		geoobject := ft["GeoObject"].(map[string]interface{})
		pointed := geoobject["Point"].(map[string]interface{})
		position := pointed["pos"].(string)
		coordinates := strings.Split(position, " ")
		lat, _ := strconv.ParseFloat(coordinates[1], 64)
		lon, _ := strconv.ParseFloat(coordinates[0], 64)

		return Coordinates{lat: lat, lon: lon}, nil
	}
	return Coordinates{0, 0}, fmt.Errorf("WrongCity")

}

func SendWeatherRequest(coords Coordinates) {
	currentTime := time.Now()
	addr := "https://api.open-meteo.com/v1/forecast?latitude=" + strconv.FormatFloat(coords.lat, 'f', 6, 64) + "&longitude=" + strconv.FormatFloat(coords.lon, 'f', 6, 64) + "&hourly=temperature_2m,rain,snowfall&start_date=" + currentTime.Format("2006-01-02") + "&end_date=" + currentTime.Format("2006-01-02") + "&time_mode=time_interval"
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		fmt.Println("Error with OpenWeather api", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error with OpenWeather api", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading json", err)
	}
	res_bytes := []byte(body)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(res_bytes, &jsonRes)
	response := jsonRes["hourly"].(map[string]interface{})
	fmt.Println(response)

}
