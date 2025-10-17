 package handlers

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"net/http"

	"github.com/example/worklog-system/internal/middleware"
	"github.com/example/worklog-system/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type TwoFactorAuthHandler struct {
	DB  *gorm.DB
	JWT *middleware.JWT
}

// Setup2FA generates a new TOTP secret for the user and returns it as a QR code.
func (h *TwoFactorAuthHandler) Setup(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "WorklogApp",
		AccountName: user.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate secret"})
		return
	}

	// Save the secret to the user's record
	if err := h.DB.Model(&user).Update("two_factor_secret", key.Secret()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save secret"})
		return
	}

	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create QR code image"})
		return
	}
	if err := png.Encode(&buf, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode QR code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
		"secret":  key.Secret(),
	})
}

type Verify2FARequest struct {
	Token string `json:"token" binding:"required"`
}

// Verify2FA enables 2FA for the user if the provided token is valid.
func (h *TwoFactorAuthHandler) Verify(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	valid := totp.Validate(req.Token, user.TwoFactorSecret)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Only enable 2FA, but do not change the required flag
	if err := h.DB.Model(&user).Update("two_factor_enabled", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enable 2FA"})
		return
	}

	// Issue a new, fully authenticated token
	token, _ := h.JWT.Sign(user.ID, string(user.Role), true)
	c.JSON(http.StatusOK, loginResp{Token: token, Role: string(user.Role)})
}
