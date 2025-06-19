package models

import (
	"time"

	"gorm.io/gorm"
)

// InterviewApplication 面试申请模型
type InterviewApplication struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"size:100;not null"`
	Email            string         `json:"email" gorm:"size:100;not null;index"`
	Phone            string         `json:"phone" gorm:"size:20;not null"`
	StudentID        string         `json:"student_id" gorm:"size:50;not null"`
	Major            string         `json:"major" gorm:"size:100;not null"`
	Grade            string         `json:"grade" gorm:"size:20;not null"`
	InterviewTime    string         `json:"interview_time" gorm:"size:100;not null"`
	Status           string         `json:"status" gorm:"type:enum('pending','interviewed','passed','rejected');default:'pending';not null;index"`
	AdminRemarks     string         `json:"admin_remarks" gorm:"type:text"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (InterviewApplication) TableName() string {
	return "interview_applications"
}

// BeforeCreate 创建前的钩子
func (ia *InterviewApplication) BeforeCreate(tx *gorm.DB) error {
	ia.CreatedAt = time.Now()
	ia.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (ia *InterviewApplication) BeforeUpdate(tx *gorm.DB) error {
	ia.UpdatedAt = time.Now()
	return nil
}

// IsPending 判断是否待面试
func (ia *InterviewApplication) IsPending() bool {
	return ia.Status == "pending"
}

// IsInterviewed 判断是否已面试
func (ia *InterviewApplication) IsInterviewed() bool {
	return ia.Status == "interviewed"
}

// IsPassed 判断是否通过
func (ia *InterviewApplication) IsPassed() bool {
	return ia.Status == "passed"
}

// IsRejected 判断是否被拒绝
func (ia *InterviewApplication) IsRejected() bool {
	return ia.Status == "rejected"
}

// InterviewApplicationResponse 面试申请响应
type InterviewApplicationResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	StudentID     string    `json:"student_id"`
	Major         string    `json:"major"`
	Grade         string    `json:"grade"`
	InterviewTime string    `json:"interview_time"`
	Status        string    `json:"status"`
	AdminRemarks  string    `json:"admin_remarks"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (ia *InterviewApplication) ToResponse() *InterviewApplicationResponse {
	return &InterviewApplicationResponse{
		ID:            ia.ID,
		Name:          ia.Name,
		Email:         ia.Email,
		Phone:         ia.Phone,
		StudentID:     ia.StudentID,
		Major:         ia.Major,
		Grade:         ia.Grade,
		InterviewTime: ia.InterviewTime,
		Status:        ia.Status,
		AdminRemarks:  ia.AdminRemarks,
		CreatedAt:     ia.CreatedAt,
		UpdatedAt:     ia.UpdatedAt,
	}
}

// InterviewApplicationUpdateRequest 面试申请更新请求
type InterviewApplicationUpdateRequest struct {
	Status       string `json:"status" validate:"required,oneof=pending interviewed passed rejected"`
	AdminRemarks string `json:"admin_remarks" validate:"omitempty"`
}

// InterviewApplicationListResponse 面试申请列表响应
type InterviewApplicationListResponse struct {
	Total int64                           `json:"total"`
	Page  int                             `json:"page"`
	Size  int                             `json:"size"`
	List  []InterviewApplicationResponse  `json:"list"`
}

// InterviewApplicationStats 面试申请统计
type InterviewApplicationStats struct {
	Total       int64 `json:"total"`
	Pending     int64 `json:"pending"`
	Interviewed int64 `json:"interviewed"`
	Passed      int64 `json:"passed"`
	Rejected    int64 `json:"rejected"`
} 