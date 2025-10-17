package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	Role         string `json:"role"`
	Is2FAAuthed  bool   `json:"is_2fa_authed"`
	jwt.RegisteredClaims
}

func NewJWT(secret string) *JWT {
	return &JWT{secret: []byte(secret)}
}

type JWT struct { secret []byte }

func (j *JWT) Sign(userID uint, role string, is2FAAuthed bool) (string, error) {
	exp := time.Now().Add(7 * 24 * time.Hour) // Default long expiry
	if !is2FAAuthed {
		exp = time.Now().Add(5 * time.Minute) // Short expiry for pre-2FA token
	}

	claims := JWTClaims{
		UserID:      userID,
		Role:        role,
		Is2FAAuthed: is2FAAuthed,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) Parse(tokenStr string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	return claims, err
}

func (j *JWT) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := j.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("is_2fa_authed", claims.Is2FAAuthed)
		c.Next()
	}
}

func (j *JWT) TFARequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool("is_2fa_authed") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "2FA required"})
			return
		}
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
