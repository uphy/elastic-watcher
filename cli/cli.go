package cli

import (
	"github.com/uphy/elastic-watcher/config"
	"github.com/urfave/cli"
)

const version = "0.0.1"

type CLI struct {
	app          *cli.App
	globalConfig *config.Config
}

func NewCLI() *CLI {
	c := &CLI{}

	app := cli.NewApp()
	app.Name = "elastic-watcher"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "url",
			Value:  "http://localhost:9200",
			EnvVar: "ES_URL",
		},
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "config",
			EnvVar: "CONFIG_FILE",
		},
		cli.BoolFlag{
			Name:   "mail-smtp-auth",
			EnvVar: "MAIL_SMTP_AUTH",
		},
		cli.BoolFlag{
			Name:   "mail-smtp-starttls-enable",
			EnvVar: "MAIL_SMTP_STARTTLS_ENABLE",
		},
		cli.BoolFlag{
			Name:   "mail-smtp-host",
			EnvVar: "MAIL_SMTP_HOST",
		},
		cli.BoolFlag{
			Name:   "mail-smtp-port",
			EnvVar: "MAIL_SMTP_PORT",
		},
		cli.BoolFlag{
			Name:   "mail-smtp-user",
			EnvVar: "MAIL_SMTP_USER",
		},
		cli.BoolFlag{
			Name:   "mail-smtp-password",
			EnvVar: "MAIL_SMTP_PASSWORD",
		},
	}
	app.Before = func(ctx *cli.Context) error {
		conf := &config.Config{}
		// read from file
		if ctx.IsSet("config") {
			cc, err := config.LoadFile(ctx.String("config"))
			if err != nil {
				return err
			}
			conf = cc
		}
		// overwrite with command options
		if ctx.IsSet("debug") {
			conf.Debug = ctx.Bool("debug")
		}
		conf.Elasticsearch.URL = ctx.String("url")
		if ctx.IsSet("mail-smtp-host") {
			account := config.Account{}
			auth := false
			if ctx.IsSet("mail-smtp-auth") {
				auth = ctx.Bool("mail-smtp-auth")
			}
			account.SMTP = config.SMTP{
				Auth: auth,
				Host: ctx.String("mail-smtp-host"),
				Port: ctx.Int("mail-smtp-port"),
			}
			if account.SMTP.Auth {
				user := ctx.String("mail-smtp-user")
				password := ctx.String("mail-smtp-password")
				account.SMTP.User = &user
				account.SMTP.Password = &password
			}
			if ctx.IsSet("mail-smtp-starttls-enable") {
				account.SMTP.StartTLS.Enable = ctx.Bool("mail-smtp-starttls-enable")
			}
			conf.Email.DefaultAccount = "default"
			conf.Email.Accounts["default"] = &account
		}
		if err := conf.Validate(); err != nil {
			return err
		}
		c.globalConfig = conf
		return nil
	}
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, c.run())
	app.Commands = append(app.Commands, c.config())
	return &CLI{
		app: app,
	}
}

func (c *CLI) Run(args []string) error {
	return c.app.Run(args)
}

func main() {

}
