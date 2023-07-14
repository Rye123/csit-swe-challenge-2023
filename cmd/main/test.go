package main

import "fmt"
import "github.com/Rye123/csit-swe-challenge-2023/internal/db"

func main() {
	fmt.Println("Sending Query.")
	flights, err := db.Flights("2023-12-10", "2023-12-20", "Abu Dhabi", -1)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(flights))
	for _, flight := range flights {
		fmt.Printf("Using %s to %s. Price: %d\n", flight.DepartureAirline, flight.ReturnAirline, flight.Price())
	}

	fmt.Println("Sent query.")
}
