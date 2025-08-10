package handlers

import (
	"net/http"
	"strconv"
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

// ListMine GET /api/timesheets/mine?page=1&page_size=10
func (h *TimesheetHandler) ListMine(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 读取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	var total int64
	if err := h.DB.Model(&models.Timesheet{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}

	var list []models.Timesheet
	if err := h.DB.Preload("Project").
		Where("user_id = ?", userID).
		Order("date DESC, id DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     list,
		"total":     total,
		"page":      page,
		"page_size": size,
	})
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

func (h *TimesheetHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var req struct {
		ProjectID *uint    `json:"project_id"` // 可选
		Date      *string  `json:"date"`       // YYYY-MM-DD，可选
		Hours     *float32 `json:"hours"`      // 可选
		Content   *string  `json:"content"`    // 可选
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	var ts models.Timesheet
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).
		First(&ts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if req.ProjectID != nil {
		ts.ProjectID = *req.ProjectID
	}
	if req.Hours != nil {
		ts.Hours = *req.Hours
	}
	if req.Content != nil {
		ts.Content = *req.Content
	}
	if req.Date != nil && *req.Date != "" {
		d, err := time.ParseInLocation("2006-01-02", *req.Date, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad date"})
			return
		}
		ts.Date = d
	}

	if err := h.DB.Save(&ts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	c.JSON(http.StatusOK, ts)
}
