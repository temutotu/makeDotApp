package response

import (
	"html/template"
	"makeDotApp/common"

	selectinput "makeDotApp/templates/input"
)

type MakeDotIndex struct {
	DotSizeSelect selectinput.SelectField
	PixelMap      [][]int
	ColorCode     []string
	PixelMapJSON  template.JS
	ColorCodeJSON template.JS
	BlockInfoMap  *map[string]common.BlockInfo
	BlockInfoJSON template.JS
	Error         *Error
}

type MakeDotIndexError struct {
	DotSizeSelect selectinput.SelectField
	PixelMapJSON  template.JS
	ColorCodeJSON template.JS
	BlockInfoJSON template.JS
	Error         *Error
}

func BuildMakeDotIndexErrorResponse(dotSizeSelect selectinput.SelectField, code int, message string) *MakeDotIndexError {
	return &MakeDotIndexError{
		DotSizeSelect: dotSizeSelect,
		PixelMapJSON:  common.ToJSONJS(make([][]int, 0)),
		ColorCodeJSON: common.ToJSONJS(make([]string, 0)),
		BlockInfoJSON: common.ToJSONJS(map[string]common.BlockInfo{}),
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}
