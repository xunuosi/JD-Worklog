package handlers

import (
	"net/http"

	"github.com/example/worklog-system/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct{ DB *gorm.DB }
type projectReq struct {
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	ContractNum string `json:"contract_num"`
	IsActive    *bool  `json:"is_active"`
}

func (h *ProjectHandler) List(c *gin.Context) {
	var ps []models.Project
	if err := h.DB.Order("id desc").Where("is_active = ?", true).Find(&ps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	c.JSON(http.StatusOK, ps)
}

func (h *ProjectHandler) AllList(c *gin.Context) {
	var ps []models.Project
	if err := h.DB.Order("id desc").Find(&ps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	c.JSON(http.StatusOK, ps)
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req projectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	// 判重：检查是否已存在同名项目
	var count int64
	h.DB.Model(&models.Project{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目名称已存在"})
		return
	}
	p := models.Project{Name: req.Name, Desc: req.Desc, ContractNum: req.ContractNum}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req projectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	var p models.Project
	if err := h.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	p.Desc = req.Desc
	p.ContractNum = req.ContractNum
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	if err := h.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
