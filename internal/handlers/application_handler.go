package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"lab-recruitment-platform/pkg/logger"
	"lab-recruitment-platform/pkg/response"
	"lab-recruitment-platform/pkg/validator"
)

// ApplicationHandler 申请处理器
type ApplicationHandler struct {
}

// NewApplicationHandler 创建申请处理器实例
func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{}
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
	var req ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 验证请求参数
	if !validator.ValidateRequest(c, &req) {
		return
	}

	// 验证验证码
	if !verifyCode(req.Email, req.VerificationCode) {
		response.BadRequest(c, "验证码错误或已过期")
		return
	}

	// 创建申请记录（这里简化处理，实际应该保存到数据库）
	logger.Infof("收到面试申请: 姓名=%s, 邮箱=%s, 专业=%s", req.Name, req.Email, req.Major)

	// 发送申请成功邮件
	if err := sendApplicationSuccessEmail(req.Email, req.Name); err != nil {
		logger.Errorf("发送申请成功邮件失败: %v", err)
	}

	response.SuccessWithMessage(c, "申请提交成功！我们会尽快联系您安排面试", gin.H{
		"application_id": fmt.Sprintf("APP_%d", time.Now().Unix()),
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
	vc, exists := verificationCodes[email]
	if !exists {
		return false
	}

	// 检查验证码是否正确且未过期
	if vc.Code == code && time.Now().Before(vc.ExpiresAt) {
		delete(verificationCodes, email) // 使用后立即删除
		return true
	}

	// 验证码错误或过期，删除
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