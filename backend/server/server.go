package server

import (
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()
	_, path, handler := NewMemoHandler()
	mux.Handle(path, handler)

	return mux
}
