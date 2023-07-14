/**
  Tests for the db package.

  To run, use:
  ```
  go test ./internal/db -v
  ```
*/

package db

import "testing"
import "math/rand"
import "strconv"

// Converts a two-digit int to a two-digit string
func itoaTwoDigit(num int) (string) {
	if num < 10 {
		return "0" + strconv.Itoa(num)
	}
	return strconv.Itoa(num)
}

func randomDateStrings() (date1, date2 string) {
	year1 := rand.Intn(5) + 2018
	year2 := year1 + rand.Intn(1)
	
	month1 := rand.Intn(11) + 1
	month2 := rand.Intn(11) + 1
	if month2 < month1 {
		year2++
	}
	day1 := rand.Intn(28) + 1
	day2 := rand.Intn(28) + 1

	// Logic to ensure date2 is AFTER date1
	if month2 == month1 {
		if day2 <= day1 && day1 == 28 {
			day2 = rand.Intn(28) + 1
			month2 += 1
			if month2 > 12 {
				month2 = 1
				year2++
			}
		} else if day2 <= day1 {
			day2 = day1 + rand.Intn(28 - day1)
		}
	}

	date1 = strconv.Itoa(year1) + "-" + itoaTwoDigit(month1) + "-" + itoaTwoDigit(day1)
	date2 = strconv.Itoa(year2) + "-" + itoaTwoDigit(month2) + "-" + itoaTwoDigit(day2)

	return date1, date2
}

func randomCity() string {
	cities := []string{"London", "Frankfurt", "Beijing", "New Delhi"}
	return cities[rand.Intn(len(cities))]
}

func TestFlights(t *testing.T) {
	departureDate, returnDate := randomDateStrings()
	destination := randomCity()
	
	flights := Flights(departureDate, returnDate, destination)

	prevFlightPrice :=  float64(0)
	
	// Test if all received flights match the relevant arguments, and assert that they are cheapest first
	for _, flight := range flights {
		if flight.departureDate != departureDate {
			t.Fatalf("Invalid Flight (departureDate). Expected %s, Given %s", departureDate, flight.departureDate)
		}
		if flight.returnDate != returnDate {
			t.Fatalf("Invalid Flight (returnDate). Expected %s, Given %s", returnDate, flight.returnDate)
		}
		if flight.city != destination {
			t.Fatalf("Invalid Flight (destination). Expected %s, Given %s", destination, flight.city)
		}
		if flight.Price() < prevFlightPrice {
			t.Fatalf("Invalid Flight (price). Not sorted in cheapest first -- Previous Price %f, Current Price %f.", prevFlightPrice, flight.Price())
		}
		prevFlightPrice = flight.Price()
	}
}

func TestHotels(t *testing.T) {
	checkInDate, checkOutDate := randomDateStrings()
	destination := randomCity()

	hotels := Hotels(checkInDate, checkOutDate, destination)
	prevHotelPrice := float64(0)

	// Test if all received hotels match the relevant arguments, and assert that they are sorted cheapest first
	for _, hotel := range hotels {
		if hotel.checkInDate != checkInDate {
			t.Fatalf("Invalid Hotel (checkInDate). Expected %s, Given %s", checkInDate, hotel.checkInDate)
		}
		if hotel.checkOutDate != checkOutDate {
			t.Fatalf("Invalid Hotel (checkOutDate). Expected %s, Given %s", checkOutDate, hotel.checkOutDate)
		}
		if hotel.city != destination {
			t.Fatalf("Invalid Hotel (destination). Expected %s, Given %s", destination, hotel.city)
		}
		if hotel.price < prevHotelPrice {
			t.Fatalf("Invalid Hotel (price). Not sorted in cheapest first -- Previous Price %f, Current Price %f.", prevHotelPrice, hotel.price)
		}
		prevHotelPrice = hotel.price
	}
}
