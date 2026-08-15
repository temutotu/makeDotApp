package middleware

import (
	commonlang "makeDotApp/common/lang"

	"github.com/gin-gonic/gin"
)

func NewLangMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		langCode := c.Query("lang")
		if langCode == "" {
			langCode = "ja"
		}
		langManager := commonlang.RegisterLangManager(langCode)
		c.Set("langManager", langManager)
		c.Next()
	}
}
