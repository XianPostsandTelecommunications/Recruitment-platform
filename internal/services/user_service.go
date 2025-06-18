package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"lab-recruitment-platform/internal/config"
	"lab-recruitment-platform/internal/models"
	"lab-recruitment-platform/pkg/logger"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		db: config.GetDB(),
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *models.UserRegisterRequest) (*models.User, error) {
	// 检查邮箱是否已存在
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 检查用户名是否已存在
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("密码加密失败: %v", err)
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if user.Role == "" {
		user.Role = "student"
	}

	if err := s.db.Create(user).Error; err != nil {
		logger.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Username != "" {
		// 检查用户名是否已被其他用户使用
		var existingUser models.User
		if err := s.db.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("用户名已存在")
		}
		user.Username = req.Username
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.StudentID != "" {
		user.StudentID = req.StudentID
	}
	if req.Major != "" {
		user.Major = req.Major
	}
	if req.Grade != "" {
		user.Grade = req.Grade
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.db.Save(user).Error; err != nil {
		logger.Errorf("更新用户失败: %v", err)
		return nil, errors.New("更新用户失败")
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(user).Error; err != nil {
		logger.Errorf("删除用户失败: %v", err)
		return errors.New("删除用户失败")
	}

	return nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(page, size int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})

	// 搜索条件
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR student_id LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("密码加密失败: %v", err)
		return errors.New("密码加密失败")
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := s.db.Save(user).Error; err != nil {
		logger.Errorf("修改密码失败: %v", err)
		return errors.New("修改密码失败")
	}

	return nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(id uint, newPassword string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("密码加密失败: %v", err)
		return errors.New("密码加密失败")
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := s.db.Save(user).Error; err != nil {
		logger.Errorf("重置密码失败: %v", err)
		return errors.New("重置密码失败")
	}

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(id uint, status string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Status = status
	if err := s.db.Save(user).Error; err != nil {
		logger.Errorf("更新用户状态失败: %v", err)
		return errors.New("更新用户状态失败")
	}

	return nil
}

// GetUserStats 获取用户统计信息
func (s *UserService) GetUserStats() (map[string]int64, error) {
	var stats = make(map[string]int64)

	// 总用户数
	var total int64
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// 活跃用户数
	var active int64
	if err := s.db.Model(&models.User{}).Where("status = ?", "active").Count(&active).Error; err != nil {
		return nil, err
	}
	stats["active"] = active

	// 学生用户数
	var students int64
	if err := s.db.Model(&models.User{}).Where("role = ?", "student").Count(&students).Error; err != nil {
		return nil, err
	}
	stats["students"] = students

	// 管理员用户数
	var admins int64
	if err := s.db.Model(&models.User{}).Where("role = ?", "admin").Count(&admins).Error; err != nil {
		return nil, err
	}
	stats["admins"] = admins

	return stats, nil
}

// VerifyPassword 验证密码
func (s *UserService) VerifyPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
} 