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
	if r.Method != http.MethodGet {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]string{"Error": "405 Method not allowed"})
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"available": true})
}

func getCrosshairs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]string{"Error": "405 Method not allowed"})
		return
	}
	crosshaircollection, collectionErr := database.ReadCrosshairCollection(db)
	if collectionErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Crosshair Collection not available"})
		return
	}
	json.NewEncoder(w).Encode(crosshaircollection)
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
