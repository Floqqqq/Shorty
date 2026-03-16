package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/yourusername/shorty/internal/config"
	"github.com/yourusername/shorty/internal/services"
)

type URLHandler struct {
	svc *services.URLService
	cfg *config.Config
}

func NewURLHandler(s *services.URLService, cfg *config.Config) *URLHandler {
	return &URLHandler{svc: s, cfg: cfg}
}

type shortenReq struct {
	URL string `json:"url" binding:"required"`
}
type shortenResp struct {
	ShortURL string `json:"short_url"`
	Code     string `json:"code"`
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req shortenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	u, err := h.svc.CreateShort(ctx, req.URL)
	if err != nil {
		if err == services.ErrInvalidURL {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Msg("create short")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	short := h.cfg.BaseURL + "/" + u.ShortCode
	c.JSON(http.StatusCreated, shortenResp{ShortURL: short, Code: u.ShortCode})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	u, err := h.svc.Resolve(ctx, code)
	if err != nil || u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	c.Redirect(http.StatusFound, u.OriginalURL)
}

func (h *URLHandler) Stats(c *gin.Context) {
	code := c.Param("code")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	u, err := h.svc.GetStats(ctx, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"short_code": u.ShortCode,
		"original":   u.OriginalURL,
		"clicks":     u.Clicks,
		"created_at": u.CreatedAt,
	})

}
func (h *URLHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	urls, err := h.svc.GetAllURLs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, urls)
}
