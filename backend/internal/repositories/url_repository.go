package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/shorty/internal/models"
)

// Экспортируемая структура
type PgUrlRepository struct {
	pool *pgxpool.Pool
}

// Интерфейс для сервиса
type URLRepository interface {
	Create(ctx context.Context, u *models.URL) error
	GetByCode(ctx context.Context, code string) (*models.URL, error)
	IncrementClicks(ctx context.Context, code string) error
	GetAll() ([]models.URL, error)
}

// Конструктор
func NewPgUrlRepository(pool *pgxpool.Pool) *PgUrlRepository {
	return &PgUrlRepository{pool: pool}
}

// Методы
func (r *PgUrlRepository) Create(ctx context.Context, u *models.URL) error {
	q := `INSERT INTO urls (original_url, short_code) VALUES ($1, $2) RETURNING id, created_at, clicks`
	return r.pool.QueryRow(ctx, q, u.OriginalURL, u.ShortCode).Scan(&u.ID, &u.CreatedAt, &u.Clicks)
}

func (r *PgUrlRepository) GetByCode(ctx context.Context, code string) (*models.URL, error) {
	var u models.URL
	q := `SELECT id, original_url, short_code, clicks, created_at FROM urls WHERE short_code=$1`
	if err := r.pool.QueryRow(ctx, q, code).Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.Clicks, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PgUrlRepository) IncrementClicks(ctx context.Context, code string) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE urls SET clicks = clicks + 1 WHERE short_code=$1`, code)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *PgUrlRepository) GetAll(ctx context.Context) ([]models.URL, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, original_url, short_code, clicks, created_at FROM urls ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var u models.URL
		if err := rows.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.Clicks, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}