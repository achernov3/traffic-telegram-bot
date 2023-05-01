package traffic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

func GetMe() {
	/*
		Возвращает информацию о боте
	*/
	resp, err := http.Get(URL + "getMe")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
}

func GetUpdates() {
	/*
		Бот получает обновления; кто написал и т.д.
		Информация передается в структуру типа Response
	*/
	var response Response
	resp, err := http.Get(URL + "getUpdates")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &response)
	lastTextUpdates := strings.Fields(response.Result[len(response.Result)-1].Message.Text)
	chatId := response.Result[len(response.Result)-1].Message.Chat.ID
	from := strings.Title(lastTextUpdates[0])
	to := strings.Title(lastTextUpdates[1])
	SendMessage(chatId, from, to)
}

func SendMessage(chatId int, from, to string) {
	//Функция отправки сообщения
	//Запрашивает два параметра: chat_id (забрать из GetUpdates()), from, to - пункт отправления и пункт назначения
	stops, duration, days, departure, arrival := GetSchedule(from, to)
	text := fmt.Sprintf("Станция отправления: %v\nСтанция прибытия: %v\nОстановки: %v\nДни: %v\nВремя в пути: %v\nОтправление: %v\nПрибытие: %v\n", from, to, stops, days, duration, departure, arrival)
	jsonData := map[string]interface{}{
		"chat_id": chatId,
		"text":    text,
	}
	answer, _ := json.Marshal(jsonData) // функция, которая делает из map json объект, который можно отправлять по сети
	resp, _ := http.Post(URL+"sendMessage", "application/json", bytes.NewBuffer(answer))
	if resp.StatusCode == 200 {
		fmt.Println("Сообщение пользователю отправилось успешно")
	}
}
