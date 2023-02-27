package utils

import (
	"encoding/json"
	"fmt"
)

func Debug(obj any) {
	bytes, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(bytes))
}
