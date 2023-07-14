/**
  Handles querying of data from the database.
*/
package db

import "math/rand"

// Represents a single flight from SG to city
type Flight struct {
	city             string // Destination City
	departureDate    string // Date of Departure from SG (YYYY-MM-DD)
	departureAirline string
	departurePrice   float64
	returnDate       string // Date of Return from Destination City (YYYY-MM-DD)
	returnAirline    string
	returnPrice      float64
}

func (f *Flight) Price() float64 {
	return f.departurePrice + f.returnPrice
}

// Represents a hotel in city
type Hotel struct {
	city         string // City of Hotel
	checkInDate  string // Date of check-in (YYYY-MM-DD)
	checkOutDate string // Date of check-out(YYYY-MM-DD)
	hotel        string
	price        float64
}


var test_hotels = []string{"A Hotel", "Hotel B", "Hotel 123"}
var test_airlines = []string{"Singapore Airlines", "Emirates", "Another Airline", "US Airways", "Scoot"}

func randomHotel() string {
	return test_hotels[rand.Intn(len(test_hotels))]
}

func randomAirline() string {
	return test_airlines[rand.Intn(len(test_airlines))]
}



// Queries and returns a list of return flights given the departureDate, returnDate and destination
func Flights(departureDate string, returnDate string, destination string) (flights []Flight) {
	// TODO: replace with db code, for now simply generates data where necessary.
	count := rand.Intn(10)
	flights = make([]Flight, count)

	for i := 0; i < count; i++ {
		flights[i] = Flight{
			city: destination,
			departureDate: departureDate,
			departureAirline: randomAirline(),
			departurePrice: float64(rand.Intn(1500) + 500),
			returnDate: returnDate,
			returnAirline: randomAirline(),
			returnPrice: float64(rand.Intn(1500) + 500),
		}
	}

	// TODO: sort by cheapest price
	return flights
}

// Queries and returns a list of hotels given the checkInDate, checkOutDate and destination
func Hotels(checkInDate string, checkOutDate string, destination string) (hotels []Hotel) {
	// TODO: replace with db code, for now simply generates data where necessary
	count := rand.Intn(15)
	hotels = make([]Hotel, count)

	for i := 0; i < count; i++ {
		hotels[i] = Hotel{
			city: destination,
			checkInDate: checkInDate,
			checkOutDate: checkOutDate,
			hotel: randomHotel(),
			price: float64(rand.Intn(1500) + 2000),
		}
	}

	// TODO: sort by cheapest price
	return hotels
}
