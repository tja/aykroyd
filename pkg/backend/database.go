package backend

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/tja/aykroyd/pkg/models"
)

// Domain is the caller-facing replica of the internal database model. It contains the state of the room domain
// object, with reference to all email forwards.
type Domain struct {
	Name     string     `json:"name"`
	Forwards []*Forward `json:"forwards"`
}

// Forward is the caller-facing replica of the internal database model. It contains the state of the leaf email
// forward object.
type Forward struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Database holds the state of the persistent storage.
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new persistent storage instance. host, database, username, and password hold relevant
// information to connect to MySQL.
func NewDatabase(host, database, username, password string) (*Database, error) {
	// Setup Gorm
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			username,
			password,
			host,
			database,
		),
	)

	if err != nil {
		return nil, err
	}

	// Update tables
	db.AutoMigrate(&models.Domain{}, &models.Forward{})

	return &Database{DB: db}, nil
}

// Close closes the persistent storage instance.
func (m *Database) Close() error {
	return m.DB.Close()
}

// Domains returns a list of all domains, including the associated email forwards.
func (m *Database) Domains() ([]*Domain, error) {
	// Fetch data
	var resp []*models.Domain

	err := m.DB.
		Set("gorm:auto_preload", true).
		Find(&resp).
		Error

	if err != nil {
		return nil, err
	}

	// Turn internal models into caller-facing structure
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

	return domains, nil
}

// CreateDomain creates a new domain with the given name. If a domain with the same name exists, an error is
// returned.
func (m *Database) CreateDomain(name string) error {
	return m.DB.
		LogMode(false).
		Create(&models.Domain{Name: name}).
		Error
}

// DeleteDomain deletes an existing domain, specified by its name. All associated forwards of the domain will
// be deleted as well. If no domain with the given name exists, an error is returned.
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

// CreateForward created a new email forward for a domain. The domain is identified by its name. The forward
// is defined by the from email and to to email. If no domain with the given name exists, an error is returned.
// If a forward with the same from address already exists, an error is returned.
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

// UpdateForward updates an existing forward for a domain. The domain is identified by its name. The forward
// is identified by the from email. Only the to email can be updated. If no domain with the given name exists,
// an error is returned. If no forward with the given from address exists, an error is returned.
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

// DeleteForward deletes an existing forward for a domain. The domain is identified by its name. The forward
// is identified by the from email. If no domain with the given name exists, an error is returned.
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
