package models

import (
	"time"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Game struct {
	ID primitive.ObjectID `bson:"_id"`
	Code string `bson:"code"`
	CreatedAt time.Time `bson:"created_at"`
	Player1Name string `bson:"player1_name"`
	Player2Name string `bson:"player2_name"`
	Player1Addr string `bson:"player1_addr"`
	Player2Addr string `bson:"player2_addr"`
	Moves []string `bson:"moves"`
}

type GameModel struct {
	DB *mongo.Collection
}

func (g *GameModel) Insert(game *Game) error {
	game.ID = primitive.NewObjectID()
	game.CreatedAt = time.Now()
	_, err := g.DB.InsertOne(context.TODO(), game)
	return err
}

func (g *GameModel) Update(game *Game) error {
	filter := bson.D{primitive.E{Key: "code", Value: game.Code}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "player1_name", Value: game.Player1Name},
		primitive.E{Key: "player2_name", Value: game.Player2Name},
		primitive.E{Key: "player1_addr", Value: game.Player1Addr},
		primitive.E{Key: "player2_addr", Value: game.Player2Addr},
		primitive.E{Key: "moves", Value: game.Moves},
	}}}
	opts := options.FindOneAndUpdate().SetUpsert(false)
	err := g.DB.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(game)
	return err
}

func (g *GameModel) Get(id string) (*Game, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := g.DB.FindOne(context.TODO(), bson.M{"_id": objectId})
	game := Game{}
	res.Decode(game)
	return &game, nil
} 
