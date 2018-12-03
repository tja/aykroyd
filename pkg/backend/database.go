package backend

type Database struct {
}

func NewDatabase() (*Database, error) {
	return &Database{}, nil
}

type EmailForward struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Domain struct {
	Name   string          `json:"name"`
	Emails []*EmailForward `json:"emails"`
}

func (m *Database) Domains() []*Domain {
	return []*Domain{}
}
