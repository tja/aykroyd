package models

type Domain struct {
	ID       uint       `gorm:"primary_key"`              // Autoincrement ID
	Name     string     `gorm:"type:varchar(255);unique"` // Domain name
	Forwards []*Forward `gorm:"foreignkey:DomainID"`      // Associated email forwardings
}
