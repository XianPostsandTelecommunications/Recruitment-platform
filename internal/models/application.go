package models

import (
	"time"

	"gorm.io/gorm"
)

// Application 申请模型
type Application struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	LabID        uint           `json:"lab_id" gorm:"not null;index"`
	Motivation   string         `json:"motivation" gorm:"type:text"`
	Skills       StringSlice    `json:"skills" gorm:"type:json"`
	Experience   string         `json:"experience" gorm:"type:text"`
	AvailableTime string        `json:"available_time" gorm:"size:100"`
	ResumeURL    string         `json:"resume_url" gorm:"size:255"`
	Status       string         `json:"status" gorm:"type:enum('pending','accepted','rejected');default:'pending';not null;index"`
	Feedback     string         `json:"feedback" gorm:"type:text"`
	ReviewedBy   *uint          `json:"reviewed_by" gorm:"index"`
	ReviewedAt   *time.Time     `json:"reviewed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	User      User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Lab       Lab  `json:"lab,omitempty" gorm:"foreignKey:LabID"`
	Reviewer  *User `json:"reviewer,omitempty" gorm:"foreignKey:ReviewedBy"`
}

// TableName 指定表名
func (Application) TableName() string {
	return "applications"
}

// BeforeCreate 创建前的钩子
func (a *Application) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (a *Application) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}

// IsPending 判断是否待审核
func (a *Application) IsPending() bool {
	return a.Status == "pending"
}

// IsAccepted 判断是否已通过
func (a *Application) IsAccepted() bool {
	return a.Status == "accepted"
}

// IsRejected 判断是否已拒绝
func (a *Application) IsRejected() bool {
	return a.Status == "rejected"
}

// IsReviewed 判断是否已审核
func (a *Application) IsReviewed() bool {
	return a.Status == "accepted" || a.Status == "rejected"
}

// ApplicationCreateRequest 申请创建请求
type ApplicationCreateRequest struct {
	LabID         uint        `json:"lab_id" validate:"required"`
	Motivation    string      `json:"motivation" validate:"required"`
	Skills        StringSlice `json:"skills" validate:"required"`
	Experience    string      `json:"experience" validate:"omitempty"`
	AvailableTime string      `json:"available_time" validate:"omitempty"`
	ResumeURL     string      `json:"resume_url" validate:"omitempty"`
}

// ApplicationUpdateRequest 申请更新请求
type ApplicationUpdateRequest struct {
	Motivation    string      `json:"motivation" validate:"omitempty"`
	Skills        StringSlice `json:"skills" validate:"omitempty"`
	Experience    string      `json:"experience" validate:"omitempty"`
	AvailableTime string      `json:"available_time" validate:"omitempty"`
	ResumeURL     string      `json:"resume_url" validate:"omitempty"`
}

// ApplicationReviewRequest 申请审核请求
type ApplicationReviewRequest struct {
	Status   string `json:"status" validate:"required,oneof=accepted rejected"`
	Feedback string `json:"feedback" validate:"omitempty"`
}

// ApplicationResponse 申请响应
type ApplicationResponse struct {
	ID            uint         `json:"id"`
	UserID        uint         `json:"user_id"`
	LabID         uint         `json:"lab_id"`
	Motivation    string       `json:"motivation"`
	Skills        StringSlice  `json:"skills"`
	Experience    string       `json:"experience"`
	AvailableTime string       `json:"available_time"`
	ResumeURL     string       `json:"resume_url"`
	Status        string       `json:"status"`
	Feedback      string       `json:"feedback"`
	ReviewedBy    *uint        `json:"reviewed_by"`
	ReviewedAt    *time.Time   `json:"reviewed_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`

	// 关联数据
	User     *UserResponse `json:"user,omitempty"`
	Lab      *LabResponse  `json:"lab,omitempty"`
	Reviewer *UserResponse `json:"reviewer,omitempty"`
}

// ToResponse 转换为响应格式
func (a *Application) ToResponse() *ApplicationResponse {
	response := &ApplicationResponse{
		ID:            a.ID,
		UserID:        a.UserID,
		LabID:         a.LabID,
		Motivation:    a.Motivation,
		Skills:        a.Skills,
		Experience:    a.Experience,
		AvailableTime: a.AvailableTime,
		ResumeURL:     a.ResumeURL,
		Status:        a.Status,
		Feedback:      a.Feedback,
		ReviewedBy:    a.ReviewedBy,
		ReviewedAt:    a.ReviewedAt,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}

	// 如果有用户信息，转换为响应格式
	if a.User.ID != 0 {
		response.User = a.User.ToResponse()
	}

	// 如果有实验室信息，转换为响应格式
	if a.Lab.ID != 0 {
		response.Lab = a.Lab.ToResponse()
	}

	// 如果有审核者信息，转换为响应格式
	if a.Reviewer != nil && a.Reviewer.ID != 0 {
		response.Reviewer = a.Reviewer.ToResponse()
	}

	return response
}

// ApplicationListResponse 申请列表响应
type ApplicationListResponse struct {
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	List  []ApplicationResponse `json:"list"`
}

// ApplicationStats 申请统计
type ApplicationStats struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Accepted  int64 `json:"accepted"`
	Rejected  int64 `json:"rejected"`
} 