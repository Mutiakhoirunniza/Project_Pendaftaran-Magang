package controllers

import (
	"bytes"
	"errors"
	"miniproject/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Mock struct to simulate the database
type MockDB struct {
	UserData              map[string]entity.User
	InternshipListingData map[uint]entity.Internship_Listing
	ApplicationFormID     uint
	ApplicationFormData   map[uint]entity.Internship_ApplicationForm
}

// Implement a mock function for checking if a user with a given email exists
func (db *MockDB) GetUserByEmail(email string) (entity.User, error) {
	user, ok := db.UserData[email]
	if !ok {
		return entity.User{}, errors.New("User not found")
	}
	return user, nil
}

// Implement a mock function for creating a new user
func (db *MockDB) CreateUser(user entity.User) error {
	db.UserData[user.Email] = user
	return nil
}

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		UserData: make(map[string]entity.User),
	}

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)

	// Define a sample user data
	userData := entity.User{
		Email:    "testuser@example.com",
		Password: "password123",
		// Add other user fields as needed
	}

	// Bind the user data to the request
	c.SetPath("/register")
	c.SetParamNames("email", "password")
	c.SetParamValues(userData.Email, userData.Password)

	// Call the RegisterUser function
	if assert.NoError(t, RegisterUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Success create new user","user":{"email":"testuser@example.com","password":"password123"}}`
		assert.Equal(t, expectedResponse, rec.Body.String())

		// Verify that the user is added to the mock database
		_, ok := mockDB.UserData[userData.Email]
		assert.True(t, ok)
	}
}

func TestLoginUserController(t *testing.T) {
	e := echo.New()
	reqBody := []byte(`{"email": "testuser@example.com", "password": "password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		UserData: make(map[string]entity.User),
	}

	// Simulate a user in the mock database
	userData := entity.User{
		Email:    "testuser@example.com",
		Password: "password123",
		// Add other user fields as needed
	}
	mockDB.UserData[userData.Email] = userData

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)

	// Call the LoginUserController function
	if assert.NoError(t, LoginUserController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Success login","user":{"id":0,"name":"","email":"testuser@example.com","token":"mocked-token"}}`
		assert.Equal(t, expectedResponse, rec.Body.String())

		// Verify that a user token is set in the context
		token := c.Get("user")
		assert.NotNil(t, token)
	}
}

func TestGetAllUsers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		UserData: make(map[string]entity.User),
	}

	// Simulate a list of users in the mock database
	user1 := entity.User{
		Email: "user1@example.com",
		// Add other user fields as needed
	}
	user2 := entity.User{
		Email: "user2@example.com",
		// Add other user fields as needed
	}
	mockDB.UserData["user1@example.com"] = user1
	mockDB.UserData["user2@example.com"] = user2

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)

	// Call the GetAllUsers function
	if assert.NoError(t, GetAllUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Success: get all users","users":[{"email":"user1@example.com"},{"email":"user2@example.com"}]}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestGetUserByID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		UserData: make(map[string]entity.User),
	}

	// Simulate a user with ID 1 in the mock database
	user := entity.User{
		Email: "user1@example.com",
		// Add other user fields as needed
	}
	mockDB.UserData["user1@example.com"] = user

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the GetUserByID function
	if assert.NoError(t, GetUserByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Success: get user by ID","user":{"email":"user1@example.com"}}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestUpdateUserByID(t *testing.T) {
	e := echo.New()
	reqBody := []byte(`{
        "username": "UpdatedUserName",
        "email": "updated@example.com",
        "password": "newpassword",
        "gender": "male",
        "phoneNumber": "123456789",
        "universityName": "UpdatedUniversity",
        "universityAddress": "UpdatedAddress",
        "major": "UpdatedMajor"
    }`)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		UserData: make(map[string]entity.User),
	}

	// Simulate an existing user without an "ID" field in the mock database
	user := entity.User{
		Username:          "ExistingUser",
		Email:             "existing@example.com",
		Password:          "password123",
		Gender:            "female",
		PhoneNumber:       "987654321",
		UniversityName:    "ExistingUniversity",
		UniversityAddress: "ExistingAddress",
		Major:             "ExistingMajor",
	}
	mockDB.UserData["existing@example.com"] = user

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the UpdateUserByID function
	if assert.NoError(t, UpdateUserByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Success update user","user":{"username":"UpdatedUserName","email":"updated@example.com","password":"newpassword","gender":"male","phoneNumber":"123456789","universityName":"UpdatedUniversity","universityAddress":"UpdatedAddress","major":"UpdatedMajor"}}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestDeleteUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		UserData: make(map[string]entity.User),
	}

	// Simulate an existing user with ID 1 in the mock database
	user := entity.User{
		Email: "user1@example.com",
		// Add other user fields as needed
	}
	mockDB.UserData["user1@example.com"] = user

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the DeleteUser function
	if assert.NoError(t, DeleteUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Success: delete user"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestGetInternshipListings(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/internship/listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		InternshipListingData: make(map[uint]entity.Internship_Listing),
	}

	// Simulate a list of internship listings in the mock database
	internship1 := entity.Internship_Listing{
		// Set the fields as needed
	}
	internship2 := entity.Internship_Listing{
		// Set the fields as needed
	}
	// Add internship listings to the mock database
	mockDB.InternshipListingData[1] = internship1
	mockDB.InternshipListingData[2] = internship2

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)

	// Call the GetInternshipListings function
	if assert.NoError(t, GetInternshipListings(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		// Adjust the expected response based on your data structure
		expectedResponse := `{"message":"Daftar magang Terbaru","listings":[{"id":1,"field1":"value1","field2":"value2"},{"id":2,"field1":"value1","field2":"value2"}]}`
		assert.Equal(t, expectedResponse, rec.Body.String())
	}
}

func TestApplyForInternship_ValidData(t *testing.T) {
	e := echo.New()
	reqBody := []byte(`{
        "CV": "sample_cv.pdf",
        "Nim": "123456",
        "GPA": 3.5,
        "EducationLevel": "Bachelor",
        "UserID": 1,
        "Status": "",
        "UserEmail": "user@example.com",
        "Username": "user123",
        "SelectedTitle": "Internship Title"
    }`)
	req := httptest.NewRequest(http.MethodPost, "/apply", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		InternshipListingData: make(map[uint]entity.Internship_Listing),
		ApplicationFormID:     1, // Set an application form ID for testing
	}

	// Simulate an existing internship listing with a quota
	listing := entity.Internship_Listing{
		Title: "Internship Title",
		Quota: 10,
	}
	mockDB.InternshipListingData[1] = listing

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)

	// Call the ApplyForInternship function
	if assert.NoError(t, ApplyForInternship(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Pendaftaran magang berhasil disimpan"}`
		assert.Equal(t, expectedResponse, rec.Body.String())
	}
}

func TestApplyForInternship_InvalidData(t *testing.T) {
	e := echo.New()
	reqBody := []byte(`{
        "CV": "sample_cv.pdf",
        "Nim": "",
        "GPA": 0,
        "EducationLevel": "",
        "UserID": 1,
        "Status": "",
        "UserEmail": "",
        "Username": "",
        "SelectedTitle": "Internship Title"
    }`)
	req := httptest.NewRequest(http.MethodPost, "/apply", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		InternshipListingData: make(map[uint]entity.Internship_Listing),
		ApplicationFormID:     1, // Set an application form ID for testing
	}

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)

	// Call the ApplyForInternship function
	if assert.NoError(t, ApplyForInternship(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// Verify the response JSON for invalid data
		// Adjust the expected response based on your application's validation rules
		expectedResponse := `{"message":"Data formulir tidak valid","invalidData":{"nim":"Nim is required","gpa":"GPA must be greater than 0","education_level":"Education level is required","username":"Username is required","user_email":"User email is required"}}`
		assert.Equal(t, expectedResponse, rec.Body.String())
	}
}

func TestApplyForInternship_QuotaFull(t *testing.T) {
	e := echo.New()
	reqBody := []byte(`{
        "CV": "sample_cv.pdf",
        "Nim": "123456",
        "GPA": 3.5,
        "EducationLevel": "Bachelor",
        "UserID": 1,
        "Status": "",
        "UserEmail": "user@example.com",
        "Username": "user123",
        "SelectedTitle": "Internship Title"
    }`)
	req := httptest.NewRequest(http.MethodPost, "/apply", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		InternshipListingData: make(map[uint]entity.Internship_Listing),
		ApplicationFormID:     1, // Set an application form ID for testing
	}

	// Simulate an existing internship listing with no available quota
	listing := entity.Internship_Listing{
		Title: "Internship Title",
		Quota: 0,
	}
	mockDB.InternshipListingData[1] = listing

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)

	// Call the ApplyForInternship function
	if assert.NoError(t, ApplyForInternship(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// Verify the response JSON for a full quota
		expectedResponse := `{"message":"Kuota pendaftaran magang sudah penuh"}`
		assert.Equal(t, expectedResponse, rec.Body.String())
	}
}

func TestCancelApplication(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/cancel-application/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		ApplicationFormData:   make(map[uint]entity.Internship_ApplicationForm),
		InternshipListingData: make(map[uint]entity.Internship_Listing),
	}

	// Simulate an existing application form
	application := entity.Internship_ApplicationForm{
		CV:                  "sample_cv.pdf",
		Nim:                 "123456",
		GPA:                 3.5,
		EducationLevel:      "Bachelor",
		UserID:              1,
		Status:              "Pending",
		UserEmail:           "user@example.com",
		Username:            "user123",
		SelectedTitle:       "Internship Title",
		InternshipListingID: 2,
		IsCanceled:          false,
	}
	mockDB.ApplicationFormData[1] = application

	// Simulate an existing internship listing with an available quota
	listing := entity.Internship_Listing{
		Title: "Internship Title",
		Quota: 5,
	}
	mockDB.InternshipListingData[2] = listing

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the CancelApplication function
	if assert.NoError(t, CancelApplication(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Formulir aplikasi berhasil dibatalkan"}`
		assert.Equal(t, expectedResponse, rec.Body.String())
	}
}

func TestGetApplicationStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/get-application-status/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a mock database
	mockDB := &MockDB{
		ApplicationFormData: make(map[uint]entity.Internship_ApplicationForm),
	}

	// Simulate an existing application form with ID 1
	application := entity.Internship_ApplicationForm{
		Status: "Pending",
	}
	mockDB.ApplicationFormData[1] = application

	// Initialize your Echo context with the mock database
	c.Set("db", mockDB)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the GetApplicationStatus function
	if assert.NoError(t, GetApplicationStatus(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify the response JSON
		expectedResponse := `{"message":"Status form aplikasi","status":"Pending"}`
		assert.Equal(t, expectedResponse, rec.Body.String())
	}
}
