package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := models.User{Username: req.Username, Password: string(hash), Role: models.RoleUser, Nickname: req.Nickname}
	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userResp{ID: u.ID, Username: u.Username, Role: string(u.Role), Nickname: u.Nickname})
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

// 重置密码 API（管理员操作）— 默认为 root
func (h *UsersHandler) ResetPassword(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"user_id"`
		NewPassword string `json:"new_password"` // 可为空
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 默认密码
	if strings.TrimSpace(req.NewPassword) == "" {
		req.NewPassword = "root"
	}

	var user models.User
	if err := h.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 可选保护：不允许重置管理员
	// if user.Role == models.RoleAdmin {
	//     c.JSON(http.StatusBadRequest, gin.H{"error": "cannot reset admin password"})
	//     return
	// }

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := h.DB.Model(&user).Update("password", string(hash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully", "default": req.NewPassword})
}

// PUT /api/admin/users/:id/nickname
func (h *UsersHandler) UpdateNickname(c *gin.Context) {
    id := c.Param("id")

    var req struct{ Nickname string `json:"nickname"` }
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Nickname) == "" {
        c.JSON(400, gin.H{"error":"invalid nickname"}); return
    }

    // 可选：禁止改管理员昵称
    // var u models.User
    // if err := h.DB.First(&u, id).Error; err != nil { c.JSON(404, gin.H{"error":"user not found"}); return }
    // if u.Role == models.RoleAdmin { c.JSON(400, gin.H{"error":"cannot edit admin"}); return }

    if err := h.DB.Model(&models.User{}).Where("id = ?", id).
        Update("nickname", req.Nickname).Error; err != nil {
        c.JSON(500, gin.H{"error":"update failed"}); return
    }
    c.JSON(200, gin.H{"ok": true})
}
