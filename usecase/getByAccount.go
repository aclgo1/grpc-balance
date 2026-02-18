package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/balance/entity"
	"github.com/google/uuid"
)

type WalletGetByAccountUC struct {
	repo entity.EntityRepository
}

func NewWalletGetByAccountUC(repo entity.EntityRepository) *WalletGetByAccountUC {
	return &WalletGetByAccountUC{repo: repo}
}

type ParamGetByAccountInput struct {
	AccountID string
}

func (p *ParamGetByAccountInput) Validate() error {
	if p.AccountID == "" {
		return errors.New("account id empty")
	}

	if _, err := uuid.Parse(p.AccountID); err != nil {
		return errors.New("invalid account uuid")
	}

	return nil
}

type ParamGetByAccountOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

func (u *WalletGetByAccountUC) Execute(ctx context.Context, in *ParamGetByAccountInput,
) (*ParamGetByAccountOutput, error) {

	entityParamGetByAccount := entity.ParamGetByAccount{
		AccountID: in.AccountID,
	}

	find, err := u.repo.GetByAccount(ctx, &entityParamGetByAccount)
	if err != nil {
		return nil, fmt.Errorf("u.repo.GetByAccount: %w", err)
	}

	out := ParamGetByAccountOutput{
		WalletID:  find.WalletID,
		AccountID: find.AccountID,
		Balance:   find.Balance,
		CreatedAT: find.CreatedAT,
		UpdatedAT: find.UpdatedAT,
	}

	return &out, nil
}
