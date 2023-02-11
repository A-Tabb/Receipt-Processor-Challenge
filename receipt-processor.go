package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []struct {
		ShortDesc string `json:"shortDescription"`
		Price     string `json:"price"`
	}
}

type ReceiptWithID struct {
	Receipt
	id string
}

type ID struct {
	Id string `json:"id"`
}

var Receipts []ReceiptWithID

var IDs []ID

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/receipts/process", processReceipts).Methods("POST")
	myRouter.HandleFunc("/receipts/{id}/points", getPoints)
	myRouter.HandleFunc("/receipts", getReceipts)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func processReceipts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: Process Receipts")
	//fmt.Fprintf(w, "returned ID")
	fileContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var receiptInfo Receipt
	err = json.Unmarshal([]byte(fileContent), &receiptInfo)
	//fmt.Printf("Receipt Information: %+v", receiptInfo)

	receiptID := ID{Id: uuid.New().String()}

	thisReceiptWithID := ReceiptWithID{
		Receipt: receiptInfo,
		id:      receiptID.Id,
	}

	byteArray, err := json.Marshal(receiptID)
	if err != nil {
		fmt.Println(err)
	}

	Receipts = append(Receipts, thisReceiptWithID)
	IDs = append(IDs, receiptID)

	fmt.Fprintf(w, string(byteArray))

}

func getReceipts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Receipts)
}

func getPoints(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println("Hello world!")
	handleRequests()
}
