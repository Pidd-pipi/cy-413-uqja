package service

import (
	"testing"

	"github.com/blueship581/mindgarden/backend/internal/util"
)

func TestJournalServicePrompts(t *testing.T) {
	s := NewJournalService(nil, util.NewLogger())
	for level := 1; level <= 10; level++ {
		prompts, err := s.Prompts(level)
		if err != nil {
			t.Fatalf("level %d: unexpected error %v", level, err)
		}
		if len(prompts) != 3 {
			t.Errorf("level %d: want 3 prompts, got %d", level, len(prompts))
		}
		for i, p := range prompts {
			if p == "" {
				t.Errorf("level %d: prompt %d is empty", level, i)
			}
		}
	}
	for _, level := range []int{-1, 0, 11} {
		if _, err := s.Prompts(level); err == nil {
			t.Errorf("level %d: want validation error, got nil", level)
		}
	}
}
