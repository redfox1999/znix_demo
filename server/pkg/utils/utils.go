package utils

import (
	"encoding/json"
)

func ToJson(v interface{}) string {
	jsonStr, _ := json.MarshalIndent(v, "", "  ")
	return string(jsonStr)
}
