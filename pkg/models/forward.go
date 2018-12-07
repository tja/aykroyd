package models

// Forward contains the database model for the leaf email forward object, associated to a Domain.
type Forward struct {
	ID       uint   `gorm:"primary_key"`              // Autoincrement ID
	DomainID uint   `gorm:""`                         // Association
	From     string `gorm:"type:varchar(255);unique"` // "From" email address (has to be unique)
	To       string `gorm:"type:text"`                // "To" email address
}
