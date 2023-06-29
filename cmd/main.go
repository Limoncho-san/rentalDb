package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Limoncho-san/rentalDb/rental"
)

func handleGetRental(w http.ResponseWriter, r *http.Request) {
	rentalID := strings.TrimPrefix(r.URL.Path, "/rentals/")
	id, err := strconv.Atoi(rentalID)
	if err != nil {
		http.Error(w, "Invalid rental ID", http.StatusBadRequest)
		return
	}

	rental, err := rental.FindRentalByIDFromDB(id)
	if err != nil {
		http.Error(w, "Failed to retrieve rental", http.StatusInternalServerError)
		return
	}
	if rental == nil {
		http.Error(w, "Rental not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rental)
}

func handleListRentals(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	priceMin, _ := strconv.Atoi(query.Get("price_min"))
	priceMax, _ := strconv.Atoi(query.Get("price_max"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	rentals, err := rental.FilterRentalsFromDB(priceMin, priceMax)
	if err != nil {
		http.Error(w, "Failed to retrieve rentals", http.StatusInternalServerError)
		return
	}

	paginatedRentals := rental.PaginateRentals(rentals, limit, offset)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedRentals)
}

func main() {
	connectionString := "postgres://postgres:password@localhost:5432/rentaldb_default?sslmode=disable"
	err := rental.InitDB(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer rental.CloseDB()

	http.HandleFunc("/rentals/", handleGetRental)
	http.HandleFunc("/rentals", handleListRentals)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
