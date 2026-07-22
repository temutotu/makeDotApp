package common

import (
	"encoding/json"
	"html/template"
)

func ToJSONJS(v any) template.JS {
	b, err := json.Marshal(v)
	if err != nil {
		return template.JS("null")
	}

	return template.JS(string(b))
}
