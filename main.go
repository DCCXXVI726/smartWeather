package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Discription struct {
	Des string `json:"description"`
}

type Weather struct {
	Temp  float64 `json:"temp"`
	Feel  float64 `json:"feels_like"`
	Press int     `json:"pressure"`
	Hum   int     `json:"humidity"`
}

type Winde struct {
	Speed int `json:"speed"`
}

type Data struct {
	Disc   []Discription `json:"weather"`
	Main   Weather       `json:"main"`
	Wind   Winde         `json:"wind"`
	MyDisc string
	Name   string `json:"name"`
}

func getWeather(city string) Data {
	var data Data

	KEY := "&appid=86c0cb2383f69fc2f22f63961ba83dc8&units=metric&lang=ru"
	APIURL := "http://api.openweathermap.org/data/2.5/weather?q="
	APIURL = APIURL + city + KEY

	req, err := http.NewRequest(http.MethodGet, APIURL, nil)
	if err != nil {
		return data
	}

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return data
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&data)
	if data.Disc != nil {
		data.MyDisc = data.Disc[0].Des
	}

	return data
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.FormValue("city")
	data := getWeather(city)

	if data.Disc == nil {
		http.Redirect(w, r, "/empty", http.StatusFound)
		return
	}

	t, err := template.ParseFiles("templates/weather.html")
	if err != nil {
		err = fmt.Errorf("can't parse weather.html: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	t.ExecuteTemplate(w, "weather", data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		err = fmt.Errorf("can't parse index.html: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = t.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/empty.html")
	if err != nil {
		err = fmt.Errorf("can't parse empty.html: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = t.ExecuteTemplate(w, "empty", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/weather", weatherHandler)
	http.HandleFunc("/empty", emptyHandler)
	http.ListenAndServe(":3030", nil)
}
