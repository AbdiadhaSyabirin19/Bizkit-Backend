package service

import (
	"errors"
	"time"

	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

// CheckIn — buat absensi baru, simpan foto masuk
func CheckIn(userID uint, photoIn string) (*model.Attendance, error) {
	// Cek sudah checkin hari ini?
	existing, err := repository.GetTodayAttendance(userID)
	if err == nil && existing != nil && existing.ID > 0 {
		return nil, errors.New("Anda sudah absen masuk hari ini")
	}

	attendance := &model.Attendance{
		UserID:  userID,
		CheckIn: time.Now(),
		PhotoIn: photoIn,
		// CheckOut sengaja tidak diset → akan NULL di database
	}

	if err := repository.CreateAttendance(attendance); err != nil {
		return nil, errors.New("Gagal menyimpan absensi: " + err.Error())
	}

	result, _ := repository.GetAttendanceByID(attendance.ID)
	return result, nil
}

// CheckOut — update absensi, simpan foto pulang
func CheckOut(attendanceID uint, userID uint, photoOut string) (*model.Attendance, error) {
	attendance, err := repository.GetAttendanceByID(attendanceID)
	if err != nil || attendance.ID == 0 {
		return nil, errors.New("Data absensi tidak ditemukan")
	}

	if attendance.UserID != userID {
		return nil, errors.New("Tidak diizinkan")
	}

	// CheckOut sekarang pointer — cek nil berarti belum checkout
	if attendance.CheckOut != nil {
		return nil, errors.New("Anda sudah absen pulang")
	}

	now := time.Now()
	attendance.CheckOut = &now
	attendance.PhotoOut = photoOut

	if err := repository.UpdateAttendance(attendance); err != nil {
		return nil, errors.New("Gagal menyimpan absensi pulang: " + err.Error())
	}

	result, _ := repository.GetAttendanceByID(attendance.ID)
	return result, nil
}

// GetTodayAttendance — absensi hari ini milik user
func GetTodayAttendance(userID uint) (*model.Attendance, error) {
	a, err := repository.GetTodayAttendance(userID)
	if err != nil || a == nil || a.ID == 0 {
		return nil, nil // belum absen hari ini, bukan error
	}
	return a, nil
}

// GetAttendanceHistory — riwayat absensi
func GetAttendanceHistory(userID uint) ([]model.Attendance, error) {
	return repository.GetAttendanceHistory(userID, 30)
}