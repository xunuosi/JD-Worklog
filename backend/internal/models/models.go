package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"uniqueIndex;size:64" json:"username"`
	Password  string `json:"-"`
	Role      Role   `gorm:"size:16" json:"role"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Project struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:128;uniqueIndex" json:"name"`
	Desc      string `gorm:"type:text" json:"desc"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time
}

type Timesheet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ProjectID uint      `gorm:"index" json:"project_id"`
	Date      time.Time `gorm:"index" json:"date"`
	Hours     float32   `json:"hours"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	Project   Project `gorm:"foreignKey:ProjectID" json:"project"`
}
