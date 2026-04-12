package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/shorty/internal/config"
	"github.com/yourusername/shorty/internal/models"
)

type mockRepo struct {
	store map[string]*models.URL
}

func (m *mockRepo) Create(ctx context.Context, u *models.URL) error {
	if m.store == nil {
		m.store = map[string]*models.URL{}
	}

	if _, ok := m.store[u.ShortCode]; ok {
		return errors.New("unique violation")
	}

	u.ID = len(m.store) + 1
	u.CreatedAt = time.Now()
	m.store[u.ShortCode] = u
	return nil
}

func (m *mockRepo) GetByCode(ctx context.Context, code string) (*models.URL, error) {
	if u, ok := m.store[code]; ok {
		copyURL := *u
		return &copyURL, nil
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) IncrementClicks(ctx context.Context, code string) error {
	if u, ok := m.store[code]; ok {
		u.Clicks++
		return nil
	}
	return errors.New("not found")
}

func (m *mockRepo) GetAll(ctx context.Context) ([]models.URL, error) {
	var urls []models.URL
	for _, u := range m.store {
		urls = append(urls, *u)
	}
	return urls, nil
}

func newTestService() *URLService {
	cfg := &config.Config{ShortCodeLength: 6}
	repo := &mockRepo{}
	return NewURLService(repo, cfg)
}

func TestCreateShort_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	u, err := svc.CreateShort(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if u == nil {
		t.Fatal("expected created URL, got nil")
	}

	if u.ShortCode == "" {
		t.Fatal("expected generated short code, got empty string")
	}

	if u.OriginalURL != "https://example.com" {
		t.Fatalf("expected original URL to be saved, got %s", u.OriginalURL)
	}
}

func TestCreateShort_InvalidURL(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.CreateShort(ctx, "not-a-valid-url")
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}

	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestResolve_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	created, err := svc.CreateShort(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	resolved, err := svc.Resolve(ctx, created.ShortCode)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	if resolved.OriginalURL != "https://example.com" {
		t.Fatalf("expected original URL https://example.com, got %s", resolved.OriginalURL)
	}

	if resolved.Clicks != 1 {
		t.Fatalf("expected clicks to be 1, got %d", resolved.Clicks)
	}
}

func TestResolve_NotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.Resolve(ctx, "unknown")
	if err == nil {
		t.Fatal("expected error for unknown short code, got nil")
	}
}

func TestGetStats_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	created, err := svc.CreateShort(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	stats, err := svc.GetStats(ctx, created.ShortCode)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if stats.ShortCode != created.ShortCode {
		t.Fatalf("expected short code %s, got %s", created.ShortCode, stats.ShortCode)
	}
}

func TestGetAllURLs_ReturnsAllItems(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateShort(ctx, "https://example.com")
	_, _ = svc.CreateShort(ctx, "https://google.com")

	urls, err := svc.GetAllURLs(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(urls) != 2 {
		t.Fatalf("expected 2 URLs, got %d", len(urls))
	}
}
