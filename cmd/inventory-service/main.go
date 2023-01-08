package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DerRomtester/csgo-app/m/v2/internal/database"
)

var Mg_db = *database.Mongo_ConnectDB()
var Pg_db = *database.Pg_ConnectDB()

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
	crosshaircollection, collectionErr := database.Mongo_ReadCrosshairCollection(&Mg_db)
	if collectionErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Error": "Crosshair Collection not available"})
		return
	}
	json.NewEncoder(w).Encode(crosshaircollection)
}

func getCrosshairsPG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]string{"Error": "405 Method not allowed"})
		return
	}
	crosshaircollection := database.Pg_ReadCrosshairCollection(&Pg_db)
	json.NewEncoder(w).Encode(crosshaircollection)
}

func main() {
	http.HandleFunc("/api/health", apiHealth)
	http.HandleFunc("/api/mg_crosshairs", getCrosshairs)
	http.HandleFunc("/api/pg_crosshairs", getCrosshairsPG)
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Fatal(srv.ListenAndServe())
}
