package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
)

type apiConfigData struct {
	ApiKey string `json:"OpenWeatherMapApiKey"`
}

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}
type Data struct {
	City string
}

func loadApiConfig(filename string) (apiConfigData, error) {
	//helps get the apiKey from .apiConfig file
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from GO.\nWelcome to the Weather Report."))
}
func query(city string) (WeatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return WeatherData{}, err
	}
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.ApiKey)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()
	var d WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}
	return d, nil
}
func fromHTMLForm(w http.ResponseWriter, r *http.Request) {
	city := r.FormValue("City")

	// data := &Data{city}

	// b, err := json.Marshal(data)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// f, err :=
	desc, err := query(city)
	if err == nil {
		return
	}
	json.NewEncoder(w).Encode(desc)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// f.Write(b)
	// f.Close()
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/process", processor)
	http.HandleFunc("/", index)

	http.ListenAndServe(":9000", nil)

}
func processor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	city := r.FormValue("cityName")

	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(data)
	// d := struct {
	// 	City string
	// }{
	// 	City: city,
	// }
	//tpl.ExecuteTemplate(w, "processor.gohtml", d)

}
func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

/*
func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
*/
