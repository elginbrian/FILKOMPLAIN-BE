package config

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type SupabaseStorage struct {
	SupabaseURL string
	SupabaseKey string
	Bucket      string
}

type UploadResponse struct {
	Key string `json:"Key"`
	URL string `json:"signedURL"`
}

func NewSupabaseStorage() *SupabaseStorage {
	return &SupabaseStorage{
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
		Bucket:      os.Getenv("SUPABASE_BUCKET"),
	}
}

func (s *SupabaseStorage) UploadFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d_%s", file.Size, filepath.Base(file.Filename))
	
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.SupabaseURL, s.Bucket, filename)
	
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("apikey", s.SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+s.SupabaseKey)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Cache-Control", "3600")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	urlEndpoint := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.SupabaseURL, s.Bucket, filename)
	
	return urlEndpoint, nil
}