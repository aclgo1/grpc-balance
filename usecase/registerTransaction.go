package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/aclgo/balance/entity"
)

type WalletRegisterTransactionUC struct {
	repo entity.EntityRepository
}

func NewRegisterTransactionUC(repo entity.EntityRepository) *WalletRegisterTransactionUC {
	return &WalletRegisterTransactionUC{repo: repo}
}

type ParamRegisterTransactionInput struct {
	ReferenceId string
	CreatedAt   time.Time
}

func (p *ParamRegisterTransactionInput) Validate() error {
	if p.ReferenceId == "" {
		return errors.New("reference id empty")
	}

	return nil
}

func (u *WalletRegisterTransactionUC) Execute(ctx context.Context, in *ParamRegisterTransactionInput,
) error {
	p := entity.ParamRegisterTransaction{
		ReferenceId: in.ReferenceId,
		CreatedAt:   time.Now(),
	}
	err := u.repo.RegisterTransaction(ctx, &p)
	if err != nil {
		return err
	}

	return nil
}
