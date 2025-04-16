package json

import (
	"encoding/json"
	"net/http"
)

func (h APIxOBJ_Handler) Write2(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")

	code := http.StatusCreated

	if h.RetCode < 1 {

		if h.RetCode != -100 {

			code = http.StatusUnauthorized

		} else {

			code = h.HttpCode
		}
	}

	w.WriteHeader(code) //http.StatusOK

	json.NewEncoder(w).Encode(h.Data)
}

func Write(w http.ResponseWriter, data interface{}, code int) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code) //http.StatusOK

	json.NewEncoder(w).Encode(data)
}
