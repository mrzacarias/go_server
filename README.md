# Go Server

Small Go Server that access two Weather APIs concurrently.

## Routes

- `/`: Hello World
- `/weather/<city>`: Returns the average temperature for a city

## Building and Running

`go build`

`./go_server`

This will open the server at [localhost:8080](http://localhost:8080/weather/boston)

## Docker

`docker-compose build && docker-compose up -d`

This will open the server at [localhost:4000](http://localhost:4000/weather/boston)
