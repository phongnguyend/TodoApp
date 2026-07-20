package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/handler"
	"github.com/todo/backend/go/internal/service"
)

// ── mock service ──────────────────────────────────────────────────────────────

type mockFileSvc struct {
	getAllFn            func(page, pageSize int) (dto.PaginatedResponse[dto.FileResponse], error)
	getByIDFn           func(id uint) (dto.FileResponse, error)
	uploadFn            func(input service.UploadInput) (dto.FileResponse, error)
	getDownloadTargetFn func(id uint) (service.DownloadTarget, error)
	deleteFn            func(id uint) error
}

func (m *mockFileSvc) GetAll(page, pageSize int) (dto.PaginatedResponse[dto.FileResponse], error) {
	return m.getAllFn(page, pageSize)
}
func (m *mockFileSvc) GetByID(id uint) (dto.FileResponse, error) {
	return m.getByIDFn(id)
}
func (m *mockFileSvc) Upload(input service.UploadInput, _ ...*uint) (dto.FileResponse, error) {
	return m.uploadFn(input)
}
func (m *mockFileSvc) GetDownloadTarget(id uint) (service.DownloadTarget, error) {
	return m.getDownloadTargetFn(id)
}
func (m *mockFileSvc) Delete(id uint) error {
	return m.deleteFn(id)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func sampleFileResponse() dto.FileResponse {
	ct := "text/plain"
	return dto.FileResponse{
		ID:          1,
		Name:        "notes.txt",
		Extension:   "txt",
		Size:        5,
		ContentType: &ct,
		CreatedAt:   time.Now(),
	}
}

func setupFileRouter(h *handler.FileHandler) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/files")
	api.GET("", h.GetAll)
	api.GET("/:id", h.GetByID)
	api.GET("/:id/download", h.Download)
	api.POST("", h.Upload)
	api.DELETE("/:id", h.Delete)
	return r
}

func doMultipartUpload(r *gin.Engine, fieldName, filename, content string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fieldName, filename)
	_, _ = io.Copy(part, bytes.NewReader([]byte(content)))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ── GetAll ────────────────────────────────────────────────────────────────────

func TestFileHandlerGetAll_Returns200(t *testing.T) {
	svc := &mockFileSvc{
		getAllFn: func(page, pageSize int) (dto.PaginatedResponse[dto.FileResponse], error) {
			assert.Equal(t, 1, page)
			assert.Equal(t, 20, pageSize)
			return dto.PaginatedResponse[dto.FileResponse]{
				Items: []dto.FileResponse{sampleFileResponse()},
				Total: 1, Page: 1, PageSize: 20, TotalPages: 1,
			}, nil
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/files", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.PaginatedResponse[dto.FileResponse]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "notes.txt", resp.Items[0].Name)
}

func TestFileHandlerGetAll_ServiceError_Returns500(t *testing.T) {
	svc := &mockFileSvc{
		getAllFn: func(page, pageSize int) (dto.PaginatedResponse[dto.FileResponse], error) {
			return dto.PaginatedResponse[dto.FileResponse]{}, assert.AnError
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/files", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestFileHandlerGetByID_Returns200(t *testing.T) {
	resp := sampleFileResponse()
	svc := &mockFileSvc{
		getByIDFn: func(id uint) (dto.FileResponse, error) {
			assert.Equal(t, uint(1), id)
			return resp, nil
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/files/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var body dto.FileResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, uint(1), body.ID)
}

func TestFileHandlerGetByID_NotFound_Returns404(t *testing.T) {
	svc := &mockFileSvc{
		getByIDFn: func(id uint) (dto.FileResponse, error) {
			return dto.FileResponse{}, service.ErrFileNotFound
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/files/99", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFileHandlerGetByID_InvalidID_Returns400(t *testing.T) {
	r := setupFileRouter(handler.NewFileHandler(&mockFileSvc{}))

	w := doRequest(r, http.MethodGet, "/api/files/abc", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Download ──────────────────────────────────────────────────────────────────

func TestDownload_NotFound_Returns404(t *testing.T) {
	svc := &mockFileSvc{
		getDownloadTargetFn: func(id uint) (service.DownloadTarget, error) {
			return service.DownloadTarget{}, service.ErrFileNotFound
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/files/99/download", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDownload_ContentMissing_Returns404(t *testing.T) {
	svc := &mockFileSvc{
		getDownloadTargetFn: func(id uint) (service.DownloadTarget, error) {
			return service.DownloadTarget{}, service.ErrFileContentMissing
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/files/1/download", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── Upload ────────────────────────────────────────────────────────────────────

func TestUpload_Returns201(t *testing.T) {
	resp := sampleFileResponse()
	svc := &mockFileSvc{
		uploadFn: func(input service.UploadInput) (dto.FileResponse, error) {
			assert.Equal(t, "notes.txt", input.OriginalName)
			return resp, nil
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doMultipartUpload(r, "file", "notes.txt", "hello")

	assert.Equal(t, http.StatusCreated, w.Code)
	var body dto.FileResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "notes.txt", body.Name)
}

func TestUpload_MissingFile_Returns400(t *testing.T) {
	r := setupFileRouter(handler.NewFileHandler(&mockFileSvc{}))

	req := httptest.NewRequest(http.MethodPost, "/api/files", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpload_TooLarge_Returns413(t *testing.T) {
	svc := &mockFileSvc{
		uploadFn: func(input service.UploadInput) (dto.FileResponse, error) {
			return dto.FileResponse{}, service.ErrFileTooLarge
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doMultipartUpload(r, "file", "big.bin", "hello")

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestUpload_ServiceError_Returns500(t *testing.T) {
	svc := &mockFileSvc{
		uploadFn: func(input service.UploadInput) (dto.FileResponse, error) {
			return dto.FileResponse{}, assert.AnError
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doMultipartUpload(r, "file", "notes.txt", "hello")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestFileHandlerDelete_Returns204(t *testing.T) {
	deleted := false
	svc := &mockFileSvc{
		deleteFn: func(id uint) error {
			assert.Equal(t, uint(1), id)
			deleted = true
			return nil
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodDelete, "/api/files/1", nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.True(t, deleted)
}

func TestFileHandlerDelete_NotFound_Returns404(t *testing.T) {
	svc := &mockFileSvc{
		deleteFn: func(id uint) error {
			return service.ErrFileNotFound
		},
	}
	r := setupFileRouter(handler.NewFileHandler(svc))

	w := doRequest(r, http.MethodDelete, "/api/files/99", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFileHandlerDelete_InvalidID_Returns400(t *testing.T) {
	r := setupFileRouter(handler.NewFileHandler(&mockFileSvc{}))

	w := doRequest(r, http.MethodDelete, "/api/files/abc", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
