package main

import "fmt"

import "github.com/Rye123/csit-swe-challenge-2023/internal/db"

func main() {
	flights := db.Flights("2023-03-09", "2023-04-01", "London")
	hotels := db.Hotels("2023-03-11", "2023-03-29", "London")

	fmt.Println(flights)
	fmt.Println(hotels)
	
}
