package db

import "context"
import "time"
import "os"
import "errors"
import "sync"
import "github.com/joho/godotenv"
import "go.mongodb.org/mongo-driver/bson"
import "go.mongodb.org/mongo-driver/mongo"
import "go.mongodb.org/mongo-driver/mongo/options"

type FlightModel struct {
	SrcCity string `bson: "srccity"`
	DestCity string `bson: "srrccity"`
	AirlineName string `bson: "airlinename"`
	Price int `bson: "price"`
	Date time.Time `bson: "date"`
}

type HotelModel struct {
	City string `bson: "city"`
	Hotel string `bson: "hotelName"`
	Price int `bson: "price"`
	Date time.Time `bson: "date"`
}

// Setup before database usage.
func mongoDB_URI() string {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic(errors.New("No MONGODB_URI environment variable found."))
	}
	return uri
}

// Queries DB for Flights
func queryFlights(departureDateStr string, returnDateStr string, destination string) (flights []Flight, err error) {
	// Convert datestrings
	departureDate, err := time.Parse(time.DateOnly, departureDateStr)
	if err != nil { return nil, err	}
	returnDate, err := time.Parse(time.DateOnly, returnDateStr)
	if err != nil { return nil, err }
	
	// Connect to DB
	mongodb_uri := mongoDB_URI()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodb_uri))
	if err != nil {
		return nil, err
	}
	
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Query for Departure and Return Flights
	collection := client.Database("minichallenge").Collection("flights")
	departFilter := bson.D{
		{"date", departureDate},
		{"destcity", destination},
	}
	returnFilter := bson.D{
		{"date", returnDate},
		{"destcity", "Singapore"},
	}
	opts := options.Find().SetSort(bson.D{{"price", 1}}) // Sort by ascending price

	departCursor, err := collection.Find(context.TODO(), departFilter, opts)
	if err != nil { return nil, err }
	returnCursor, err := collection.Find(context.TODO(), returnFilter, opts)
	if err != nil { return nil, err }

	var departFlights []FlightModel
	var returnFlights []FlightModel
	if err = departCursor.All(context.TODO(), &departFlights); err != nil {
		return nil, err
	}
	if err = returnCursor.All(context.TODO(), &returnFlights); err != nil {
		return nil, err
	}

	if len(departFlights) == 0 || len(returnFlights) == 0 {
		return []Flight{}, nil
	}
	lowestPrice := departFlights[0].Price + returnFlights[0].Price

	// Combine to give two-way flights
	var wg sync.WaitGroup
	flightChan := make(chan Flight)
	flights = make([]Flight, 0)
	for _, departFlight := range departFlights {
		wg.Add(1)
		go func(departFlight FlightModel, flightChan chan Flight) {
			defer wg.Done()
			for _, returnFlight := range returnFlights {
				if departFlight.Price + returnFlight.Price != lowestPrice {
					continue
				}
				flight := Flight{
					City: departFlight.DestCity,
					DepartureDate: departureDateStr,
					DepartureAirline: departFlight.AirlineName,
					DeparturePrice: departFlight.Price,
					ReturnDate: returnDateStr,
					ReturnAirline: returnFlight.AirlineName,
					ReturnPrice: returnFlight.Price,
				}
				flightChan <- flight
			}
		}(departFlight, flightChan)
	}

	// Consume incoming values
	go func() {
		for flight := range flightChan {
			flights = append(flights, flight)
		}
	}()

	wg.Wait()
	close(flightChan)
		
	return flights, nil
}
