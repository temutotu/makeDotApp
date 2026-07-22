package main

import (
	"html/template"
	"path/filepath"
	"strconv"

	"makeDotApp/common"
	constants "makeDotApp/const"
	handler "makeDotApp/handler"
	response "makeDotApp/network/response"
	selectinput "makeDotApp/templates/input"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

type mainPageData struct {
	DotSizeSelect selectinput.SelectField
	PixelMap      [][]int
	ColorCode     []string
	PixelMapJSON  template.JS
	ColorCodeJSON template.JS
	BlockInfoMap  *map[string]common.BlockInfo
	BlockInfoJSON template.JS
	Error         *response.Error
}

func main() {
	r := gin.Default()
	partFiles, err := filepath.Glob("templates/parts/*.tmpl")
	if err != nil {
		panic(err)
	}

	common.SessionInit()
	store := memstore.NewStore([]byte(common.SessionKey))
	r.Use(sessions.Sessions(common.SessionName, store))

	templateFiles := append([]string{"templates/main.tmpl"}, partFiles...)
	tmpl := template.Must(template.ParseFiles(templateFiles...))
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "./static")
	r.Static("/img", "./img")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/main", func(c *gin.Context) {
		common.StartSession(c)

		selectedDotSize := strconv.Itoa(common.GetSize(c))
		if value := c.Query("dotSize"); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil {
				selectedDotSize = strconv.Itoa(parsed)
			}
		}

		dotSizeOptions := make([]selectinput.SelectOption, 0, len(constants.DotSizes))
		for _, size := range constants.DotSizes {
			value := strconv.Itoa(size)
			dotSizeOptions = append(dotSizeOptions, selectinput.SelectOption{
				Value: value,
				Label: value,
			})
		}

		dotSizeSelect := selectinput.SelectField{
			ID:            "dotSize",
			Name:          "dotSize",
			Label:         "ドットサイズ",
			Options:       dotSizeOptions,
			SelectedValue: selectedDotSize,
		}

		var pixelMap [][]int
		p := common.GetPixelMap(c)
		if p == nil {
			pixelMap = make([][]int, 0)
		} else {
			pixelMap = *p
		}

		var colorCode []string
		cc := common.GetColorCode(c)
		if cc == nil {
			colorCode = make([]string, 0)
		} else {
			colorCode = *cc
		}

		var blockInfoMap *map[string]common.BlockInfo
		var blockInfoMapJSON template.JS
		if pixelMap == nil || len(colorCode) == 0 {
			blockInfoMap = nil
			blockInfoMapJSON = "{}"
		} else {
			blockInfoMap = common.GetBlockInfo(filepath.Join("resource", "blockInfo", "blockInfo.csv"), &colorCode)
			blockInfoMapJSON = common.ToJSONJS(*blockInfoMap)
		}

		c.HTML(200, "main.tmpl", mainPageData{
			DotSizeSelect: dotSizeSelect,
			PixelMap:      pixelMap,
			ColorCode:     colorCode,
			PixelMapJSON:  common.ToJSONJS(pixelMap),
			ColorCodeJSON: common.ToJSONJS(colorCode),
			BlockInfoMap:  blockInfoMap,
			BlockInfoJSON: blockInfoMapJSON,
			Error:         nil,
		})

	})

	r.POST("/makeDot", func(c *gin.Context) {
		handler.MakeDotHandler(c)
	})

	r.Run(":80") // http://localhost:80
}
