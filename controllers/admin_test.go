package controllers

import (
	"encoding/json"
	"miniproject/entity"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLoginAdminController(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a mock database and set it in the Echo context
	mockDB := &MockDB{
		UserData: make(map[string]entity.User),
	}
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", mockDB)
			return next(c)
		}
	})

	// Initialize the database or database connection here
	// You need to ensure that the database or connection is properly initialized.

	// Create a request and response recorder
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Initialize the mock database with user data
	userData := entity.User{
		Email:    "admin@example.com",
		Password: "your_password",
	}
	mockDB.UserData[userData.Email] = userData

	// Call the LoginAdminController function
	if assert.NoError(t, LoginAdminController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Add your assertion checks for the response here
		// For example, check the response body for expected content.
	}
}

func TestGetAdminByID(t *testing.T) {
	// Membuat objek Echo untuk pengujian
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi GetAdminByID dengan ID yang valid
	if assert.NoError(t, GetAdminByID(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusOK, rec.Code)

		// Anda juga dapat memeriksa konten JSON dalam respons
		// Misalnya, jika Anda ingin memastikan bahwa respons berisi "Success" dalam pesan.
		expectedJSON := `{"message":"Success","admin":{}}`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}

	// Memanggil fungsi GetAdminByID dengan ID yang tidak valid
	req = httptest.NewRequest(http.MethodGet, "/admin/invalid_id", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, GetAdminByID(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Anda dapat memeriksa pesan dalam respons.
		expectedMessage := "ID Admin tidak valid"
		assert.Equal(t, expectedMessage, rec.Body.String())
	}

	// Memanggil fungsi GetAdminByID dengan ID yang tidak ditemukan
	// Pastikan ID yang Anda berikan benar-benar tidak ada dalam database.
	req = httptest.NewRequest(http.MethodGet, "/admin/999", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, GetAdminByID(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// Anda dapat memeriksa pesan dalam respons.
		expectedMessage := "Admin tidak ditemukan"
		assert.Equal(t, expectedMessage, rec.Body.String())
	}
}
func TestUpdateAdminControllerInvalidID(t *testing.T) {
	// Membuat objek Echo untuk pengujian
	e := echo.New()

	// Membuat body request JSON yang valid untuk pengujian
	requestBody := `{
		"username": "new_username",
		"email": "new_email@example.com",
		"password": "new_password"
	}`

	// Memanggil fungsi UpdateAdminController dengan ID yang tidak valid
	req := httptest.NewRequest(http.MethodPut, "/admin/invalid_id", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, UpdateAdminController(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse respons JSON
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			// Memeriksa pesan kesalahan dalam respons JSON
			assert.Equal(t, "Invalid request body", response["error"])
		}
	}
}
func TestCreateInternshipListing(t *testing.T) {
	// Membuat objek Echo untuk pengujian
	e := echo.New()
	requestBody := `{
		"title": "Magang 2023",
		"quota": 10,
		"description": "Deskripsi lowongan magang"
	}`

	req := httptest.NewRequest(http.MethodPost, "/create-listing", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi CreateInternshipListing
	if assert.NoError(t, CreateInternshipListing(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Anda juga dapat memeriksa konten JSON dalam respons
		// Misalnya, jika Anda ingin memastikan bahwa respons berisi "Lowongan magang berhasil dibuat" dalam pesan.
		expectedJSON := `{"message":"Lowongan magang berhasil dibuat","listing":{}}`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}

	// Memanggil fungsi CreateInternshipListing dengan kuota yang tidak valid
	requestBody = `{
		"title": "Magang 2023",
		"quota": -1, // Kuota tidak valid
		"description": "Deskripsi lowongan magang"
	}`

	req = httptest.NewRequest(http.MethodPost, "/create-listing", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, CreateInternshipListing(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Anda dapat memeriksa pesan dalam respons.
		expectedMessage := "Kuota harus lebih dari 0"
		assert.Equal(t, expectedMessage, rec.Body.String())
	}
}
func TestUpdateInternshipListingByIDInvalidID(t *testing.T) {
	// Membuat objek Echo untuk pengujian
	e := echo.New()

	// Membuat body request JSON yang valid untuk pengujian
	requestBody := `{
        "title": "Magang 2023 (Diperbarui)",
        "quota": 15,
        "description": "Deskripsi lowongan magang (Diperbarui)"
    }`

	// Memanggil fungsi UpdateInternshipListingByID dengan ID yang tidak valid
	req := httptest.NewRequest(http.MethodPut, "/update-listing/invalid_id", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, UpdateInternshipListingByID(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse respons JSON
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			// Memeriksa pesan kesalahan dalam respons JSON
			assert.Equal(t, "Failed to parse request body", response["error"])
		}
	}
}

func TestDeleteInternshipListingByID(t *testing.T) {
	// Membuat objek Echo untuk pengujian
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/delete-listing/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi DeleteInternshipListingByID dengan ID yang valid
	if assert.NoError(t, DeleteInternshipListingByID(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedJSON := `{"message":"Pendaftaran magang berhasil dihapus"}`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}

	// Memanggil fungsi DeleteInternshipListingByID dengan ID yang tidak valid
	req = httptest.NewRequest(http.MethodDelete, "/delete-listing/invalid_id", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, DeleteInternshipListingByID(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Anda dapat memeriksa pesan dalam respons.
		expectedMessage := "Gagal menghapus pendaftaran magang"
		assert.Equal(t, expectedMessage, rec.Body.String())
	}
}
func TestSelectCandidatesByGPAID(t *testing.T) {
	// Membuat objek Echo untuk pengujian
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/select-candidates/1", nil) // Ganti URL dengan yang sesuai
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi SelectCandidatesByGPAID
	if assert.NoError(t, SelectCandidatesByGPAID(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusOK, rec.Code)

		// Anda juga dapat memeriksa konten JSON dalam respons
		// Misalnya, jika Anda ingin memastikan bahwa respons berisi "Candidates selected based on GPA range" dalam pesan.
		expectedJSON := `{"message":"Candidates selected based on GPA range"}`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}
func TestSendEmailHandler(t *testing.T) {
	// Membuat objek Echo untuk pengujian
	e := echo.New()

	// Membuat form request dengan parameter yang sesuai
	form := url.Values{}
	form.Set("userEmail", "user@example.com")
	form.Set("username", "TestUser")
	form.Set("status", "accepted") // Pastikan status sesuai dengan yang diharapkan

	req := httptest.NewRequest(http.MethodPost, "/send-email", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi SendEmailHandler
	if assert.NoError(t, SendEmailHandler(c)) {
		// Memeriksa respons HTTP
		assert.Equal(t, http.StatusOK, rec.Code)

		// Anda juga dapat memeriksa konten JSON dalam respons
		// Misalnya, jika Anda ingin memastikan bahwa respons berisi "Email sent successfully" dalam pesan.
		expectedJSON := `{"message":"Email sent successfully"}`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}
func TestViewAllCandidates(t *testing.T) {
	// Membuat objek Echo untuk menguji handler
	e := echo.New()

	// Membuat HTTP request palsu dengan metode GET ke endpoint "/candidates"
	req := httptest.NewRequest(http.MethodGet, "/candidates", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi ViewAllCandidates dengan konteks palsu
	if assert.NoError(t, ViewAllCandidates(c)) {
		// Memeriksa kode status HTTP yang diharapkan
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
