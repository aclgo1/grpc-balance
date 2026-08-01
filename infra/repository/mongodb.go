package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aclgo/balance/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	collection            *mongo.Collection
	transactionCollection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection, transactionCollection *mongo.Collection) entity.EntityRepository {
	return &MongoRepository{collection: collection,
		transactionCollection: transactionCollection,
	}
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

	filter := bson.M{
		"_id": id,
	}

	if param.Balance < 0 {
		requiredBalance := param.Balance * -1
		filter["balance"] = bson.M{"$gte": requiredBalance}
	}

	update := bson.M{
		"$set": bson.M{"updated_at": param.UpdatedAT},
		"$inc": bson.M{"balance": param.Balance},
	}

	pout := paramRepositoryMongoOutput{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = m.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&pout)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("m.collection.FindOneAndUpdate: wallet not searched %w", err)
		}

		return nil, fmt.Errorf("m.collection.FindOneAndUpdate: %w", err)
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

func(m *MongoRepository)ProcessLedgerEntry(ctx context.Context, param *entity.ParamLedgerEntry) (*entity.ParamUpdateOutput, error){
	session,err := m.collection.Database().Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("start session: %w",err)
	}

	defer session.EndSession(ctx)

	var pout paramRepositoryMongoOutput

	callback := func(sessCtx mongo.SessionContext)(any,error){
		id, err := primitive.ObjectIDFromHex(param.WalletId)
		if err != nil {
			return nil,err
		}

		tx := bson.M{
			"reference_id": param.ReferenceId,
			"wallet_id":    id,
			"amount":       param.Amount,
			"created_at":   param.CreatedAt,
		}

		_,err = m.transactionCollection.InsertOne(sessCtx,tx)
		if err != nil {
			if mongo.IsDuplicateKeyError(err){
				return nil, errors.New("transaction already processed")
			}

			return nil, fmt.Errorf("insert ledger entry: %w", err)
		}

		filter := bson.M{"_id":id}
		if param.Amount < 0{
			requiredBalance := param.Amount *-1
			filter["balance"] = bson.M{"$gte":requiredBalance}
		}

		update := bson.M{
			"$set": bson.M{"updated_at": param.CreatedAt},
			"$inc": bson.M{"balance": param.Amount},
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = m.collection.FindOneAndUpdate(sessCtx, filter,update,opts).Decode(pout)
		if err != nil {
			if err == mongo.ErrNoDocuments{
				return nil, fmt.Errorf("insufficient funds or wallet not found")
			}

			return nil,err
		}
		return nil,nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, err
	}

	out := entity.ParamUpdateOutput{
		WalletID:  pout.WalletID.Hex(),
		AccountID: pout.AccountID,
		Balance:   pout.Balance,
		CreatedAT: pout.CreatedAT,
		UpdatedAT: pout.UpdatedAT,
	}

	return &out,nil
}

func(m *MongoRepository) AuditWallet(ctx context.Context, param *entity.ParamAuditWallet)(*entity.AuditReport,error){
	client := m.collection.Database().Client()
	session, err := client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("error init session mongodb: %w",err)
	}

	defer session.EndSession(ctx)

	txnOpts := options.Transaction().SetReadConcern(readconcern.Snapshot()).
	SetReadPreference(readpref.Primary())

	var report entity.AuditReport

	callback := func(sessCtx mongo.SessionContext)(any,error){
		objId,err := primitive.ObjectIDFromHex(param.WalletId)
		if err != nil {
			return nil, fmt.Errorf("id wallet invalid: %w",err)
		}

		var wallet WalletDocument
		err = m.collection.FindOne(sessCtx, bson.M{"_id": objId}).Decode(&wallet)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments){
				return nil, fmt.Errorf("wallet not searched: %w",err)
			}

			return nil, fmt.Errorf("error read wallet: %w",err)
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "wallet_id", Value: objId}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}}	

		cursor, err := m.transactionCollection.Aggregate(sessCtx, mongo.Pipeline{matchStage, groupStage})
		if err != nil {
			return nil, fmt.Errorf("error calculate sum ledger: %w",err)
		}

		defer cursor.Close(sessCtx)

		var ledgerTotal int64 = 0 
		if cursor.Next(sessCtx){
			var sumResult LedgerSumResult
			if err := cursor.Decode(&sumResult); err != nil {
				return nil, fmt.Errorf("error decodify sum ledger: %w",err)
			}

			ledgerTotal = sumResult.Total
		}

		diff := wallet.Balance - ledgerTotal
		report = entity.AuditReport{
			WalletID:      param.WalletId,
			WalletBalance: wallet.Balance,
			LedgerTotal:   ledgerTotal,
			Difference:    diff,
			IsConsistent:  diff == 0,
		}

		return nil,nil
	}

	_,err = session.WithTransaction(ctx, callback,txnOpts)
	if err != nil{
		return nil, fmt.Errorf("failed transaction auditory: %w",err)
	}

	return &report,nil

}



// func (m *MongoRepository) RegisterTransaction(ctx context.Context, param *entity.ParamRegisterTransaction) error {
// 	tx := bson.M{
// 		"reference_id": param.ReferenceId,
// 		"created_at":   param.CreatedAt,
// 	}

// 	_, err := m.transactionCollection.InsertOne(ctx, tx)
// 	if mongo.IsDuplicateKeyError(err) {
// 		return errors.New("transaction already processed")
// 	}

// 	return err
// }

func (m *MongoRepository) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "reference_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := m.transactionCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("m.transactionCollection.Indexes().CreateOne: %w", err)
	}

	return nil
}
