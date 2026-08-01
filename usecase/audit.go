package usecase

import (
	"context"
	"errors"

	"github.com/aclgo/balance/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type WalletAuditUC struct {
	repo entity.EntityRepository
} 

type ParamsAudit struct {
	WalletId string
}

func (p *ParamsAudit) Validate()error{
	if p.WalletId == ""{
		return errors.New("id empty")
	}

	if _, err := primitive.ObjectIDFromHex(p.WalletId); err != nil {
		return errors.New("invalid object id")
	}

	return nil
}

type ParamsAuditOutput struct{
    WalletID      string
    WalletBalance int64
    LedgerTotal   int64
    Difference    int64
    IsConsistent  bool
}

func(u *WalletAuditUC)Execute(ctx context.Context, param *ParamsAudit)(*ParamsAuditOutput,error){

	pm := entity.ParamAuditWallet{
		WalletId: param.WalletId,
	}

	a, err := u.repo.AuditWallet(ctx,&pm)
	if err != nil {
		return nil,err
	}

	o := ParamsAuditOutput{
		WalletID      :a.WalletID,
	    WalletBalance :a.WalletBalance,
	    LedgerTotal  :a.LedgerTotal,
	    Difference   :a.Difference,
	    IsConsistent :a.IsConsistent,
	}

	return &o,nil
}