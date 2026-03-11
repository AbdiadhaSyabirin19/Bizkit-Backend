package handler

import (
	"net/http"
	"strconv"

	"bizkit-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// POST /shifts/open
func OpenShift(c *gin.Context) {
	var req service.OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid"})
		return
	}

	userID := c.MustGet("user_id").(uint)
	shift, err := service.OpenShift(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Shift berhasil dibuka",
		"data":    shift,
	})
}

// POST /shifts/close
func CloseShift(c *gin.Context) {
	var req service.CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid"})
		return
	}

	userID := c.MustGet("user_id").(uint)
	shift, err := service.CloseShift(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shift berhasil ditutup",
		"data":    shift,
	})
}

// GET /shifts/active
func GetActiveShift(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	shift, err := service.GetActiveShift(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data shift"})
		return
	}
	if shift == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Tidak ada shift aktif", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": shift})
}

// GET /shifts/history
func GetShiftHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	shifts, err := service.GetShiftHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil riwayat shift"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": shifts})
}

// GET /shifts/:id/summary
func GetShiftSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	summary, err := service.GetShiftSummary(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": summary})
}

func UpdateShiftNotes(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request tidak valid"})
		return
	}

	shift, err := service.UpdateShiftNotes(uint(id), req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Keterangan berhasil diperbarui", "data": shift})
}