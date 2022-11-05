//go:generate mockery --name=Repository
package contracts

import (
	entity "github.com/tpp/msf/domain/repository"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/context"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Contract, int64, error)
	Get(ctx context.Context, contractID uint64) (*model.Contract, error)
}

type repo struct {
	base.Repository
}

func (r *repo) List(ctx context.Context, limit, offset int, sort string, filters base.Filters) ([]*model.Contract, int64, error) {
	var contracts []*entity.Contract
	var total int64

	err := r.DB(ctx).
		Model(&entity.Contract{}).
		Scopes(r.Filter(filters)).
		Count(&total).
		Scopes(r.Paginate(limit, offset, sort)).
		Preload("SuplyVendor").
		Find(&contracts).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	for _, c := range contracts {
		c.UpdateVendorName()
	}

	return contracts, total, nil
}

func (r *repo) Get(ctx context.Context, contractID uint64) (*model.Contract, error) {
	var contract *entity.Contract
	err := r.DB(ctx).
		Model(&entity.Contract{}).
		Preload("SuplyVendor").
		Take(&contract, contractID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, base.ErrorNotFound
		}
		r.Error(ctx).Err(err).Msg("GetContractError")
		return nil, err
	}

	contract.UpdateVendorName()
	return contract, nil
}

func New() Repository {
	baseRepo := base.NewBaseRepository("contracts")

	baseRepo.RegisterActualKey("supply_vendor_name", "orgs.name")
	baseRepo.RegisterJoinScope("supply_vendor_name", func(d *gorm.DB) *gorm.DB { return d.Joins(`INNER JOIN "orgs" ON "supply_vendor_id" = "orgs"."id"`) })

	return &repo{baseRepo}
}
