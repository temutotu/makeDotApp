package ja

import messagekey "makeDotApp/const/lang/messageKey"

var ErrorMessages = map[int]string{
	messagekey.ERROR_FAILED_TO_PARSE_MULTIPART_FORM: "リクエストに失敗しました。もう一度お試しください。",
	messagekey.ERROR_NO_IMAGE_FILES_PROVIDED:        "画像ファイルが選択されていません。",
	messagekey.ERROR_FAILED_TO_OPEN_IMAGE_FILE:      "画像ファイルの読み込みに失敗しました。もう一度お試しください。",
	messagekey.ERROR_FAILED_TO_DECODE_IMAGE:         "画像ファイルのデコードに失敗しました。PNGまたはJPEG形式の画像ファイルを選択してください。",
	messagekey.ERROR_UNSUPPORTED_IMAGE_FORMAT:       "サポートされていない画像形式です。PNGまたはJPEG形式の画像ファイルを選択してください。",
	messagekey.ERROR_INVALID_DOT_SIZE:               "ドットサイズが不正です。",
	messagekey.ERROR_MAKE_DOT_FAILED:                "ドット絵の生成に失敗しました。",
	messagekey.ERROR_FAILED_TO_GET_BLOCK_INFO:       "ブロック情報の取得に失敗しました。",
}
