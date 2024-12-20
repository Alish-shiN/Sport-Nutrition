package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"onlinestore/mongoDB"
	"onlinestore/product"
	// "go.mongodb.org/mongo-driver/mongo"
	// "time"
)


type RequestData struct {	
	Message string `json:"message"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}


func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	requestData := RequestData{
		Message : "false",
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)

	if err != nil {
		response := ResponseData{
			Status:  "Fail",
			Message: "Invalid JSON message",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if requestData.Message == "false" {
		response := ResponseData{
			Status:  "Fail",
			Message: "Invalid JSON message",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println("Client's message:", requestData.Message)

	response := ResponseData{
		Status:  "success",
		Message: "Data is recieved.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	response := ResponseData{
		Status:  "Success",
		Message: "Data is recieved.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}



func main() {
	client, err := mongoDB.ConnectToMongoDB()
	if err != nil {
		log.Fatal()
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB is closed")
	} ()

	// database := client.Database("OnlineStore")
	// collection := database.Collection("Products")

	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/get", handleGet)

	http.HandleFunc("/", product.HomePage)
	http.HandleFunc("/create", product.CreateProduct)
	http.HandleFunc("/update", product.UpdateProduct)
	http.HandleFunc("/delete", product.DeleteProduct)

	log.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal()
	}
}