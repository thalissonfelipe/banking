package account

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	collection *mongo.Collection
}
