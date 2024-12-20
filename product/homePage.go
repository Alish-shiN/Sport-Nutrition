package product

import (
	"context"
	"html/template"
	"net/http"
	
	"onlinestore/mongoDB"

	"go.mongodb.org/mongo-driver/bson"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
    // Получаем продукты из базы данных
    collection := mongoDB.GetCollection()
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, "Error fetching products", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    var products []Item
    if err := cursor.All(context.TODO(), &products); err != nil {
        http.Error(w, "Error reading products", http.StatusInternalServerError)
        return
    }

    // Рендерим шаблон
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, struct {
        Products []Item
    }{
        Products: products,
    })
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}
