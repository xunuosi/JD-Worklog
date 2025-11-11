package handlers

import (
	"net/http"
	"time"

	"github.com/example/worklog-system/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateWorkPlanRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Content   string `json:"content"`
}

func CreateWorkPlan(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateWorkPlanRequest
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

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		workPlan := models.WorkPlan{
			UserID:    userID.(uint),
			ProjectID: req.ProjectID,
			StartDate: startDate,
			EndDate:   endDate,
			Content:   req.Content,
		}

		if err := db.Create(&workPlan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work plan"})
			return
		}
		c.JSON(http.StatusOK, workPlan)
	}
}

type GetWorkPlansRequest struct {
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

type GetWorkPlansByProjectRequest struct {
	ProjectID uint    `json:"project_id"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

func GetWorkPlans(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		role := models.Role(roleStr.(string))

		var req GetWorkPlansRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var workPlans []models.WorkPlan
		query := db.Model(&models.WorkPlan{})

		if role == models.RoleUser {
			query = query.Where("user_id = ?", userID)
		}

		now := time.Now()
		var parsedStartDate, parsedEndDate time.Time
		var err error

		if req.StartDate != nil && *req.StartDate != "" {
			parsedStartDate, err = time.Parse("2006-01-02", *req.StartDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
				return
			}
		} else {
			parsedStartDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		}

		if req.EndDate != nil && *req.EndDate != "" {
			parsedEndDate, err = time.Parse("2006-01-02", *req.EndDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
				return
			}
		} else {
			parsedEndDate = parsedStartDate.AddDate(0, 1, -1)
		}

		if parsedEndDate.Sub(parsedStartDate).Hours()/24 > 62 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Date range cannot exceed 2 months"})
			return
		}

		query = query.Where("start_date <= DATE(?) AND end_date >= DATE(?)", parsedEndDate, parsedStartDate)

		if err := query.Preload("User").Preload("Project").Find(&workPlans).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get work plans"})
			return
		}
		c.JSON(http.StatusOK, workPlans)
	}
}

func GetWorkPlansByProject(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		role := models.Role(roleStr.(string))

		var req GetWorkPlansByProjectRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var workPlans []models.WorkPlan
		query := db.Model(&models.WorkPlan{})

		query = query.Where("project_id = ?", req.ProjectID)

		if role == models.RoleUser {
			query = query.Where("user_id = ?", userID)
		}

		now := time.Now()
		var parsedStartDate, parsedEndDate time.Time
		var err error

		if req.StartDate != nil && *req.StartDate != "" {
			parsedStartDate, err = time.Parse("2006-01-02", *req.StartDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
				return
			}
		} else {
			parsedStartDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		}

		if req.EndDate != nil && *req.EndDate != "" {
			parsedEndDate, err = time.Parse("2006-01-02", *req.EndDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
				return
			}
		} else {
			parsedEndDate = parsedStartDate.AddDate(0, 1, -1)
		}

		if parsedEndDate.Sub(parsedStartDate).Hours()/24 > 62 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Date range cannot exceed 2 months"})
			return
		}

		query = query.Where("start_date <= DATE(?) AND end_date >= DATE(?)", parsedEndDate, parsedStartDate)

		if err := query.Preload("User").Preload("Project").Find(&workPlans).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get work plans by project"})
			return
		}
		c.JSON(http.StatusOK, workPlans)
	}
}

type UpdateWorkPlanRequest struct {
	ID        uint   `json:"id" binding:"required"`
	ProjectID uint   `json:"project_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Content   string `json:"content"`
}

func UpdateWorkPlan(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateWorkPlanRequest
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

		var workPlan models.WorkPlan
		if err := db.First(&workPlan, req.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work plan not found"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		role := models.Role(roleStr.(string))

		if role != models.RoleAdmin && workPlan.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		workPlan.ProjectID = req.ProjectID
		workPlan.StartDate = startDate
		workPlan.EndDate = endDate
		workPlan.Content = req.Content
		workPlan.UpdatedAt = time.Now()

		if err := db.Save(&workPlan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work plan"})
			return
		}
		c.JSON(http.StatusOK, workPlan)
	}
}

type DeleteWorkPlanRequest struct {
	ID uint `json:"id"`
}

func DeleteWorkPlan(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteWorkPlanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var workPlan models.WorkPlan
		if err := db.First(&workPlan, req.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work plan not found"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		role := models.Role(roleStr.(string))

		if role != models.RoleAdmin && workPlan.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		if err := db.Delete(&workPlan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete work plan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Work plan deleted successfully"})
	}
}
