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
