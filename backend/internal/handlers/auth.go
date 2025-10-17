package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/example/worklog-system/internal/middleware"
	"github.com/example/worklog-system/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	JWT *middleware.JWT
}

type loginReq struct{ Username, Password string }
type loginResp struct {
	Token             string `json:"token"`
	Role              string `json:"role"`
	TwoFactorRequired bool   `json:"two_factor_required"`
	Force2FASetup     bool   `json:"force_2fa_setup"`
}

type login2FAReq struct {
	Token string `json:"token" binding:"required"`
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

	if u.TwoFactorRequiredByAdmin && !u.TwoFactorEnabled {
		// Issue a temporary token to force setup
		token, _ := h.JWT.Sign(u.ID, string(u.Role), false)
		c.JSON(http.StatusOK, loginResp{Token: token, Role: string(u.Role), Force2FASetup: true})
		return
	}

	if u.TwoFactorEnabled {
		// Issue a temporary token for 2FA verification
		token, _ := h.JWT.Sign(u.ID, string(u.Role), false)
		c.JSON(http.StatusOK, loginResp{Token: token, Role: string(u.Role), TwoFactorRequired: true})
	} else {
		// Issue a final token
		token, _ := h.JWT.Sign(u.ID, string(u.Role), true)
		c.JSON(http.StatusOK, loginResp{Token: token, Role: string(u.Role)})
	}
}

func (h *AuthHandler) Login2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req login2FAReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	var u models.User
	if err := h.DB.First(&u, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if !u.TwoFactorEnabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "2FA not enabled"})
		return
	}

	valid := totp.Validate(req.Token, u.TwoFactorSecret)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid 2FA token"})
		return
	}

	// Issue a final token
	token, _ := h.JWT.Sign(u.ID, string(u.Role), true)
	c.JSON(http.StatusOK, loginResp{Token: token, Role: string(u.Role)})
}
