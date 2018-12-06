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

	err := m.DB.
		Set("gorm:auto_preload", true).
		Find(&resp).
		Error

	if err != nil {
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

func (m *Database) CreateDomain(name string) error {
	// Store new domain
	return m.DB.
		LogMode(false).
		Create(&models.Domain{Name: name}).
		Error
}

func (m *Database) DeleteDomain(name string) error {
	// Get domain
	var domain models.Domain

	err := m.DB.
		Where(&models.Domain{Name: name}).
		First(&domain).
		Error

	if err != nil {
		return err
	}

	// Delete all forwards
	err = m.DB.
		Where(&models.Forward{DomainID: domain.ID}).
		Delete(&models.Forward{}).
		Error

	if err != nil {
		return err
	}

	// Delete domain itself
	return m.DB.
		Delete(&domain).
		Error
}

func (m *Database) CreateForward(name string, from string, to string) error {
	// Get domain
	var domain models.Domain

	err := m.DB.
		Where(&models.Domain{Name: name}).
		First(&domain).
		Error

	if err != nil {
		return err
	}

	// Store new forward
	return m.DB.
		LogMode(false).
		Model(&domain).
		Association("Forwards").
		Append(&models.Forward{From: from, To: to}).
		Error
}

func (m *Database) UpdateForward(name string, from string, to string) error {
	// Get domain
	var domain models.Domain

	err := m.DB.
		Where(&models.Domain{Name: name}).
		First(&domain).
		Error

	if err != nil {
		return err
	}

	// Find forward
	var forward models.Forward

	err = m.DB.
		Where(&models.Forward{From: from, DomainID: domain.ID}).
		First(&forward).
		Error

	if err != nil {
		return err
	}

	// Update existing forward
	return m.DB.
		LogMode(false).
		Model(&forward).
		Update(&models.Forward{To: to}).
		Error
}

func (m *Database) DeleteForward(name string, from string) error {
	// Get domain
	var domain models.Domain

	err := m.DB.
		Where(&models.Domain{Name: name}).
		First(&domain).
		Error

	if err != nil {
		return err
	}

	// Delete forward
	return m.DB.
		Where(&models.Forward{From: from, DomainID: domain.ID}).
		Delete(&models.Forward{}).
		Error
}
