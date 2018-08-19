package actions

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"net/mail"
	"net/smtp"

	"github.com/uphy/elastic-watcher/pkg/context"
)

type (
	SendEmailAction struct {
		Email *Email `json:"email"`
	}

	Email struct {
		Account string                `json:"account,omitempty"`
		From    Address               `json:"from"`
		To      Addresses             `json:"to"`
		CC      Addresses             `json:"cc"`
		Subject context.TemplateValue `json:"subject"`
		Body    context.TemplateValue `json:"body"`
	}

	Address struct {
		context.TemplateValue
	}
	Addresses struct {
		context.TemplateValues
	}
)

func (a Address) parse(ctx context.ExecutionContext) (*mail.Address, error) {
	rendered, err := a.TemplateValue.String(ctx)
	if err != nil {
		return nil, err
	}
	return mail.ParseAddress(rendered)
}

func (a Addresses) parse(ctx context.ExecutionContext) ([]mail.Address, error) {
	values, err := a.TemplateValues.StringSlice(ctx)
	if err != nil {
		return nil, err
	}

	var addr []mail.Address
	for _, v := range values {
		parsed, err := mail.ParseAddress(v)
		if err != nil {
			return nil, err
		}
		addr = append(addr, *parsed)
	}
	return addr, nil
}

func joinMailAddress(a []mail.Address) string {
	buf := new(bytes.Buffer)
	for i, aa := range a {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(aa.Address)
	}
	return buf.String()
}

func (l *SendEmailAction) Run(ctx context.ExecutionContext) error {
	return l.run(ctx, false)
}
func (l *SendEmailAction) DryRun(ctx context.ExecutionContext) error {
	return l.run(ctx, true)
}
func (l *SendEmailAction) run(ctx context.ExecutionContext, dryRun bool) error {
	if ctx.GlobalConfig().Email == nil {
		return errors.New("no email config")
	}
	account := ctx.GlobalConfig().Email.GetAccount(l.Email.Account)
	if account == nil {
		accountName := l.Email.Account
		if len(accountName) == 0 {
			accountName = "(default)"
		}
		return fmt.Errorf("no such account %s; please set the smtp account in the global configuration", accountName)
	}

	to, err := l.Email.To.parse(ctx)
	if err != nil {
		return fmt.Errorf("failed to render `to`: %v", err)
	}
	cc, err := l.Email.CC.parse(ctx)
	if err != nil {
		return fmt.Errorf("failed to render `cc`: %v", err)
	}
	from, err := l.Email.From.parse(ctx)
	if err != nil {
		return fmt.Errorf("failed to render `from`: %v", err)
	}
	subject, err := l.Email.Subject.String(ctx)
	if err != nil {
		return fmt.Errorf("failed to render `subject`: %v", err)
	}
	body, err := l.Email.Body.String(ctx)
	if err != nil {
		return fmt.Errorf("failed to render `body`: %v", err)
	}

	msg := new(bytes.Buffer)
	msg.WriteString("From: " + from.String())
	msg.WriteString("\r\n")
	msg.WriteString("To: " + joinMailAddress(to))
	msg.WriteString("\r\n")
	if len(cc) > 0 {
		msg.WriteString("Cc: " + joinMailAddress(cc))
		msg.WriteString("\r\n")
	}
	msg.WriteString("Subject: " + subject)
	msg.WriteString("\r\n\r\n")
	msg.WriteString(body)

	recipients := []string{}
	for _, a := range to {
		recipients = append(recipients, a.Address)
	}
	for _, a := range cc {
		recipients = append(recipients, a.Address)
	}
	var auth smtp.Auth
	if account.SMTP.Auth {
		auth = smtp.PlainAuth("", *account.SMTP.User, *account.SMTP.Password, account.SMTP.Host)
	}
	addr := fmt.Sprintf("%s:%d", account.SMTP.Host, account.SMTP.Port)
	if dryRun {
		logger := ctx.Logger()
		logger.Infof("addr: %s", addr)
		logger.Infof("auth: %#v", auth)
		logger.Infof("from: %v", from.Address)
		logger.Infof("to: %v", recipients)
		for _, s := range strings.Split(msg.String(), "\r\n") {
			logger.Info(s)
		}
		return nil
	}
	return smtp.SendMail(addr, auth, from.Address, recipients, msg.Bytes())
}
