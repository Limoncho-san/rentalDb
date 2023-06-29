package rental

import (
	"database/sql"
	"log"
)

// Rental struct definition

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

// Price struct definition

type Price struct {
	Day int `json:"day"`
}

// Location struct definition

type Location struct {
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zip     string  `json:"zip"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

// User struct definition

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var db *sql.DB

// FindRentalByIDFromDB retrieves a rental by ID from the database

func FindRentalByIDFromDB(id int) (*Rental, error) {
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
		return nil, nil
	} else if err != nil {
		log.Fatal(err)
	}

	return &rental, nil
}

// FilterRentalsFromDB filters rentals based on price range from the database

func FilterRentalsFromDB(priceMin, priceMax int) ([]Rental, error) {
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

// PaginateRentals returns a slice of rentals based on limit and offset

func PaginateRentals(rentals []Rental, limit, offset int) []Rental {
	if offset >= len(rentals) {
		return []Rental{}
	}

	end := offset + limit
	if end > len(rentals) {
		end = len(rentals)
	}

	return rentals[offset:end]
}
