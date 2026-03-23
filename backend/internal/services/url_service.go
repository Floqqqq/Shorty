package services

import (
	"context"
	"errors"

	"net/url"

	"github.com/yourusername/shorty/internal/config"
	"github.com/yourusername/shorty/internal/models"
)

var ErrInvalidURL = errors.New("invalid URL")

// Интерфейс репозитория
type URLRepository interface {
	Create(ctx context.Context, u *models.URL) error
	GetByCode(ctx context.Context, code string) (*models.URL, error)
	IncrementClicks(ctx context.Context, code string) error
	GetAll(ctx context.Context) ([]models.URL, error)
}

// Сервис для работы с короткими ссылками
type URLService struct {
	repo URLRepository
	cfg  *config.Config
}

// Передаем реализацию репозитория при создании сервиса
func NewURLService(repo URLRepository, cfg *config.Config) *URLService {
	return &URLService{
		repo: repo, // теперь любой объект, реализующий URLRepository
		cfg:  cfg,
	}
}

// Проверка на корректность URL
func isValidURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

// Создать короткую ссылку
func (s *URLService) CreateShort(ctx context.Context, original string) (*models.URL, error) {
	if !isValidURL(original) {
		return nil, ErrInvalidURL
	}

	u := &models.URL{
		OriginalURL: original,
	}

	for {
		u.ShortCode = generateShortCode(s.cfg.ShortCodeLength)
		err := s.repo.Create(ctx, u)
		if err != nil {
			if err.Error() == "unique violation" {
				continue // повторяем генерацию
			}
			return nil, err
		}
		break
	}
	return u, nil
}

// Получить исходный URL по коду
func (s *URLService) Resolve(ctx context.Context, code string) (*models.URL, error) {
	u, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if err := s.repo.IncrementClicks(ctx, code); err != nil {
		return nil, err
	}

	u.Clicks++
	return u, nil
}

// Получить статистику по короткой ссылке
func (s *URLService) GetStats(ctx context.Context, code string) (*models.URL, error) {
	return s.repo.GetByCode(ctx, code)
}

// Получить все ссылки

func (s *URLService) GetAllURLs(ctx context.Context) ([]models.URL, error) {
	return s.repo.GetAll(ctx)
}
