package main

import (
  "encoding/json"
  "net/http"
  "strings"
)

// ================ CONSTS AND TYPES ================

const api_key = "ea91772698ddeb376160487657804955"

type weatherData struct {
  Name string `json:"name"`
  Main struct {
    Kelvin float64 `json:"temp"`
  } `json:"main"`
}

// ================ MAIN ================

func main() {
  http.HandleFunc("/", hello)

  http.HandleFunc("/weather/", weather)

  http.ListenAndServe(":8080", nil)
}

// ================ HANDLERS (ROUTES?) ================

func hello(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello World!"))
}

func weather(w http.ResponseWriter, r *http.Request) {
  city := strings.SplitN(r.URL.Path, "/", 3)[2]

  data, err := query(city)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  json.NewEncoder(w).Encode(data)
}

// ================ AUXILIARY ================

func query(city string) (weatherData, error) {
  resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + api_key + "&q=" + city)
  if err !=  nil {
    return weatherData{}, err
  }

  defer resp.Body.Close()

  var d weatherData

  if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
    return weatherData{}, err
  }

  return d, nil
}
