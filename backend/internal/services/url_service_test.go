package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/shorty/internal/config"
	"github.com/yourusername/shorty/internal/models"
)

// Мок-репозиторий для тестов
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
		return u, nil
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
// добавляем метод GetAll в mockRepo
func (m *mockRepo) GetAll(ctx context.Context) ([]models.URL, error) {
	var urls []models.URL
	for _, u := range m.store {
		urls = append(urls, *u)
	}
	return urls, nil
}

// Тест создания и разрешения короткой ссылки
func TestCreateAndResolve(t *testing.T) {
	cfg := &config.Config{ShortCodeLength: 6}
	mr := &mockRepo{}
	svc := NewURLService(mr, cfg)

	ctx := context.Background()
	u, err := svc.CreateShort(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if u.ShortCode == "" {
		t.Fatalf("empty code")
	}

	res, err := svc.Resolve(ctx, u.ShortCode)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if res.OriginalURL != "https://example.com" {
		t.Fatalf("unexpected original URL: %v", res.OriginalURL)
	}
	if res.Clicks != 1 {
		t.Fatalf("clicks not incremented: %v", res.Clicks)
	}
}