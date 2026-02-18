package service

import (
	"context"
	"fmt"

	"github.com/aclgo/balance/proto"
	"github.com/aclgo/balance/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcService struct {
	CreateUC       *usecase.WalletCreateUC
	CreditUC       *usecase.WalletCreditUC
	DebitUC        *usecase.WalletDebitUC
	GetByAccountUC *usecase.WalletGetByAccountUC
	proto.UnimplementedWalletServiceServer
}

func NewGrpcService(createUC *usecase.WalletCreateUC,
	creditUC *usecase.WalletCreditUC,
	debitUC *usecase.WalletDebitUC,
	getByAccountUC *usecase.WalletGetByAccountUC,
) *GrpcService {
	return &GrpcService{
		CreateUC:       createUC,
		CreditUC:       creditUC,
		DebitUC:        debitUC,
		GetByAccountUC: getByAccountUC,
	}
}

func (s *GrpcService) Create(ctx context.Context, in *proto.ParamCreateWalletRequest,
) (*proto.ParamCreateWalletResponse, error) {

	p := usecase.ParamCreateInput{
		AccountID: in.AccountID,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	created, err := s.CreateUC.Execute(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.CreateUC.Execute: %w", err)
	}

	out := proto.ParamCreateWalletResponse{
		WalletID:  created.WalletID,
		AccountID: created.AccountID,
		Balance:   created.Balance,
		CreatedAT: timestamppb.New(created.CreatedAT),
		UpdatedAT: timestamppb.New(created.UpdatedAT),
	}

	return &out, nil
}
func (s *GrpcService) Credit(ctx context.Context, in *proto.ParamCreditWalletRequest,
) (*proto.ParamCreditWalletResponse, error) {

	p := usecase.ParamCreditInput{
		WalletID: in.WalletID,
		Amount:   in.Amount,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	updated, err := s.CreditUC.Execute(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.CreditUC.Execute: %w", err)
	}

	out := proto.ParamCreditWalletResponse{
		WalletID:  updated.WalletID,
		AccountID: updated.AccountID,
		Balance:   updated.Balance,
		CreatedAT: timestamppb.New(updated.CreatedAT),
		UpdatedAT: timestamppb.New(updated.UpdatedAT),
	}

	return &out, nil
}
func (s *GrpcService) Debit(ctx context.Context, in *proto.ParamDebitWalletRequest,
) (*proto.ParamDebitWalletResponse, error) {

	p := usecase.ParamDebitInput{
		WalletID: in.WalletID,
		Amount:   in.Amount,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	updated, err := s.DebitUC.Execute(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.DebitUC.Execute: %w", err)
	}

	out := proto.ParamDebitWalletResponse{
		WalletID:  updated.WalletID,
		AccountID: updated.AccountID,
		Balance:   updated.Balance,
		CreatedAT: timestamppb.New(updated.CreatedAT),
		UpdatedAT: timestamppb.New(updated.UpdatedAT),
	}

	return &out, nil
}

func (s *GrpcService) GetWalletByAccount(ctx context.Context, in *proto.ParamGetWalletByAccountRequest,
) (*proto.ParamgGetWalletByAccountResponse, error) {

	p := usecase.ParamGetByAccountInput{
		AccountID: in.AccountID,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	find, err := s.GetByAccountUC.Execute(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.GetByAccountUC.Execute: %w", err)
	}

	out := proto.ParamgGetWalletByAccountResponse{
		WalletID:  find.WalletID,
		AccountID: find.AccountID,
		Balance:   find.Balance,
		CreatedAT: timestamppb.New(find.CreatedAT),
		UpdatedAT: timestamppb.New(find.UpdatedAT),
	}

	return &out, nil
}
