package gofkass

import (
	"fmt"
	"errors"
	"encoding/json"
)

func I2bytes(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json err:", err)
		return []byte("")
	}
	return b
}

func FromJson(data string) (map[string]interface{}, error) {
	variable := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &variable)
	return variable, err
}

func I2String(data interface{}) (string, error) {
	if a, ok := (data).(string); ok {
		return a, nil
	}
	return "", errors.New("type assertion to string failed")
}