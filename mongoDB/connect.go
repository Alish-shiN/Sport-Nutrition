package mongoDB

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectToMongoDB() (*mongo.Client, error) {
    const uri = "mongodb://localhost:27017/"

    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.TODO(), clientOpts)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
    }

    Client = client
    fmt.Println("Successful connection.")
    return client, nil
}

func GetCollection() *mongo.Collection {
    if Client == nil {
        panic("MongoDB client is not initialized")
    }
    return Client.Database("OnlineStore").Collection("Products")
}
