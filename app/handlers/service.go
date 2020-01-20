package handlers

import (
	"github.com/lisa-bella97/tech-db-forum/pkg/database"
	"github.com/lisa-bella97/tech-db-forum/pkg/network"
	"net/http"
)

func Clear(w http.ResponseWriter, r *http.Request) {
	database.Clear()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Очистка базы успешно завершена"))
}

func Status(w http.ResponseWriter, r *http.Request) {
	status := database.GetStatus()
	network.WriteResponse(w, http.StatusOK, status)
}
