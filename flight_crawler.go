package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	FlyFrom = "CUU"
	DateFrom = "30/01/2022"
	DateTo = "07/02/2022"
)

type Flight struct {
	SearchId string `json:"search_id"`
	Time     int    `json:"time"`
	Currency string `json:"currency"`
	FxRate   int    `json:"fx_rate"`
	Data     []struct {
		Id           string `json:"id"`
		FlyFrom      string `json:"flyFrom"`
		FlyTo        string `json:"flyTo"`
		CityFrom     string `json:"cityFrom"`
		CityCodeFrom string `json:"cityCodeFrom"`
		CityTo       string `json:"cityTo"`
		CityCodeTo   string `json:"cityCodeTo"`
		CountryFrom  struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"countryFrom"`
		CountryTo struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"countryTo"`
		TypeFlights  []string    `json:"type_flights"`
		NightsInDest interface{} `json:"nightsInDest"`
		Quality      float64     `json:"quality"`
		Distance     float64     `json:"distance"`
		Duration     struct {
			Departure int `json:"departure"`
			Return    int `json:"return"`
			Total     int `json:"total"`
		} `json:"duration"`
		Price      float64 `json:"price"`
		Conversion struct {
			EUR float64 `json:"EUR"`
		} `json:"conversion"`
		BagsPrice struct {
			Field1 float64 `json:"1"`
		} `json:"bags_price"`
		Baglimit struct {
			HandHeight         int `json:"hand_height"`
			HandLength         int `json:"hand_length"`
			HandWeight         int `json:"hand_weight"`
			HandWidth          int `json:"hand_width"`
			HoldDimensionsSum  int `json:"hold_dimensions_sum"`
			HoldHeight         int `json:"hold_height"`
			HoldLength         int `json:"hold_length"`
			HoldWeight         int `json:"hold_weight"`
			HoldWidth          int `json:"hold_width"`
			PersonalItemHeight int `json:"personal_item_height"`
			PersonalItemLength int `json:"personal_item_length"`
			PersonalItemWeight int `json:"personal_item_weight"`
			PersonalItemWidth  int `json:"personal_item_width"`
		} `json:"baglimit"`
		Availability struct {
			Seats int `json:"seats"`
		} `json:"availability"`
		Routes   [][]string `json:"routes"`
		Airlines []string   `json:"airlines"`
		Route    []struct {
			Id                  string      `json:"id"`
			CombinationId       string      `json:"combination_id"`
			FlyFrom             string      `json:"flyFrom"`
			FlyTo               string      `json:"flyTo"`
			CityFrom            string      `json:"cityFrom"`
			CityCodeFrom        string      `json:"cityCodeFrom"`
			CityTo              string      `json:"cityTo"`
			CityCodeTo          string      `json:"cityCodeTo"`
			Airline             string      `json:"airline"`
			FlightNo            int         `json:"flight_no"`
			OperatingCarrier    string      `json:"operating_carrier"`
			OperatingFlightNo   string      `json:"operating_flight_no"`
			FareBasis           string      `json:"fare_basis"`
			FareCategory        string      `json:"fare_category"`
			FareClasses         string      `json:"fare_classes"`
			FareFamily          string      `json:"fare_family"`
			Return              int         `json:"return"`
			BagsRecheckRequired bool        `json:"bags_recheck_required"`
			ViConnection        bool        `json:"vi_connection"`
			Guarantee           bool        `json:"guarantee"`
			LastSeen            time.Time   `json:"last_seen"`
			RefreshTimestamp    time.Time   `json:"refresh_timestamp"`
			Equipment           interface{} `json:"equipment"`
			VehicleType         string      `json:"vehicle_type"`
			LocalArrival        time.Time   `json:"local_arrival"`
			UtcArrival          time.Time   `json:"utc_arrival"`
			LocalDeparture      time.Time   `json:"local_departure"`
			UtcDeparture        time.Time   `json:"utc_departure"`
		} `json:"route"`
		BookingToken                string        `json:"booking_token"`
		DeepLink                    string        `json:"deep_link"`
		TrackingPixel               string        `json:"tracking_pixel"`
		FacilitatedBookingAvailable bool          `json:"facilitated_booking_available"`
		PnrCount                    int           `json:"pnr_count"`
		HasAirportChange            bool          `json:"has_airport_change"`
		TechnicalStops              int           `json:"technical_stops"`
		ThrowAwayTicketing          bool          `json:"throw_away_ticketing"`
		HiddenCityTicketing         bool          `json:"hidden_city_ticketing"`
		VirtualInterlining          bool          `json:"virtual_interlining"`
		Transfers                   []interface{} `json:"transfers"`
		LocalArrival                time.Time     `json:"local_arrival"`
		UtcArrival                  time.Time     `json:"utc_arrival"`
		LocalDeparture              time.Time     `json:"local_departure"`
		UtcDeparture                time.Time     `json:"utc_departure"`
	} `json:"data"`
	Results      int `json:"_results"`
	SearchParams struct {
		FlyFromType string `json:"flyFrom_type"`
		ToType      string `json:"to_type"`
		Seats       struct {
			Passengers int `json:"passengers"`
			Adults     int `json:"adults"`
			Children   int `json:"children"`
			Infants    int `json:"infants"`
		} `json:"seats"`
	} `json:"search_params"`
	AllAirlines         []interface{} `json:"all_airlines"`
	AllStopoverAirports []interface{} `json:"all_stopover_airports"`
	Del                 int           `json:"del"`
	CurrencyRate        int           `json:"currency_rate"`
	Connections         []interface{} `json:"connections"`
	Refresh             []interface{} `json:"refresh"`
	RefTasks            []interface{} `json:"ref_tasks"`
	SortVersion         int           `json:"sort_version"`
}

func buildTequilaAPICall(iataCode string, price float64) string {
	b := strings.Builder{}
	b.WriteString("https://tequila-api.kiwi.com/v2/search?")
	b.WriteString("fly_from=")
	b.WriteString(FlyFrom)
	b.WriteString("&date_from=")
	b.WriteString(DateFrom)
	b.WriteString("&date_to=")
	b.WriteString(DateTo)
	b.WriteString("&fly_to=")
	b.WriteString(iataCode)
	b.WriteString("&price_to=")
	b.WriteString(strconv.Itoa(int(price)))
	b.WriteString("&limit=1")
	return b.String()
}

type FlightPrices struct {
	CurrentPrice float64
	FoundPrice float64
}

func scheduleFlightTask(wg *sync.WaitGroup, c *http.Client, n SmsNotifier, t FlightTask) {
	defer wg.Done()
	prices := fetchFlights(c, t)
	notifyIfLowerPriceFound(n, prices.CurrentPrice, prices.FoundPrice, t.Destination)
}

func fetchFlights(c *http.Client, t FlightTask) FlightPrices {
	apiEndpoint := buildTequilaAPICall(t.IATACode, t.TrackPrice)
	req, _ := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	req.Header.Set("apiKey", os.Getenv("TEQUILA_API_KEY"))
	res, errReq := c.Do(req)
	if errReq != nil {
		return FlightPrices{}
	}
	defer res.Body.Close()

	body, errRead := ioutil.ReadAll(res.Body)
	if errRead != nil {
		return FlightPrices{}
	}

	flight := Flight{}
	if err := json.Unmarshal(body, &flight); err != nil || len(flight.Data) < 1 {
		return FlightPrices{}
	}
	return FlightPrices{
		CurrentPrice: t.TrackPrice,
		FoundPrice:   flight.Data[0].Price,
	}
}
