package models

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"` // UUID primary key

	Name              string `json:"name" validate:"required,max=15"`                            // Required, max length 15
	Description       string `json:"description" validate:"max=3000"`                            // Optional, max length 3000
	Type              string `gorm:"type:company_type;not null" json:"type" validate:"required"` // Required
	Registered        bool   `json:"registered" validate:"required"`                             // Required
	AmountOfEmployees int    `json:"amount_of_employees" validate:"required,min=1"`              // Required, at least 1

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateCompanyRequest struct {
	Name              string `json:"name" validate:"required,max=15"`                             // Required, max length 15
	Description       string `json:"description" validate:"max=3000"`                             // Optional, max length 3000
	Type              string `gorm:"type:company_type;not null"  json:"type" validate:"required"` // Required
	Registered        bool   `json:"registered" validate:"required"`                              // Required
	AmountOfEmployees int    `json:"amount_of_employees" validate:"required,min=1"`               // Required, at least 1
}

type UpdateCompanyRequest struct {
	Description       *string `json:"description,omitempty" validate:"max=3000"`                   // Optional; max length 3000
	Type              *string `gorm:"type:company_type;not null"  json:"type" validate:"required"` // Required
	Registered        *bool   `json:"registered" validate:"required"`                              // Required
	AmountOfEmployees *int    `json:"amount_of_employees" validate:"required,min=1"`               // Required, at least 1
}
