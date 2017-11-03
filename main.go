package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strings"
  "time"
)

// ================ CONSTS AND TYPES ================

const ow_key = "ea91772698ddeb376160487657804955"
const wu_key = "68175cc1c0c4261d"

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
  mw := multiWeatherProvider{
    openWeatherMap{key: ow_key},
    weatherUnderground{key: wu_key},
  }

  begin := time.Now()
  city := strings.SplitN(r.URL.Path, "/", 3)[2]

  temp, err := mw.temperature(city)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  json.NewEncoder(w).Encode(map[string]interface{}{
    "city": city,
    "temp": temp,
    "took": time.Since(begin).String(),
  })
}

// ================ AUXILIARY ================

type weatherProvider interface {
  temperature(city string) (float64, error)
}

type openWeatherMap struct{
  key string
}

func (w openWeatherMap) temperature(city string) (float64, error) {
  resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.key + "&q=" + city)
  if err !=  nil {
    return 0, err
  }

  defer resp.Body.Close()

  var d struct {
    Main struct {
      Kelvin float64 `json:"temp"`
    } `json:"main"`
  }

  if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
    return 0, err
  }

  log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
  return d.Main.Kelvin, nil
}

type weatherUnderground struct {
  key string
}

func (w weatherUnderground) temperature(city string) (float64, error) {
  resp, err := http.Get("http://api.wunderground.com/api/" + w.key + "/conditions/q/" + city + ".json")
  if err != nil {
    return 0, err
  }

  defer resp.Body.Close()

  var d struct {
    Observation struct {
      Celsius float64 `json:"temp_c"`
    } `json:"current_observation"`
  }

  if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
    return 0, err
  }

  kelvin := d.Observation.Celsius + 273.15
  log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
  return kelvin, nil
}

type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
  sum := 0.0

  for _, provider := range w {
    k, err  := provider.temperature(city)
    if err != nil {
      return 0, err
    }

    sum += k
  }

  return sum / float64(len(w)), nil
}
