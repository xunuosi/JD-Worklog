package seed

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"gorm.io/gorm"
	"github.com/example/worklog-system/internal/models"
)

func hash(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(h[:])
}

func Run(db *gorm.DB) {
	var cnt int64
	db.Model(&models.User{}).Count(&cnt)
	if cnt == 0 {
		log.Println("seeding users and projects...")
		_ = db.Create(&models.User{Username: "admin", Password: hash("admin123"), Role: models.RoleAdmin}).Error
		_ = db.Create(&models.User{Username: "alice", Password: hash("alice123"), Role: models.RoleUser}).Error
		_ = db.Create(&models.Project{Name: "演示项目A", Desc: "示例描述"}).Error
	}
}
