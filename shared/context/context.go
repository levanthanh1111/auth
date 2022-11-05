package context

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/tpp/msf/model"
)

type set interface {
	// WithDBTx to save the database transactions
	WithDBTx(*gorm.DB) Context
	// WithUser to save the user information
	WithUser(*model.User) Context
	// WithRole to save user's role
	WithRole(model.RoleName) Context
	// WithReqID to save user's role
	WithReqID(string) Context
}

type get interface {
	// ReqID
	ReqID() string
	// User
	User() *model.User
	// Role
	Role() model.RoleName
}

// Context wrapped golang based context
type Context interface {
	context.Context
	set
	get
}

// CancelFunc wrapped from base context
type CancelFunc = context.CancelFunc

type ctxKey struct{ value string }

var (
	ctxDBKey    = &ctxKey{"database transaction"}
	ctxUserKey  = &ctxKey{"user"}
	ctxRoleKey  = &ctxKey{"role"}
	ctxReqIDKey = &ctxKey{"request id"}
)

type appContext struct {
	context.Context
}

func (ctx *appContext) WithDBTx(db *gorm.DB) Context {
	ctx.Context = context.WithValue(ctx.Context, ctxDBKey, db)
	return ctx
}

func (ctx *appContext) WithUser(user *model.User) Context {
	ctx.Context = context.WithValue(ctx.Context, ctxUserKey, user)
	return ctx
}

func (ctx *appContext) WithRole(r model.RoleName) Context {
	ctx.Context = context.WithValue(ctx.Context, ctxRoleKey, r)
	return ctx
}

func (ctx *appContext) WithReqID(s string) Context {
	ctx.Context = context.WithValue(ctx.Context, ctxReqIDKey, s)
	return ctx
}

func valueFromCtx[V string | *model.User | model.RoleName | *gorm.DB](ctx context.Context, key *ctxKey) V {
	v, _ := ctx.Value(key).(V)
	return v
}

func (ctx *appContext) ReqID() string {
	return valueFromCtx[string](ctx, ctxReqIDKey)
}

func (ctx *appContext) User() *model.User {
	return valueFromCtx[*model.User](ctx, ctxUserKey)
}

func (ctx *appContext) Role() model.RoleName {
	return valueFromCtx[model.RoleName](ctx, ctxRoleKey)
}

// DBTxFromContext get database from context
func DBTxFromContext(ctx Context) *gorm.DB {
	return valueFromCtx[*gorm.DB](ctx, ctxDBKey)
}

// FromBaseContext wrapping base context
func FromBaseContext(ctx context.Context) Context {
	switch v := ctx.(type) {
	case *appContext:
		return v
	default:
		return &appContext{Context: v}
	}
}

// WithDeadline wrapped base context.WithDeadline
func WithDeadline(parent context.Context, d time.Time) (Context, CancelFunc) {
	switch v := parent.(type) {
	case *appContext:
		ctx, cf := context.WithDeadline(v.Context, d)
		v.Context = ctx
		return v, cf
	default:
		ctx, cf := context.WithDeadline(v, d)
		return &appContext{Context: ctx}, cf
	}
}

// WithCancel wrapped base context.WithCancel
func WithCancel(parent context.Context) (Context, CancelFunc) {
	switch v := parent.(type) {
	case *appContext:
		ctx, cf := context.WithCancel(v.Context)
		v.Context = ctx
		return v, cf
	default:
		ctx, cf := context.WithCancel(v)
		return &appContext{Context: ctx}, cf
	}
}

// WithTimeout wrapped base context.WithTimeout
func WithTimeout(parent context.Context, d time.Duration) (Context, CancelFunc) {
	switch v := parent.(type) {
	case *appContext:
		ctx, cf := context.WithTimeout(v.Context, d)
		v.Context = ctx
		return v, cf
	default:
		ctx, cf := context.WithTimeout(v, d)
		return &appContext{Context: ctx}, cf
	}
}

// WithValue wrapped base context.WithValue
func WithValue(parent context.Context, key any, val any) Context {
	switch v := parent.(type) {
	case *appContext:
		v.Context = context.WithValue(v.Context, key, val)
		return v
	default:
		ctx := context.WithValue(v, key, val)
		return &appContext{Context: ctx}
	}
}

// TODO wrapped base context.TODO
func TODO() Context {
	return &appContext{Context: context.TODO()}
}

// Background wrapped base context.Background
func Background() Context {
	return &appContext{Context: context.Background()}
}
