/**
  Handles querying of data from the database.
*/
package db

import "math/rand"
import "errors"
import "time"

// Represents a single flight from SG to city
type Flight struct {
	City             string  `json: "city"`             // Destination City
	DepartureDate    string  `json: "departureDate"`    // Date of Departure from SG (YYYY-MM-DD)
	DepartureAirline string  `json: "departureAirline"`
	DeparturePrice   float64 `json: "departurePrice"`
	ReturnDate       string  `json: "returnDate"`       // Date of Return from Destination City (YYYY-MM-DD)
	ReturnAirline    string  `json: "returnAirline"`
	ReturnPrice      float64 `json: "returnPrice"`
}

func (f *Flight) Price() float64 {
	return f.DeparturePrice + f.ReturnPrice
}

// Represents a hotel in city
type Hotel struct {
	City         string   `json: "city"`         // City of Hotel
	CheckInDate  string   `json: "checkInDate"`  // Date of check-in (YYYY-MM-DD)
	CheckOutDate string   `json: "checkOutDate"` // Date of check-out(YYYY-MM-DD)
	Hotel        string   `json: "hotel"`
	Price        float64  `json: "price"`
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
func Flights(departureDate string, returnDate string, destination string, limit int) (flights []Flight, err error) {
	// Validate dates
	if !isValidDate(departureDate) {
		return nil, errors.New("Flights: Invalid departureDate.")
	}
	if !isValidDate(returnDate) {
		return nil, errors.New("Flights: Invalid returnDate.")
	}

	// Validate limit
	if limit < -1 {
		return nil, errors.New("Flights: Invalid limit.")
	}
	
	// TODO: replace with db code, for now simply generates data where necessary.
	if limit == -1 {
		limit = 10
	}
	count := rand.Intn(limit)
	flights = make([]Flight, count)

	for i := 0; i < count; i++ {
		flights[i] = Flight{
			City: destination,
			DepartureDate: departureDate,
			DepartureAirline: randomAirline(),
			DeparturePrice: float64(rand.Intn(1500) + 500),
			ReturnDate: returnDate,
			ReturnAirline: randomAirline(),
			ReturnPrice: float64(rand.Intn(1500) + 500),
		}
	}

	// TODO: sort by cheapest price
	return flights, nil
}

// Queries and returns a list of hotels given the checkInDate, checkOutDate and destination
func Hotels(checkInDate string, checkOutDate string, destination string, limit int) (hotels []Hotel, err error) {
	// Validate dates
	if !isValidDate(checkInDate) {
		return nil, errors.New("Hotels: Invalid checkInDate.")
	}
	
	if !isValidDate(checkOutDate) {
		return nil, errors.New("Hotels: Invalid checkOutDate.")
	}
	
	// Validate limit
	if limit < -1 {
		return nil, errors.New("Hotels: Invalid limit.")
	}
	
	// TODO: replace with db code, for now simply generates data where necessary.
	if limit == -1 {
		limit = 15
	}
	
	count := rand.Intn(limit)
	hotels = make([]Hotel, count)

	for i := 0; i < count; i++ {
		hotels[i] = Hotel{
			City: destination,
			CheckInDate: checkInDate,
			CheckOutDate: checkOutDate,
			Hotel: randomHotel(),
			Price: float64(rand.Intn(1500) + 2000),
		}
	}

	// TODO: sort by cheapest price
	return hotels, nil
}
