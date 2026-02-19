package entity

import "context"

type EntityRepository interface {
	Create(context.Context, *ParamCreate) (*ParamCreateOutput, error)
	Update(context.Context, *ParamUpdate) (*ParamUpdateOutput, error)
	Get(context.Context, *ParamGet) (*ParamGetOutput, error)
	GetByAccount(context.Context, *ParamGetByAccount) (*ParamGetByAccountOutput, error)
	EnsureIndexes(ctx context.Context) error
	RegisterTransaction(ctx context.Context, param *ParamRegisterTransaction) error
}
