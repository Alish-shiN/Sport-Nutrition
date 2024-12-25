package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"onlinestore/mongoDB"
	"onlinestore/product"
)

type RequestData struct {
	Message string `json:"message"` // Структура для парсинга входящих данных в JSON
}

type ResponseData struct {
	Status  string `json:"status"`  // Статус обработки запроса
	Message string `json:"message"` // Сообщение ответа
}

// Обработчик POST-запроса
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { // Проверка метода запроса
		http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Инициализация структуры запроса со значением по умолчанию
	requestData := RequestData{
		Message: "false",
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData) // Декодирование JSON из тела запроса

	if err != nil {
		// Возвращаем ошибку, если JSON некорректный
		response := ResponseData{
			Status:  "Fail",
			Message: "Invalid JSON message",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if requestData.Message == "false" { // Проверяем, что поле Message содержит валидное значение
		response := ResponseData{
			Status:  "Fail",
			Message: "Invalid JSON message",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println("Client's message:", requestData.Message) // Логируем сообщение клиента

	// Формируем успешный ответ
	response := ResponseData{
		Status:  "success",
		Message: "Data is received.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Обработчик GET-запроса
func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" { // Проверка метода запроса
		http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Формируем успешный ответ
	response := ResponseData{
		Status:  "Success",
		Message: "Data is received.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Подключение к MongoDB
	client, err := mongoDB.ConnectToMongoDB()
	if err != nil {
		log.Fatal() // Завершаем выполнение, если не удалось подключиться
	}

	defer func() {
		// Закрытие соединения при завершении работы программы
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB is closed")
	}()

	db := client.Database("OnlineStore") // Выбираем базу данных

	// Регистрация маршрутов
	http.HandleFunc("/post", handlePost) // Обработчик для POST-запросов
	http.HandleFunc("/get", handleGet)   // Обработчик для GET-запросов

	http.HandleFunc("/", product.HomePage) // Главная страница
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		// Создание продукта с передачей подключения к базе
		product.CreateProduct(w, r, db)
	})
	http.HandleFunc("/update", product.UpdateProduct) // Обновление продукта

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// Удаление продукта с передачей подключения к базе
		product.DeleteProduct(w, r, db)
	})

	http.HandleFunc("/getByID", product.GetProductByID)

	fmt.Println("Starting server on http://localhost:8080")   // Уведомление о запуске сервера
	if err := http.ListenAndServe(":8080", nil); err != nil { // Запуск HTTP-сервера
		log.Fatal()
	}
}
