package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"makeDotApp/common"
	constants "makeDotApp/const"
	handler "makeDotApp/handler"
	"makeDotApp/middleware"
	response "makeDotApp/network/response"
	selectinput "makeDotApp/templates/input"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type mainPageData struct {
	DotSizeSelect selectinput.SelectField
	PixelMap      [][]int
	ColorCode     []string
	PixelMapJSON  template.JS
	ColorCodeJSON template.JS
	BlockInfoMap  *map[string]common.BlockInfo
	BlockInfoJSON template.JS
	OGPImageURL   string
	Error         *response.Error
}

func main() {
	r := gin.Default()
	r.Use(middleware.NewLangMiddleware())
	r.MaxMultipartMemory = 8 << 20 // 8MB
	r.Use(middleware.NewGlobalRateLimitMiddleware(rate.Limit(30), 60))
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
	r.GET("/ogp-thumbnail", handler.OGPThumbnailHandler)

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
			OGPImageURL:   ogpThumbnailURL(c),
			Error:         nil,
		})

	})

	r.GET("/makeDot", func(c *gin.Context) {
		c.Redirect(302, "/main")
	})

	makeDotProtected := r.Group("/makeDot")
	makeDotProtected.Use(
		middleware.NewConcurrencyLimitMiddleware(8),
		middleware.NewMaxBodyBytesMiddleware(10<<20), // 10MB
	)
	makeDotProtected.POST("", func(c *gin.Context) {
		handler.MakeDotHandler(c)
	})

	srv := &http.Server{
		Addr:              ":80",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func ogpThumbnailURL(c *gin.Context) string {
	host := strings.TrimSpace(c.Request.Host)
	scheme := "http"
	if c.Request.TLS != nil || strings.EqualFold(strings.TrimSpace(strings.Split(c.GetHeader("X-Forwarded-Proto"), ",")[0]), "https") {
		scheme = "https"
	}

	return scheme + "://" + host + "/ogp-thumbnail"
}
