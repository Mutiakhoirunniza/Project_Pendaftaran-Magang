package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	payload := `{
        "Username": "rara",
        "Email": "rara@gmail.com",
        "Password": "rara12"
    }`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, RegisterUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Success create new user")
	}
}

func TestRegisterUserError(t *testing.T) {
	e := echo.New()
	payload := `{
        "Username": "rara",
        "Email": "rara@gmail.com",
        "Password": "rara12"
    }`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.Error(t, RegisterUser(c)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Contains(t, rec.Body.String(), "User already exists")
	}
}

func TestLoginUserController(t *testing.T) {
	e := echo.New()
	payload := `{
        "Username": "rara",
        "Email": "rara@gmail.com",
    }`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, LoginUserController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Success login")
	}
}

func TestGetAllUsers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetAllUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetUserByID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, GetUserByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestGetUserByIDError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	if assert.Error(t, GetUserByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "User not found")
	}
}

func TestUpdateUserByID(t *testing.T) {
	e := echo.New()
	payload := `{
        "Username": "raras",
        "Email": "raras@gmail.com",
        "Password": "raras12"
    }`
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, UpdateUserByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Success update user")
	}
}
func TestUpdateUserByIDError(t *testing.T) {
	e := echo.New()
	payload := `{
        "Username": "updateduser",
        "Email": "updateduser@example.com",
        "Password": "newpassword"
    }`
	req := httptest.NewRequest(http.MethodPut, "/users/999", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	if assert.Error(t, UpdateUserByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "User not found")
	}
}

func TestDeleteUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, DeleteUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "success delete user")
	}
}

func TestGetInternshipListings(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/internship-listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetInternshipListings(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestApplyForInternship(t *testing.T) {
	e := echo.New()
	payload := `{
        "SelectedTitle": "Software Engineer ",
        "CV": "path/to/cv.pdf",
        "Nim": "123456",
        "GPA": 3.5,
        "EducationLevel": "S1",
        "UserID": 1,
        "UserEmail": "raras@gmail.com",
        "Username": "raras"
    }`
	req := httptest.NewRequest(http.MethodPost, "/apply-for-internship", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ApplyForInternship(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Pendaftaran magang berhasil disimpan")
	}
}

func TestCancelApplication(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/cancel-application/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, CancelApplication(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Formulir aplikasi berhasil dibatalkan")
	}
}
func TestCancelApplicationError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/cancel-application/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	if assert.Error(t, CancelApplication(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Application not found")
	}
}

func TestGetApplicationStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/application-status/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, GetApplicationStatus(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Status form aplikasi")
	}
}
