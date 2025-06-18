package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// StringSlice 字符串切片类型，用于JSON存储
type StringSlice []string

// Value 实现 driver.Valuer 接口
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("cannot scan non-string value into StringSlice")
	}
}

// Lab 实验室模型
type Lab struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:100;not null" validate:"required,min=2,max=100"`
	Description    string         `json:"description" gorm:"type:text"`
	Requirements   string         `json:"requirements" gorm:"type:text"`
	MaxMembers     int            `json:"max_members" gorm:"default:10;not null" validate:"required,min=1,max=100"`
	CurrentMembers int            `json:"current_members" gorm:"default:0;not null"`
	ContactEmail   string         `json:"contact_email" gorm:"size:100" validate:"omitempty,email"`
	ContactPhone   string         `json:"contact_phone" gorm:"size:20"`
	Location       string         `json:"location" gorm:"size:200"`
	Tags           StringSlice    `json:"tags" gorm:"type:json"`
	CoverImage     string         `json:"cover_image" gorm:"size:255"`
	Status         string         `json:"status" gorm:"type:enum('active','inactive');default:'active';not null"`
	CreatedBy      uint           `json:"created_by" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Creator      User           `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Applications []Application  `json:"applications,omitempty" gorm:"foreignKey:LabID"`
}

// TableName 指定表名
func (Lab) TableName() string {
	return "labs"
}

// BeforeCreate 创建前的钩子
func (l *Lab) BeforeCreate(tx *gorm.DB) error {
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (l *Lab) BeforeUpdate(tx *gorm.DB) error {
	l.UpdatedAt = time.Now()
	return nil
}

// IsActive 判断是否激活
func (l *Lab) IsActive() bool {
	return l.Status == "active"
}

// IsFull 判断是否已满员
func (l *Lab) IsFull() bool {
	return l.CurrentMembers >= l.MaxMembers
}

// GetAvailableSlots 获取可用名额
func (l *Lab) GetAvailableSlots() int {
	available := l.MaxMembers - l.CurrentMembers
	if available < 0 {
		return 0
	}
	return available
}

// LabCreateRequest 实验室创建请求
type LabCreateRequest struct {
	Name           string      `json:"name" validate:"required,min=2,max=100"`
	Description    string      `json:"description" validate:"required"`
	Requirements   string      `json:"requirements" validate:"required"`
	MaxMembers     int         `json:"max_members" validate:"required,min=1,max=100"`
	ContactEmail   string      `json:"contact_email" validate:"required,email"`
	ContactPhone   string      `json:"contact_phone" validate:"omitempty"`
	Location       string      `json:"location" validate:"omitempty"`
	Tags           StringSlice `json:"tags" validate:"omitempty"`
	CoverImage     string      `json:"cover_image" validate:"omitempty"`
}

// LabUpdateRequest 实验室更新请求
type LabUpdateRequest struct {
	Name           string      `json:"name" validate:"omitempty,min=2,max=100"`
	Description    string      `json:"description" validate:"omitempty"`
	Requirements   string      `json:"requirements" validate:"omitempty"`
	MaxMembers     int         `json:"max_members" validate:"omitempty,min=1,max=100"`
	ContactEmail   string      `json:"contact_email" validate:"omitempty,email"`
	ContactPhone   string      `json:"contact_phone" validate:"omitempty"`
	Location       string      `json:"location" validate:"omitempty"`
	Tags           StringSlice `json:"tags" validate:"omitempty"`
	CoverImage     string      `json:"cover_image" validate:"omitempty"`
	Status         string      `json:"status" validate:"omitempty,oneof=active inactive"`
}

// LabResponse 实验室响应
type LabResponse struct {
	ID             uint         `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Requirements   string       `json:"requirements"`
	MaxMembers     int          `json:"max_members"`
	CurrentMembers int          `json:"current_members"`
	ContactEmail   string       `json:"contact_email"`
	ContactPhone   string       `json:"contact_phone"`
	Location       string       `json:"location"`
	Tags           StringSlice  `json:"tags"`
	CoverImage     string       `json:"cover_image"`
	Status         string       `json:"status"`
	CreatedBy      uint         `json:"created_by"`
	Creator        *UserResponse `json:"creator,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (l *Lab) ToResponse() *LabResponse {
	response := &LabResponse{
		ID:             l.ID,
		Name:           l.Name,
		Description:    l.Description,
		Requirements:   l.Requirements,
		MaxMembers:     l.MaxMembers,
		CurrentMembers: l.CurrentMembers,
		ContactEmail:   l.ContactEmail,
		ContactPhone:   l.ContactPhone,
		Location:       l.Location,
		Tags:           l.Tags,
		CoverImage:     l.CoverImage,
		Status:         l.Status,
		CreatedBy:      l.CreatedBy,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
	}

	// 如果有创建者信息，转换为响应格式
	if l.Creator.ID != 0 {
		response.Creator = l.Creator.ToResponse()
	}

	return response
}

// LabListResponse 实验室列表响应
type LabListResponse struct {
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
	List  []LabResponse `json:"list"`
} 