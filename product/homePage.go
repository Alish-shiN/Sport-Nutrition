package product

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"onlinestore/mongoDB"

	"go.mongodb.org/mongo-driver/bson"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

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
		cursor.Decode(&item)
		products = append(products, item)
	}

	tmpl.Execute(w, struct {
		Products []Item
	}{
		Products: products,
	})
}
