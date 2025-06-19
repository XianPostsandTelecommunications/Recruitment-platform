package handlers

import (
	"github.com/gin-gonic/gin"
	"lab-recruitment-platform/internal/middleware"
	"lab-recruitment-platform/internal/models"
	"lab-recruitment-platform/internal/services"
	"lab-recruitment-platform/pkg/logger"
	"lab-recruitment-platform/pkg/response"
	"lab-recruitment-platform/pkg/validator"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *services.UserService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: services.NewUserService(),
	}
}



// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=gin.H}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 验证请求参数
	if !validator.ValidateRequest(c, &req) {
		return
	}

	// 根据邮箱获取用户
	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		logger.Warnf("用户登录失败，邮箱不存在: %s", req.Email)
		response.Unauthorized(c, "邮箱或密码错误")
		return
	}

	// 验证密码
	if !h.userService.VerifyPassword(user, req.Password) {
		logger.Warnf("用户登录失败，密码错误: %s", req.Email)
		response.Unauthorized(c, "邮箱或密码错误")
		return
	}

	// 检查用户状态
	if !user.IsActive() {
		response.Forbidden(c, "账户已被禁用")
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		logger.Errorf("生成令牌失败: %v", err)
		response.InternalServerError(c, "生成令牌失败")
		return
	}

	// 返回用户信息和令牌
	response.SuccessWithMessage(c, "登录成功", gin.H{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.UserResponse}
// @Failure 401 {object} response.Response
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v", err)
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user.ToResponse())
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserUpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=models.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 验证请求参数
	if !validator.ValidateRequest(c, &req) {
		return
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		logger.Errorf("更新用户信息失败: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", user.ToResponse())
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 验证请求参数
	if !validator.ValidateRequest(c, &req) {
		return
	}

	err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		logger.Errorf("修改密码失败: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 刷新JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=gin.H}
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v", err)
		response.NotFound(c, "用户不存在")
		return
	}

	// 生成新的JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		logger.Errorf("生成令牌失败: %v", err)
		response.InternalServerError(c, "生成令牌失败")
		return
	}

	response.SuccessWithMessage(c, "令牌刷新成功", gin.H{
		"token": token,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 在实际应用中，可以将令牌加入黑名单
	// 这里简单返回成功响应
	response.SuccessWithMessage(c, "登出成功", nil)
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
} 