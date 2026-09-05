package common

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func OGPThumbnailURL(c *gin.Context) string {
	host := strings.TrimSpace(c.Request.Host)
	scheme := "https"
	if host == "localhost" {
		scheme = "http"
	}

	return scheme + "://" + host + "/ogp-thumbnail"
}
