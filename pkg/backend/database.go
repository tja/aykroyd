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
