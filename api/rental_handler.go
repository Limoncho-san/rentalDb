package api

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"

	model "github.com/Limoncho-san/rentalDb/models"
	"github.com/gin-gonic/gin"
)

type RentalHandler struct {
	service *RentalService
}

func NewRentalHandler(db *gorm.DB) *RentalHandler {
	return &RentalHandler{service: NewRentalService(db)}
}

func (h *RentalHandler) GetRental(c *gin.Context) {
	rentalID := c.Param("id")
	id, err := strconv.Atoi(rentalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rental ID"})
		return
	}

	rental, err := h.service.GetRentalByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if rental == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rental not found"})
		return
	}

	c.JSON(http.StatusOK, rental)
}

func (h *RentalHandler) ListRentals(c *gin.Context) {
	priceMinStr := c.Query("price_min")
	priceMaxStr := c.Query("price_max")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	idsStr := c.Query("ids")
	nearStr := c.Query("near")
	sort := c.Query("sort")

	priceMin, err := strconv.Atoi(priceMinStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price_min parameter"})
		return
	}

	priceMax, err := strconv.Atoi(priceMaxStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price_max parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	var ids []int
	if idsStr != "" {
		idsStrArr := strings.Split(idsStr, ",")
		ids = make([]int, len(idsStrArr))
		for i, idStr := range idsStrArr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ids parameter"})
				return
			}
			ids[i] = id
		}
	}

	var near []float64
	if nearStr != "" {
		nearStrArr := strings.Split(nearStr, ",")
		if len(nearStrArr) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid near parameter"})
			return
		}
		for _, val := range nearStrArr {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid near parameter"})
				return
			}
			near = append(near, f)
		}
	}

	filter := &model.RentalFilter{
		PriceMin: priceMin,
		PriceMax: priceMax,
		Limit:    limit,
		Offset:   offset,
		IDs:      ids,
		Near:     near,
		Sort:     sort,
	}

	rentals, total, err := h.service.ListRentals(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	response := &ListRentalsResponse{
		Total:   total,
		Rentals: rentals,
	}

	c.JSON(http.StatusOK, response)
}
