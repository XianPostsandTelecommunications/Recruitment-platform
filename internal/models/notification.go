package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification 通知模型
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Type      string         `json:"type" gorm:"type:enum('system','application','lab');default:'system';not null"`
	IsRead    bool           `json:"is_read" gorm:"default:false;not null"`
	RelatedID *uint          `json:"related_id" gorm:"index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate 创建前的钩子
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (n *Notification) BeforeUpdate(tx *gorm.DB) error {
	n.UpdatedAt = time.Now()
	return nil
}

// IsSystemNotification 判断是否为系统通知
func (n *Notification) IsSystemNotification() bool {
	return n.Type == "system"
}

// IsApplicationNotification 判断是否为申请通知
func (n *Notification) IsApplicationNotification() bool {
	return n.Type == "application"
}

// IsLabNotification 判断是否为实验室通知
func (n *Notification) IsLabNotification() bool {
	return n.Type == "lab"
}

// MarkAsRead 标记为已读
func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

// NotificationCreateRequest 通知创建请求
type NotificationCreateRequest struct {
	UserID    uint   `json:"user_id" validate:"required"`
	Title     string `json:"title" validate:"required,max=200"`
	Content   string `json:"content" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=system application lab"`
	RelatedID *uint  `json:"related_id" validate:"omitempty"`
}

// NotificationUpdateRequest 通知更新请求
type NotificationUpdateRequest struct {
	Title   string `json:"title" validate:"omitempty,max=200"`
	Content string `json:"content" validate:"omitempty"`
	IsRead  *bool  `json:"is_read" validate:"omitempty"`
}

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	RelatedID *uint     `json:"related_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联数据
	User *UserResponse `json:"user,omitempty"`
}

// ToResponse 转换为响应格式
func (n *Notification) ToResponse() *NotificationResponse {
	response := &NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Content:   n.Content,
		Type:      n.Type,
		IsRead:    n.IsRead,
		RelatedID: n.RelatedID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	// 如果有用户信息，转换为响应格式
	if n.User.ID != 0 {
		response.User = n.User.ToResponse()
	}

	return response
}

// NotificationListResponse 通知列表响应
type NotificationListResponse struct {
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	List  []NotificationResponse `json:"list"`
}

// NotificationStats 通知统计
type NotificationStats struct {
	Total   int64 `json:"total"`
	Unread  int64 `json:"unread"`
	Read    int64 `json:"read"`
} 