package scf

import (
	"encoding/json"
)

func Json(data interface{}) string {
	if val, err := json.Marshal(data); err != nil {
		return ""
	} else {
		return string(val)
	}
}

func Map(data string) map[string]interface{} {
	var mp map[string]interface{}
	if err := json.Unmarshal([]byte(data), &mp); err != nil {
		return nil
	}
	return mp
}
