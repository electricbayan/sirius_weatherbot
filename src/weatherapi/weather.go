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

type WeatherForecast struct {
	AverageTemperature float64
	RainStart          int
	RainStop           int
	IsRain             bool
	CurrentTemp        float64
}

func SendGeocoderRequest(city string) (float64, float64, error) {
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

		return lat, lon, nil
	}
	return 0, 0, fmt.Errorf("WrongCity")

}

func SendWeatherRequest(lat float64, lon float64) WeatherForecast {
	currentTime := time.Now()
	addr := "https://api.open-meteo.com/v1/forecast?latitude=" + strconv.FormatFloat(lat, 'f', 6, 64) + "&longitude=" + strconv.FormatFloat(lon, 'f', 6, 64) + "&hourly=temperature_2m,rain,snowfall&start_date=" + currentTime.Format("2006-01-02") + "&end_date=" + currentTime.Format("2006-01-02") + "&time_mode=time_interval&timezone=Europe%2FMoscow"
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
	temperature := response["temperature_2m"].([]interface{})

	// fmt.Println(response)
	count := 0
	sum := 0.0
	for _, temp := range temperature {
		count += 1
		sum += temp.(float64)
	}
	avg_temp := sum / float64(count)

	rain := response["rain"].([]interface{})
	rain_started := false
	rain_start_time := 0
	rain_stop_time := 0
	for i, rainy := range rain {
		rainy = rainy.(float64)
		if rain_started {
			if rainy == 0.0 {
				rain_stop_time = i
			}
		} else if rainy != 0.0 {
			rain_started = true
			rain_start_time = i
		}
	}
	cur_time := time.Now()
	if cur_time.Hour()-2 > 0 {
		cur_temp := temperature[cur_time.Hour()-3].(float64)
		return WeatherForecast{
			AverageTemperature: avg_temp,
			RainStart:          rain_start_time,
			RainStop:           rain_stop_time,
			IsRain:             rain_started,
			CurrentTemp:        cur_temp,
		}
	} else {
		currentTime = currentTime.Add(-time.Hour * 24)
		addr := "https://api.open-meteo.com/v1/forecast?latitude=" + strconv.FormatFloat(lat, 'f', 6, 64) + "&longitude=" + strconv.FormatFloat(lon, 'f', 6, 64) + "&hourly=temperature_2m,rain,snowfall&start_date=" + currentTime.Format("2006-01-02") + "&end_date=" + currentTime.Format("2006-01-02") + "&time_mode=time_interval&timezone=Europe%2FMoscow"
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
		temperature := response["temperature_2m"].([]interface{})

		cur_temp := temperature[cur_time.Hour()-3+24].(float64)
		return WeatherForecast{
			AverageTemperature: avg_temp,
			RainStart:          rain_start_time,
			RainStop:           rain_stop_time,
			IsRain:             rain_started,
			CurrentTemp:        cur_temp,
		}
	}
}
