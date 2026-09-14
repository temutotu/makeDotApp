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
	SelectedSize uint16
	PixelMap     [][]uint8
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

	return int(data.SelectedSize)
}

func SetSize(c *gin.Context, size int) {
	session, data, err := getSessionData(c)
	if err != nil {
		return
	}

	data.SelectedSize = uint16(size)
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

	pixelMap := make([][]int, len(data.PixelMap))
	for rowIndex, row := range data.PixelMap {
		pixelMap[rowIndex] = make([]int, len(row))
		for columnIndex, value := range row {
			pixelMap[rowIndex][columnIndex] = int(value)
		}
	}

	return &pixelMap
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
		data.PixelMap = make([][]uint8, len(*pixelMap))
		for rowIndex, row := range *pixelMap {
			data.PixelMap[rowIndex] = make([]uint8, len(row))
			for columnIndex, value := range row {
				data.PixelMap[rowIndex][columnIndex] = uint8(value)
			}
		}
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
