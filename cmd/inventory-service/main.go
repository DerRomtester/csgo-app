package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DerRomtester/csgo-app/m/v2/internal/database"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.HandleFunc("/api/health", apiHealth).Methods("GET")
	router.HandleFunc("/api/crosshairs", getCrosshairs).Methods("GET")
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	log.Fatal(srv.ListenAndServe())
}
