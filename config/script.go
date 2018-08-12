package config

type Script struct {
	Lang   *string `json:"lang,omitempty"`
	Source string  `json:"source"`
}
