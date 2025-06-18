package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lab-recruitment-platform/internal/config"
	"lab-recruitment-platform/internal/models"
	"lab-recruitment-platform/pkg/logger"
	"lab-recruitment-platform/pkg/response"
)

// Claims JWT声明
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.GlobalConfig.JWT.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lab-recruitment-platform",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "令牌格式错误")
			c.Abort()
			return
		}

		// 提取令牌
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析令牌
		claims, err := ParseToken(tokenString)
		if err != nil {
			logger.Warnf("令牌解析失败: %v", err)
			response.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		// 检查令牌是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			response.Unauthorized(c, "令牌已过期")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行JWT认证
		AuthMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// 检查用户角色
		role, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "用户信息不完整")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有令牌，继续处理
			c.Next()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// 令牌格式错误，继续处理（不强制要求）
			c.Next()
			return
		}

		// 提取令牌
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析令牌
		claims, err := ParseToken(tokenString)
		if err != nil {
			// 令牌无效，继续处理（不强制要求）
			logger.Warnf("可选认证令牌解析失败: %v", err)
			c.Next()
			return
		}

		// 检查令牌是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			// 令牌过期，继续处理（不强制要求）
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUserEmail 获取当前用户邮箱
func GetCurrentUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	return email.(string), true
}

// GetCurrentUserRole 获取当前用户角色
func GetCurrentUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(string), true
}

// IsCurrentUserAdmin 判断当前用户是否为管理员
func IsCurrentUserAdmin(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == "admin"
}

// IsCurrentUserStudent 判断当前用户是否为学生
func IsCurrentUserStudent(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == "student"
} 