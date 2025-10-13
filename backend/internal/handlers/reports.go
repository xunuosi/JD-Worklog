package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/xuri/excelize/v2"
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

// 导出 XLSX：两个 sheet，total_*（项目+用户汇总），detail_*（明细，列顺序：p.name, u.nickname, t.hours, t.content）
func (h *ReportHandler) ProjectExportXLSX(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	uidStr := c.Query("user_id")

	from, err1 := time.ParseInLocation("2006-01-02", fromStr, time.Local)
	to, err2 := time.ParseInLocation("2006-01-02", toStr, time.Local)
	if err1 != nil || err2 != nil {
		c.String(http.StatusBadRequest, "bad date")
		return
	}

	// 生成 sheet 名与文件名
	rng := fmt.Sprintf("%s_%s", from.Format("2006-01-02"), to.Format("2006-01-02"))
	totalSheet := "total_" + rng
	detailSheet := "detail_" + rng
	fname := fmt.Sprintf("worklog_%s.xlsx", rng)

	f := excelize.NewFile()
	// 默认 Sheet1 -> total_*
	defaultSheet := f.GetSheetName(0)
	_ = f.SetSheetName(defaultSheet, totalSheet)
	// 新建 detail_*
	_, _ = f.NewSheet(detailSheet)

	// ===== 1) total sheet：项目总工时（不含 nickname）=====
	_ = f.SetSheetRow(totalSheet, "A1", &[]string{"project_name", "contract_num", "total_hours"})

	type totalRow struct {
		ProjectName string
		ContractNum string
		TotalHours  float64
	}
	totals := []totalRow{}
	if uidStr != "" && uidStr != "0" {
		h.DB.Raw(`
            SELECT p.name AS project_name,
                   p.contract_num,
                   SUM(t.hours) AS total_hours
            FROM timesheets t
            JOIN projects p ON t.project_id = p.id
            WHERE DATE(t.date) BETWEEN ? AND ? AND t.user_id = ?
            GROUP BY p.name, p.contract_num
            ORDER BY total_hours DESC
        `, from, to, uidStr).Scan(&totals)
	} else {
		h.DB.Raw(`
            SELECT p.name AS project_name,
                   p.contract_num,
                   SUM(t.hours) AS total_hours
            FROM timesheets t
            JOIN projects p ON t.project_id = p.id
            WHERE DATE(t.date) BETWEEN ? AND ?
            GROUP BY p.name, p.contract_num
            ORDER BY total_hours DESC
        `, from, to).Scan(&totals)
	}
	for i, r := range totals {
		cell := fmt.Sprintf("A%d", i+2)
		_ = f.SetSheetRow(totalSheet, cell, &[]any{r.ProjectName, r.ContractNum, r.TotalHours})
	}

	// ===== 2) detail sheet：明细（按你要求的列顺序）=====
	// 列顺序：p.name, u.nickname, t.hours, t.content
	_ = f.SetSheetRow(detailSheet, "A1", &[]string{"project_name", "contract_num", "nickname", "hours", "content"})

	type detailRow struct {
		ProjectName string
		ContractNum string
		Nickname    string
		Hours       float64
		Content     string
	}
	details := []detailRow{}
	if uidStr != "" && uidStr != "0" {
		h.DB.Raw(`
            SELECT p.name AS project_name,
                   p.contract_num,
                   COALESCE(u.nickname, u.username) AS nickname,
                   t.hours,
                   t.content
            FROM timesheets t
            JOIN projects p ON t.project_id = p.id
            JOIN users    u ON t.user_id   = u.id
            WHERE DATE(t.date) BETWEEN ? AND ? AND t.user_id = ?
            ORDER BY t.date DESC, t.id DESC
        `, from, to, uidStr).Scan(&details)
	} else {
		h.DB.Raw(`
            SELECT p.name AS project_name,
                   p.contract_num,
                   COALESCE(u.nickname, u.username) AS nickname,
                   t.hours,
                   t.content
            FROM timesheets t
            JOIN projects p ON t.project_id = p.id
            JOIN users    u ON t.user_id   = u.id
            WHERE DATE(t.date) BETWEEN ? AND ?
            ORDER BY t.date DESC, t.id DESC
        `, from, to).Scan(&details)
	}
	for i, r := range details {
		cell := fmt.Sprintf("A%d", i+2)
		_ = f.SetSheetRow(detailSheet, cell, &[]any{r.ProjectName, r.ContractNum, r.Nickname, r.Hours, r.Content})
	}

	// 响应下载
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fname)
	c.Header("File-Name", fname)
	_ = f.Write(c.Writer)
}
