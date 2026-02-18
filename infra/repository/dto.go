package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type paramRepositoryMongoInput struct {
	WalletID  primitive.ObjectID `bson:"_id,omitempty"`
	AccountID string             `bson:"account_id"`
	Balance   float64            `bson:"balance"`
	CreatedAT time.Time          `bson:"created_at"`
	UpdatedAT time.Time          `bson:"updated_at"`
}

type paramRepositoryMongoOutput struct {
	WalletID  primitive.ObjectID `bson:"_id,omitempty"`
	AccountID string             `bson:"account_id"`
	Balance   float64            `bson:"balance"`
	CreatedAT time.Time          `bson:"created_at"`
	UpdatedAT time.Time          `bson:"updated_at"`
}
