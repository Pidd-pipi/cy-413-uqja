package service

import (
	"github.com/blueship581/mindgarden/backend/internal/constants"
	"github.com/blueship581/mindgarden/backend/internal/dto"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/blueship581/mindgarden/backend/internal/repository"
	"io"
	"log/slog"
	"testing"
)

type fakeJournalRepo struct {
	items map[uint]*model.Journal
	next  uint
}

func (f *fakeJournalRepo) Create(v *model.Journal) error {
	f.next++
	v.ID = f.next
	f.items[v.ID] = v
	return nil
}
func (f *fakeJournalRepo) List(uid uint, level int) ([]model.Journal, error) {
	out := []model.Journal{}
	for _, v := range f.items {
		if v.UserID == uid && (level <= 0 || v.MoodLevel == level) {
			out = append(out, *v)
		}
	}
	return out, nil
}
func (f *fakeJournalRepo) ByID(id, uid uint) (*model.Journal, error) {
	v, ok := f.items[id]
	if !ok || v.UserID != uid {
		return nil, repository.ErrNotFound
	}
	return v, nil
}
func (f *fakeJournalRepo) Update(v *model.Journal) error { f.items[v.ID] = v; return nil }
func (f *fakeJournalRepo) Delete(v *model.Journal) error { delete(f.items, v.ID); return nil }

func newJournalService() (*JournalService, *fakeJournalRepo) {
	repo := &fakeJournalRepo{items: map[uint]*model.Journal{}}
	return NewJournalService(repo, slog.New(slog.NewTextHandler(io.Discard, nil))), repo
}

func TestJournalServicePrompts(t *testing.T) {
	s, _ := newJournalService()
	cases := []struct {
		name  string
		level int
		band  string
	}{
		{"low mood", 1, "低落"},
		{"low edge", 3, "低落"},
		{"mid", 5, "平稳"},
		{"high edge", 8, "明媚"},
		{"high", 10, "明媚"},
		{"zero falls back to mid", 0, "平稳"},
		{"out of range falls back to mid", 99, "平稳"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := s.Prompts(tt.level)
			if got.Band != tt.band {
				t.Fatalf("band=%s want %s", got.Band, tt.band)
			}
			if len(got.Prompts) != constants.JournalPromptCount {
				t.Fatalf("prompts=%d want %d", len(got.Prompts), constants.JournalPromptCount)
			}
		})
	}
}

func TestJournalServicePromptPersistedPerEntry(t *testing.T) {
	s, _ := newJournalService()
	first, e := s.Create(7, dto.JournalRequest{Title: "第一篇", Content: "正文一", MoodLevel: 2, Prompt: "提示A"})
	if e != nil {
		t.Fatalf("create first: %v", e)
	}
	second, e := s.Create(7, dto.JournalRequest{Title: "第二篇", Content: "正文二", MoodLevel: 9, Prompt: "提示B"})
	if e != nil {
		t.Fatalf("create second: %v", e)
	}
	// 只改第一篇的提示，第二篇必须原样保留
	if _, e = s.Update(7, first.ID, dto.JournalRequest{Title: "第一篇", Content: "改写", MoodLevel: 2, Prompt: "提示A改"}); e != nil {
		t.Fatalf("update first: %v", e)
	}
	kept, e := s.List(7, 0)
	if e != nil || len(kept) != 2 {
		t.Fatalf("list got %v %v", kept, e)
	}
	byID := map[uint]model.Journal{}
	for _, j := range kept {
		byID[j.ID] = j
	}
	if byID[first.ID].Prompt != "提示A改" || byID[first.ID].Content != "改写" {
		t.Fatalf("first entry not updated: %+v", byID[first.ID])
	}
	if byID[second.ID].Prompt != "提示B" || byID[second.ID].Content != "正文二" {
		t.Fatalf("second entry must stay untouched: %+v", byID[second.ID])
	}
	// 只删第一篇，第二篇仍在
	if e = s.Delete(7, first.ID); e != nil {
		t.Fatalf("delete first: %v", e)
	}
	kept, _ = s.List(7, 0)
	if len(kept) != 1 || kept[0].ID != second.ID || kept[0].Prompt != "提示B" {
		t.Fatalf("delete must only affect one entry: %+v", kept)
	}
}
