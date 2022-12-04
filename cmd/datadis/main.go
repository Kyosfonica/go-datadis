package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rubiojr/go-datadis"
	"github.com/rubiojr/go-datadis/cmd/storage"
)

// Fetch datadis last day consumption
func main() {
	sqlite, err := storage.NewSqlite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sqlite.Close()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := datadis.NewClient()
	err = client.Login(os.Getenv("DATADIS_USERNAME"), os.Getenv("DATADIS_PASSWORD"))
	if err != nil {
		panic(err)
	}

	s, err := client.Supplies()
	if err != nil {
		panic(err)
	}

	now := time.Now()
	year, month, day := now.Date()
	dateFrom := time.Date(year, month-2, day, 0, 0, 0, 0, now.UTC().Location())
	dateTo := time.Date(year, month, day, 0, 0, 0, 0, now.UTC().Location())
	data, err := client.ConsumptionData(&s[0], dateFrom, dateTo)
	for _, d := range data {
		fmt.Println("CUPS: ", d.Cups)
		fmt.Println("Date: ", d.Date)
		fmt.Println("Time: ", d.Time)
		fmt.Printf("Consumption: %f KWh\n", d.Consumption)
		fmt.Println("Obtained Method: ", d.ObtainMethod)

		err = sqlite.Writer(&d)
		if err != nil {
			panic(err)
		}
	}
}
