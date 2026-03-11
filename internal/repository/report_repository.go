package repository

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
	"time"
)

func GetShiftsByPeriod(start, end time.Time) ([]model.Shift, error) {
	var shifts []model.Shift
	result := config.DB.
		Preload("User").
		Where("start_time BETWEEN ? AND ?", start, end).
		Order("start_time DESC").
		Find(&shifts)
	return shifts, result.Error
}