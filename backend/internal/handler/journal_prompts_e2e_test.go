package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blueship581/mindgarden/backend/internal/dto"
	"github.com/blueship581/mindgarden/backend/internal/middleware"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/blueship581/mindgarden/backend/internal/service"
	"github.com/blueship581/mindgarden/backend/internal/util"
	"github.com/gin-gonic/gin"
)

type fakeJournalRepo struct{ saved []model.Journal }

func (f *fakeJournalRepo) Create(v *model.Journal) error { f.saved = append(f.saved, *v); return nil }
func (f *fakeJournalRepo) List(uint, int) ([]model.Journal, error) {
	return f.saved, nil
}
func (f *fakeJournalRepo) ByID(id, uid uint) (*model.Journal, error) {
	for i := range f.saved {
		if f.saved[i].ID == id {
			return &f.saved[i], nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (f *fakeJournalRepo) Update(v *model.Journal) error { return nil }
func (f *fakeJournalRepo) Delete(v *model.Journal) error { return nil }

func setupJournalRouter(repo *fakeJournalRepo) (*gin.Engine, string) {
	gin.SetMode(gin.TestMode)
	logger := util.NewLogger()
	svc := service.NewJournalService(repo, logger)
	h := NewJournalHandler(svc, logger)
	r := gin.New()
	r.Use(middleware.ErrorHandler(logger))
	auth := middleware.Auth("secret", "issuer", logger)
	p := r.Group("/api/v1/journals", auth)
	p.GET("", h.List)
	p.GET("/prompts", h.Prompts)
	p.POST("", h.Create)
	p.PUT("/:id", h.Update)
	p.DELETE("/:id", h.Delete)
	token, _ := util.CreateToken(7, "user", "secret", "issuer")
	return r, token
}

func TestJournalPromptsEndToEnd(t *testing.T) {
	repo := &fakeJournalRepo{}
	r, token := setupJournalRouter(repo)

	// 1. prompts for mood_level=3
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/journals/prompts?mood_level=3", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("prompts status = %d, body = %s", w.Code, w.Body.String())
	}
	var resp dto.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	arr, _ := resp.Data.([]any)
	if len(arr) != 3 {
		t.Fatalf("want 3 prompts, got %v", resp.Data)
	}

	// 2. invalid level rejected
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/journals/prompts?mood_level=99", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid level status = %d, body = %s", w.Code, w.Body.String())
	}

	// 3. unauthenticated rejected
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/journals/prompts?mood_level=3", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("no-auth status = %d", w.Code)
	}

	// 4. create journal with prompt, then list returns it
	body := `{"title":"t","content":"c","prompt":"写一句安慰","mood_level":3,"weather":"雨","is_private":true}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/journals", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", w.Code, w.Body.String())
	}
	if len(repo.saved) != 1 || repo.saved[0].Prompt != "写一句安慰" {
		t.Fatalf("prompt not persisted: %+v", repo.saved)
	}
	if repo.saved[0].UserID != 7 {
		t.Fatalf("user scoping lost: %+v", repo.saved[0])
	}
}
