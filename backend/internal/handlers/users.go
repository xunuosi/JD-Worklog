package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/example/worklog-system/internal/models"
)

type UsersHandler struct{ DB *gorm.DB }

type createUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type userResp struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Nickname string `json:"nickname"`
}

// POST /api/admin/users  -> 创建普通用户（禁止创建 admin），用户名判重
func (h *UsersHandler) Create(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username/password required"})
		return
	}
	// 判重
	var cnt int64
	h.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&cnt)
	if cnt > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}
	u := models.User{Username: req.Username, Password: hash(req.Password), Role: models.RoleUser, Nickname: req.Username}
	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userResp{ID: u.ID, Username: u.Username, Role: string(u.Role)})
}

// GET /api/admin/users -> 简单列表（不返回密码）
func (h *UsersHandler) List(c *gin.Context) {
	var users []models.User
	if err := h.DB.Order("id desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	resp := make([]userResp, 0, len(users))
	for _, u := range users {
		resp = append(resp, userResp{ID: u.ID, Username: u.Username, Role: string(u.Role), Nickname: u.Nickname})
	}
	c.JSON(http.StatusOK, resp)
}

// DELETE /api/admin/users/:id -> 软删除普通用户
func (h *UsersHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// 禁止删除管理员账号
	var u models.User
	if err := h.DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if u.Role == models.RoleAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不允许删除管理员账号"})
		return
	}

	// 可选：禁止删除自己（避免误操作）
	if uid := c.GetUint("user_id"); uid == u.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不允许删除当前登录账号"})
		return
	}

	if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
