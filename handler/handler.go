package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"urls/service"
)

// ShortenURL
// @Param request body handler.ShortenURL.req true "query params"
// @Router /url/shorten [post]
// ShortenURL generates a short URL and stores the original and short URL in the MySQL database and Redis cache.
func ShortenURL(c *gin.Context) {
	type req struct {
		OriginalURL string `json:"original_url" binding:"required"`
	}
	var req1 req
	if err := c.ShouldBindJSON(&req1); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := service.Shorten(req1.OriginalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

// ExpandURL
// @Param	short	path	string	true	"short"
// @Router /url/expand/{short} [get]
// ExpandURL expands a short URL by looking up the original URL in the MySQL database and Redis cache.
func ExpandURL(c *gin.Context) {
	short := c.Param("shorten")

	originalURL, err := service.Expand(short)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"original_url": originalURL})
}

// GetHits
// @Param	short	path	string	true	"short"
// @Router /url/hits/{short} [get]
func GetHits(c *gin.Context) {
	short := c.Param("shorten")

	hits, err := service.GetURLHits(short)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hits": hits})
}
