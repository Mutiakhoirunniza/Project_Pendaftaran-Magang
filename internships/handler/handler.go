package handler

import (
	"fmt"
	"miniproject/internships/usecase"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type InternshipHandler struct {
	InternshipUsecase usecase.InternshipApplicationUsecase 
}

type InternshipResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func NewInternshipHandler(usecase usecase.InternshipApplicationUsecase) *InternshipHandler {
	return &InternshipHandler{
		InternshipUsecase: usecase,
	}
}

func (h *InternshipHandler) SubmitApplication(c echo.Context) error {
	var requestData map[string]interface{}
	err := c.Bind(&requestData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid JSON format")
	}

	name, okName := requestData["name"].(string)
	email, okEmail := requestData["email"].(string)
	userInput, okUserInput := requestData["userInput"].(string) // Input tambahan untuk pendaftaran magang

	if !okName || !okEmail || !okUserInput {
		return c.JSON(http.StatusBadRequest, "Invalid request format")
	}

	// Ganti pesan yang dibuat oleh user ke sesuatu yang relevan untuk pendaftaran magang
	applicationDetails := fmt.Sprintf("Pendaftaran magang oleh %s (email: %s) - Pendaftaran di%s Hanya memerlukan identitas diri kamu dan Kuota, Karena nantinya kamu akan mengisi formulir yang telah disediakan oleh pt kami.", name, email, userInput)

	answer, err := h.InternshipUsecase.SubmitApplication(applicationDetails, name, email, os.Getenv("chatbot"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error dalam pengajuan pendaftaran magang")
	}

	responseData := InternshipResponse{
		Status: "success",
		Data:   answer,
	}

	return c.JSON(http.StatusOK, responseData)
}
