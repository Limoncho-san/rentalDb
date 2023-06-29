package api

import (
	model "github.com/Limoncho-san/rentalDb/models"
	"github.com/jinzhu/gorm"
)

type RentalService struct {
	db *gorm.DB
}

func NewRentalService(db *gorm.DB) *RentalService {
	return &RentalService{db: db}
}

func (s *RentalService) GetRentalByID(id int) (*model.Rental, error) {
	rental := &model.Rental{}
	err := s.db.Preload("User").First(rental, id).Error
	if err != nil {
		return nil, err
	}
	return rental, nil
}

func (s *RentalService) ListRentals(filter *model.RentalFilter) ([]model.Rental, int, error) {
	rentals := []model.Rental{}
	query := s.db.Preload("User").
		Where("price_day >= ?", filter.PriceMin).
		Where("price_day <= ?", filter.PriceMax)

	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", filter.IDs)
	}

	if len(filter.Near) == 2 {
		query = query.Where("earth_box(ll_to_earth(?, ?), ?) @> ll_to_earth(location.lat, location.lng)",
			filter.Near[0], filter.Near[1], 100) // Assuming the distance is in miles
	}

	var total int
	if err := query.Model(&model.Rental{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Sort == "price" {
		query = query.Order("price_day ASC")
	}

	err := query.Limit(filter.Limit).Offset(filter.Offset).Find(&rentals).Error
	if err != nil {
		return nil, 0, err
	}

	return rentals, total, nil
}
