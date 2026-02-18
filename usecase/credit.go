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

type WalletCreditUC struct {
	repo entity.EntityRepository
	mu   *sync.Mutex
}

func NewWalletCreditUC(repo entity.EntityRepository, mu *sync.Mutex) *WalletCreditUC {
	return &WalletCreditUC{repo: repo, mu: mu}
}

type ParamCreditInput struct {
	WalletID string
	Amount   float64
}

func (p *ParamCreditInput) Validate() error {
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

type ParamCreditOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

func (u *WalletCreditUC) Execute(ctx context.Context, in *ParamCreditInput,
) (*ParamCreditOutput, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	entityParamGet := entity.ParamGet{
		WalletID: in.WalletID,
	}

	oldWallet, err := u.repo.Get(ctx, &entityParamGet)
	if err != nil {
		return nil, fmt.Errorf("u.repo.Get: %w", err)
	}

	newBalance := oldWallet.Balance + in.Amount

	entityParamUpdate := entity.ParamUpdate{
		WalletID:  oldWallet.WalletID,
		Balance:   newBalance,
		UpdatedAT: time.Now(),
	}

	updatedWallet, err := u.repo.Update(ctx, &entityParamUpdate)
	if err != nil {
		return nil, fmt.Errorf("u.repo.Update: %w", err)
	}

	out := ParamCreditOutput{
		WalletID:  updatedWallet.WalletID,
		AccountID: updatedWallet.AccountID,
		Balance:   updatedWallet.Balance,
		CreatedAT: updatedWallet.CreatedAT,
		UpdatedAT: updatedWallet.UpdatedAT,
	}

	return &out, nil
}
