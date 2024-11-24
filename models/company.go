package models

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"` // UUID primary key

	Name              string `gorm:"type:varchar(15);not null;unique" json:"name"` // Name (15 characters, unique, required)
	Description       string `gorm:"type:varchar(3000)" json:"description"`        // Description (optional, up to 3000 characters)
	AmountOfEmployees int    `gorm:"not null" json:"amount_of_employees"`          // Amount of Employees (required)
	Registered        bool   `gorm:"not null" json:"registered"`                   // Registered (boolean, required)
	Type              string `gorm:"type:company_type;not null" json:"type"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateCompanyRequest struct {
	Name              string `gorm:"type:varchar(15);not null;unique" json:"name"` // Name (15 characters, unique, required)
	Description       string `gorm:"type:varchar(3000)" json:"description"`        // Description (optional, up to 3000 characters)
	AmountOfEmployees int    `gorm:"not null" json:"amount_of_employees"`          // Amount of Employees (required)
	Registered        bool   `gorm:"not null" json:"registered"`                   // Registered (boolean, required)
	Type              string `gorm:"type:company_type;not null" json:"type"`       // Type (required, enum)
}

type UpdateCompanyRequest struct {
	Description       *string `json:"description,omitempty"`         // Optional; use pointer to allow distinguishing between "not provided" and "empty"
	AmountOfEmployees *int    `json:"amount_of_employees,omitempty"` // Optional; pointer for the same reason
	Registered        *bool   `json:"registered,omitempty"`          // Optional; pointer for the same reason
	Type              *string `json:"type,omitempty"`                // Optional; pointer to validate presence
}
