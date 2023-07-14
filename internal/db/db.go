/**
  Handles querying of data from the database.
*/
package db

import "math/rand"
import "errors"
import "time"
import "strings"

// Represents a single flight from SG to city
type Flight struct {
	City             string  `json: "city"`          // Destination City
	DepartureDate    string  `json: "departureDate"` // Date of Departure from SG (YYYY-MM-DD)
	DepartureAirline string  `json: "departureAirline"`
	DeparturePrice   int     `json: "departurePrice"`
	ReturnDate       string  `json: "returnDate"` // Date of Return from Destination City (YYYY-MM-DD)
	ReturnAirline    string  `json: "returnAirline"`
	ReturnPrice      int     `json: "returnPrice"`
}

func (f *Flight) Price() int {
	return f.DeparturePrice + f.ReturnPrice
}

// Represents a hotel in city
type Hotel struct {
	City         string  `json: "city"`         // City of Hotel
	CheckInDate  string  `json: "checkInDate"`  // Date of check-in (YYYY-MM-DD)
	CheckOutDate string  `json: "checkOutDate"` // Date of check-out(YYYY-MM-DD)
	Hotel        string  `json: "hotel"`
	Price        int     `json: "price"`
}

var test_hotels = []string{"A Hotel", "Hotel B", "Hotel 123"}
var test_airlines = []string{"Singapore Airlines", "Emirates", "Another Airline", "US Airways", "Scoot"}

func randomHotel() string {
	return test_hotels[rand.Intn(len(test_hotels))]
}

func randomAirline() string {
	return test_airlines[rand.Intn(len(test_airlines))]
}

func isValidDate(dateStr string) bool {
	_, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return false
	}
	return true
}	

// Queries and returns a list of return flights given the departureDate, returnDate and destination
func Flights(departureDate string, returnDate string, destination string) (flights []Flight, err error) {
	// Validate dates
	if !isValidDate(departureDate) {
		return nil, errors.New("Flights: Invalid departureDate.")
	}
	if !isValidDate(returnDate) {
		return nil, errors.New("Flights: Invalid returnDate.")
	}

	// Standardise destination string
	destination = strings.Title(strings.ToLower(destination))

	// Query DB
	flights, err = queryFlights(departureDate, returnDate, destination)
	if err != nil {
		return nil, err
	}
	
	return flights, nil
}

// Queries and returns a list of hotels given the checkInDate, checkOutDate and destination
func Hotels(checkInDate string, checkOutDate string, destination string) (hotels []Hotel, err error) {
	// Validate dates
	if !isValidDate(checkInDate) {
		return nil, errors.New("Hotels: Invalid checkInDate.")
	}

	if !isValidDate(checkOutDate) {
		return nil, errors.New("Hotels: Invalid checkOutDate.")
	}

	// Standardise destination string
	destination = strings.Title(strings.ToLower(destination))

	// Query DB
	hotels, err = queryHotels(checkInDate, checkOutDate, destination)
	if err != nil {
		return nil, err
	}
	
	return hotels, nil
}
