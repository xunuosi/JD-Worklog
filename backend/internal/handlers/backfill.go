package handlers

import (
	"net/http"
	"time"

	"github.com/example/worklog-system/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BackfillRequest struct {
	ProjectID uint    `json:"project_id" binding:"required"`
	UserID    uint    `json:"user_id" binding:"required"`
	TotalDays float32 `json:"total_days" binding:"required"`
	StartDate string  `json:"start_date" binding:"required"`
	EndDate   string  `json:"end_date" binding:"required"`
	Content   string  `json:"content"`
}

func BackfillTimesheets(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BackfillRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		layout := "2006-01-02"
		startDate, err := time.Parse(layout, req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
			return
		}
		endDate, err := time.Parse(layout, req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
			return
		}

		adminID := c.GetUint("user_id")

		if startDate.After(endDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start date cannot be after end date"})
			return
		}

		var existingTimesheets []models.Timesheet
		db.Where("user_id = ? AND date >= ? AND date <= ?", req.UserID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Find(&existingTimesheets)

		existingDates := make(map[time.Time]bool)
		for _, ts := range existingTimesheets {
			// Normalize to UTC midnight
			utcDate := time.Date(ts.Date.Year(), ts.Date.Month(), ts.Date.Day(), 0, 0, 0, 0, time.UTC)
			existingDates[utcDate] = true
		}

		var availableDays []time.Time
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			// d is already UTC midnight from parsing
			if !existingDates[d] {
				availableDays = append(availableDays, d)
			}
		}

		if float32(len(availableDays)) < req.TotalDays {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not enough available days in the selected date range"})
			return
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			backfillLog := models.BackfillLog{
				AdminID:   adminID,
				UserID:    req.UserID,
				ProjectID: req.ProjectID,
				TotalDays: req.TotalDays,
				StartDate: startDate,
				EndDate:   endDate,
			}
			if err := tx.Create(&backfillLog).Error; err != nil {
				return err
			}

			for i := 0; i < int(req.TotalDays); i++ {
				timesheet := models.Timesheet{
					UserID:       req.UserID,
					ProjectID:    req.ProjectID,
					Date:         availableDays[i],
					Hours:        8,
					Content:      req.Content,
					BackfillLogID: &backfillLog.ID,
				}
				if err := tx.Create(&timesheet).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to backfill timesheets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "timesheets backfilled successfully"})
	}
}

func GetBackfillHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var history []models.BackfillLog
		if err := db.Preload("Operator").Preload("User").Preload("Project").Order("created_at desc").Find(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch backfill history"})
			return
		}
		c.JSON(http.StatusOK, history)
	}
}

func DeleteBackfill(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("backfill_log_id = ?", id).Delete(&models.Timesheet{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&models.BackfillLog{}, id).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete backfill"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "backfill deleted successfully"})
	}
}
