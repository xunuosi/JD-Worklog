package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/example/worklog-system/internal/middleware"
	"github.com/example/worklog-system/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	JWT *middleware.JWT
}

type loginReq struct{ Username, Password string }
type loginResp struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func hash(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(h[:])
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	var u models.User
	if err := h.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if u.Password != hash(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := h.JWT.Sign(u.ID, string(u.Role))
	c.JSON(http.StatusOK, loginResp{Token: token, Role: string(u.Role)})
}
