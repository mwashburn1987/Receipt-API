package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/mux"
)

type ReceiptIdResponse struct {
	Id string `json:"id,omitempty"`
}

type ReceiptPointsResponse struct {
	Points int `json:"points,omitempty"`
}

type item struct {
	ShortDescription string `json:"shortDescription,omitempty"`
	Price            string `json:"price,omitempty"`
}
type Receipt struct {
	Retailer     string `json:"retailer,omitempty"`
	PurchaseDate string `json:"purchaseDate,omitempty"`
	PurchaseTime string `json:"purchaseTime,omitempty"`
	Items        []item `json:"items,omitempty"`
	Total        string `json:"total,omitempty"`
}

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ received: ProcessReceipt")

	//Read receipt from body
	body, _ := io.ReadAll(r.Body)

	//unmarshal receipt into struct vals
	var receipt Receipt
	json.Unmarshal(body, &receipt)

	// create time from receipt string
	dttm := receipt.PurchaseDate + " " + receipt.PurchaseTime
	layout := "2006-01-02 15:04"

	t, err := time.Parse(layout, dttm)
	if err != nil {
		log.Println("Error parsing time", err)
	}

	fmt.Println(receipt)

	var points int
	// start calculating points for each rule
	points += calculateNamePoints(receipt.Retailer)
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		log.Println("Unable to parse float from total string", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	points += calculatePointsForTotal(totalFloat)
	itemPoints, err := calculatePointsOfItems(receipt.Items)
	if err != nil {
		log.Println("unable to get points of supplied items", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	points += itemPoints
	points += calculatePointsByDate(t)

	recResp := ReceiptPointsResponse{
		Points: points,
	}

	// create an id for our receipt with incrementing by length
	receiptID := strconv.Itoa(len(ReceiptResponses) + 1)

	// add our point total receipt to map with created id
	ReceiptResponses[receiptID] = recResp

	idRecResp := ReceiptIdResponse{
		Id: receiptID,
	}
	resp, err := json.Marshal(idRecResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQ received: GetPoints")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Println("missing id field in request")
		w.Header().Set("x-missing-field", "\"id\"")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	receiptPoints, ok := ReceiptResponses[id]
	if !ok {
		log.Println("Unable to find receipt by that id")
		w.Header().Set("x-receipt-not-found", id)
		w.WriteHeader(http.StatusBadRequest)
	}

	data, err := json.Marshal(receiptPoints)
	if err != nil {
		log.Printf("unable to marshal receipt points response with id: %s", id)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// separated each point rule by function, could be consolidated if we knew we would never
// receive another rule
func calculateNamePoints(retailerName string) (points int) {
	// 1 point for each alphanumeric character in retailer name
	points = 0
	for _, r := range retailerName {
		if unicode.IsLetter(r) || unicode.IsLetter(r) {
			points++
		}
	}
	return points
}

func calculatePointsForTotal(total float64) (points int) {
	points = 0

	//50 points if total is round dollar amount with no cents
	if int(total*100)%100 == 0 {
		points += 50
	}

	if int(total*100)%25 == 0 {
		points += 25
	}

	return points
}

func calculatePointsOfItems(items []item) (points int, err error) {
	points = 0

	points += (len(items) / 2) * 5

	for _, item := range items {
		priceFloat, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			log.Println("unable to parse float from item price", err)
			return 0, err
		}
		trimmedItem := strings.TrimSpace(item.ShortDescription)
		if len(trimmedItem)%3 == 0 {
			points += int(math.Ceil(priceFloat * .2))
		}
	}
	return points, nil
}

func calculatePointsByDate(dttm time.Time) (points int) {
	points = 0

	if dttm.Day()%2 == 1 {
		points += 6
	}
	start := time.Date(dttm.Year(), dttm.Month(), dttm.Day(), 14, 0, 0, 0, dttm.Location())
	end := time.Date(dttm.Year(), dttm.Month(), dttm.Day(), 16, 0, 0, 0, dttm.Location())
	if dttm.After(start) && dttm.Before(end) {
		points += 10
	}
	return points
}
