package db

import (
	"log/slog"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Game {
	player1Name string
	player2Name string
	player1Addr string
	player2Addr string
	moves []string
	moves []uint16
}

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	
	if err != nil {
		slog.Error("could not connect", "err", err)
		return
	}
	
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	
	if err != nil {
		slog.Error("could not ping", "err", err)
		return
	}
	
	slog.Info("Connected to MongoDB!")
}
