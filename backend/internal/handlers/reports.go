package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	DB *gorm.DB
}

type ProjectHours struct {
	ProjectID   uint    `json:"project_id"`
	ProjectName string  `json:"project_name"`
	TotalHours  float64 `json:"total_hours"`
}

type reportReq struct {
	From string `json:"from"`
	To   string `json:"to"`
} // YYYY-MM-DD

func (h *ReportHandler) ProjectTotals(c *gin.Context) {
	var req reportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	from, err1 := time.Parse("2006-01-02", req.From)
	to, err2 := time.Parse("2006-01-02", req.To)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad date"})
		return
	}
	var rows []ProjectHours
	h.DB.Raw(`
		SELECT p.id as project_id, p.name as project_name, SUM(t.hours) as total_hours
		FROM timesheets t JOIN projects p ON t.project_id = p.id
		WHERE t.date BETWEEN ? AND ?
		GROUP BY p.id, p.name
		ORDER BY total_hours DESC
	`, from, to).Scan(&rows)
	c.JSON(http.StatusOK, rows)
}

// CSV 导出: GET /api/admin/reports/project-totals.csv?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *ReportHandler) ProjectTotalsCSV(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	from, err1 := time.Parse("2006-01-02", fromStr)
	to, err2 := time.Parse("2006-01-02", toStr)
	if err1 != nil || err2 != nil {
		c.String(http.StatusBadRequest, "bad date")
		return
	}
	var rows []ProjectHours
	h.DB.Raw(`
		SELECT p.id as project_id, p.name as project_name, SUM(t.hours) as total_hours
		FROM timesheets t JOIN projects p ON t.project_id = p.id
		WHERE t.date BETWEEN ? AND ?
		GROUP BY p.id, p.name
		ORDER BY total_hours DESC
	`, from, to).Scan(&rows)

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=project_totals.csv")
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()
	writer.Write([]string{"project_id", "project_name", "total_hours"})
	for _, r := range rows {
		writer.Write([]string{
			strconv.FormatUint(uint64(r.ProjectID), 10),
			r.ProjectName,
			strconv.FormatFloat(r.TotalHours, 'f', 2, 64),
		})
	}
}
