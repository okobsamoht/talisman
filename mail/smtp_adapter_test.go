package mail

import (
	"testing"

	"github.com/okobsamoht/talisman/config"
	"github.com/okobsamoht/talisman/types"
)

func Test_smtp(t *testing.T) {
	config.TConfig = &config.Config{
		SMTPServer:   "smtp.163.com",
		MailUsername: "user@163.com",
		MailPassword: "password",
	}

	s := NewSMTPAdapter()
	object := types.M{
		"text":    "text from talisman",
		"to":      "user@163.com",
		"subject": "talisman send",
	}
	s.SendMail(object)
}
