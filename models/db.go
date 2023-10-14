package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Model struct {
	Games GameModel
}

func CreateModel(c1 *mongo.Collection) *Model {
	return &Model {
		Games: GameModel {
			DB: c1,
		},
	}
}
