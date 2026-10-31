# go-yaweather

Unofficial Go client for Yandex Weather forecasts and location search. No API key required.

## Install

```sh
go get github.com/mishamyrt/go-yaweather
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mishamyrt/go-yaweather"
)

func main() {
	location := yaweather.Location{Latitude: 55.7558, Longitude: 37.6173}
	forecast, err := location.GetForecast(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, day := range forecast {
		fmt.Println(day.Date.Format("2006-01-02"))
		fmt.Println(day)
	}
}
```

Each day includes morning, day, evening, and night forecasts: temperature (°C),
feels-like temperature, conditions, wind, humidity, and pressure.

Use `yaweather.Find(ctx, "Moscow")` to search for locations.
See [examples](examples) for more.

## License

[MIT](./LICENSE)

## Legal

This project is not affiliated with, endorsed by, or sponsored by Yandex.
