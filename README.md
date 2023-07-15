# csit-swe-challenge-2023
A REST API for [CSIT's 2023 Mini Challenge](https://se-mini-challenge.csit-events.sg/).

## Problem
The challenge was to create a REST API providing the required routes:
- `GET /flight?departureDate=YYYY-MM-DD&returnDate=YYYY-MM-DD&destination=city`
  - This would return a JSON array, with each element having the following shape:
    ```json
    {
      "City": "City 2",
      "Departure Date": "2023-12-10",
      "Departure Airline": "Airline 1",
      "Departure Price": 1000,
      "Return Date": "2023-12-16",
      "Return Airline": "Airline 2",
      "Return Price": 1000
    }
    ```
  - This array would contain the **cheapest** flights to and from the given destination, with the source city being Singapore.
- `GET /hotel?checkInDate=YYYY-MM-DD&checkoutDate=YYYY-MM-DD&destination=city`
  - This would return a JSON array, with each element having the following shape:
    ```json
    {
      "City": "City 2",
      "Check In Date": "2023-12-10",
      "Check Out Date": "2023-12-16",
      "Hotel": "Hotel A",
      "Price": 1000
    }
    ```
  - This array would contain the **cheapest** hotel stays from the check-in date to the check-out date at the given destination.
A MongoDB database was provided, where there were two collections:
- `minichallenge.flights`, with the following truncated JSON shape:
  ```json
  {
    "airlinename": "Airline",
    "srccity": "City 1",
    "destcity": "City 2",
    "price": 1000,
    "date": 2023-12-10T00:00:00.000+00:00
  }
  ```
- `minichallenge.hotels`, with the following truncated JSON shape:
  ```json
  {
    "city": "City 1",
    "hotelName": "Hotel A",
    "price": 1000,
    "date": 2023-12-10T00:00:00.000+00:00
  }
  ```

## Solving the Problem
This was done in Go, because I'm relatively new to the language and wanted to use it more.

The main challenge was transforming the queried data for both flights and hotels into the requested forms.
- **Flights**: The chosen solution was to simply query for departure flights with the relevant cities and departure date (return date for return flights), and combine the departure flights with the return flights and sort by the cheapest total price.
- **Hotels**: A MongoDB aggregate was used.
  1. We first select only the relevant hotel records -- where the stay date was between the check-in and check-out dates and the hotel was in the correct city.
  2. Then, we group by the hotel and date of stay, taking the lowest price for that specific group. This allowed us to ensure each (hotel, date of stay) record would have the cheapest price.
  3. We then group by the hotel, taking the sum of the cheapest prices. Since each (hotel, date of stay) record was already between the check-in date and check-out date (in Step (1.)), this gave the **cheapest total price** for the given period of stay.
  4. We do some data manipulation to output the records in a nice format for Go to process.

## Installation
### Docker
If you have Docker installed, you can pull from the [Docker repository](https://hub.docker.com/r/ryedralisk/csit-swe-challenge-2023) and run it.

### Local
Ensure that there exists a `.env` file in the root directory (and the `./internal/db` directory if you wish to run tests).
- This `.env` file must contain a `MONGODB_URI` variable, with the URI discoverable on the CSIT challenge site. Presumably, this would be taken down after the challenge has ended.

To run tests:
```bash
make test
```

To build:
```bash
make build
```
- This would create an executable `server` in the root directory, which can be directly run.

To run (without building):
```bash
make run
```
