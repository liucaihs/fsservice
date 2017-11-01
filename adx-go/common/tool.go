package common

import (
	"encoding/json"
	"os"
)

func GetEnvDef(name string, def_str string) string {
	value := os.Getenv(name)
	if len(value) > 0 {
		return value
	}
	return def_str
}

func Json2map(req []byte) (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal(req, &result); err != nil {
		return nil, err
	}
	return result, nil
}
