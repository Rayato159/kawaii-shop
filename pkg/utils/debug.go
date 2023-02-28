package utils

import (
	"encoding/json"
	"fmt"
)

func Debug(obj any) {
	bytes, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(bytes))
}

func Output(obj any) []byte {
	bytes, _ := json.Marshal(obj)
	return bytes
}
