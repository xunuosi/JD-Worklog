package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/example/worklog-system/internal/models"
)

type TimesheetHandler struct{ DB *gorm.DB }
type tsReq struct {
	ProjectID uint    `json:"project_id"`
	Date      string  `json:"date"` // YYYY-MM-DD
	Hours     float32 `json:"hours"`
	Content   string  `json:"content"`
}

func (h *TimesheetHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req tsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad date"})
		return
	}
	ts := models.Timesheet{UserID: userID, ProjectID: req.ProjectID, Date: d, Hours: req.Hours, Content: req.Content}
	if err := h.DB.Create(&ts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ts)
}

func (h *TimesheetHandler) ListMine(c *gin.Context) {
	userID := c.GetUint("user_id")
	var list []models.Timesheet
	if err := h.DB.Preload("Project").Where("user_id = ?", userID).Order("date desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *TimesheetHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	res := h.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Timesheet{})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
