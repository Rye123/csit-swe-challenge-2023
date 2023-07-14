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
	HotelName string `bson: "hotelName"`
	Price int `bson: "price"`
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

func queryHotels(checkInDateStr string, checkOutDateStr string, destination string) (hotels []Hotel, err error) {
	// Convert datestrings
	checkInDate, err := time.Parse(time.DateOnly, checkInDateStr)
	if err != nil { return nil, err }
	checkOutDate, err := time.Parse(time.DateOnly, checkOutDateStr)
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

	// Query
	collection := client.Database("minichallenge").Collection("hotels")
	agg := bson.A{
		bson.D{
			{"$match", // Filter to only the relevant days
				bson.D{
					{"date",
						bson.D{
							{"$gte", checkInDate},
							{"$lte", checkOutDate},
						},
					},
					{"city", destination},
				},
			},
		},
		bson.D{
			{"$group", // Group by {hotelName, date} and get cheapest price that day
				bson.D{
					{"_id",
						bson.D{
							{"hotelName", "$hotelName"},
							{"date", "$date"},
						},
					},
					{"lowestPriceForDay", bson.D{{"$min", "$price"}}},
				},
			},
		},
		bson.D{
			{"$group", // Group by hotelName, get sum of prices for that hotel
				bson.D{
					{"_id", "$_id.hotelName"},
					{"totalPrice", bson.D{{"$sum", "$lowestPriceForDay"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", false},
					{"hotelName", "$_id"},
					{"price", "$totalPrice"},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{{"price", 1}},
			},
		},
	}
	
	cursor, err := collection.Aggregate(context.TODO(), agg)
	if err != nil {
		return nil, err
	}

	var hotelModels []HotelModel
	if err = cursor.All(context.TODO(), &hotelModels); err != nil {
		return nil, err
	}
	
	if len(hotelModels) == 0 {
		return []Hotel{}, nil
	}
	
	// Return as Hotels
	hotels = make([]Hotel, 0)
	lowestPrice := hotelModels[0].Price
	
	for _, hotelModel := range hotelModels {
		if hotelModel.Price > lowestPrice {
			continue
		}
		hotel := Hotel{
			City: destination,
			CheckInDate: checkInDateStr,
			CheckOutDate: checkOutDateStr,
			Hotel: hotelModel.HotelName,
			Price: hotelModel.Price,
		}
		hotels = append(hotels, hotel)
	}
	return hotels, nil
}
