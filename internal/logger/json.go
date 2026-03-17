package logger

import (
	"encoding/json"
	"log"
)

func DebugJson(v any) {
	data, _ := json.MarshalIndent(v, "", " ")
	log.Println(string(data))
}
