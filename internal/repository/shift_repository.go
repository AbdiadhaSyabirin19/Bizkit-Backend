package repository

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
)

func GetActiveShift(userID uint) (*model.Shift, error) {
    var shift model.Shift
    err := config.DB.
        Preload("User").
        Where("user_id = ? AND end_time IS NULL", userID).
        Order("created_at DESC").
        First(&shift).Error
    return &shift, err
}

func GetActiveShiftAny() (*model.Shift, error) {
    var shift model.Shift
    err := config.DB.
        Preload("User").
        Where("end_time IS NULL").
        Order("created_at DESC").
        First(&shift).Error
    return &shift, err
}

func CreateShift(shift *model.Shift) error {
	return config.DB.Create(shift).Error
}

func UpdateShift(shift *model.Shift) error {
	return config.DB.Save(shift).Error
}

func GetShiftByID(id uint) (*model.Shift, error) {
	var shift model.Shift
	err := config.DB.Preload("User").First(&shift, id).Error
	return &shift, err
}

func GetShiftHistory(userID uint, limit int) ([]model.Shift, error) {
	var shifts []model.Shift
	q := config.DB.Preload("User").Order("created_at DESC")
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&shifts).Error
	return shifts, err
}