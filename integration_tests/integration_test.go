//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Define constants
const baseURL = "http://localhost:8001/api/v1"

// Generate a unique company name to avoid conflicts
func generateCompanyName() string {
	// Generate a shortened company name by truncating the Unix timestamp
	return fmt.Sprintf("C-%d", time.Now().UnixNano()%10000000) // Truncate the nano timestamp to limit chars <15
}

// Helper function to make POST requests and decode the response
func postRequest(url string, body interface{}) (map[string]interface{}, int, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, resp.StatusCode, nil
}

// getJwtToken registers or logs in a user and returns the JWT token
func getJwtToken(username, role string) (string, error) {
	loginURL := fmt.Sprintf("%s/users/login", baseURL)
	registerURL := fmt.Sprintf("%s/users/register", baseURL)

	// Login data
	userData := map[string]interface{}{
		"username": username,
		"password": "securepassword",
		"role":     role,
	}

	// Attempt to log in
	response, statusCode, err := postRequest(loginURL, userData)
	if err != nil && statusCode != http.StatusUnauthorized && statusCode != http.StatusNotFound {
		return "", fmt.Errorf("login failed: %w", err)
	}

	// If login succeeds, return the token
	if statusCode == http.StatusOK {
		if token, ok := response["token"].(string); ok && token != "" {
			return token, nil
		}
		return "", fmt.Errorf("JWT token not found in login response")
	}

	// If login fails due to user not existing, register the user
	_, statusCode, err = postRequest(registerURL, userData)
	if err != nil || (statusCode != http.StatusOK && statusCode != http.StatusCreated) {
		return "", fmt.Errorf("registration failed with status: %d, error: %w", statusCode, err)
	}

	// Retry login after registration
	response, statusCode, err = postRequest(loginURL, userData)
	if err != nil || statusCode != http.StatusOK {
		return "", fmt.Errorf("login failed after registration with status: %d, error: %w", statusCode, err)
	}

	// Extract and return the JWT token
	if token, ok := response["token"].(string); ok && token != "" {
		return token, nil
	}

	return "", fmt.Errorf("JWT token not found after registration")
}

// Create a new company
func createCompany(jwtToken string, CName string, CType string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/companies", baseURL)
	client := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", jwtToken)
	companyData := map[string]interface{}{
		"name":                CName,
		"description":         "A newly registered company",
		"amount_of_employees": 100,
		"registered":          true,
		"type":                CType,
	}

	companyBody, err := json.Marshal(companyData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(companyBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create company with status: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// Get a specific company by ID
func getCompany(companyID, jwtToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/companies/%s", baseURL, companyID)
	client := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", jwtToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get company with status: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// Update a company
func updateCompany(companyID, jwtToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/companies/%s", baseURL, companyID)
	client := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", jwtToken)
	companyData := map[string]interface{}{
		"description":         "Updated company description",
		"amount_of_employees": 150,
		"registered":          false,
		"type":                "NonProfit",
	}
	companyBody, err := json.Marshal(companyData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(companyBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update company with status: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// Delete a company
func deleteCompany(companyID, jwtToken string) error {
	url := fmt.Sprintf("%s/companies/%s", baseURL, companyID)
	client := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", jwtToken)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete company with status: %d", resp.StatusCode)
	}

	return nil
}

func setupTestEnvironment(t *testing.T) string {
	// Get JWT token for admin
	jwtToken, err := getJwtToken("superuser", "admin")
	assert.NoError(t, err, "Failed to authenticate and fetch JWT token")
	return jwtToken
}

func cleanupCompany(t *testing.T, companyID string, jwtToken string) {
	// Attempt to delete the company, ignore errors during cleanup
	err := deleteCompany(companyID, jwtToken)
	if err != nil {
		t.Logf("Cleanup failed for company ID: %s. Error: %s", companyID, err)
	}
}

func TestAPI_Positive_Cases(t *testing.T) {
	// Setup test environment
	adminToken := setupTestEnvironment(t)

	t.Run("CreateCompany", func(t *testing.T) {
		companyResponse, err := createCompany(adminToken, generateCompanyName(), "Corporations")
		assert.NoError(t, err, "Failed to create company")

		companyID := companyResponse["data"].(map[string]interface{})["id"].(string)
		defer cleanupCompany(t, companyID, adminToken)

		// Validate response
		assert.NotEmpty(t, companyID, "Company ID should not be empty")
		assert.Equal(t, "Corporations", companyResponse["data"].(map[string]interface{})["type"], "Unexpected company type")
	})

	t.Run("GetCompany", func(t *testing.T) {
		// Create a company to fetch
		companyResponse, err := createCompany(adminToken, generateCompanyName(), "NonProfit")
		assert.NoError(t, err, "Failed to create company")

		companyID := companyResponse["data"].(map[string]interface{})["id"].(string)
		defer cleanupCompany(t, companyID, adminToken)

		// Fetch company
		fetchedCompany, err := getCompany(companyID, adminToken)
		assert.NoError(t, err, "Failed to fetch company")
		assert.Equal(t, companyID, fetchedCompany["data"].(map[string]interface{})["id"].(string), "Fetched company ID mismatch")
	})

	t.Run("UpdateCompany", func(t *testing.T) {
		// Create a company to update
		companyResponse, err := createCompany(adminToken, generateCompanyName(), "Cooperative")
		assert.NoError(t, err, "Failed to create company")

		companyID := companyResponse["data"].(map[string]interface{})["id"].(string)
		defer cleanupCompany(t, companyID, adminToken)

		// Update company
		updatedCompany, err := updateCompany(companyID, adminToken)
		assert.NoError(t, err, "Failed to update company")
		assert.Equal(t, "Updated company description", updatedCompany["data"].(map[string]interface{})["description"], "Description update mismatch")
		assert.Equal(t, 150.0, updatedCompany["data"].(map[string]interface{})["amount_of_employees"], "Employee count update mismatch")
	})

	t.Run("DeleteCompany", func(t *testing.T) {
		// Create a company to delete
		companyResponse, err := createCompany(adminToken, generateCompanyName(), "Sole Proprietorship")
		assert.NoError(t, err, "Failed to create company")

		companyID := companyResponse["data"].(map[string]interface{})["id"].(string)

		// Delete company
		err = deleteCompany(companyID, adminToken)
		assert.NoError(t, err, "Failed to delete company")
	})
}

func TestAPI_Negative_Cases(t *testing.T) {
	// Setup test environment
	adminToken := setupTestEnvironment(t)

	t.Run("GetNonExistentCompany", func(t *testing.T) {
		// Try to fetch a non-existent company
		_, err := getCompany("nonexistent-id", adminToken)
		assert.Error(t, err, "Expected error when fetching a non-existent company")
	})

	t.Run("CreateCompany_Invalid_Name_length", func(t *testing.T) {
		// Try to fetch a non-existent company
		_, err := createCompany(adminToken, "invlida-name-length-with-more-than-15-chard", "Corporations")
		assert.Error(t, err, "Expected error when fetching a non-existent company")
	})

	t.Run("CreateCompany_Invalid_Company_Type", func(t *testing.T) {
		// Try to create a company with an invalid type
		_, err := createCompany(adminToken, generateCompanyName(), "Corporations-inv")
		assert.Error(t, err, "Expected error when creating a company with an invalid type")
	})
	t.Run("CreateCompany_InsufficientPrivileges", func(t *testing.T) {
		// Get a reader token
		readerToken, err := getJwtToken("readeruser", "reader")
		assert.NoError(t, err, "Failed to authenticate reader")
		// Create a company as reader
		_, err = createCompany(readerToken, generateCompanyName(), "Corporations")
		assert.Error(t, err, "Expected error when creating company with insufficient privileges")
	})

	t.Run("UpdateCompany_InsufficientPrivileges", func(t *testing.T) {
		// Create a company as admin for further tests
		companyResponse, err := createCompany(adminToken, generateCompanyName(), "Corporations")
		assert.NoError(t, err, "Failed to create company as admin")
		companyID := companyResponse["data"].(map[string]interface{})["id"].(string)
		defer cleanupCompany(t, companyID, adminToken)

		// Get a reader token
		readerToken, err := getJwtToken("readeruser", "reader")
		assert.NoError(t, err, "Failed to authenticate reader")

		// Try to update the company with insufficient privileges
		_, err = updateCompany(companyID, readerToken)
		assert.Error(t, err, "Expected error when updating company with insufficient privileges")
	})

	t.Run("DeleteCompany_InsufficientPrivileges", func(t *testing.T) {
		// Create a company as admin for further tests
		companyResponse, err := createCompany(adminToken, generateCompanyName(), "Corporations")
		assert.NoError(t, err, "Failed to create company as admin")
		companyID := companyResponse["data"].(map[string]interface{})["id"].(string)
		defer cleanupCompany(t, companyID, adminToken)

		// Get a reader token
		readerToken, err := getJwtToken("readeruser", "reader")
		assert.NoError(t, err, "Failed to authenticate reader")

		// Try to delete the company with insufficient privileges
		err = deleteCompany(companyID, readerToken)
		assert.Error(t, err, "Expected error when deleting company with insufficient privileges")
	})
}
