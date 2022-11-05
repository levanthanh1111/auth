package utils

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tpp/msf/shared/context"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestTeardown func()

type TestSetup[V any] func(*testing.T, *require.Assertions, V) TestTeardown

func NewDBMock() (*gorm.DB, sqlmock.Sqlmock, func(), error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, nil, err
	}
	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gDB, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      false,
		}),
	})
	gDB = gDB.Debug()

	return gDB, mock, func() {
		db.Close()
	}, err
}

func NewContextWithDB() (context.Context, sqlmock.Sqlmock, func(), error) {
	gDB, mock, cnl, err := NewDBMock()
	if err != nil {
		return nil, nil, nil, err
	}
	ctx := context.Background()
	ctx.WithDBTx(gDB)
	return ctx, mock, cnl, err
}
