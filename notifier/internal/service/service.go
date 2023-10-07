package service

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/miha3009/market/notifier/pkg/config"
)

type MailService interface {
	SendMail(to, subject, body string)
}

type MailServiceImpl struct {
	logger *log.Logger
	cfg    config.EmailConfig
}

func NewMailService(logger *log.Logger, cfg config.EmailConfig) MailService {
	return &MailServiceImpl{
		logger: logger,
		cfg:    cfg,
	}
}

func (s *MailServiceImpl) SendMail(to, subject, body string) {
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n\r\n"+
		"%s\r\n", s.cfg.From, to, subject, body))

	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	err := smtp.SendMail(addr, auth, s.cfg.From, []string{to}, msg)

	if err != nil {
		log.Fatal(err)
	}
}
