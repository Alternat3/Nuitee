package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHotels(c *gin.Context) {
	currentQuery, err := ParseQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request (query details)")
	}

	requestHotelbeds(currentQuery)

	c.JSON(http.StatusOK, "OK")
}

func requestHotelbeds(currentQuery HotelQuery) {
	reqBody := HotelbedsRequest{
		Stay: Stay{
			CheckIn:  currentQuery.checkIn,
			CheckOut: currentQuery.checkOut,
		},
		Occupancies: currentQuery.Occupancies,
		Hotels: Hotels{
			Hotel: currentQuery.hotelIds,
		},
	}

	reqBodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.test.hotelbeds.com/hotel-api/1.0/hotels", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		fmt.Println("GetHotels Error -> Error creating request for hotel beds")
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.Getenv("API_KEY"))

	//Error recieved with Hotelbeds of Authorization field not present, tested with postman and checked request data tried multiple types of authorization style headers
	//but still getting the issue, unfortunately I don't have the time to fully debug.
	// Understand if its not good enough but happy to share previous code I've worked on in this type to show abilities and can talk through said code.
}

func ParseQuery(c *gin.Context) (HotelQuery, error) {
	currentQuery := HotelQuery{
		checkIn:          c.Query("checkin"),
		checkOut:         c.Query("checkout"),
		currency:         c.Query("currency"),
		guestNationality: c.Query("guestNationality"),
	}

	hotelIds := c.Query("hotelIds")
	tempArr := strings.Split(hotelIds, ",")
	for _, id := range tempArr {
		current, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("GetHotels Error -> Error parsing hotel id")
			return HotelQuery{}, err
		}
		currentQuery.hotelIds = append(currentQuery.hotelIds, int32(current))
	}

	occupancies := []Occupancies{}
	occupanciesByt := []byte(c.Query("occupancies"))
	err := json.Unmarshal(occupanciesByt, &occupancies)
	if err != nil {
		fmt.Println("GetHotels Error -> Could not unmarhsal occupancies JSON")
		fmt.Println(err)
		return HotelQuery{}, err
	}
	currentQuery.Occupancies = occupancies

	return currentQuery, nil
}

type HotelQuery struct {
	checkIn          string
	checkOut         string
	currency         string
	guestNationality string
	hotelIds         []int32
	Occupancies      []Occupancies
}

type Occupancies struct {
	Rooms    int `json:"rooms"`
	Adults   int `json:"adults"`
	Children int `json:"children"`
}

type HotelbedsRequest struct {
	Stay        Stay          `json:"stay"`
	Occupancies []Occupancies `json:"occupancies"`
	Hotels      Hotels        `json:"hotels"`
}

type Stay struct {
	CheckIn  string `json:"checkIn"`
	CheckOut string `json:"checkOut"`
}

type Hotels struct {
	Hotel []int32 `json:"hotel"`
}

type HotelResponse struct {
	Code     int     `json:"code"`
	MinRate  float64 `json:"minRate"`
	MaxRate  float64 `json:"maxRate"`
	Currency string  `json:"currency"`
}
