package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aclgo/balance/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletDebitUC struct {
	repo entity.EntityRepository
	mu   *sync.Mutex
}

func NewWalletDebitUC(repo entity.EntityRepository, mu *sync.Mutex) *WalletDebitUC {
	return &WalletDebitUC{repo: repo, mu: mu}
}

type ParamDebitInput struct {
	WalletID string
	Amount   float64
}

func (p *ParamDebitInput) Validate() error {
	if p.WalletID == "" {
		return errors.New("wallet id empty")
	}

	if _, err := primitive.ObjectIDFromHex(p.WalletID); err != nil {
		return errors.New("invalid object id")
	}

	if p.Amount <= 0 {
		return fmt.Errorf("amount is %v", p.Amount)
	}

	return nil
}

type ParamDebitOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

func (u *WalletDebitUC) Execute(ctx context.Context, in *ParamDebitInput,
) (*ParamDebitOutput, error) {

	entityParamUpdate := entity.ParamUpdate{
		WalletID:  in.WalletID,
		Balance:   -in.Amount,
		UpdatedAT: time.Now(),
	}

	updatedWallet, err := u.repo.Update(ctx, &entityParamUpdate)
	if err != nil {
		return nil, fmt.Errorf("u.repo.Update: %w", err)
	}

	out := ParamDebitOutput{
		WalletID:  updatedWallet.WalletID,
		AccountID: updatedWallet.AccountID,
		Balance:   updatedWallet.Balance,
		CreatedAT: updatedWallet.CreatedAT,
		UpdatedAT: updatedWallet.UpdatedAT,
	}

	return &out, nil
}
