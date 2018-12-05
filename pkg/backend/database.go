package backend

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/tja/postfix-web/pkg/models"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(connection string) (*Database, error) {
	// Setup Gorm
	db, err := gorm.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Domain{}, &models.Forward{})

	return &Database{DB: db}, nil
}

func (m *Database) Close() {
	m.DB.Close()
}

type Forward struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Domain struct {
	Name     string     `json:"name"`
	Forwards []*Forward `json:"forwards"`
}

func (m *Database) ListDomains() []*Domain {
	// Fetch data
	resp := []*models.Domain{}
	if err := m.DB.Set("gorm:auto_preload", true).Find(&resp).Error; err != nil {
		return []*Domain{}
	}

	// Convert
	domains := make([]*Domain, 0, len(resp))
	for _, r := range resp {
		forwards := make([]*Forward, 0, len(r.Forwards))
		for _, f := range r.Forwards {
			forwards = append(forwards, &Forward{
				From: f.From,
				To:   f.To,
			})
		}

		domains = append(domains, &Domain{
			Name:     r.Name,
			Forwards: forwards,
		})
	}

	return domains
}
