package base

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/go-mail/mail/v2"
	"github.com/tpp/msf/external-adapter/mailer"
	"github.com/tpp/msf/shared/log"
)

type Usecase interface {
	Logger
	// additional method helper for usecase
	SendEmail(accessToken string, to string, templateFile string, subject string) error
}

type usecase struct {
	Logger
	// additional helper for usecase
	*mailer.Mailer
}

//go:embed templates/*
var templateFS embed.FS

type Repo struct {
	Subject       string
	ResetPassword string
	ActiveUser    string
	To            string
	Map           map[string]string
}

func (m *usecase) SendEmail(accessToken string, to string, templateFile string, subject string) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	newData := Repo{
		ResetPassword: m.Config.ResetPath + accessToken,
		ActiveUser:    m.Config.ActivePath + accessToken,
		To:            to,
		Subject:       subject,
		//Map: map[string]any,
	}
	err = tmpl.ExecuteTemplate(&body, templateFile, newData)
	if err != nil {
		return err
	}
	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetHeader("From", m.Sender)
	msg.AddAlternative("text/html", body.String())
	//msg.SetBody("text/html", "")

	return m.Dialer.DialAndSend(msg)

}

func NewBaseUsecase(usecaseName string) Usecase {
	return &usecase{
		Logger: newBaseLogger(log.Logger.With().Str("layer", fmt.Sprintf("usecase:%s", usecaseName)).Logger()),
		Mailer: mailer.GetMailInstance(),
	}
}
