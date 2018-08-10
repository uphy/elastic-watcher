package actions

import (
	"bytes"
	"fmt"

	"net/mail"
	"net/smtp"

	"github.com/uphy/elastic-watcher/watcher/context"
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

	Address   context.TemplateValue
	Addresses []Address
)

func (a Address) parse(ctx context.ExecutionContext) (*mail.Address, error) {
	rendered, err := context.TemplateValue(a).String(ctx)
	if err != nil {
		return nil, err
	}
	return mail.ParseAddress(rendered)
}

func (a Addresses) parse(ctx context.ExecutionContext) ([]mail.Address, error) {
	var addr []mail.Address
	for _, v := range a {
		parsed, err := v.parse(ctx)
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

	msg := "From:" + from.String() + "\r\n" +
		"To:" + joinMailAddress(to) + "\r\n" +
		"Cc:" + joinMailAddress(cc) + "\r\n" +
		"Subject:" + subject + "\r\n" +
		"\r\n" + body

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
	return smtp.SendMail(fmt.Sprintf("%s:%d", account.SMTP.Host, account.SMTP.Port), auth, from.Address, recipients, []byte(msg))
}
