package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "kcal-api")
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	food := mux.Vars(r)["food"]
	measure := mux.Vars(r)["measure"]
	qtd, _ := strconv.ParseFloat(mux.Vars(r)["qtd"], 64)

	qtdByUnit := GetInfo(measure, food)
	result := map[string]interface{}{
		"food":      food,
		"measure":   measure,
		"qtd":       qtd,
		"qtdByUnit": qtdByUnit,
		"total":     qtd * qtdByUnit,
	}
	json.NewEncoder(w).Encode(result)
}

func saveInfo(w http.ResponseWriter, r *http.Request) {
	var newFood NewFood

	err := json.NewDecoder(r.Body).Decode(&newFood)

	if err != nil || !newFood.Valid() {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = SaveFood(newFood)

	if err != nil {
		errJson := make(map[string]string)
		errJson["error"] = "Already exist"
		json.NewEncoder(w).Encode(errJson)
	} else {
		json.NewEncoder(w).Encode(newFood)
	}

}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/info/{food}/{measure}/{qtd}", getInfo).Methods("GET")
	router.HandleFunc("/info", saveInfo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8444", router))
}
