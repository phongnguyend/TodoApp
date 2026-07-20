package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/handler"
	"github.com/todo/backend/go/internal/service"
)

type attachmentSvcMock struct {
	service.TodoItemAttachmentService
	getAll func(uint) ([]dto.TodoItemAttachmentResponse, error)
	create func(uint, uint) (dto.TodoItemAttachmentResponse, error)
	delete func(uint, uint) error
}

func (m *attachmentSvcMock) GetAll(id uint) ([]dto.TodoItemAttachmentResponse, error) {
	return m.getAll(id)
}
func (m *attachmentSvcMock) Create(todo, file uint, _ ...*uint) (dto.TodoItemAttachmentResponse, error) {
	return m.create(todo, file)
}
func (m *attachmentSvcMock) Delete(todo, id uint) error { return m.delete(todo, id) }

func attachmentRouter(s service.TodoItemAttachmentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewTodoItemAttachmentHandler(s)
	r.GET("/api/todo-items/:id/attachments", h.GetAll)
	r.POST("/api/todo-items/:id/attachments", h.Create)
	r.DELETE("/api/todo-items/:id/attachments/:attachmentId", h.Delete)
	return r
}
func attachmentRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestAttachmentHandlerGetAllReturns200(t *testing.T) {
	svc := &attachmentSvcMock{getAll: func(id uint) ([]dto.TodoItemAttachmentResponse, error) {
		assert.Equal(t, uint(10), id)
		return []dto.TodoItemAttachmentResponse{{ID: 1}}, nil
	}}
	w := attachmentRequest(attachmentRouter(svc), http.MethodGet, "/api/todo-items/10/attachments", "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[{"id":1,"todoItemId":0,"fileId":0,"createdAt":"0001-01-01T00:00:00Z","createdByUserId":null,"updatedAt":null,"updatedByUserId":null}]`, w.Body.String())
}
func TestAttachmentHandlerCreateReturns201(t *testing.T) {
	svc := &attachmentSvcMock{create: func(todo, file uint) (dto.TodoItemAttachmentResponse, error) {
		assert.Equal(t, uint(10), todo)
		assert.Equal(t, uint(5), file)
		return dto.TodoItemAttachmentResponse{ID: 2, TodoItemID: todo, FileID: file}, nil
	}}
	w := attachmentRequest(attachmentRouter(svc), http.MethodPost, "/api/todo-items/10/attachments", `{"fileId":5}`)
	assert.Equal(t, http.StatusCreated, w.Code)
}
func TestAttachmentHandlerRejectsInvalidBody(t *testing.T) {
	w := attachmentRequest(attachmentRouter(&attachmentSvcMock{}), http.MethodPost, "/api/todo-items/10/attachments", `{"fileId":0}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestAttachmentHandlerRejectsInvalidPathID(t *testing.T) {
	w := attachmentRequest(attachmentRouter(&attachmentSvcMock{}), http.MethodGet, "/api/todo-items/nope/attachments", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestAttachmentHandlerMapsNotFound(t *testing.T) {
	svc := &attachmentSvcMock{getAll: func(id uint) ([]dto.TodoItemAttachmentResponse, error) { return nil, service.ErrNotFound }}
	w := attachmentRequest(attachmentRouter(svc), http.MethodGet, "/api/todo-items/99/attachments", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestAttachmentHandlerDeleteReturns204(t *testing.T) {
	svc := &attachmentSvcMock{delete: func(todo, id uint) error { assert.Equal(t, uint(10), todo); assert.Equal(t, uint(3), id); return nil }}
	w := attachmentRequest(attachmentRouter(svc), http.MethodDelete, "/api/todo-items/10/attachments/3", "")
	assert.Equal(t, http.StatusNoContent, w.Code)
}
