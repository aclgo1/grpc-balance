package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/balance/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type WalletReconcileUC struct {
	repo entity.EntityRepository
}

type ParamReconcile struct {
	WalletID string
	Operator string
}

func (p *ParamReconcile) Validate()error{
	if p.WalletID == ""{
		return errors.New("id empty")
	}

	if _, err := primitive.ObjectIDFromHex(p.WalletID); err != nil {
		return errors.New("invalid object id")
	}

	if p.Operator == ""{
		return errors.New("operator reconcile empty")
	}

	return nil
}

func (u *WalletReconcileUC)Execute(ctx context.Context, param *ParamReconcile)(error){
	report, err := u.repo.AuditWallet(ctx,&entity.ParamAuditWallet{WalletId: param.WalletID})
	if err != nil {
		return fmt.Errorf("failed audit to reconcile: %w",err)
	}

	if report.IsConsistent{
		return nil
	}

	adjustmentAmount := -report.Difference

	refID := fmt.Sprintf("RECONCILE_%s_%s_%d", param.WalletID,param.Operator, time.Now().UnixNano())
	
	entry := &entity.ParamLedgerEntry{
		ReferenceId: refID,
		WalletId:    param.WalletID,
		Amount:      adjustmentAmount, 
		CreatedAt:   time.Now(),
	}

	_, err = u.repo.ProcessLedgerEntry(ctx,entry)
	if err != nil {
		return fmt.Errorf("error aplicate reconcile: %w",err)
	}

	return nil
}