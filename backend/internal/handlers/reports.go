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
	From   string `json:"from"`
	To     string `json:"to"`
	UserID *uint  `json:"user_id"` // ← 可为空；为空或 0 表示全体
} // YYYY-MM-DD

func (h *ReportHandler) ProjectTotals(c *gin.Context) {
	var req reportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	from, err1 := time.ParseInLocation("2006-01-02", req.From, time.Local)
	to, err2 := time.ParseInLocation("2006-01-02", req.To, time.Local)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad date"})
		return
	}

	var rows []ProjectHours
	if req.UserID != nil && *req.UserID != 0 {
		h.DB.Raw(`
            SELECT p.id as project_id, p.name as project_name, SUM(t.hours) as total_hours
            FROM timesheets t JOIN projects p ON t.project_id = p.id
            WHERE DATE(t.date) BETWEEN ? AND ? AND t.user_id = ?
            GROUP BY p.id, p.name
            ORDER BY total_hours DESC
        `, from, to, *req.UserID).Scan(&rows)
	} else {
		h.DB.Raw(`
            SELECT p.id as project_id, p.name as project_name, SUM(t.hours) as total_hours
            FROM timesheets t JOIN projects p ON t.project_id = p.id
            WHERE DATE(t.date) BETWEEN ? AND ?
            GROUP BY p.id, p.name
            ORDER BY total_hours DESC
        `, from, to).Scan(&rows)
	}
	c.JSON(http.StatusOK, rows)
}

// CSV 导出同样支持 user_id
func (h *ReportHandler) ProjectTotalsCSV(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	uidStr := c.Query("user_id")
	from, err1 := time.ParseInLocation("2006-01-02", fromStr, time.Local)
	to, err2 := time.ParseInLocation("2006-01-02", toStr, time.Local)
	if err1 != nil || err2 != nil {
		c.String(http.StatusBadRequest, "bad date")
		return
	}

	var rows []ProjectHours
	if uidStr != "" && uidStr != "0" {
		h.DB.Raw(`
            SELECT p.id as project_id, p.name as project_name, SUM(t.hours) as total_hours
            FROM timesheets t JOIN projects p ON t.project_id = p.id
            WHERE DATE(t.date) BETWEEN ? AND ? AND t.user_id = ?
            GROUP BY p.id, p.name
            ORDER BY total_hours DESC
        `, from, to, uidStr).Scan(&rows)
	} else {
		h.DB.Raw(`
            SELECT p.id as project_id, p.name as project_name, SUM(t.hours) as total_hours
            FROM timesheets t JOIN projects p ON t.project_id = p.id
            WHERE DATE(t.date) BETWEEN ? AND ?
            GROUP BY p.id, p.name
            ORDER BY total_hours DESC
        `, from, to).Scan(&rows)
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=project_totals.csv")
	w := csv.NewWriter(c.Writer)
	defer w.Flush()
	_ = w.Write([]string{"project_id", "project_name", "total_hours"})
	for _, r := range rows {
		_ = w.Write([]string{strconv.FormatUint(uint64(r.ProjectID), 10), r.ProjectName, strconv.FormatFloat(r.TotalHours, 'f', 2, 64)})
	}
}
