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

func init() {
	gin.SetMode(gin.TestMode)
}

// ── mock service ──────────────────────────────────────────────────────────────

type mockSvc struct {
	getAllFn        func(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error)
	getIncompleteFn func(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error)
	getByIDFn       func(id uint) (dto.TodoItemResponse, error)
	createFn        func(req dto.CreateTodoItemRequest) (dto.TodoItemResponse, error)
	updateFn        func(id uint, req dto.UpdateTodoItemRequest) (dto.TodoItemResponse, error)
	markCompleteFn  func(id uint) (dto.TodoItemResponse, error)
	deleteFn        func(id uint) error
	importCSVFn     func(r io.Reader) (dto.ImportResult, error)
	exportCSVFn     func() (string, error)
	importExcelFn   func(r io.Reader) (dto.ImportResult, error)
	exportExcelFn   func() ([]byte, error)
}

func (m *mockSvc) GetAll(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error) {
	return m.getAllFn(page, pageSize)
}
func (m *mockSvc) GetIncomplete(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error) {
	return m.getIncompleteFn(page, pageSize)
}
func (m *mockSvc) GetByID(id uint) (dto.TodoItemResponse, error) {
	return m.getByIDFn(id)
}
func (m *mockSvc) Create(req dto.CreateTodoItemRequest, _ ...*uint) (dto.TodoItemResponse, error) {
	return m.createFn(req)
}
func (m *mockSvc) Update(id uint, req dto.UpdateTodoItemRequest, _ ...*uint) (dto.TodoItemResponse, error) {
	return m.updateFn(id, req)
}
func (m *mockSvc) MarkComplete(id uint, _ ...*uint) (dto.TodoItemResponse, error) {
	return m.markCompleteFn(id)
}
func (m *mockSvc) Delete(id uint) error {
	return m.deleteFn(id)
}
func (m *mockSvc) ImportCSV(r io.Reader, _ ...*uint) (dto.ImportResult, error) {
	return m.importCSVFn(r)
}
func (m *mockSvc) ExportCSV() (string, error) {
	return m.exportCSVFn()
}
func (m *mockSvc) ImportExcel(r io.Reader, _ ...*uint) (dto.ImportResult, error) {
	return m.importExcelFn(r)
}
func (m *mockSvc) ExportExcel() ([]byte, error) {
	return m.exportExcelFn()
}

// ── helpers ───────────────────────────────────────────────────────────────────

func sampleResponse() dto.TodoItemResponse {
	desc := "a description"
	return dto.TodoItemResponse{
		ID:          1,
		Title:       "Buy milk",
		Description: &desc,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}
}

func setupRouter(h *handler.TodoItemHandler) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/todo-items")
	api.GET("", h.GetAll)
	api.GET("/incomplete", h.GetIncomplete)
	api.GET("/:id", h.GetByID)
	api.POST("", h.Create)
	api.PUT("/:id", h.Update)
	api.PATCH("/:id/complete", h.MarkComplete)
	api.DELETE("/:id", h.Delete)
	api.POST("/import/csv", h.ImportCSV)
	api.GET("/export/csv", h.ExportCSV)
	api.POST("/import/excel", h.ImportExcel)
	api.GET("/export/excel", h.ExportExcel)
	return r
}

func doRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	r.ServeHTTP(w, req)
	return w
}

// ── GetAll ────────────────────────────────────────────────────────────────────

func TestGetAll_Returns200(t *testing.T) {
	svc := &mockSvc{
		getAllFn: func(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error) {
			assert.Equal(t, 1, page)
			assert.Equal(t, 20, pageSize)
			return dto.PaginatedResponse[dto.TodoItemResponse]{
				Items: []dto.TodoItemResponse{sampleResponse()},
				Total: 1, Page: 1, PageSize: 20, TotalPages: 1,
			}, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.PaginatedResponse[dto.TodoItemResponse]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "Buy milk", resp.Items[0].Title)
}

func TestGetAll_ServiceError_Returns500(t *testing.T) {
	svc := &mockSvc{
		getAllFn: func(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error) {
			return dto.PaginatedResponse[dto.TodoItemResponse]{}, assert.AnError
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── GetIncomplete ─────────────────────────────────────────────────────────────

func TestGetIncomplete_Returns200(t *testing.T) {
	svc := &mockSvc{
		getIncompleteFn: func(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error) {
			return dto.PaginatedResponse[dto.TodoItemResponse]{
				Items: []dto.TodoItemResponse{sampleResponse()},
				Total: 1, Page: 1, PageSize: 20, TotalPages: 1,
			}, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items/incomplete", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.PaginatedResponse[dto.TodoItemResponse]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Items, 1)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestGetByID_Returns200(t *testing.T) {
	resp := sampleResponse()
	svc := &mockSvc{
		getByIDFn: func(id uint) (dto.TodoItemResponse, error) {
			assert.Equal(t, uint(1), id)
			return resp, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var body dto.TodoItemResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, uint(1), body.ID)
}

func TestGetByID_NotFound_Returns404(t *testing.T) {
	svc := &mockSvc{
		getByIDFn: func(id uint) (dto.TodoItemResponse, error) {
			return dto.TodoItemResponse{}, service.ErrNotFound
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items/99", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetByID_InvalidID_Returns400(t *testing.T) {
	r := setupRouter(handler.NewTodoItemHandler(&mockSvc{}))

	w := doRequest(r, http.MethodGet, "/api/todo-items/abc", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCreate_Returns201(t *testing.T) {
	resp := sampleResponse()
	svc := &mockSvc{
		createFn: func(req dto.CreateTodoItemRequest) (dto.TodoItemResponse, error) {
			assert.Equal(t, "Buy milk", req.Title)
			return resp, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodPost, "/api/todo-items", map[string]string{"title": "Buy milk"})

	assert.Equal(t, http.StatusCreated, w.Code)
	var body dto.TodoItemResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Buy milk", body.Title)
}

func TestCreate_MissingTitle_Returns400(t *testing.T) {
	r := setupRouter(handler.NewTodoItemHandler(&mockSvc{}))

	w := doRequest(r, http.MethodPost, "/api/todo-items", map[string]string{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_ServiceError_Returns500(t *testing.T) {
	svc := &mockSvc{
		createFn: func(req dto.CreateTodoItemRequest) (dto.TodoItemResponse, error) {
			return dto.TodoItemResponse{}, assert.AnError
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodPost, "/api/todo-items", map[string]string{"title": "Test"})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdate_Returns200(t *testing.T) {
	resp := sampleResponse()
	resp.Title = "Updated"
	svc := &mockSvc{
		updateFn: func(id uint, req dto.UpdateTodoItemRequest) (dto.TodoItemResponse, error) {
			assert.Equal(t, uint(1), id)
			return resp, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodPut, "/api/todo-items/1", map[string]string{"title": "Updated"})

	assert.Equal(t, http.StatusOK, w.Code)
	var body dto.TodoItemResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Updated", body.Title)
}

func TestUpdate_NotFound_Returns404(t *testing.T) {
	svc := &mockSvc{
		updateFn: func(id uint, req dto.UpdateTodoItemRequest) (dto.TodoItemResponse, error) {
			return dto.TodoItemResponse{}, service.ErrNotFound
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodPut, "/api/todo-items/99", map[string]string{"title": "X"})

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdate_InvalidID_Returns400(t *testing.T) {
	r := setupRouter(handler.NewTodoItemHandler(&mockSvc{}))

	w := doRequest(r, http.MethodPut, "/api/todo-items/abc", map[string]string{"title": "X"})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── MarkComplete ──────────────────────────────────────────────────────────────

func TestMarkComplete_Returns200(t *testing.T) {
	resp := sampleResponse()
	resp.IsCompleted = true
	svc := &mockSvc{
		markCompleteFn: func(id uint) (dto.TodoItemResponse, error) {
			assert.Equal(t, uint(1), id)
			return resp, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodPatch, "/api/todo-items/1/complete", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var body dto.TodoItemResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.True(t, body.IsCompleted)
}

func TestMarkComplete_NotFound_Returns404(t *testing.T) {
	svc := &mockSvc{
		markCompleteFn: func(id uint) (dto.TodoItemResponse, error) {
			return dto.TodoItemResponse{}, service.ErrNotFound
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodPatch, "/api/todo-items/99/complete", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMarkComplete_InvalidID_Returns400(t *testing.T) {
	r := setupRouter(handler.NewTodoItemHandler(&mockSvc{}))

	w := doRequest(r, http.MethodPatch, "/api/todo-items/abc/complete", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDelete_Returns204(t *testing.T) {
	deleted := false
	svc := &mockSvc{
		deleteFn: func(id uint) error {
			assert.Equal(t, uint(1), id)
			deleted = true
			return nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodDelete, "/api/todo-items/1", nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.True(t, deleted)
}

func TestDelete_NotFound_Returns404(t *testing.T) {
	svc := &mockSvc{
		deleteFn: func(id uint) error {
			return service.ErrNotFound
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodDelete, "/api/todo-items/99", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDelete_InvalidID_Returns400(t *testing.T) {
	r := setupRouter(handler.NewTodoItemHandler(&mockSvc{}))

	w := doRequest(r, http.MethodDelete, "/api/todo-items/abc", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── ImportCSV ─────────────────────────────────────────────────────────────────

func multipartCSVRequest(csvContent string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "todo_items.csv")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write([]byte(csvContent)); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req := httptest.NewRequest(http.MethodPost, "/api/todo-items/import/csv", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestImportCSV_Returns200(t *testing.T) {
	svc := &mockSvc{
		importCSVFn: func(r io.Reader) (dto.ImportResult, error) {
			return dto.ImportResult{Imported: 2, Failed: 0, Errors: []dto.ImportRowError{}}, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	req, err := multipartCSVRequest("title,description,is_completed\nBuy milk,,false\n")
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body dto.ImportResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 2, body.Imported)
	assert.Equal(t, 0, body.Failed)
}

func TestImportCSV_MissingFile_Returns400(t *testing.T) {
	r := setupRouter(handler.NewTodoItemHandler(&mockSvc{}))

	w := doRequest(r, http.MethodPost, "/api/todo-items/import/csv", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestImportCSV_ServiceError_Returns500(t *testing.T) {
	svc := &mockSvc{
		importCSVFn: func(r io.Reader) (dto.ImportResult, error) {
			return dto.ImportResult{}, assert.AnError
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	req, err := multipartCSVRequest("title,description,is_completed\nBuy milk,,false\n")
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── ExportCSV ─────────────────────────────────────────────────────────────────

func TestExportCSV_Returns200(t *testing.T) {
	svc := &mockSvc{
		exportCSVFn: func() (string, error) {
			return "id,title,description,is_completed,created_at,updated_at\n", nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items/export/csv", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	assert.Contains(t, w.Header().Get("Content-Disposition"), "todo_items.csv")
	assert.Contains(t, w.Body.String(), "id,title,description,is_completed,created_at,updated_at")
}

func TestExportCSV_ServiceError_Returns500(t *testing.T) {
	svc := &mockSvc{
		exportCSVFn: func() (string, error) {
			return "", assert.AnError
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items/export/csv", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── ImportExcel ───────────────────────────────────────────────────────────────

func multipartExcelRequest(content []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "todo_items.xlsx")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(content); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req := httptest.NewRequest(http.MethodPost, "/api/todo-items/import/excel", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestImportExcel_Returns200(t *testing.T) {
	svc := &mockSvc{
		importExcelFn: func(r io.Reader) (dto.ImportResult, error) {
			return dto.ImportResult{Imported: 2, Failed: 0, Errors: []dto.ImportRowError{}}, nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	req, err := multipartExcelRequest([]byte("fake xlsx bytes"))
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var respBody dto.ImportResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &respBody))
	assert.Equal(t, 2, respBody.Imported)
	assert.Equal(t, 0, respBody.Failed)
}

func TestImportExcel_MissingFile_Returns400(t *testing.T) {
	r := setupRouter(handler.NewTodoItemHandler(&mockSvc{}))

	w := doRequest(r, http.MethodPost, "/api/todo-items/import/excel", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestImportExcel_ServiceError_Returns500(t *testing.T) {
	svc := &mockSvc{
		importExcelFn: func(r io.Reader) (dto.ImportResult, error) {
			return dto.ImportResult{}, assert.AnError
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	req, err := multipartExcelRequest([]byte("fake xlsx bytes"))
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── ExportExcel ───────────────────────────────────────────────────────────────

func TestExportExcel_Returns200(t *testing.T) {
	svc := &mockSvc{
		exportExcelFn: func() ([]byte, error) {
			return []byte("fake xlsx bytes"), nil
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items/export/excel", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	assert.Contains(t, w.Header().Get("Content-Disposition"), "todo_items.xlsx")
	assert.Equal(t, "fake xlsx bytes", w.Body.String())
}

func TestExportExcel_ServiceError_Returns500(t *testing.T) {
	svc := &mockSvc{
		exportExcelFn: func() ([]byte, error) {
			return nil, assert.AnError
		},
	}
	r := setupRouter(handler.NewTodoItemHandler(svc))

	w := doRequest(r, http.MethodGet, "/api/todo-items/export/excel", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
