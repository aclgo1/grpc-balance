package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/balance/entity"
	"github.com/google/uuid"
)

type WalletCreateUC struct {
	repo entity.EntityRepository
}

func NewWalletCreateUC(repo entity.EntityRepository) *WalletCreateUC {
	return &WalletCreateUC{repo: repo}
}

type ParamCreateInput struct {
	AccountID string
}

func (p *ParamCreateInput) Validate() error {
	if p.AccountID == "" {
		return errors.New("account id empty")
	}

	if _, err := uuid.Parse(p.AccountID); err != nil {
		return errors.New("invalid account uuid")
	}

	return nil
}

type ParamCreateOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

func (u *WalletCreateUC) Execute(ctx context.Context, in *ParamCreateInput,
) (*ParamCreateOutput, error) {

	entityParamCreate := entity.ParamCreate{
		AccountID: in.AccountID,
		Balance:   0,
		CreatedAT: time.Now(),
		UpdatedAT: time.Now(),
	}

	entityParamGet := entity.ParamGetByAccount{AccountID: in.AccountID}

	_, err := u.repo.GetByAccount(ctx, &entityParamGet)
	if err == nil {
		return nil, fmt.Errorf("account id exists")
	}

	created, err := u.repo.Create(ctx, &entityParamCreate)
	if err != nil {
		return nil, fmt.Errorf("u.repo.Create: %w", err)
	}

	out := ParamCreateOutput{
		WalletID:  created.WalletID,
		AccountID: created.AccountID,
		Balance:   created.Balance,
		CreatedAT: created.CreatedAT,
		UpdatedAT: created.UpdatedAT,
	}

	return &out, nil
}
