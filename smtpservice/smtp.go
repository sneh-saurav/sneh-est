package smtp

import (
	"encoding/json"
	"os"

	"gopkg.in/gomail.v2"
)

var localConf *config

const defConfPath = "../.././config/smtp.json"

type config struct {
	From      string `json:"from"`
	Host      string `json:"host"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Templates string `json:"template_path"`
}

// Configure will configure the smtp server
func Configure() error {
	f, err := os.Open(defConfPath)
	if err != nil {
		return err
	}
	defer f.Close()

	conf := config{}
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		return err
	}

	localConf = &conf
	return nil
}

// SMTP represents an SMTP connection
type SMTP struct {
	from string
}

// SetFrom will override the default from address
func (s *SMTP) SetFrom(email string) {
	s.from = email
}

// Send will send an email using the SMTP parameters
func (s SMTP) Send(subject, body string, attach []string, to string) error {
	if localConf == nil {
		if err := Configure(); err != nil {
			return err
		}
	}
	dialer := gomail.NewDialer(localConf.Host, localConf.Port, localConf.User, localConf.Password)
	con, err := dialer.Dial()
	if err != nil {
		return err
	}
	from := localConf.From
	if s.from != "" {
		from = s.from
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetAddressHeader("To", to, to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if len(attach) > 0 {
		for _, p := range attach {
			m.Attach(p)
		}
	}
	return gomail.Send(con, m)
}
