package services

import (
	"errors"

	"gorm.io/gorm"
	"lab-recruitment-platform/internal/config"
	"lab-recruitment-platform/internal/models"
	"lab-recruitment-platform/pkg/logger"
)

// InterviewApplicationService 面试申请服务
type InterviewApplicationService struct {
	db *gorm.DB
}

// NewInterviewApplicationService 创建面试申请服务实例
func NewInterviewApplicationService() *InterviewApplicationService {
	return &InterviewApplicationService{
		db: config.GetDB(),
	}
}

// CreateApplication 创建面试申请
func (s *InterviewApplicationService) CreateApplication(name, email, phone, studentID, major, grade, interviewTime string) (*models.InterviewApplication, error) {
	// 检查邮箱是否已经申请过
	var existingApp models.InterviewApplication
	if err := s.db.Where("email = ?", email).First(&existingApp).Error; err == nil {
		return nil, errors.New("该邮箱已提交过申请，请勿重复申请")
	}

	// 创建新申请
	application := &models.InterviewApplication{
		Name:          name,
		Email:         email,
		Phone:         phone,
		StudentID:     studentID,
		Major:         major,
		Grade:         grade,
		InterviewTime: interviewTime,
		Status:        "pending",
	}

	if err := s.db.Create(application).Error; err != nil {
		logger.Errorf("创建面试申请失败: %v", err)
		return nil, errors.New("创建面试申请失败")
	}

	logger.Infof("面试申请创建成功: ID=%d, 姓名=%s, 邮箱=%s", application.ID, application.Name, application.Email)
	return application, nil
}

// GetApplicationByID 根据ID获取面试申请
func (s *InterviewApplicationService) GetApplicationByID(id uint) (*models.InterviewApplication, error) {
	var application models.InterviewApplication
	if err := s.db.First(&application, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("面试申请不存在")
		}
		return nil, err
	}
	return &application, nil
}

// GetApplicationByEmail 根据邮箱获取面试申请
func (s *InterviewApplicationService) GetApplicationByEmail(email string) (*models.InterviewApplication, error) {
	var application models.InterviewApplication
	if err := s.db.Where("email = ?", email).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("面试申请不存在")
		}
		return nil, err
	}
	return &application, nil
}

// ListApplications 获取面试申请列表
func (s *InterviewApplicationService) ListApplications(page, size int, status, name string) (*models.InterviewApplicationListResponse, error) {
	var applications []models.InterviewApplication
	var total int64

	query := s.db.Model(&models.InterviewApplication{})

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 姓名搜索（支持模糊匹配）
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]models.InterviewApplicationResponse, len(applications))
	for i, app := range applications {
		list[i] = *app.ToResponse()
	}

	return &models.InterviewApplicationListResponse{
		Total: total,
		Page:  page,
		Size:  size,
		List:  list,
	}, nil
}

// UpdateApplication 更新面试申请状态
func (s *InterviewApplicationService) UpdateApplication(id uint, status, adminRemarks string) (*models.InterviewApplication, error) {
	var application models.InterviewApplication
	if err := s.db.First(&application, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("面试申请不存在")
		}
		return nil, err
	}

	// 更新状态和备注
	application.Status = status
	application.AdminRemarks = adminRemarks

	if err := s.db.Save(&application).Error; err != nil {
		logger.Errorf("更新面试申请失败: %v", err)
		return nil, errors.New("更新面试申请失败")
	}

	logger.Infof("面试申请更新成功: ID=%d, 状态=%s", application.ID, application.Status)
	return &application, nil
}

// DeleteApplication 删除面试申请
func (s *InterviewApplicationService) DeleteApplication(id uint) error {
	if err := s.db.Delete(&models.InterviewApplication{}, id).Error; err != nil {
		logger.Errorf("删除面试申请失败: %v", err)
		return errors.New("删除面试申请失败")
	}

	logger.Infof("面试申请删除成功: ID=%d", id)
	return nil
}

// GetApplicationStats 获取面试申请统计
func (s *InterviewApplicationService) GetApplicationStats() (*models.InterviewApplicationStats, error) {
	var stats models.InterviewApplicationStats

	// 总数
	s.db.Model(&models.InterviewApplication{}).Count(&stats.Total)

	// 各状态数量
	s.db.Model(&models.InterviewApplication{}).Where("status = ?", "pending").Count(&stats.Pending)
	s.db.Model(&models.InterviewApplication{}).Where("status = ?", "interviewed").Count(&stats.Interviewed)
	s.db.Model(&models.InterviewApplication{}).Where("status = ?", "passed").Count(&stats.Passed)
	s.db.Model(&models.InterviewApplication{}).Where("status = ?", "rejected").Count(&stats.Rejected)

	return &stats, nil
} 