package api

import model "github.com/Limoncho-san/rentalDb/models"

type ListRentalsResponse struct {
	Total   int            `json:"total"`
	Rentals []model.Rental `json:"rentals"`
}
