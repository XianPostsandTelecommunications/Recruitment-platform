package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=20"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email"`
	Password  string         `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
	Role      string         `json:"role" gorm:"type:enum('student','admin');default:'student';not null"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Phone     string         `json:"phone" gorm:"size:20"`
	StudentID string         `json:"student_id" gorm:"size:20"`
	Major     string         `json:"major" gorm:"size:100"`
	Grade     string         `json:"grade" gorm:"size:20"`
	Status    string         `json:"status" gorm:"type:enum('active','inactive');default:'active';not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Applications []Application `json:"applications,omitempty" gorm:"foreignKey:UserID"`
	Labs         []Lab         `json:"labs,omitempty" gorm:"foreignKey:CreatedBy"`
	Notifications []Notification `json:"notifications,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// IsAdmin 判断是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsStudent 判断是否为学生
func (u *User) IsStudent() bool {
	return u.Role == "student"
}

// IsActive 判断是否激活
func (u *User) IsActive() bool {
	return u.Status == "active"
}



// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Username  string `json:"username" validate:"omitempty,min=3,max=20"`
	Phone     string `json:"phone" validate:"omitempty,len=11"`
	StudentID string `json:"student_id" validate:"omitempty"`
	Major     string `json:"major" validate:"omitempty"`
	Grade     string `json:"grade" validate:"omitempty"`
	Avatar    string `json:"avatar" validate:"omitempty"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Phone     string    `json:"phone"`
	StudentID string    `json:"student_id"`
	Major     string    `json:"major"`
	Grade     string    `json:"grade"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		StudentID: u.StudentID,
		Major:     u.Major,
		Grade:     u.Grade,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
} 