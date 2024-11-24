//go:build !integration
// +build !integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pchawandi/xm-company/database"
	"github.com/pchawandi/xm-company/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func setupUnitTestEnv() (*database.MockDatabase, *gin.Context, CompanyRepository) {
	mockDB := &database.MockDatabase{}
	ctx := context.Background()
	repo := NewCompanyRepository(ctx, mockDB, &zap.Logger{})

	// Create a mock Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("appCtx", repo)

	return mockDB, c, repo
}

func TestCreateCompany_Unit(t *testing.T) {
	mockDB, c, repo := setupUnitTestEnv()

	// Mock input
	input := models.CreateCompanyRequest{
		Name:              "Test Company",
		Description:       "Test Description",
		Type:              "Software",
		Registered:        true,
		AmountOfEmployees: 50,
	}
	company := models.Company{
		ID:                uuid.New(),
		Name:              input.Name,
		Description:       input.Description,
		Type:              input.Type,
		Registered:        input.Registered,
		AmountOfEmployees: input.AmountOfEmployees,
	}

	// Mock request body
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	body, _ := json.Marshal(input)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Mock database behavior
	mockDB.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Company)
		*arg = company
	}).Return(&gorm.DB{})

	// Execute method
	repo.Create(c)

	// Verify response
	assert.Equal(t, http.StatusCreated, c.Writer.Status())
	mockDB.AssertExpectations(t)
}

func TestGetCompany_Unit(t *testing.T) {
	mockDB, c, repo := setupUnitTestEnv()

	u := uuid.New()
	// Mock URL parameter
	c.Params = []gin.Param{{Key: "id", Value: u.String()}}

	// Mock database behavior
	mockDB.On("First", mock.AnythingOfType("*models.Company"), mock.Anything).Run(func(args mock.Arguments) {
		company := args.Get(0).(*models.Company)
		*company = models.Company{
			ID:                u,
			Name:              "Test Company",
			Description:       "Test Description",
			Type:              "Software",
			Registered:        true,
			AmountOfEmployees: 50,
		}
	}).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Execute the Get method
	repo.Get(c)

	// Verify the response
	assert.Equal(t, http.StatusOK, c.Writer.Status())
	mockDB.AssertExpectations(t)
}

// TODO: gorm.DB is a struct, so having issue in
// mocking even wrapping around GormDatabase

// func TestPatchCompany_Unit(t *testing.T) {
// 	mockDB, c, repo := setupUnitTestEnv()

// 	u := uuid.New()
// 	// Mock URL parameter
// 	c.Params = []gin.Param{{Key: "id", Value: u.String()}}

// 	// Mock request body
// 	c.Request = httptest.NewRequest(http.MethodPatch, "/companies/"+u.String(), strings.NewReader(`{
//         "name": "Updated Company Name",
//         "description": "Updated Description"
//     }`))
// 	c.Request.Header.Set("Content-Type", "application/json")

// 	// Mock database behavior
// 	mockDB.On("First", mock.AnythingOfType("*models.Company"), mock.Anything).Run(func(args mock.Arguments) {
// 		company := args.Get(0).(*models.Company)
// 		*company = models.Company{
// 			ID:                u,
// 			Name:              "Old Company Name",
// 			Description:       "Old Description",
// 			Type:              "Old Type",
// 			Registered:        true,
// 			AmountOfEmployees: 10,
// 		}
// 	}).Return(mockDB)
// 	mockDB.On("Model", mock.AnythingOfType("*models.Company")).Return(&gorm.DB{})
// 	mockDB.On("Updates", mock.Anything).Return(mockDB)
// 	mockDB.On("Error").Return(nil)

// 	// Execute the Patch method
// 	repo.Patch(c)

// 	// Verify the response
// 	assert.Equal(t, http.StatusOK, c.Writer.Status())
// 	mockDB.AssertExpectations(t)
// }

func TestDeleteCompany_Unit(t *testing.T) {
	mockDB, c, repo := setupUnitTestEnv()

	// Mock input
	u := uuid.New()
	company := models.Company{ID: u}
	c.Params = gin.Params{{Key: "id", Value: u.String()}}

	// Mock database behavior
	mockDB.On("Error").Return(nil).Once()
	mockDB.On("First", mock.AnythingOfType("*models.Company"), mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Company)
		*arg = company
	}).Return(mockDB).Once()
	mockDB.On("Delete", mock.AnythingOfType("*models.Company"), mock.Anything).Return(&gorm.DB{}).Once()

	// Execute method
	repo.Delete(c)

	// Verify response
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())
	mockDB.AssertExpectations(t)
}
