package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/example/worklog-system/internal/middleware"
	"github.com/example/worklog-system/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

func hashSHA256(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(h[:])
}

// 判断是否为 bcrypt 哈希
func isBcryptHash(s string) bool {
	return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
}

// 校验明文密码与存储哈希（兼容 bcrypt 与旧 sha256）
func checkPassword(stored, plain string) bool {
	if isBcryptHash(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil
	}
	return stored == hashSHA256(plain)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	var u models.User
	if err := h.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	if !checkPassword(u.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := h.JWT.Sign(u.ID, string(u.Role))
	c.JSON(http.StatusOK, loginResp{Token: token, Role: string(u.Role)})
}
