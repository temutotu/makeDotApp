package lang

import (
	"makeDotApp/const/lang/ja"
)

type JaManager struct {
	ErrorMessages map[int]string
}

func NewJaManager() *JaManager {
	return &JaManager{
		ErrorMessages: ja.ErrorMessages,
	}
}

func (m *JaManager) GetErrorMessage(key int) string {
	if msg, ok := m.ErrorMessages[key]; ok {
		return msg
	}
	return ""
}
