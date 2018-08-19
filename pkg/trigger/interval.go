package trigger

type Interval string

func (c *Interval) Cron() ([]Cron, error) {
	return []Cron{newCronSpecial("@every " + string(*c))}, nil
}
