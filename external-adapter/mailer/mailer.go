package mailer

import (
	"time"

	"github.com/go-mail/mail/v2"
)

type Config struct {
	Timeout    time.Duration `config:"timeout"`
	Host       string        `config:"host"`
	Port       int           `config:"port"`
	Username   string        `config:"username"`
	Password   string        `config:"password"`
	Sender     string        `config:"sender"`
	ResetPath  string        `config:"reset_path"`
	ActivePath string        `config:"active_path"`
}
type Mailer struct {
	Dialer *mail.Dialer
	Config *Config
	Sender string
}

var mailer *Mailer

func NewMailer(config *Config) *Mailer {
	dialer := mail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	dialer.Timeout = config.Timeout

	mailer = &Mailer{
		Dialer: dialer,
		Sender: config.Sender,
		Config: config,
	}
	return mailer

}
func GetMailInstance() *Mailer {
	if mailer == nil {
		panic("mailer instance is not initialized")
	}
	return mailer
}
