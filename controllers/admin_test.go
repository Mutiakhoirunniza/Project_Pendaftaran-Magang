package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"miniproject/entity"
	"miniproject/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAdmin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Membuat objek Admin untuk pengujian
	admin := entity.Admin{
		Email: "caca@gmail.com",
	}

	// Marshal objek Admin menjadi JSON
	adminJSON, err := json.Marshal(admin)
	if err != nil {
		t.Fatal(err)
	}

	// Menggunakan ioutil.NopCloser untuk membungkus JSON sebagai ReadCloser
	req.Body = ioutil.NopCloser(bytes.NewReader(adminJSON))

	// Memanggil fungsi RegisterAdmin
	err = RegisterAdmin(c)

	// Memeriksa respons HTTP
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Memeriksa respons JSON
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if assert.NoError(t, err) {
		assert.Equal(t, "Success create new user", response["message"])
	}
}

func TestLoginAdminController(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Membuat data admin palsu untuk digunakan dalam pengujian
	fakeAdmin := entity.Admin{
		Username: "caca",
		Password: "caca12",
		Email:    "caca@gmail.com",
	}

	// Menambahkan data palsu ke body request
	body, _ := json.Marshal(fakeAdmin)
	req.Body = ioutil.NopCloser(bytes.NewReader(body)) 

	// Memanggil fungsi LoginAdminController
	if assert.NoError(t, LoginAdminController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Memeriksa isi respon JSON
		expectedResponse := map[string]interface{}{
			"message": "Success login",
			"user": map[string]interface{}{
				"ID":       1, 
				"Username": "caca",
				"Email":    "caca@gmail.com",
				"Token":    "your_token", 
			},
		}
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestGetAdminByID(t *testing.T) {
	// Membuat instansi Echo dan request palsu
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Menambahkan token palsu ke konteks untuk simulasi
	fakeToken, err := middleware.CreateToken(1, "testadmin") 
	assert.NoError(t, err)                                  
	c.Set("user", fakeToken)

	// Memanggil fungsi GetAdminByID
	if assert.NoError(t, GetAdminByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Memeriksa isi respon JSON
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "Success", response["message"])
			admin, ok := response["admin"].(map[string]interface{})
			if assert.True(t, ok) {
				assert.Equal(t, "admin_username", admin["Username"]) 
				assert.Equal(t, "admin_email", admin["Email"])      
			}
		}
	}
}

func TestGetAdminByIDInvalidID(t *testing.T) {
	// Membuat instansi Echo dan request palsu dengan ID yang tidak valid
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/invalid_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Menambahkan token palsu ke konteks untuk simulasi
	fakeToken, err := middleware.CreateToken(1, "testadmin") 
	assert.NoError(t, err)                                   
	c.Set("user", fakeToken)

	// Memanggil fungsi GetAdminByID dengan ID yang tidak valid
	if assert.NoError(t, GetAdminByID(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Memeriksa isi respon JSON
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ID Admin tidak valid", response)
		}
	}
}

func TestGetAdminByIDNotFound(t *testing.T) {
	// Membuat instansi Echo dan request palsu dengan ID yang tidak ada dalam basis data
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Menambahkan token palsu ke konteks untuk simulasi
	fakeToken, err := middleware.CreateToken(1, "testadmin") 
	assert.NoError(t, err)                                   
	c.Set("user", fakeToken)

	// Memanggil fungsi GetAdminByID dengan ID yang tidak ditemukan
	if assert.NoError(t, GetAdminByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// Memeriksa isi respon JSON
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "Admin tidak ditemukan", response)
		}
	}
}

func TestUpdateAdminController(t *testing.T) {
	e := echo.New()

	// Membuat permintaan HTTP palsu untuk tes
	reqBody := `{"Username":"cacar", "Email":"cacar@gmail.com", "Password":"cacar12"}`
	req := httptest.NewRequest(http.MethodPut, "/admins/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi UpdateAdminController
	if err := UpdateAdminController(c); err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Memeriksa respons HTTP
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, but got: %d", http.StatusOK, rec.Code)
	}

	// Anda juga dapat memeriksa respons JSON yang dihasilkan
	expectedResponse := `{"message": "success update admin"`
	if !strings.Contains(rec.Body.String(), expectedResponse) {
		t.Fatalf("Expected response to contain: %s", expectedResponse)
	}
}

func TestUpdateAdminControllerInvalidID(t *testing.T) {
	e := echo.New()

	// Membuat permintaan HTTP palsu dengan ID yang tidak valid
	reqBody := `{"Username":"newUsername", "Email":"newEmail", "Password":"newPassword"}`
	req := httptest.NewRequest(http.MethodPut, "/admins/invalidID", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil fungsi UpdateAdminController
	if err := UpdateAdminController(c); err == nil {
		t.Fatalf("Expected an error, but got none")
	}

	// Memeriksa status code 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, but got: %d", http.StatusBadRequest, rec.Code)
	}
}


func TestCreateInternshipListing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/create-internship", nil)
	rec := httptest.NewRecorder()
	e.NewContext(req, rec)

	// Test case 1: Successful Creation
	requestBody1 := map[string]interface{}{
		"Quota": 10,
	}
	reqBody1, _ := json.Marshal(requestBody1)
	req1 := httptest.NewRequest(http.MethodPost, "/create-internship", bytes.NewReader(reqBody1))
	req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)
	CreateInternshipListing(c1)
	assert.Equal(t, http.StatusCreated, rec1.Code)

	// Test case 2: Unauthorized
	req2 := httptest.NewRequest(http.MethodPost, "/create-internship", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	CreateInternshipListing(c2)
	assert.Equal(t, http.StatusUnauthorized, rec2.Code)

	// Test case 3: Bad Request
	requestBody3 := map[string]interface{}{
		"Quota": -1,
	}
	reqBody3, _ := json.Marshal(requestBody3)
	req3 := httptest.NewRequest(http.MethodPost, "/create-internship", bytes.NewReader(reqBody3))
	req3.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	CreateInternshipListing(c3)
	assert.Equal(t, http.StatusBadRequest, rec3.Code)
}


func TestUpdateInternshipListingByID(t *testing.T) {
	// Inisialisasi Echo framework
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/admin/InternshipListing/:id", strings.NewReader(`{"field1": "value1", "field2": "value2"}`)) 
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulasikan token dan data Admin yang valid
	token := "your-valid-token"
	username := "caca"
	c.Set("user", map[string]interface{}{"token": token, "username": username})
	err := UpdateInternshipListingByID(c)

	// Lakukan pengujian
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}