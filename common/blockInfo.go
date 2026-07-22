package common

import (
	"encoding/csv"
	"os"
)

type BlockInfo struct {
	ColorCode string
	Name      string
}

func GetBlockInfo(filePath string, colorCodeValue *[]string) *map[string]BlockInfo {
	f, err := os.Open(filePath)
	if err != nil {
		println(err)
		return nil
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		println(err)
		return nil
	}

	blockInfoMap := make(map[string]BlockInfo)

	for i, rec := range records {
		if i == 0 {
			continue // Skip header row
		}

		colorCode := rec[0]
		blockName := rec[1]

		blockInfoMap[colorCode] = BlockInfo{
			ColorCode: colorCode,
			Name:      blockName,
		}
	}

	return &blockInfoMap
}
