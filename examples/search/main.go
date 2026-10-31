package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mishamyrt/go-yaweather"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: search <location>")
		os.Exit(1)
	}
	cityName := os.Args[1]
	ctx := context.TODO()
	results, _ := yaweather.Find(ctx, cityName)
	var target yaweather.Location
	for _, result := range results {
		if result.Title == cityName {
			target = result.Location
			break
		}
	}
	if target.IsZero() {
		fmt.Println("Location is not found")
		fmt.Println("Possible locations:")
		for _, result := range results {
			fmt.Println(result.Title)
		}
		os.Exit(1)
	}
	forecast, err := target.GetForecast(ctx)
	if err != nil {
		panic(err)
	}
	for _, day := range forecast {
		fmt.Println(day.Date.Format("Mon Jan _2"))
		fmt.Println(day.String())
		fmt.Println()
	}
}
