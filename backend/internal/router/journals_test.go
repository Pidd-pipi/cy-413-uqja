package router

import (
	"encoding/json"
	"github.com/blueship581/mindgarden/backend/internal/dto"
	"github.com/blueship581/mindgarden/backend/internal/handler"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/blueship581/mindgarden/backend/internal/repository"
	"github.com/blueship581/mindgarden/backend/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubJournalRepo struct{ saved *model.Journal }

func (s *stubJournalRepo) Create(v *model.Journal) error { s.saved = v; v.ID = 1; return nil }
func (s *stubJournalRepo) List(uint, int) ([]model.Journal, error) {
	return []model.Journal{}, nil
}
func (s *stubJournalRepo) ByID(id, uid uint) (*model.Journal, error) {
	return nil, repository.ErrNotFound
}
func (s *stubJournalRepo) Update(*model.Journal) error { return nil }
func (s *stubJournalRepo) Delete(*model.Journal) error { return nil }

func setupJournalRouter(repo *stubJournalRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewJournalHandler(service.NewJournalService(repo, logger), logger)
	auth := func(c *gin.Context) { c.Set("userID", uint(1)); c.Next() }
	RegisterJournals(r.Group("/v1"), h, auth)
	return r
}

func TestJournalPromptsEndpoint(t *testing.T) {
	r := setupJournalRouter(&stubJournalRepo{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/journals/prompts?mood_level=2", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var resp dto.Response
	if e := json.Unmarshal(w.Body.Bytes(), &resp); e != nil {
		t.Fatalf("decode: %v", e)
	}
	data, _ := json.Marshal(resp.Data)
	var prompts dto.JournalPromptsResponse
	if e := json.Unmarshal(data, &prompts); e != nil {
		t.Fatalf("decode data: %v", e)
	}
	if prompts.Band != "低落" || len(prompts.Prompts) != 3 {
		t.Fatalf("unexpected prompts payload: %+v", prompts)
	}
}

func TestJournalCreateCarriesPrompt(t *testing.T) {
	repo := &stubJournalRepo{}
	r := setupJournalRouter(repo)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/journals", strings.NewReader(`{"title":"t","content":"c","mood_level":4,"prompt":"今天有哪三个瞬间值得被记下来，哪怕很普通？"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if repo.saved == nil || repo.saved.Prompt == "" {
		t.Fatalf("prompt not persisted: %+v", repo.saved)
	}
}
