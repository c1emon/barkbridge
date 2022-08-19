package utils

import "encoding/json"

func PrettyMarshal(data any) string {
	b, _ := json.MarshalIndent(data, "", "    ")
	return string(b)
}
