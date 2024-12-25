package product

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"onlinestore/mongoDB"

	"html/template"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}

func getNextID(collection *mongo.Collection, key string) (int, error) {
	var result struct {
		Seq int `bson:"seq"`
	}

	filter := bson.M{"key": key}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Seq, nil
}

func CreateProduct(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	counterCollection := db.Collection("counters")
	itemCollection := db.Collection("Products")

	newID, err := getNextID(counterCollection, "item_id")
	if err != nil {
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	item := Item{
		ID:          newID,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       atoi(r.FormValue("price")),
		Discount:    atoi(r.FormValue("discount")),
		Quantity:    atoi(r.FormValue("quantity")),
	}

	//collection := mongoDB.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = itemCollection.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	//_, err := collection.InsertOne(ctx, item)
	//if err != nil {
	//	http.Error(w, "Failed to create product", http.StatusInternalServerError)
	//	return
	//}
	cursor, err := itemCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve items", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var items []Item
	if err = cursor.All(context.Background(), &items); err != nil {
		http.Error(w, "Failed to decode items", http.StatusInternalServerError)
		return
	}

	// Обновляем id для всех элементов
	for i, item := range items {
		_, err := itemCollection.UpdateOne(context.Background(),
			bson.M{"id": item.ID}, // _id используется MongoDB для уникальности записи
			bson.M{"$set": bson.M{"id": i + 1}})
		if err != nil {
			http.Error(w, "Failed to update product IDs", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	collection := mongoDB.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	id := atoi(r.URL.Query().Get("id"))
	if id == 0 {
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

	data := PageData{
		Products: []Item{item}, 
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}


func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	id := atoi(r.FormValue("id"))
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

func DeleteProduct(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	r.ParseForm()

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id := atoi(r.FormValue("id"))
	filter := bson.M{"id": id}

	itemCollection := db.Collection("Products")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := itemCollection.DeleteOne(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	cursor, err := itemCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve items", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var items []Item
	if err = cursor.All(context.Background(), &items); err != nil {
		http.Error(w, "Failed to decode items", http.StatusInternalServerError)
		return
	}

	// Обновляем id для всех элементов
	for i, item := range items {
		_, err := itemCollection.UpdateOne(context.Background(),
			bson.M{"id": item.ID}, // _id используется MongoDB для уникальности записи
			bson.M{"$set": bson.M{"id": i + 1}})
		if err != nil {
			http.Error(w, "Failed to update product IDs", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
