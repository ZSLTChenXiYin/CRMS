package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ZSLTChenXiYin/CRMS/crms-backend/api/service"
	"github.com/ZSLTChenXiYin/CRMS/crms-backend/com"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 注册用户
//
// 请求体参数：form
//
//	email: string
//	password: string
func PostRegister(c *gin.Context) {
	// 处理请求参数
	email := c.PostForm("email")
	password := c.PostForm("password")

	// 密码加密
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "注册失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 创建用户
	err = service.NewUserService(com.Database).CreateUser(email, string(password_hash))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "注册失败",
				"error":   "邮箱已存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "注册失败",
				"error":   err.Error(),
			})
			log.Panicln("[ERROR]", err)
		}
		return
	}

	// 返回注册结果
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// 删除用户
//
// 请求头参数：
//
//	Authorization: Bearer <token>
func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "功能暂未实现",
	})
}

// 查询用户个人信息
//
// 请求头参数：
//
//	Authorization: Bearer <token>
func GetUserInfo(c *gin.Context) {
	// 获取用户令牌
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "查询用户个人信息失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 获取用户信息
	user, err := service.NewUserService(com.Database).GetUserByID(user_token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询用户个人信息失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"message":       "查询用户个人信息成功",
		"email":         user.Email,
		"created_at":    user.CreatedAt,
		"last_login_at": user.LastLoginAt,
	})
}

// 修改用户个人信息
//
// 请求头参数：
//
//	Authorization: Bearer <token>
func PutUserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "功能暂未实现",
	})
}

// 登录用户
//
// 请求体参数：form
//
//	email: string
//	password: string
func PostLogin(c *gin.Context) {
	// 处理登录请求参数
	email := c.PostForm("email")
	password := c.PostForm("password")

	exists, err := com.Redis.Exists(com.Context, email).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "登录失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	if exists == 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "登录失败",
			"error":   "登录请求太过频繁，请5秒后再试",
		})
		return
	}

	err = com.Redis.Set(com.Context, email, 1, 5*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "登录失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 创建用户服务
	user_service := service.NewUserService(com.Database)

	// 获取用户
	user, err := user_service.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "登录失败",
				"error":   "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "登录失败",
				"error":   err.Error(),
			})
			log.Panicln("[ERROR]", err)
		}
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "登录失败",
			"error":   "密码错误",
		})
		return
	}

	// 生成令牌
	expiration_time := time.Now().Add(time.Hour * 24 * 7)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"user_id": user.ID,
		"expire":  expiration_time.Unix(),
	}).SignedString(com.Config.JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "登录失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 更新用户最后登录时间
	err = user_service.UpdateUserLastLoginById(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "登录失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 更新用户过期时间
	if user.ExpiredAt.Before(time.Now()) {
		err = user_service.UpdateUserExpiredById(user.ID, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "登录失败",
				"error":   err.Error(),
			})
			log.Panicln("[ERROR]", err)
			return
		}
	}

	// 返回登录成功信息
	c.JSON(http.StatusOK, gin.H{
		"message":       "登录成功",
		"token":         token,
		"last_login_at": user.LastLoginAt,
	})
}

// 登出用户
//
// 请求头参数：
//
//	Authorization: Bearer <token>
func PostLogout(c *gin.Context) {
	// 获取用户令牌
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "注销失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 更新向前过期时间
	err = service.NewUserService(com.Database).UpdateUserExpiredById(user_token.UserID, time.Unix(user_token.Expire, 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "注销失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "注销成功",
	})
}

// 修改用户个人密码
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	old_password: string
//	new_password: string
func PutResetPassword(c *gin.Context) {
	// 获取用户令牌
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "修改用户个人密码失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 处理用户请求
	old_password := c.PostForm("old_password")
	new_password := c.PostForm("new_password")

	// 获取用户
	user_service := service.NewUserService(com.Database)

	user, err := user_service.GetUserByID(user_token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改用户个人密码失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(old_password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改用户个人密码失败",
			"error":   "旧密码错误",
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 新密码加密
	password_hash, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改用户个人密码失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 修改用户密码
	err = user_service.UpdateUserPasswordHashById(user_token.UserID, string(password_hash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改用户个人密码失败",
			"error":   err.Error(),
		})
		log.Panicln("[ERROR]", err)
		return
	}

	// 创建操作日志
	additional_information := map[string]any{}
	service.RecordUserBasicInformation(c, additional_information)

	// 异步记录操作日志
	go func() {
		operation_log := &service.OperationLog{
			UserID:        user_token.UserID,
			ResourceType:  "user",
			ResourceID:    user_token.UserID,
			Action:        "update",
			ActionDetails: "修改用户个人密码",
			CreatedAt:     time.Now(),
		}

		additional_infomation_json, err := json.Marshal(additional_information)
		if err != nil {
			log.Println("[操作日志记录错误]", err)
		} else {
			operation_log.AdditionalInformation = additional_infomation_json
		}

		err = service.NewOperationLogService(com.Database).CreateOperationLog(operation_log)
		if err != nil {
			log.Println("[操作日志记录错误]", err)

			json_data, err := json.Marshal(operation_log)
			if err != nil {
				log.Println("[操作日志记录错误]", err)
			}

			log.Println("[操作日志备份]", string(json_data))
			return
		}
	}()

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "修改用户个人密码成功",
	})
}
