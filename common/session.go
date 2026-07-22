package common

import (
	"encoding/gob"
	"errors"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SessionKey  = "secret-key"
	SessionName = "session-name"
)

type SessionData struct {
	SelectedSize int
	PixelMap     [][]int
	ColorCode    []string
}

func SessionInit() {
	gob.Register(&SessionData{})
}

func getSessionData(c *gin.Context) (sessions.Session, *SessionData, error) {
	if c == nil {
		return nil, nil, errors.New("invalid gin context")
	}

	session := sessions.Default(c)
	raw := session.Get(SessionKey)
	if raw == nil {
		return session, &SessionData{}, nil
	}

	if data, ok := raw.(*SessionData); ok {
		return session, data, nil
	}

	if data, ok := raw.(SessionData); ok {
		copied := data
		return session, &copied, nil
	}

	return session, &SessionData{}, nil
}

func StartSession(c *gin.Context) error {
	session, data, err := getSessionData(c)
	if err != nil {
		return err
	}

	if session.Get(SessionKey) == nil {
		session.Set(SessionKey, *data)
		if err := session.Save(); err != nil {
			log.Printf("failed to save session: %v", err)
			return err
		}
		return nil
	}

	return nil
}

func GetSize(c *gin.Context) int {
	_, data, err := getSessionData(c)
	if err != nil {
		return 32
	}
	if data.SelectedSize <= 0 {
		return 32
	}

	return data.SelectedSize
}

func SetSize(c *gin.Context, size int) {
	session, data, err := getSessionData(c)
	if err != nil {
		return
	}

	data.SelectedSize = size
	session.Set(SessionKey, *data)
	if err := session.Save(); err != nil {
		log.Printf("failed to save session: %v", err)
	}
}

func GetPixelMap(c *gin.Context) *[][]int {
	_, data, err := getSessionData(c)
	if err != nil {
		return nil
	}

	if data.PixelMap == nil {
		return nil
	}

	return &data.PixelMap
}

func GetColorCode(c *gin.Context) *[]string {
	_, data, err := getSessionData(c)
	if err != nil {
		return nil
	}

	if data.ColorCode == nil {
		return nil
	}

	return &data.ColorCode
}

func SetPixelMap(c *gin.Context, pixelMap *[][]int, colorCode *[]string) error {
	session, data, err := getSessionData(c)
	if err != nil {
		return err
	}

	if pixelMap != nil {
		data.PixelMap = *pixelMap
	}
	if colorCode != nil {
		data.ColorCode = *colorCode
	}
	session.Set(SessionKey, *data)
	if err := session.Save(); err != nil {
		pixelMapLen := len(data.PixelMap)
		colorCodeLen := len(data.ColorCode)
		log.Printf("failed to save session in SetPixelMap (pixelMapLen=%d colorCodeLen=%d): %v", pixelMapLen, colorCodeLen, err)
		return err
	}

	return nil
}
