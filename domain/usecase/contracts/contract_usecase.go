package contracts

import (
	"github.com/tpp/msf/domain/repository/contracts"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/context"
)

type Usecase interface {
	List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Contract, int64, error)
	GetContract(ctx context.Context, contractId uint64) (*model.Contract, error)
}

type usecase struct {
	base.Usecase
	repo contracts.Repository
}

func (u *usecase) List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Contract, int64, error) {
	return u.repo.List(ctx, limit, offset, sort, filters)
}

func (u *usecase) GetContract(ctx context.Context, contractId uint64) (*model.Contract, error) {
	return u.repo.Get(ctx, contractId)
}

func New() Usecase {
	return &usecase{
		Usecase: base.NewBaseUsecase("contracts"),
		repo:    contracts.New(),
	}
}
