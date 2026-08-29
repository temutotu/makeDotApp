package lang

import "github.com/gin-gonic/gin"

type LangManager interface {
	GetErrorMessage(key int) string
}

func RegisterLangManager(langStr string) LangManager {
	switch langStr {
	case "ja":
		return NewJaManager()
	default:
		return NewJaManager()
	}
}

func GetLangManagerFromContext(c *gin.Context) LangManager {
	value, exists := c.Get("langManager")
	if !exists {
		return NewJaManager()
	}
	langManager, ok := value.(LangManager)
	if !ok {
		return NewJaManager()
	}
	return langManager
}
