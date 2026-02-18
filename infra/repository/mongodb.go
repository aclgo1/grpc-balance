package repository

import (
	"context"
	"fmt"

	"github.com/aclgo/balance/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) entity.EntityRepository {
	return &MongoRepository{collection: collection}
}

func (m *MongoRepository) Create(ctx context.Context, param *entity.ParamCreate,
) (*entity.ParamCreateOutput, error) {

	pin := paramRepositoryMongoInput{
		AccountID: param.AccountID,
		Balance:   param.Balance,
		CreatedAT: param.CreatedAT,
		UpdatedAT: param.UpdatedAT,
	}

	result, err := m.collection.InsertOne(ctx, &pin)
	if err != nil {
		return nil, fmt.Errorf("m.collection.InsertOne: %w", err)
	}

	insertID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("result.InsertID.(primitive.ObjectID): %w", err)
	}

	filter := bson.M{"_id": insertID}
	pout := paramRepositoryMongoOutput{}

	err = m.collection.FindOne(ctx, filter).Decode(&pout)
	if err != nil {
		return nil, fmt.Errorf("m.collection.FindOne: %w", err)
	}

	out := entity.ParamCreateOutput{
		WalletID:  pout.WalletID.Hex(),
		AccountID: pout.AccountID,
		Balance:   pout.Balance,
		CreatedAT: pout.CreatedAT,
		UpdatedAT: pout.UpdatedAT,
	}

	return &out, nil
}
func (m *MongoRepository) Update(ctx context.Context, param *entity.ParamUpdate,
) (*entity.ParamUpdateOutput, error) {
	id, err := primitive.ObjectIDFromHex(param.WalletID)
	if err != nil {
		return nil, fmt.Errorf("primitive.ObjectIDFromHex: %w", err)
	}

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{"balance": param.Balance, "updated_at": param.UpdatedAT},
	}

	res, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("m.collection.UpdateOne: %w", err)
	}

	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("wallet no searched")
	}

	pout := paramRepositoryMongoOutput{}

	if err := m.collection.FindOne(ctx, filter).Decode(&pout); err != nil {
		return nil, fmt.Errorf("m.collection.FindOne: %w", err)
	}

	out := entity.ParamUpdateOutput{
		WalletID:  pout.WalletID.Hex(),
		AccountID: pout.AccountID,
		Balance:   pout.Balance,
		CreatedAT: pout.CreatedAT,
		UpdatedAT: pout.UpdatedAT,
	}

	return &out, nil
}

func (m *MongoRepository) Get(ctx context.Context, param *entity.ParamGet,
) (*entity.ParamGetOutput, error) {

	id, err := primitive.ObjectIDFromHex(param.WalletID)
	if err != nil {
		return nil, fmt.Errorf("primitive.ObjectIDFromHex: %w", err)
	}

	filter := bson.M{"_id": id}

	pout := paramRepositoryMongoOutput{}

	if err := m.collection.FindOne(ctx, filter).Decode(&pout); err != nil {
		return nil, fmt.Errorf("m.collection.FindOne: %w", err)
	}

	out := entity.ParamGetOutput{
		WalletID:  pout.WalletID.Hex(),
		AccountID: pout.AccountID,
		Balance:   pout.Balance,
		CreatedAT: pout.CreatedAT,
		UpdatedAT: pout.UpdatedAT,
	}

	return &out, nil
}

func (m *MongoRepository) GetByAccount(ctx context.Context, param *entity.ParamGetByAccount,
) (*entity.ParamGetByAccountOutput, error) {

	filter := bson.M{"account_id": param.AccountID}

	pout := paramRepositoryMongoOutput{}

	if err := m.collection.FindOne(ctx, filter).Decode(&pout); err != nil {
		return nil, fmt.Errorf("m.collection.FindOne: %w", err)
	}

	out := entity.ParamGetByAccountOutput{
		WalletID:  pout.WalletID.Hex(),
		AccountID: pout.AccountID,
		Balance:   pout.Balance,
		CreatedAT: pout.CreatedAT,
		UpdatedAT: pout.UpdatedAT,
	}

	return &out, nil
}
