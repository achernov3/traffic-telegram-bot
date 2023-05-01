package traffic

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const TOKEN string = "902aafd0-5401-4b29-a8b3-723948b3a5c1"

var stationCodes map[string]string = map[string]string{
	"Бердск":              "s9610394",
	"Речной Вокзал":       "s9610287",
	"Центр":               "s9610286",
	"Новосибирск-главный": "s9610189",
}

type SchedulerTable struct {
	Search struct {
		From struct {
			Title string `json:"title"`
		} `json:"from"`
		To struct {
			Title string `json:"title"`
		} `json:"to"`
		Date any `json:"date"`
	} `json:"search"`
	Segments []struct {
		Stops     string  `json:"stops"`
		Duration  float64 `json:"duration"`
		Days      string  `json:"days"`
		Departure string  `json:"departure"`
		Arrival   string  `json:"arrival"`
		StartDate string  `json:"start_date"`
	} `json:"segments"`
}

func GetSchedule(from, to string) (string, float64, string, string, string) {
	//Функция отправки расписания пользователю
	//Вызывается с двумя параметрами: from и to
	fromDestination := getStationCode(from)
	toDestination := getStationCode(to)
	var stops string
	var duration float64
	var days string
	var departure string
	var arrival string
	var startDate string
	timeNow := time.Now()
	year := strconv.Itoa(timeNow.Year())
	month := strconv.Itoa(int(timeNow.Month()))
	day := strconv.Itoa(timeNow.Day())
	resp, err := http.Get("https://api.rasp.yandex.net/v3.0/search/?apikey=" + TOKEN + "&format=json&from=" + fromDestination + "&to=" + toDestination + "&lang=ru_RU&page=1&date=" + year + "-" + month + "-" + day)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var scheduler SchedulerTable
	json.Unmarshal(body, &scheduler)

	for i := 0; i != len(scheduler.Segments); i++ {
		stops = scheduler.Segments[i].Stops
		duration = scheduler.Segments[i].Duration
		days = scheduler.Segments[i].Days
		departure = scheduler.Segments[i].Departure[11:19]
		arrival = scheduler.Segments[i].Arrival[11:19]
		startDate = scheduler.Segments[i].StartDate
		departureYear, _ := strconv.Atoi(startDate[:4])
		departureMonth := time.Time.Month(timeNow)
		departureDay, _ := strconv.Atoi(startDate[8:10])
		departureHour, _ := strconv.Atoi(departure[0:2])
		departureMinute, _ := strconv.Atoi(departure[3:5])
		departureDate := time.Date(departureYear, departureMonth, departureDay, departureHour, departureMinute, 0, 0, time.Local)
		if timeNow.Before(departureDate) {
			break
		}
	}
	return stops, duration, days, departure, arrival
}

func getStationCode(destination string) string {
	var code string
	for station, code := range stationCodes {
		if station == destination {
			return code
		}
	}
	return code
}
