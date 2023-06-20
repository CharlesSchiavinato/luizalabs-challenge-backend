package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
)

type NotFound struct{}

func NewNotFound() *NotFound {
	return &NotFound{}
}

func (*NotFound) NotFound(rw http.ResponseWriter, req *http.Request) {
	errorNotFound := model.NotFound("URL")
	rw.WriteHeader(http.StatusNotFound)
	json.NewEncoder(rw).Encode(errorNotFound)
}
