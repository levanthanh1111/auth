package main

import (
	"github.com/tpp/msf/application"
	"github.com/tpp/msf/config"
	"github.com/tpp/msf/external-adapter/db"
	mailer "github.com/tpp/msf/external-adapter/mailer"
	"github.com/tpp/msf/shared/log"
)

type Config struct {
	Service    *application.Config `config:"service"`
	Logger     *log.Config         `config:"logger"`
	DB         *db.Config          `config:"db"`
	AuthServer any                 `config:"auth_server"`
	Mailer     *mailer.Config      `config:"mailer"`
}

func main() {
	var cfg *Config
	config.Load(&cfg)

	// init logger
	log.New(cfg.Logger)

	// init db
	db.New(cfg.DB)

	//init mailer server
	mailer.NewMailer(cfg.Mailer)
	//err := sender.Sendmail("email_template.html", "anhtramvu97@gmail.com", "anhtramvu97@gmail.com", nil)

	application.Run(cfg.Service)
}
