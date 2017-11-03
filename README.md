# Go Server

Small Go Server that access two Weather APIs concurrently.

## Routes

- `/`: Hello World
- `/weather/<city>`: Returns the average temperature for a city

## Building and Running

`go build`

`./go_server`

This will open the server at localhost:8080

## Testing

`http://localhost:8080/weather/boston`
