package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/pchawandi/xm-company/database"
	"github.com/pchawandi/xm-company/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New() // Initialize the validator

// CompanyRepository defines the repository interface
type CompanyRepository interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Patch(c *gin.Context)
	Delete(c *gin.Context)
}

type companyRepository struct {
	DB     database.Database
	Logger *zap.Logger
	Ctx    context.Context
}

// NewCompanyRepository initializes a new repository
func NewCompanyRepository(ctx context.Context, db database.Database, logger *zap.Logger) CompanyRepository {
	return &companyRepository{
		DB:     db,
		Logger: logger,
		Ctx:    ctx,
	}
}

// Create handles creating a new company
func (r *companyRepository) Create(c *gin.Context) {
	var input models.CreateCompanyRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		r.respondWithError(c, http.StatusBadRequest, "invalid input", err)
		return
	}

	if err := validate.Struct(&input); err != nil {
		r.respondWithError(c, http.StatusBadRequest, "validation failed", err)
		return
	}

	company := models.Company{
		Name:              input.Name,
		Description:       input.Description,
		Type:              input.Type,
		Registered:        input.Registered,
		AmountOfEmployees: input.AmountOfEmployees,
	}

	if err := r.DB.Create(&company).Error; err != nil {
		r.Logger.Error("failed to create company", zap.Error(err))
		r.respondWithError(c, http.StatusInternalServerError, "could not create company", nil)
		return
	}

	r.respondWithSuccess(c, http.StatusCreated, company)
}

// Get fetches a company by ID
func (r *companyRepository) Get(c *gin.Context) {
	var company models.Company
	id := c.Param("id")

	if err := r.DB.First(&company, "id = ?", id).Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.respondWithError(c, http.StatusNotFound, "company record not found", nil)
		} else {
			r.Logger.Error("error fetching company", zap.String("id", id), zap.Error(err))
			r.respondWithError(c, http.StatusInternalServerError, "could not fetch company", nil)
		}
		return
	}

	r.respondWithSuccess(c, http.StatusOK, company)
}

// Patch updates an existing company
func (r *companyRepository) Patch(c *gin.Context) {
	var company models.Company
	var input models.UpdateCompanyRequest
	id := c.Param("id")

	if err := r.DB.First(&company, "id = ?", id).Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.respondWithError(c, http.StatusNotFound, "company record not found", nil)
		} else {
			r.Logger.Error("error fetching company", zap.String("id", id), zap.Error(err))
			r.respondWithError(c, http.StatusInternalServerError, "could not fetch company", nil)
		}
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		r.respondWithError(c, http.StatusBadRequest, "invalid input", err)
		return
	}

	updates, err := r.buildUpdateMap(input)
	if err != nil {
		r.Logger.Error("failed to update company", zap.Error(err))
		r.respondWithError(c, http.StatusBadRequest, "could not update company", err)
		return
	}

	if err := r.DB.Model(&company).Updates(updates).Error; err != nil {
		r.Logger.Error("failed to update company", zap.Error(err))
		r.respondWithError(c, http.StatusInternalServerError, "could not update company", nil)
		return
	}

	r.respondWithSuccess(c, http.StatusOK, company)
}

// Delete removes a company by ID
func (r *companyRepository) Delete(c *gin.Context) {
	var company models.Company
	id := c.Param("id")

	if err := r.DB.First(&company, "id = ?", id).Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.respondWithError(c, http.StatusNotFound, "company record not found", nil)
		} else {
			r.Logger.Error("error fetching company", zap.String("id", id), zap.Error(err))
			r.respondWithError(c, http.StatusInternalServerError, "could not fetch company", nil)
		}
		return
	}

	if err := r.DB.Delete(&company).Error; err != nil {
		r.Logger.Error("failed to delete company", zap.String("id", id), zap.Error(err))
		r.respondWithError(c, http.StatusInternalServerError, "could not delete company", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper Functions
func (r *companyRepository) respondWithError(c *gin.Context, status int, message string, err error) {
	if err != nil {
		r.Logger.Error(message, zap.Error(err))
	}
	c.JSON(status, gin.H{"error": message})
}

func (r *companyRepository) respondWithSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

func (r *companyRepository) buildUpdateMap(input models.UpdateCompanyRequest) (map[string]interface{}, error) {
	updates := make(map[string]interface{})
	if input.Description != nil {
		if len(*input.Description) > 3000 {
			return nil, errors.New("max characters allowed for description field is 3000")
		}
		updates["description"] = *input.Description
	}
	if input.AmountOfEmployees != nil {
		if *input.AmountOfEmployees <= 0 {
			return nil, errors.New("minimum value for amount of employes is 1")
		}
		updates["amount_of_employees"] = *input.AmountOfEmployees
	}
	if input.Registered != nil {
		updates["registered"] = *input.Registered
	}
	if input.Type != nil {
		// Proper value will be validated by postgres
		updates["type"] = *input.Type
	}
	if len(updates) == 0 {
		return nil, errors.New("no fields requested for update")
	}

	return updates, nil
}
