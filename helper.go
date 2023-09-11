package main

import "encoding/json"

func SprintInterfaceJson(a interface{}) string {
	r, _ := json.MarshalIndent(a, "", "  ")
	return string(r)
}
