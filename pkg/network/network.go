package network

import (
	"encoding/json"
	"fmt"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, modelError *models.ModelError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(modelError.ErrorCode)
	marshalBody, err := json.Marshal(&modelError)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(marshalBody)
}

func WriteResponse(w http.ResponseWriter, code int, body interface{ MarshalJSON() ([]byte, error) }) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	marshalBody, err := body.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(marshalBody)
}
