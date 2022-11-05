package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host            string        `config:"host"`
	Port            int           `config:"port"`
	User            string        `config:"user"`
	Password        string        `config:"pass"`
	DBName          string        `config:"db_name"`
	MaxIdleConn     int           `config:"max_idle_conn"`
	MaxOpenConn     int           `config:"max_open_conn"`
	MaxConnLifeTime time.Duration `config:"max_conn_life_time"`
	LogmodeLevel    string        `config:"logmode_level"`
	DisableLogColor bool          `config:"disable_log_color"`
}

var mapStringLogmodeLevel = map[string]logger.LogLevel{
	"silent": logger.Silent,
	"error":  logger.Error,
	"warn":   logger.Warn,
	"info":   logger.Info,
}

var db *gorm.DB

func GetDBInstance() *gorm.DB {
	if db == nil {
		panic("database instance is not initialized")
	}
	return db
}

func New(c *Config) *gorm.DB {
	if db != nil {
		return db
	}

	if c == nil {
		panic("missing database configuration")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)

	level := mapStringLogmodeLevel[c.LogmodeLevel]
	if level == 0 {
		level = logger.Silent
	}

	var err error
	db, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      level,
				Colorful:      !c.DisableLogColor,
			}),
		})
	if err != nil {
		panic(err)
	}

	if c.MaxConnLifeTime*time.Second < time.Hour {
		c.MaxConnLifeTime = 3600
	}

	if c.MaxIdleConn == 0 {
		c.MaxIdleConn = 3
	}

	if c.MaxOpenConn == 0 {
		c.MaxOpenConn = 10
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(c.MaxConnLifeTime * time.Second)

	return db
}
