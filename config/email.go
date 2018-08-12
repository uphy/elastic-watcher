package config

type (
	Email struct {
		DefaultAccount string              `json:"default_account,omitempty"`
		Accounts       map[string]*Account `json:"account,omitempty"`
	}
	Account struct {
		Profile *string `json:"profile"`
		SMTP    SMTP    `json:"smtp"`
	}
	SMTP struct {
		Auth     bool     `json:"auth"`
		StartTLS StartTLS `json:"starttls,omitempty"`
		Host     string   `json:"host"`
		Port     int      `json:"port"`
		User     *string  `json:"user,omitempty"`
		Password *string  `json:"password,omitempty"`
	}
	StartTLS struct {
		Enable   bool  `json:"enable"`
		Required *bool `json:"required,omitempty"`
	}
)

func (e *Email) GetDefaultAccount() *Account {
	return e.GetAccount("")
}

func (e *Email) GetAccount(name string) *Account {
	if name == "" {
		return e.Accounts[e.DefaultAccount]
	}
	return e.Accounts[name]
}
