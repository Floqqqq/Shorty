package services

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/yourusername/shorty/internal/config"
	"github.com/yourusername/shorty/internal/models"
	// <- добавляем импорт
)

var ErrInvalidURL = errors.New("invalid URL")

// Интерфейс репозитория
type URLRepository interface {
	Create(ctx context.Context, u *models.URL) error
	GetByCode(ctx context.Context, code string) (*models.URL, error)
	IncrementClicks(ctx context.Context, code string) error
	GetAll(ctx context.Context) ([]models.URL, error)}

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

// Генерация короткого кода
func (s *URLService) generateCode() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, s.cfg.ShortCodeLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Создать короткую ссылку
func (s *URLService) CreateShort(ctx context.Context, original string) (*models.URL, error) {
	if original == "" {
		return nil, ErrInvalidURL
	}

	u := &models.URL{
		OriginalURL: original,
	}

	for {
		u.ShortCode = s.generateCode()
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
	// Увеличиваем клики
	_ = s.repo.IncrementClicks(ctx, code)
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