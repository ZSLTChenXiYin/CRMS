package service

import (
	"errors"

	"github.com/ZSLTChenXiYin/CRMS/crms-backend/com"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserToken struct {
	UserID uint
	Expire int64
}

func AnalyzeToken(c *gin.Context) (user_token *UserToken, err error) {
	// 获取令牌
	value, exists := c.Get("signed_token")
	if !exists {
		return nil, errors.New("未登录")
	}
	signed_token := value.(string)

	// 解析令牌
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(signed_token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(com.Config.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("非法令牌")
	}

	// 获取令牌信息
	user_token = &UserToken{
		UserID: uint((*claims)["user_id"].(float64)),
		Expire: int64((*claims)["expire"].(float64)),
	}

	return user_token, nil
}

func RecordUserBasicInformation(c *gin.Context, information map[string]any) {
	information["ip"] = c.ClientIP()
	information["user_agent"] = c.GetHeader("User-Agent")
	information["referer"] = c.GetHeader("Referer")
}
