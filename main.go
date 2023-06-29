package main

import (
	"github.com/Limoncho-san/rentalDb/api"
	"github.com/Limoncho-san/rentalDb/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	rentalHandler := api.NewRentalHandler(db)

	r.GET("/rentals/:id", rentalHandler.GetRental)
	r.GET("/rentals", rentalHandler.ListRentals)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
