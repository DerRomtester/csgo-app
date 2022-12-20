package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DerRomtester/csgo-app/m/v2/internal/database"
	"github.com/gorilla/mux"
)

func getCrosshairs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	crosshairs := database.ReadCrosshairCollection(database.ConnectDB())
	json.NewEncoder(w).Encode(crosshairs)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/crosshairs", getCrosshairs).Methods("GET")
	log.Fatal(http.ListenAndServe(":6000", router))

}
