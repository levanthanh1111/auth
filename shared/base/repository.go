package base

import (
	"fmt"
	"strings"

	"github.com/tpp/msf/shared/context"
	"github.com/tpp/msf/shared/log"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type Scope func(*gorm.DB) *gorm.DB

var (
	noneScope        Scope = func(db *gorm.DB) *gorm.DB { return db }
	mapQueryOperator       = map[Operator]string{
		OperatorEqual:              "=",
		OperatorNotEqual:           "IS NOT",
		OperatorLessThan:           "<",
		OperatorLessThanOrEqual:    "<=",
		OperatorGreaterThan:        ">",
		OperatorGreaterThanOrEqual: ">=",
		OperatorLike:               "like",
	}
)

type Repository interface {
	Logger
	// DB get database transaction instance from context
	DB(ctx context.Context) *gorm.DB
	// additional method helper for repository
	Filter(filters []*Filter) Scope
	Paginate(limit int, offset int, sort string) Scope
	RegisterActualKey(presentKey, actualKey string)
	RegisterJoinScope(presentKey string, scopes ...Scope)
}

type repo struct {
	Logger
	// additional helper for repository
	mapJoinScopes map[string][]*Scope
	mapActualKeys map[string]string
}

func (r *repo) RegisterActualKey(presentKey, actualKey string) {
	r.mapActualKeys[presentKey] = actualKey
}

func (r *repo) RegisterJoinScope(presentKey string, scopes ...Scope) {
	r.mapJoinScopes[presentKey] = make([]*Scope, len(scopes))
	for idx, scope := range scopes {
		r.mapJoinScopes[presentKey][idx] = &scope
	}
}

func (r *repo) Filter(filters []*Filter) Scope {
	if len(filters) == 0 {
		return noneScope
	}

	var queryStr string
	var queryFilters, values = []string{}, []any{}

	var joinScopes = make([]*Scope, 0)

	for _, filter := range filters {
		operator := mapQueryOperator[filter.Operator]
		if operator == "" || filter.Key == "" {
			continue
		}

		var value string
		if operator == "like" {
			value = fmt.Sprintf("%%%s%%", filter.Value)
		} else {
			value = filter.Value
		}

		if len(r.mapJoinScopes[filter.Key]) > 0 {
			for _, scope := range r.mapJoinScopes[filter.Key] {
				if scope != nil && !slices.Contains(joinScopes, scope) {
					joinScopes = append(joinScopes, scope)
				}
			}
		}

		key := filter.Key
		if r.mapActualKeys[key] != "" {
			key = r.mapActualKeys[key]
		}

		queryFilters = append(queryFilters, fmt.Sprintf("%s %s ?", key, operator))
		values = append(values, value)
	}

	queryStr = strings.Join(queryFilters, " AND ")

	return func(db *gorm.DB) *gorm.DB {
		for _, scopePtr := range joinScopes {
			db = (*scopePtr)(db)
		}
		return db.Where(queryStr, values...)
	}
}

func (r *repo) Paginate(limit int, offset int, sort string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit).Order(sort)
	}
}

func (r *repo) DB(ctx context.Context) *gorm.DB {
	return context.DBTxFromContext(ctx)
}

func NewBaseRepository(repositoryName string) Repository {
	return &repo{
		Logger:        newBaseLogger(log.Logger.With().Str("layer", fmt.Sprintf("repository:%s", repositoryName)).Logger()),
		mapActualKeys: map[string]string{},
		mapJoinScopes: map[string][]*Scope{},
	}
}
