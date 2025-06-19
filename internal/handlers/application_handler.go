package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"lab-recruitment-platform/internal/services"
	"lab-recruitment-platform/pkg/logger"
	"lab-recruitment-platform/pkg/response"
	"lab-recruitment-platform/pkg/validator"
)

// ApplicationHandler 申请处理器
type ApplicationHandler struct {
	interviewService *services.InterviewApplicationService
}

// NewApplicationHandler 创建申请处理器实例
func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		interviewService: services.NewInterviewApplicationService(),
	}
}

// VerificationCode 验证码存储
type VerificationCode struct {
	Code      string
	Email     string
	ExpiresAt time.Time
}

// 内存存储验证码
var verificationCodes = make(map[string]*VerificationCode)

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ApplyRequest 申请请求
type ApplyRequest struct {
	Name             string `json:"name" validate:"required"`
	Email            string `json:"email" validate:"required,email"`
	Phone            string `json:"phone" validate:"required"`
	StudentID        string `json:"student_id" validate:"required"`
	Major            string `json:"major" validate:"required"`
	Grade            string `json:"grade" validate:"required"`
	InterviewTime    string `json:"interview_time" validate:"required"`
	VerificationCode string `json:"verification_code" validate:"required"`
}

// SendCode 发送验证码
// @Summary 发送验证码
// @Description 向指定邮箱发送验证码
// @Tags 申请
// @Accept json
// @Produce json
// @Param request body SendCodeRequest true "邮箱信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /send-code [post]
func (h *ApplicationHandler) SendCode(c *gin.Context) {
	// 检查Content-Type
	contentType := c.GetHeader("Content-Type")
	if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
		response.BadRequest(c, "Content-Type必须为application/json")
		return
	}

	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 验证请求参数
	if !validator.ValidateRequest(c, &req) {
		return
	}

	// 生成6位随机验证码
	code, err := generateVerificationCode()
	if err != nil {
		logger.Errorf("生成验证码失败: %v", err)
		response.InternalServerError(c, "生成验证码失败")
		return
	}

	// 存储验证码
	verificationCodes[req.Email] = &VerificationCode{
		Code:      code,
		Email:     req.Email,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟有效期
	}

	// 发送邮件
	if err := sendVerificationEmail(req.Email, code); err != nil {
		logger.Errorf("发送验证码邮件失败: %v", err)
		// 发送失败，删除验证码
		delete(verificationCodes, req.Email)
		response.InternalServerError(c, "邮件发送失败，请稍后重试")
		return
	}

	response.SuccessWithMessage(c, "验证码已发送到邮箱，请注意查收", nil)
}

// Apply 提交申请
// @Summary 提交面试申请
// @Description 提交面试申请信息
// @Tags 申请
// @Accept json
// @Produce json
// @Param request body ApplyRequest true "申请信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /apply [post]
func (h *ApplicationHandler) Apply(c *gin.Context) {
	logger.Infof("开始处理面试申请")
	
	// 检查Content-Type
	contentType := c.GetHeader("Content-Type")
	logger.Infof("Content-Type: %s", contentType)
	if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
		logger.Errorf("Content-Type验证失败: %s", contentType)
		response.BadRequest(c, "Content-Type必须为application/json")
		return
	}

	var req ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("JSON绑定失败: %v", err)
		response.BadRequest(c, "请求参数错误")
		return
	}
	logger.Infof("请求参数: name=%s, email=%s, verification_code=%s", req.Name, req.Email, req.VerificationCode)

	// 验证请求参数
	if !validator.ValidateRequest(c, &req) {
		logger.Errorf("参数验证失败")
		return
	}
	logger.Infof("参数验证通过")

	// 验证验证码（特殊处理测试验证码）
	if req.VerificationCode == "999999" {
		logger.Infof("使用测试验证码，跳过验证: email=%s", req.Email)
	} else if !verifyCode(req.Email, req.VerificationCode) {
		logger.Errorf("验证码验证失败: email=%s, code=%s", req.Email, req.VerificationCode)
		response.BadRequest(c, "验证码错误或已过期")
		return
	}
	logger.Infof("验证码验证通过")

	// 保存申请到数据库
	application, err := h.interviewService.CreateApplication(
		req.Name, req.Email, req.Phone, req.StudentID, 
		req.Major, req.Grade, req.InterviewTime,
	)
	if err != nil {
		logger.Errorf("保存面试申请失败: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	// 发送申请成功邮件
	if err := sendApplicationSuccessEmail(req.Email, req.Name); err != nil {
		logger.Errorf("发送申请成功邮件失败: %v", err)
	}

	response.SuccessWithMessage(c, "🎉 申请提交成功！我们会尽快联系您安排面试", gin.H{
		"application_id": application.ID,
		"application":    application.ToResponse(),
	})
}

// generateVerificationCode 生成6位随机验证码
func generateVerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// verifyCode 验证验证码
func verifyCode(email, code string) bool {
	logger.Infof("验证码验证开始: email=%s, code=%s", email, code)
	
	// 测试环境支持：验证码 999999 用于测试
	if code == "999999" {
		logger.Infof("使用测试验证码进行验证: %s", email)
		return true
	}

	vc, exists := verificationCodes[email]
	if !exists {
		logger.Infof("验证码不存在: email=%s", email)
		return false
	}

	// 检查验证码是否正确且未过期
	if vc.Code == code && time.Now().Before(vc.ExpiresAt) {
		logger.Infof("验证码验证成功: email=%s", email)
		delete(verificationCodes, email) // 使用后立即删除
		return true
	}

	// 验证码错误或过期，删除
	logger.Infof("验证码错误或过期: email=%s, expected=%s, received=%s", email, vc.Code, code)
	delete(verificationCodes, email)
	return false
}

// sendVerificationEmail 发送验证码邮件
func sendVerificationEmail(email, code string) error {
	// SMTP配置
	host := "smtp.qq.com"
	port := 587
	user := "1785260184@qq.com"
	password := "yccnobsfrnoncaci"

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "实验室面试申请验证码")

	// HTML邮件内容
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>验证码</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 20px; border-radius: 0 0 8px 8px; }
        .code { background: #1890ff; color: white; padding: 10px 20px; font-size: 24px; font-weight: bold; text-align: center; border-radius: 4px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧪 EPI实验室面试申请</h1>
        </div>
        <div class="content">
            <p>您好！</p>
            <p>您正在申请EPI实验室面试，请使用以下验证码完成验证：</p>
            <div class="code">%s</div>
            <p><strong>验证码有效期：5分钟</strong></p>
            <p>如果这不是您的操作，请忽略此邮件。</p>
            <p>感谢您的关注！</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复</p>
        </div>
    </div>
</body>
</html>
`, code)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(host, port, user, password)
	return d.DialAndSend(m)
}

// ListApplications 获取面试申请列表（管理员接口）
// @Summary 获取面试申请列表
// @Description 获取面试申请列表，支持分页、状态过滤和姓名搜索
// @Tags 管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param status query string false "状态过滤" Enums(pending,interviewed,passed,rejected)
// @Param name query string false "姓名搜索"
// @Success 200 {object} response.Response{data=models.InterviewApplicationListResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/applications [get]
func (h *ApplicationHandler) ListApplications(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")
	name := c.Query("name")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 获取申请列表
	result, err := h.interviewService.ListApplications(page, size, status, name)
	if err != nil {
		logger.Errorf("获取面试申请列表失败: %v", err)
		response.InternalServerError(c, "获取申请列表失败")
		return
	}

	response.Success(c, result)
}

// GetApplication 获取单个面试申请详情（管理员接口）
// @Summary 获取面试申请详情
// @Description 根据ID获取面试申请详情
// @Tags 管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "申请ID"
// @Success 200 {object} response.Response{data=models.InterviewApplicationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的申请ID")
		return
	}

	application, err := h.interviewService.GetApplicationByID(uint(id))
	if err != nil {
		logger.Errorf("获取面试申请详情失败: %v", err)
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, application.ToResponse())
}

// UpdateApplication 更新面试申请状态（管理员接口）
// @Summary 更新面试申请状态
// @Description 更新面试申请的状态和备注
// @Tags 管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "申请ID"
// @Param request body models.InterviewApplicationUpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=models.InterviewApplicationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/applications/{id} [put]
func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的申请ID")
		return
	}

	// 检查Content-Type
	contentType := c.GetHeader("Content-Type")
	if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
		response.BadRequest(c, "Content-Type必须为application/json")
		return
	}

	var req struct {
		Status       string `json:"status" validate:"required,oneof=pending interviewed passed rejected"`
		AdminRemarks string `json:"admin_remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 验证请求参数
	if !validator.ValidateRequest(c, &req) {
		return
	}

	application, err := h.interviewService.UpdateApplication(uint(id), req.Status, req.AdminRemarks)
	if err != nil {
		logger.Errorf("更新面试申请失败: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "申请状态更新成功", application.ToResponse())
}

// DeleteApplication 删除面试申请（管理员接口）
// @Summary 删除面试申请
// @Description 删除指定的面试申请
// @Tags 管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "申请ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的申请ID")
		return
	}

	err = h.interviewService.DeleteApplication(uint(id))
	if err != nil {
		logger.Errorf("删除面试申请失败: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "申请删除成功", nil)
}

// GetApplicationStats 获取面试申请统计（管理员接口）
// @Summary 获取面试申请统计
// @Description 获取各状态的申请数量统计
// @Tags 管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.InterviewApplicationStats}
// @Failure 401 {object} response.Response
// @Router /admin/applications/stats [get]
func (h *ApplicationHandler) GetApplicationStats(c *gin.Context) {
	stats, err := h.interviewService.GetApplicationStats()
	if err != nil {
		logger.Errorf("获取面试申请统计失败: %v", err)
		response.InternalServerError(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// sendApplicationSuccessEmail 发送申请成功邮件
func sendApplicationSuccessEmail(email, name string) error {
	host := "smtp.qq.com"
	port := 587
	user := "1785260184@qq.com"
	password := "yccnobsfrnoncaci"

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "EPI实验室面试申请已收到")

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>申请成功</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 20px; border-radius: 0 0 8px 8px; }
        .success { background: #52c41a; color: white; padding: 15px; text-align: center; border-radius: 4px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧪 EPI实验室</h1>
        </div>
        <div class="content">
            <p>亲爱的 %s：</p>
            <div class="success">
                <h3>🎉 您的面试申请已成功提交！</h3>
            </div>
            <p>我们已收到您的面试申请，将在3个工作日内联系您安排具体的面试时间。</p>
            <p>请保持手机畅通，注意查收邮件和电话通知。</p>
            <p>感谢您对EPI实验室的关注，期待与您见面！</p>
        </div>
    </div>
</body>
</html>
`, name)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(host, port, user, password)
	return d.DialAndSend(m)
} 