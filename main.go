package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	mrand "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// 面试申请者结构
type Applicant struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	StudentID   string    `json:"student_id"`
	Major       string    `json:"major"`
	Grade       string    `json:"grade"`
	InterviewTime string  `json:"interview_time"`
	Status      string    `json:"status"` // pending, first_pass, second_pass, passed, rejected
	FirstRemark string    `json:"first_remark"`
	SecondRemark string   `json:"second_remark"`
	ThirdRemark string    `json:"third_remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 管理员结构
type Admin struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // 不返回密码
	Role     string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 登录响应结构
type LoginResponse struct {
	Token string `json:"token"`
	User  Admin  `json:"user"`
}

// 申请请求结构
type ApplicationRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	StudentID     string `json:"student_id" binding:"required"`
	Major         string `json:"major" binding:"required"`
	Grade         string `json:"grade" binding:"required"`
	InterviewTime string `json:"interview_time" binding:"required"`
	VerificationCode string `json:"verification_code" binding:"required"`
}

// 状态更新请求
type StatusUpdateRequest struct {
	Status       string `json:"status" binding:"required"`
	FirstRemark  string `json:"first_remark"`
	SecondRemark string `json:"second_remark"`
	ThirdRemark  string `json:"third_remark"`
}

// API响应结构
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 验证码结构
type VerificationCode struct {
	Code      string
	ExpiresAt time.Time
}

// 内存存储
var applicants = make(map[int]*Applicant)
var nextID = 1

// 验证码存储（实际项目中应该使用Redis）
var verificationCodes = make(map[string]*VerificationCode)

// 管理员账户（实际项目中应该存储在数据库中）
var adminAccount = &Admin{
	ID:        1,
	Email:     "1234567@qq.com",
	Password:  "epi666",
	Role:      "admin",
	CreatedAt: time.Now(),
}

// JWT token存储（实际项目中应该使用Redis）
var adminTokens = make(map[string]time.Time)

// 生成6位随机验证码
func generateVerificationCode() string {
	// 生成100000-999999之间的随机数，确保是6位数字
	code := 100000 + mrand.Intn(900000)
	return fmt.Sprintf("%d", code)
}

// 发送验证码邮件
func sendVerificationEmail(email, code string) error {
	config := getEmailConfig()
	
	// 测试模式：直接返回成功，不实际发送邮件
	if config.SMTPHost == "test" {
		fmt.Printf("测试模式 - 验证码 %s 已发送到 %s\n", code, email)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.SMTPUser)
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
					<h1>🧪 实验室面试申请</h1>
				</div>
				<div class="content">
					<p>您好！</p>
					<p>您正在申请实验室面试，请使用以下验证码完成验证：</p>
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

	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)
	return d.DialAndSend(m)
}

// 验证验证码
func verifyCode(email, code string) bool {
	if vc, exists := verificationCodes[email]; exists {
		if vc.Code == code && time.Now().Before(vc.ExpiresAt) {
			delete(verificationCodes, email) // 使用后立即删除
			return true
		}
		// 验证码过期或错误，删除
		delete(verificationCodes, email)
	}
	return false
}

// 生成JWT token
func generateToken(email string) string {
	header := `{"alg":"HS256","typ":"JWT"}`
	payload := fmt.Sprintf(`{"email":"%s","exp":%d}`, email, time.Now().Add(24*time.Hour).Unix())
	
	headerB64 := base64.RawURLEncoding.EncodeToString([]byte(header))
	payloadB64 := base64.RawURLEncoding.EncodeToString([]byte(payload))
	
	signature := hmac.New(sha256.New, []byte("your-secret-key"))
	signature.Write([]byte(headerB64 + "." + payloadB64))
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature.Sum(nil))
	
	return headerB64 + "." + payloadB64 + "." + signatureB64
}

// 验证JWT token
func verifyToken(tokenString string) (string, bool) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", false
	}
	
	// 检查token是否在存储中且未过期
	if expiry, exists := adminTokens[tokenString]; exists && time.Now().Before(expiry) {
		// 解析payload获取email
		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return "", false
		}
		
		// 简单的JSON解析，实际项目中应该使用proper JSON parser
		payload := string(payloadBytes)
		if strings.Contains(payload, adminAccount.Email) {
			return adminAccount.Email, true
		}
	}
	
	return "", false
}

// 管理员登录验证
func verifyAdminLogin(email, password string) bool {
	return email == adminAccount.Email && password == adminAccount.Password
}

func main() {
	// 初始化随机数种子（Go 1.20+中自动初始化，无需手动调用rand.Seed）
	
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		config := getEmailConfig()
		emailStatus := "已配置"
		if config.SMTPHost == "test" {
			emailStatus = "测试模式"
		}
		
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "实验室面试申请系统运行正常",
			"email_status": emailStatus,
			"time":    time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 管理员登录
		api.POST("/auth/login", func(c *gin.Context) {
			var req LoginRequest

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "请求参数错误: " + err.Error(),
				})
				return
			}

			// 验证管理员账户
			if !verifyAdminLogin(req.Email, req.Password) {
				c.JSON(http.StatusUnauthorized, ApiResponse{
					Code:    401,
					Message: "邮箱或密码错误",
				})
				return
			}

			// 生成JWT token
			token := generateToken(req.Email)
			adminTokens[token] = time.Now().Add(24 * time.Hour)

			// 返回登录成功响应
			c.JSON(http.StatusOK, ApiResponse{
				Code:    200,
				Message: "登录成功",
				Data: LoginResponse{
					Token: token,
					User:  *adminAccount,
				},
			})
		})

		// 发送验证码
		api.POST("/send-code", func(c *gin.Context) {
			var req struct {
				Email string `json:"email" binding:"required,email"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "请求参数错误: " + err.Error(),
				})
				return
			}

			// 检查是否已经申请过
			for _, applicant := range applicants {
				if applicant.Email == req.Email {
					c.JSON(http.StatusConflict, ApiResponse{
						Code:    409,
						Message: "该邮箱已提交过申请",
					})
					return
				}
			}

			// 生成验证码
			code := generateVerificationCode()
			expiresAt := time.Now().Add(5 * time.Minute) // 5分钟有效期

			// 存储验证码
			verificationCodes[req.Email] = &VerificationCode{
				Code:      code,
				ExpiresAt: expiresAt,
			}

			// 发送邮件
			if err := sendVerificationEmail(req.Email, code); err != nil {
				// 发送失败，删除验证码
				delete(verificationCodes, req.Email)
				c.JSON(http.StatusInternalServerError, ApiResponse{
					Code:    500,
					Message: "邮件发送失败，请稍后重试",
				})
				return
			}

			config := getEmailConfig()
			message := "验证码已发送到邮箱，请注意查收"
			if config.SMTPHost == "test" {
				message = fmt.Sprintf("测试模式 - 验证码：%s（5分钟内有效）", code)
			}

			c.JSON(http.StatusOK, ApiResponse{
				Code:    200,
				Message: message,
			})
		})

		// 提交面试申请
		api.POST("/apply", func(c *gin.Context) {
			var req ApplicationRequest

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "请求参数错误: " + err.Error(),
				})
				return
			}

			// 验证验证码
			if !verifyCode(req.Email, req.VerificationCode) {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "验证码错误或已过期",
				})
				return
			}

			// 检查是否已经申请过
			for _, applicant := range applicants {
				if applicant.Email == req.Email {
					c.JSON(http.StatusConflict, ApiResponse{
						Code:    409,
						Message: "该邮箱已提交过申请",
					})
					return
				}
			}

			// 创建申请记录
			now := time.Now()
			applicant := &Applicant{
				ID:            nextID,
				Name:          req.Name,
				Email:         req.Email,
				Phone:         req.Phone,
				StudentID:     req.StudentID,
				Major:         req.Major,
				Grade:         req.Grade,
				InterviewTime: req.InterviewTime,
				Status:        "pending",
				CreatedAt:     now,
				UpdatedAt:     now,
			}

			applicants[nextID] = applicant
			nextID++

			c.JSON(http.StatusOK, ApiResponse{
				Code:    200,
				Message: "申请提交成功",
				Data:    applicant,
			})
		})

		// 获取所有申请者列表（管理员接口）
		api.GET("/applicants", func(c *gin.Context) {
			applicantList := make([]*Applicant, 0, len(applicants))
			for _, applicant := range applicants {
				applicantList = append(applicantList, applicant)
			}

			c.JSON(http.StatusOK, ApiResponse{
				Code:    200,
				Message: "获取申请者列表成功",
				Data:    applicantList,
			})
		})

		// 更新申请状态（管理员接口）
		api.PUT("/applicants/:id/status", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "无效的申请ID",
				})
				return
			}

			applicant, exists := applicants[id]
			if !exists {
				c.JSON(http.StatusNotFound, ApiResponse{
					Code:    404,
					Message: "申请记录不存在",
				})
				return
			}

			var req StatusUpdateRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "请求参数错误: " + err.Error(),
				})
				return
			}

			// 更新状态和备注
			applicant.Status = req.Status
			applicant.FirstRemark = req.FirstRemark
			applicant.SecondRemark = req.SecondRemark
			applicant.ThirdRemark = req.ThirdRemark
			applicant.UpdatedAt = time.Now()

			c.JSON(http.StatusOK, ApiResponse{
				Code:    200,
				Message: "状态更新成功",
				Data:    applicant,
			})
		})

		// 获取申请详情
		api.GET("/applicants/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "无效的申请ID",
				})
				return
			}

			applicant, exists := applicants[id]
			if !exists {
				c.JSON(http.StatusNotFound, ApiResponse{
					Code:    404,
					Message: "申请记录不存在",
				})
				return
			}

			c.JSON(http.StatusOK, ApiResponse{
				Code:    200,
				Message: "获取申请详情成功",
				Data:    applicant,
			})
		})

		// 删除申请记录（管理员接口）
		api.DELETE("/applicants/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Code:    400,
					Message: "无效的申请ID",
				})
				return
			}

			if _, exists := applicants[id]; !exists {
				c.JSON(http.StatusNotFound, ApiResponse{
					Code:    404,
					Message: "申请记录不存在",
				})
				return
			}

			delete(applicants, id)

			c.JSON(http.StatusOK, ApiResponse{
				Code:    200,
				Message: "申请记录删除成功",
			})
		})
	}

	r.Run(":8080")
}