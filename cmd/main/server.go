package main

import "fmt"
import "net/http"

import "github.com/Rye123/csit-swe-challenge-2023/internal/db"
import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/flight", getFlights)
	router.GET("/hotel", getHotels)

	router.Run(":8080")
	fmt.Println("Server running on http://localhost:8080")
}

// Serves URL: /flight?departureDate=2023-12-10&returnDate=2023-12-16&destination=Frankfurt
func getFlights(c *gin.Context) {
	departureDate := c.Query("departureDate")
	returnDate := c.Query("returnDate")
	destination := c.Query("destination")

	flights, err := db.Flights(departureDate, returnDate, destination)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, nil)
	} else {
		c.JSON(http.StatusOK, flights[0])
	}
}

// Serves URL: /hotel?checkInDate=2023-12-10&checkOutDate=2023-12-16&destination=Frankfurt
func getHotels(c *gin.Context) {
	checkInDate := c.Query("checkInDate")
	checkOutDate := c.Query("checkOutDate")
	destination := c.Query("destination")

	hotels, err := db.Hotels(checkInDate, checkOutDate, destination)
	fmt.Println(hotels)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, nil)
	} else {
		c.JSON(http.StatusOK, hotels)
	}
}
