package advertisement

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	createPath       = "/advertisement"
	defaultUserAgent = "Chrome/143.0.0.0"
)

type CreateRequest struct {
	Title       string
	Description string
	Price       string
	Quantity    string
	PhotoPath   string
}

func Create(t *testing.T, token string, req CreateRequest, expectedStatusCode int) string {
	t.Helper()

	statusCode, body, err := createRaw(token, req)

	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, statusCode)

	return body
}

func createRaw(token string, req CreateRequest) (int, string, error) {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		return 0, "", fmt.Errorf("API_URL is empty (check .env / .env.override)")
	}

	// multipart body
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	_ = writer.WriteField("title", req.Title)
	_ = writer.WriteField("description", req.Description)

	if req.Price != "" {
		_ = writer.WriteField("price", req.Price)
	}
	if req.Quantity != "" {
		_ = writer.WriteField("quantity", req.Quantity)
	}

	if req.PhotoPath == "" {
		_ = writer.Close()
		return 0, "", fmt.Errorf("PhotoPath is empty")
	}

	absPhotoPath, err := resolvePathFromRepoRoot(req.PhotoPath)
	if err != nil {
		_ = writer.Close()
		return 0, "", err
	}

	if err := addFile(writer, "photos", absPhotoPath); err != nil {
		_ = writer.Close()
		return 0, "", err
	}

	if err := writer.Close(); err != nil {
		return 0, "", err
	}

	url := apiURL + createPath

	httpReq, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return 0, "", err
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", defaultUserAgent)
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return resp.StatusCode, "", readErr
	}

	return resp.StatusCode, string(bodyBytes), nil
}

func addFile(w *multipart.Writer, fieldName, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open photo file: %s: %w", filePath, err)
	}
	defer f.Close()

	part, err := w.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, f)
	return err
}

func resolvePathFromRepoRoot(relPath string) (string, error) {
	relPath = filepath.Clean(relPath)

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	root, err := findGoModRoot(wd)
	if err != nil {
		return "", err
	}

	abs := filepath.Join(root, relPath)

	if _, err := os.Stat(abs); err != nil {
		return "", fmt.Errorf("photo file not found: %s (resolved from repo root: %s): %w", relPath, abs, err)
	}

	return abs, nil
}

func findGoModRoot(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found вверх по дереву от: %s", start)
		}
		dir = parent
	}
}
