package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"bizkit-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// POST /attendances/checkin
func CheckIn(c *gin.Context) {
	userID := c.GetUint("user_id")

	photoPath, err := saveAttendancePhoto(c, "photo", userID, "in")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	attendance, err := service.CheckIn(userID, photoPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil absen masuk", "data": attendance})
}

// POST /attendances/:id/checkout
func CheckOut(c *gin.Context) {
	userID := c.GetUint("user_id")

	idParam := c.Param("id")
	attendanceID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	photoPath, err := saveAttendancePhoto(c, "photo", userID, "out")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	attendance, err := service.CheckOut(uint(attendanceID), userID, photoPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil absen pulang", "data": attendance})
}

// GET /attendances/today
func GetTodayAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")

	attendance, err := service.GetTodayAttendance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data absensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendance})
}

// GET /attendances/history
func GetAttendanceHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	list, err := service.GetAttendanceHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil riwayat absensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ── Helper: simpan foto ke disk, return path/URL ──────────────────────────
func saveAttendancePhoto(c *gin.Context, field string, userID uint, suffix string) (string, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return "", fmt.Errorf("foto wajib diupload")
	}

	uploadDir := "uploads/attendance"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload")
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%d_%s_%d%s", userID, suffix, time.Now().Unix(), ext)

	// Selalu pakai forward slash untuk URL (bukan filepath.Join yang bisa jadi backslash di Windows)
	savePath := uploadDir + "/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", fmt.Errorf("gagal menyimpan foto: %s", err.Error())
	}

	// Return sebagai URL path: /uploads/attendance/filename.jpg
	return "/uploads/attendance/" + filename, nil
}