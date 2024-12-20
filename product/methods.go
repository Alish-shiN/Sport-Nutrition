package product

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"onlinestore/mongoDB"
	
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func atoi(s string) int {
    value, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return value
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	item := Item{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       atoi(r.FormValue("price")),
		Discount:    atoi(r.FormValue("discount")),
		Quantity:    atoi(r.FormValue("quantity")),
	}

	collection := mongoDB.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func GetProducts(w http.ResponseWriter, r *http.Request) {
	collection := mongoDB.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []Item
	for cursor.Next(ctx) {
		var item Item
		if err := cursor.Decode(&item); err != nil {
			http.Error(w, "Failed to decode product", http.StatusInternalServerError)
			return
		}
		products = append(products, item)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}


func GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	collection := mongoDB.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item Item
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&item)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	id := r.FormValue("id")
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        r.FormValue("name"),
			"description": r.FormValue("description"),
			"price":       atoi(r.FormValue("price")),
			"discount":    atoi(r.FormValue("discount")),
			"quantity":    atoi(r.FormValue("quantity")),
		},
	}

	collection := mongoDB.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	filter := bson.M{"id": id}

	collection := mongoDB.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
