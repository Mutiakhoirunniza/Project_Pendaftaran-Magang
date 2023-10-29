package controllers

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

type InternshipsResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type InternshipsUsecase interface {
	ApplyInternships(userInput, openAIKey string) (string, error) //tambahkan lagi
}

type internshipsUsecase struct{}

func NewInternshipsUsecase() InternshipsUsecase {
	return &internshipsUsecase{}
}

func (uc *internshipsUsecase) ApplyInternships(userInput, openAIKey string) (string, error) {
	ctx := context.Background()
	client := openai.NewClient(openAIKey)
	model := openai.GPT3Dot5Turbo
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Halo, perkenalkan saya sistem untuk Tata Cara Pendaftaran & Tips lolos wawancara",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		},
	}

	resp, err := uc.getCompletionFromMessages(ctx, client, messages, model)
	if err != nil {
		return "", err
	}
	answer := resp.Choices[0].Message.Content
	return answer, nil
}

func (uc *internshipsUsecase) getCompletionFromMessages(
	ctx context.Context,
	client *openai.Client,
	messages []openai.ChatCompletionMessage,
	model string,
) (openai.ChatCompletionResponse, error) {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	return resp, err
}

func ApplyInternship(c echo.Context,  internshipsUsecase InternshipsUsecase) error {
	var requestData map[string]interface{}
	err := c.Bind(&requestData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid JSON format"})
	}

	userInput, ok := requestData["message"].(string)
	if !ok || userInput == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid or missing 'message' in the request"})
	}

	answer, err := internshipsUsecase.ApplyInternships(userInput, os.Getenv("chatbot"))
	if err != nil {
		errorMessage := "Gagal menghasilkan panduan pendaftaran magang"
		if strings.Contains(err.Error(), "rate limits exceeded") {
			errorMessage = "Batas tingkat permintaan terlampaui. Silakan coba lagi nanti."
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": errorMessage})
	}

	responseData := InternshipsResponse{
		Status: "success",
		Data:   answer,
	}

	return c.JSON(http.StatusOK, responseData)
}
