package configs

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	path = "customer"
)

type MinioConfig struct {
	Client    *minio.Client
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Path      string
	Secure    bool
}

func NewMinioConfig() (*MinioConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	minioSecure := os.Getenv("MINIO_SSL_MODE")

	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCES_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: minioSecure == "true",
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &MinioConfig{
		Client:    minioClient,
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCES_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Bucket:    os.Getenv("MINIO_BUCKET_NAME"),
		Path:      path,
	}, nil
}

func (mc *MinioConfig) UplaodFile(identifier int64, base64Data string) (string, error) {
	ctx := context.Background()

	contentType, err := getContentTypeOfFile(base64Data)

	if err != nil {
		Logger.Error("Error upload proses -> ", err)
		return "", fmt.Errorf("unable to determine content type: %w", err)
	}

	data := RemoveDataURLPrefix(base64Data)

	fileName := generateObjectName(data)

	object := fmt.Sprintf("%s/%d/%s%s", mc.Path, identifier, fileName, getFileExtension(base64Data))

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))

	_, err = mc.Client.PutObject(ctx, mc.Bucket, object, reader, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		Logger.Error("Error upload proses -> ", err)
	}

	url := fmt.Sprintf("%s", object)

	return url, nil
}

func (mc *MinioConfig) GetObjectName(identifier int64, base64Data string) string {

	data := RemoveDataURLPrefix(base64Data)

	fileName := generateObjectName(data)

	object := fmt.Sprintf("%s/%d/%s%s", mc.Path, identifier, fileName, getFileExtension(base64Data))
	return object
}

func (mc *MinioConfig) GetFile(fileURL string) (*minio.Object, error) {
	ctx := context.Background()
	object, err := mc.Client.GetObject(ctx, mc.Bucket, fileURL, minio.GetObjectOptions{})
	if err != nil {
		Logger.Error("errored get object from minio -> ", err)
		return nil, err
	}

	defer object.Close()

	return object, nil
}

func (mc *MinioConfig) GetUrlPublishData(fileUrl string) (string, error) {
	ctx := context.Background()
	presignUrl, err := mc.Client.PresignedGetObject(ctx, mc.Bucket, fileUrl, time.Hour, nil)
	if err != nil {
		Logger.Error("presigned url generate error -> ", err.Error())
	}
	return presignUrl.String(), nil
}

func getContentTypeOfFile(base64Data string) (string, error) {
	if strings.Contains(base64Data, "data:") {
		parts := strings.SplitN(base64Data, ";", 2)
		if len(parts) > 0 && strings.Contains(parts[0], "data:") {
			return parts[0][5:], nil
		}
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(data[:512]), nil
}

func RemoveDataURLPrefix(data string) string {
	if idx := strings.Index(data, ","); idx != -1 {
		return data[idx+1:]
	}
	return data
}

func getFileExtension(contentType string) string {
	switch {
	case strings.Contains(contentType, "image/jpg"):
		return ".jpg"
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/webp"):
		return ".webp"
	case strings.Contains(contentType, "text/plain"):
		return ".txt"
	case strings.Contains(contentType, "application/pdf"):
		return ".pdf"
	case strings.Contains(contentType, "video/mp4"):
		return ".mp4"
	case strings.Contains(contentType, "audio/mpeg"):
		return ".mp3"
	default:
		return ".bin"
	}
}

func generateObjectName(base64Data string) string {
	hash := sha1.New()
	io.WriteString(hash, base64Data)
	return hex.EncodeToString(hash.Sum(nil))
}
