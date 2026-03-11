package service

import (
	"errors"
	"time"

	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

type OpenShiftRequest struct {
    CashIn float64 `json:"cash_in"`
    Notes  string  `json:"notes"`
}

type CloseShiftRequest struct {
	CashOut float64 `json:"cash_out"`
	Notes   string  `json:"notes"`
}

// Ganti bagian OpenShift — EndTime tidak perlu diset:
func OpenShift(req OpenShiftRequest, userID uint) (*model.Shift, error) {
    existing, err := repository.GetActiveShift(userID)
    if err == nil && existing.ID > 0 {
        return nil, errors.New("Anda masih memiliki shift yang sedang berjalan")
    }

    shift := &model.Shift{
        UserID:    userID,
        StartTime: time.Now(),
        CashIn:    req.CashIn,
        Notes:     req.Notes,
    }

    if err := repository.CreateShift(shift); err != nil {
        return nil, errors.New("Gagal membuka shift")
    }

    result, _ := repository.GetShiftByID(shift.ID)
    return result, nil
}

// Ganti bagian CloseShift:
func CloseShift(req CloseShiftRequest, userID uint) (*model.Shift, error) {
    shift, err := repository.GetActiveShift(userID)
    if err != nil || shift.ID == 0 {
        return nil, errors.New("Tidak ada shift aktif untuk ditutup")
    }

    now := time.Now()
    endTime := now  // endTime adalah time.Time biasa

    sales, _ := repository.GetSalesByTimeRange(shift.StartTime, endTime, shift.UserID)
    var totalSales float64
    for _, s := range sales {
        totalSales += s.GrandTotal
    }

    shift.EndTime    = &now
    shift.CashOut    = req.CashOut
    shift.Difference = (shift.CashIn + totalSales) - req.CashOut
    if req.Notes != "" {
        shift.Notes = req.Notes
    }

    if err := repository.UpdateShift(shift); err != nil {
        return nil, errors.New("Gagal menutup shift")
    }

    result, _ := repository.GetShiftByID(shift.ID)
    return result, nil
}

// Ganti bagian GetShiftSummary:
func GetShiftSummary(shiftID uint) (*ShiftSummary, error) {
    shift, err := repository.GetShiftByID(shiftID)
    if err != nil {
        return nil, errors.New("Shift tidak ditemukan")
    }

    endTime := time.Now()
    if shift.EndTime != nil {
        endTime = *shift.EndTime  // ← dereference pointer
    }

    sales, _ := repository.GetSalesByTimeRange(shift.StartTime, endTime, shift.UserID)

    var totalSales float64
    for _, s := range sales {
        totalSales += s.GrandTotal
    }

    return &ShiftSummary{
        Shift:      shift,
        TotalSales: totalSales,
        TotalTrx:   len(sales),
    }, nil
}

func GetActiveShift(userID uint) (*model.Shift, error) {
	shift, err := repository.GetActiveShift(userID)
	if err != nil || shift.ID == 0 {
		return nil, nil // tidak ada shift aktif, bukan error
	}
	return shift, nil
}

func GetShiftHistory(userID uint) ([]model.Shift, error) {
	return repository.GetShiftHistory(userID, 20)
}

type ShiftSummary struct {
	Shift       *model.Shift  `json:"shift"`
	TotalSales  float64       `json:"total_sales"`
	TotalTrx    int           `json:"total_trx"`
}

func UpdateShiftNotes(shiftID uint, notes string) (*model.Shift, error) {
	shift, err := repository.GetShiftByID(shiftID)
	if err != nil || shift.ID == 0 {
		return nil, errors.New("Shift tidak ditemukan")
	}

	shift.Notes = notes
	if err := repository.UpdateShift(shift); err != nil {
		return nil, errors.New("Gagal memperbarui keterangan")
	}

	return shift, nil
}