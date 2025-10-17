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
	Nickname  string `json:"nickname"`
	Password  string `json:"-"`
	Role      Role   `gorm:"size:16" json:"role"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Project struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:128;uniqueIndex" json:"name"`
	Desc        string `gorm:"type:text" json:"desc"`
	ContractNum string `gorm:"size:128" json:"contract_num"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time
}

type Timesheet struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	ProjectID    uint      `gorm:"index" json:"project_id"`
	Date         time.Time `gorm:"type:date;index" json:"date"`
	Hours        float32   `json:"hours"`
	Content      string    `gorm:"type:text" json:"content"`
	BackfillLogID *uint     `gorm:"index" json:"backfill_log_id,omitempty"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	User         User `gorm:"foreignKey:UserID" json:"user"`
	Project      Project `gorm:"foreignKey:ProjectID" json:"project"`
}

type BackfillLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ProjectID uint      `gorm:"index" json:"project_id"`
	TotalDays float32   `json:"total_days"`
	StartDate time.Time `gorm:"type:date" json:"start_date"`
	EndDate   time.Time `gorm:"type:date" json:"end_date"`
	CreatedAt time.Time
	Operator  User `gorm:"foreignKey:AdminID" json:"operator"`
	User      User `gorm:"foreignKey:UserID" json:"user"`
	Project   Project `gorm:"foreignKey:ProjectID" json:"project"`
}

