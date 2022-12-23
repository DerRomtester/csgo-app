package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DerRomtester/csgo-app/m/v2/internal/database"
)

var db = database.ConnectDB()

func apiHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"alive": true})
}

func getCrosshairs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	crosshairs := database.ReadCrosshairCollection(db)
	json.NewEncoder(w).Encode(crosshairs)
}

func main() {
	http.HandleFunc("/api/health", apiHealth)
	http.HandleFunc("/api/crosshairs", getCrosshairs)
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Fatal(srv.ListenAndServe())
}
