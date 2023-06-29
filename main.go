package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Rental struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Type            string   `json:"type"`
	Make            string   `json:"make"`
	Model           string   `json:"model"`
	Year            int      `json:"year"`
	Length          float64  `json:"length"`
	Sleeps          int      `json:"sleeps"`
	PrimaryImageURL string   `json:"primary_image_url"`
	Price           Price    `json:"price"`
	Location        Location `json:"location"`
	User            User     `json:"user"`
}

type Price struct {
	Day int `json:"day"`
}

type Location struct {
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zip     string  `json:"zip"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var db *sql.DB

func handleGetRental(w http.ResponseWriter, r *http.Request) {
	rentalID := strings.TrimPrefix(r.URL.Path, "/rentals/")
	id, err := strconv.Atoi(rentalID)
	if err != nil {
		http.Error(w, "Invalid rental ID", http.StatusBadRequest)
		return
	}

	rental := findRentalByIDFromDB(id, db)
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

	rentals, err := filterRentalsFromDB(priceMin, priceMax, db)
	if err != nil {
		http.Error(w, "Failed to retrieve rentals", http.StatusInternalServerError)
		return
	}

	paginatedRentals := paginateRentals(rentals, limit, offset)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedRentals)
}

func findRentalByIDFromDB(id int, db *sql.DB) *Rental {
	query := "SELECT id, name, description, type, make, model, year, length, sleeps, primary_image_url, price_day, location_city, location_state, location_zip, location_country, location_lat, location_lng, user_id, user_first_name, user_last_name FROM rentals WHERE id = $1"
	row := db.QueryRow(query, id)

	var rental Rental
	err := row.Scan(
		&rental.ID,
		&rental.Name,
		&rental.Description,
		&rental.Type,
		&rental.Make,
		&rental.Model,
		&rental.Year,
		&rental.Length,
		&rental.Sleeps,
		&rental.PrimaryImageURL,
		&rental.Price.Day,
		&rental.Location.City,
		&rental.Location.State,
		&rental.Location.Zip,
		&rental.Location.Country,
		&rental.Location.Lat,
		&rental.Location.Lng,
		&rental.User.ID,
		&rental.User.FirstName,
		&rental.User.LastName,
	)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	return &rental
}

func filterRentalsFromDB(priceMin, priceMax int, db *sql.DB) ([]Rental, error) {
	query := "SELECT id, name, description, type, make, model, year, length, sleeps, primary_image_url, price_day, location_city, location_state, location_zip, location_country, location_lat, location_lng, user_id, user_first_name, user_last_name FROM rentals WHERE price_day >= $1 AND price_day <= $2"
	rows, err := db.Query(query, priceMin, priceMax)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []Rental
	for rows.Next() {
		var rental Rental
		err := rows.Scan(
			&rental.ID,
			&rental.Name,
			&rental.Description,
			&rental.Type,
			&rental.Make,
			&rental.Model,
			&rental.Year,
			&rental.Length,
			&rental.Sleeps,
			&rental.PrimaryImageURL,
			&rental.Price.Day,
			&rental.Location.City,
			&rental.Location.State,
			&rental.Location.Zip,
			&rental.Location.Country,
			&rental.Location.Lat,
			&rental.Location.Lng,
			&rental.User.ID,
			&rental.User.FirstName,
			&rental.User.LastName,
		)
		if err != nil {
			return nil, err
		}
		rentals = append(rentals, rental)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rentals, nil
}

func paginateRentals(rentals []Rental, limit, offset int) []Rental {
	if offset >= len(rentals) {
		return []Rental{}
	}

	end := offset + limit
	if end > len(rentals) {
		end = len(rentals)
	}

	return rentals[offset:end]
}

func testDBConnection() {
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", "postgres://username:password@localhost:5432/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection successful")
}

func main() {
	// Test the database connection
	//testDBConnection()
	// Open a connection to the PostgreSQL database
	var err error
	db, err = sql.Open("postgres", "postgres://username:password@localhost:5432/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/rentals/", handleGetRental)
	http.HandleFunc("/rentals", func(w http.ResponseWriter, r *http.Request) {
		handleListRentals(w, r)
	})

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
