package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w *http.ResponseWriter, datas *map[string]interface{}) {
	(*w).Header().Add("Content-Type", "application/json")

	if datas != nil {
		j, err := json.Marshal(datas)
		if err == nil {
			(*w).Write(j)
			return
		}
	}
}

func jsonError(w http.ResponseWriter, e error) {
	failure := map[string]interface{}{
		"success": false,
		"error":   e.Error(),
	}
	writeJSON(&w, &failure)
}

func jsonSuccess(w http.ResponseWriter, data map[string]interface{}) {
	failure := map[string]interface{}{
		"success": true,
		"result":  data,
	}
	writeJSON(&w, &failure)
}
