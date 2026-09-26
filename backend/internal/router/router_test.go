package router

import (
	"testing"

	"github.com/blueship581/mindgarden/backend/internal/config"
	"github.com/blueship581/mindgarden/backend/internal/handler"
	"github.com/blueship581/mindgarden/backend/internal/util"
	"github.com/gin-gonic/gin"
)

func TestNewRegistersJournalPromptsRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := util.NewLogger()
	h := Handlers{
		User:       handler.NewUserHandler(nil, nil, logger, "secret", "issuer"),
		Mood:       handler.NewMoodHandler(nil, logger),
		Assessment: handler.NewAssessmentHandler(nil, logger),
		Journal:    handler.NewJournalHandler(nil, logger),
	}
	engine := New(config.Config{CORSOrigin: "http://localhost:18413", JWTSecret: "secret", JWTIssuer: "issuer"}, h, logger)
	found := map[string]bool{}
	for _, r := range engine.Routes() {
		found[r.Method+" "+r.Path] = true
	}
	for _, want := range []string{"GET /v1/journals/prompts", "GET /api/v1/journals/prompts"} {
		if !found[want] {
			t.Errorf("route %s not registered", want)
		}
	}
}
