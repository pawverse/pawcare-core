package utils

import "encoding/json"

func MapToStruct[T any](m map[string]string) (T, error) {
	var result T
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(jsonStr, &result); err != nil {
		return result, err
	}

	return result, nil
}
