package validator

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"lab-recruitment-platform/pkg/response"
)

var validate *validator.Validate

// InitValidator 初始化验证器
func InitValidator() {
	validate = validator.New()
	
	// 注册自定义验证器
	registerCustomValidators()
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidateVar 验证单个字段
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

// ValidateRequest 验证请求并返回错误信息
func ValidateRequest(c *gin.Context, req interface{}) bool {
	if err := validate.Struct(req); err != nil {
		errors := make([]string, 0)
		
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			param := err.Param()
			
			// 获取字段的中文名称
			fieldName := getFieldName(req, field)
			
			// 根据验证标签生成错误信息
			message := getErrorMessage(fieldName, tag, param)
			errors = append(errors, message)
		}
		
		response.ValidationError(c, strings.Join(errors, "; "))
		return false
	}
	
	return true
}

// registerCustomValidators 注册自定义验证器
func registerCustomValidators() {
	// 可以在这里添加自定义验证器
	// 例如：validate.RegisterValidation("custom_rule", customValidationFunc)
}

// getFieldName 获取字段的中文名称
func getFieldName(s interface{}, field string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	f, ok := t.FieldByName(field)
	if !ok {
		return field
	}
	
	// 从json标签获取字段名
	jsonTag := f.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		// 移除omitempty等选项
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}
		return jsonTag
	}
	
	// 从validate标签获取字段名
	validateTag := f.Tag.Get("validate")
	if validateTag != "" {
		// 这里可以解析validate标签中的自定义名称
		// 例如：validate:"required,label=用户名"
		return field
	}
	
	return field
}

// getErrorMessage 根据验证标签生成错误信息
func getErrorMessage(fieldName, tag, param string) string {
	switch tag {
	case "required":
		return fieldName + "不能为空"
	case "email":
		return fieldName + "格式不正确"
	case "min":
		return fieldName + "长度不能少于" + param + "个字符"
	case "max":
		return fieldName + "长度不能超过" + param + "个字符"
	case "len":
		return fieldName + "长度必须为" + param + "个字符"
	case "oneof":
		return fieldName + "必须是以下值之一: " + param
	case "numeric":
		return fieldName + "必须是数字"
	case "alpha":
		return fieldName + "只能包含字母"
	case "alphanum":
		return fieldName + "只能包含字母和数字"
	case "url":
		return fieldName + "必须是有效的URL"
	case "file":
		return fieldName + "必须是有效的文件"
	case "image":
		return fieldName + "必须是有效的图片文件"
	case "datetime":
		return fieldName + "必须是有效的日期时间格式"
	case "date":
		return fieldName + "必须是有效的日期格式"
	case "time":
		return fieldName + "必须是有效的时间格式"
	case "gt":
		return fieldName + "必须大于" + param
	case "gte":
		return fieldName + "必须大于等于" + param
	case "lt":
		return fieldName + "必须小于" + param
	case "lte":
		return fieldName + "必须小于等于" + param
	case "eq":
		return fieldName + "必须等于" + param
	case "ne":
		return fieldName + "不能等于" + param
	default:
		return fieldName + "验证失败"
	}
}

// ValidatePagination 验证分页参数
func ValidatePagination(c *gin.Context) (page, size int, valid bool) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	
	// 验证页码
	if err := validate.Var(pageStr, "numeric,min=1"); err != nil {
		response.BadRequest(c, "页码必须是大于0的数字")
		return 0, 0, false
	}
	
	// 验证每页大小
	if err := validate.Var(sizeStr, "numeric,min=1,max=100"); err != nil {
		response.BadRequest(c, "每页大小必须是1-100之间的数字")
		return 0, 0, false
	}
	
	page, _ = strconv.Atoi(pageStr)
	size, _ = strconv.Atoi(sizeStr)
	
	return page, size, true
}

// ValidateID 验证ID参数
func ValidateID(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	
	if err := validate.Var(idStr, "required,numeric,min=1"); err != nil {
		response.BadRequest(c, "ID参数无效")
		return 0, false
	}
	
	id, _ := strconv.Atoi(idStr)
	return uint(id), true
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	return validate.Var(email, "required,email") == nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) bool {
	// 密码至少6位，包含字母和数字
	return validate.Var(password, "required,min=6,alphanum") == nil
}

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) bool {
	// 简单的手机号验证，可以根据需要调整
	return validate.Var(phone, "required,len=11,numeric") == nil
} 