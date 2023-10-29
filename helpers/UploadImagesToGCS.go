package helpers

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// func UploadImageToGCS(ctx context.Context, imageData []byte, imageName string) (string, error) {
// 	credentialsFile := os.Getenv("GOOGLE_CLOUD_CREDENTIALS_PATH")

// 	// Periksa apakah variabel lingkungan GOOGLE_CLOUD_CREDENTIALS_PATH telah diatur
// 	if credentialsFile == "" {
// 		return "", fmt.Errorf("Variabel lingkungan GOOGLE_CLOUD_CREDENTIALS_PATH tidak diatur")
// 	}

// 	// Setel variabel lingkungan GOOGLE_APPLICATION_CREDENTIALS
// 	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialsFile)
// 	if err != nil {
// 		return "", fmt.Errorf("Gagal mengatur variabel lingkungan: %v", err)
// 	}

// 	client, err := storage.NewClient(ctx)
// 	if err != nil {
// 		return "", fmt.Errorf("Gagal membuat klien Storage: %v", err)
// 	}
// 	defer client.Close()

// 	bucketName := "krisnadwipayana"
// 	object := client.Bucket(bucketName).Object(imageName)
// 	wc := object.NewWriter(ctx)
// 	wc.ContentType = "application/octet-stream"

// 	if _, err := io.Copy(wc, bytes.NewReader(imageData)); err != nil {
// 		wc.Close()
// 		return "", fmt.Errorf("Gagal menyalin data ke GCS: %v", err)
// 	}

// 	if err := wc.Close(); err != nil {
// 		return "", fmt.Errorf("Gagal menutup penulis objek: %v", err)
// 	}

// 	cvPath := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, imageName)
// 	return cvPath, nil
// }

var (
	DEFAULT_GCS_LINK string = "https://storage.googleapis.com/krisnadwipayana/"
)

type ClientUploader struct {
	storageClient *storage.Client
	bucketName    string
	cvPath        string
}

var clientUploader *ClientUploader

func GetStorageClient() *ClientUploader {
	if clientUploader == nil {
		client, err := storage.NewClient(context.Background(), option.WithoutAuthentication())
		if err != nil {
			// fmt.Println("Failed to create client: %v", err)
		}

		clientUploader = &ClientUploader{
			storageClient: client,
			bucketName:    os.Getenv("bucketname"),
			cvPath:        "cv/",
		}

		return clientUploader
	}
	return clientUploader
}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(file multipart.File, objectName string) (fileLocation string, err error) {
	ctx := context.Background()

	// Upload an object with storage.Writer.
	wc := c.storageClient.Bucket(c.bucketName).Object(c.cvPath + objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	return DEFAULT_GCS_LINK + wc.Name, nil
}

func (c *ClientUploader) DeleteFile(objectName string) error {
	ctx := context.Background()

	wc := c.storageClient.Bucket(c.bucketName).Object(strings.Replace(objectName, DEFAULT_GCS_LINK, "", 1))
	if err := wc.Delete(ctx); err != nil {
		return err
	}

	return nil
}
