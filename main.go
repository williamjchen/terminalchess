package main

import (
	"log/slog"
	"context"
	"os"

	"github.com/williamjchen/terminalchess/server"	

	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func connect() *mongo.Collection {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_STRING")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("chess").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	slog.Info("Pinged your deployment. You successfully connected to MongoDB!")

	c1 := client.Database("chess").Collection("game_data")
	return c1
}

func main() {
	c1 := connect()
	s, err := server.NewServer("./.ssh/term_info_ed25519", "0.0.0.0", 2324, c1)
	if err != nil {
		return
	}
	s.Start()
}

