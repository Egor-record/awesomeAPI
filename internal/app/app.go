package app

import (
	"awesomeAPI/internal/delivery/http"
)

func Start() {
	handler := http.Handler{}
	router := handler.Init()
	router.Run()
}
