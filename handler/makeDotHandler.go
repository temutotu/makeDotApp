package handler

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"makeDotApp/common"
	"makeDotApp/common/lang"
	constants "makeDotApp/const"
	messageKey "makeDotApp/const/lang/messageKey"
	response "makeDotApp/network/response"
	selectinput "makeDotApp/templates/input"

	makeDot "github.com/temutotu/makeDot"

	"github.com/gin-gonic/gin"
)

func MakeDotHandler(c *gin.Context) {
	langManager := lang.GetLangManagerFromContext(c)

	dotSizeStr := c.PostForm("dotSize")
	dotSize, err := strconv.Atoi(dotSizeStr)
	if err != nil {
		c.HTML(400, "main.tmpl", response.BuildMakeDotIndexErrorResponse(selectinput.SelectField{
			ID:            "dotSize",
			Name:          "dotSize",
			Label:         "ドットサイズ",
			Options:       buildDotSizeOptions(),
			SelectedValue: "32",
		}, 400, "Invalid dot size"))
		return
	}

	dotSizeSelect := selectinput.SelectField{
		ID:            "dotSize",
		Name:          "dotSize",
		Label:         "ドットサイズ",
		Options:       buildDotSizeOptions(),
		SelectedValue: strconv.Itoa(dotSize),
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.HTML(400, "main.tmpl", response.BuildMakeDotIndexErrorResponse(dotSizeSelect, 400, langManager.GetErrorMessage(messageKey.ERROR_FAILED_TO_PARSE_MULTIPART_FORM)))
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.HTML(400, "main.tmpl", response.BuildMakeDotIndexErrorResponse(dotSizeSelect, 400, langManager.GetErrorMessage(messageKey.ERROR_NO_IMAGE_FILES_PROVIDED)))
		return
	}

	file := files[0]
	src, err := file.Open()
	if err != nil {
		c.HTML(500, "main.tmpl", response.BuildMakeDotIndexErrorResponse(dotSizeSelect, 500, langManager.GetErrorMessage(messageKey.ERROR_FAILED_TO_OPEN_IMAGE_FILE)))
		return
	}
	defer src.Close()

	src.Seek(0, 0) // ポインタをリセット
	srcImg, format, err := image.Decode(src)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		c.HTML(400, "main.tmpl", response.BuildMakeDotIndexErrorResponse(dotSizeSelect, 400, langManager.GetErrorMessage(messageKey.ERROR_FAILED_TO_DECODE_IMAGE)))
		return
	}

	if format != "jpeg" && format != "png" {
		c.HTML(400, "main.tmpl", response.BuildMakeDotIndexErrorResponse(dotSizeSelect, 400, langManager.GetErrorMessage(messageKey.ERROR_UNSUPPORTED_IMAGE_FORMAT)))
		return
	}

	dstImg := image.NewRGBA(srcImg.Bounds())

	conf := &makeDot.MakeDotConfig{
		OutputPath:    filepath.Join("resource", "output.png"),
		NumImgPath:    filepath.Join("resource", "numImage") + string(filepath.Separator),
		BlockInfoPath: filepath.Join("resource", "blockInfo", "blockInfo.csv"),
		PixcelSize:    dotSize,
	}

	pixelMap, colorCode, err := makeDot.MakeDotMap(srcImg, dstImg, conf)
	if err != nil {
		c.HTML(500, "main.tmpl", response.BuildMakeDotIndexErrorResponse(dotSizeSelect, 500, langManager.GetErrorMessage(messageKey.ERROR_MAKE_DOT_FAILED)))
		return
	}

	pixelMapValue := make([][]int, 0)
	if pixelMap != nil {
		pixelMapValue = *pixelMap
	}

	colorCodeValue := make([]string, 0)
	if colorCode != nil {
		colorCodeValue = *colorCode
	}

	blockInfoMap := common.GetBlockInfo(conf.BlockInfoPath, &colorCodeValue)
	if blockInfoMap == nil {
		c.HTML(500, "main.tmpl", response.BuildMakeDotIndexErrorResponse(dotSizeSelect, 500, langManager.GetErrorMessage(messageKey.ERROR_FAILED_TO_GET_BLOCK_INFO)))
		return
	}

	common.SetSize(c, dotSize)
	if err := common.SetPixelMap(c, pixelMap, colorCode); err != nil {
		log.Printf("SetPixelMap failed: %v", err)
	}

	c.HTML(200, "main.tmpl", response.MakeDotIndex{
		DotSizeSelect: dotSizeSelect,
		PixelMap:      pixelMapValue,
		ColorCode:     colorCodeValue,
		PixelMapJSON:  common.ToJSONJS(pixelMapValue),
		ColorCodeJSON: common.ToJSONJS(colorCodeValue),
		BlockInfoMap:  blockInfoMap,
		BlockInfoJSON: common.ToJSONJS(blockInfoMap),
		OGPImageURL:   common.OGPThumbnailURL(c),
	})
}

func OGPThumbnailHandler(c *gin.Context) {
	outputPath := filepath.Join("resource", "thumbnail", "thumbnail.png")
	if _, err := os.Stat(outputPath); err != nil {
		if os.IsNotExist(err) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Cache-Control", "no-store")
	c.File(outputPath)
}

func buildDotSizeOptions() []selectinput.SelectOption {
	options := make([]selectinput.SelectOption, 0, len(constants.DotSizes))
	for _, size := range constants.DotSizes {
		value := strconv.Itoa(size)
		options = append(options, selectinput.SelectOption{
			Value: value,
			Label: value,
		})
	}

	return options
}
