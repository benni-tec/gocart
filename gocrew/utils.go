package gocrew

import (
	"encoding/json"
)

func structToMap(obj any) (map[string]any, error) {
	str, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var mp map[string]any
	err = json.Unmarshal(str, &mp)
	if err != nil {
		return nil, err
	}

	return mp, nil
}
