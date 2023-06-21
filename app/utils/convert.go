package utils

import (
	"encoding/json"
	"fmt"
)

func ConvertJsonDataAnyToString(input map[string]any) (map[string]string, error) {
	var data map[string]string
	dataMap, ok := input["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Error: 'data' is not a map[string]interface{}")
	}

	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
