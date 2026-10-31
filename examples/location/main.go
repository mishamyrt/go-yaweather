package main

import (
	"context"
	"fmt"

	"github.com/mishamyrt/go-yaweather"
)

func main() {
	loc := &yaweather.Location{
		Latitude:  43.90551758,
		Longitude: 42.71571732,
	}
	ctx := context.TODO()
	forecast, err := loc.GetForecast(ctx)
	if err != nil {
		panic(err)
	}
	for _, day := range forecast {
		fmt.Println(day.Date.Format("Mon Jan _2"))
		fmt.Println(day.String())
		fmt.Println()
	}
}
