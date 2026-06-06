package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenli03/humq/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4000, "msg": "参数错误"})
		return
	}

	var user database.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "用户名或密码错误"})
		return
	}

	token, err := GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5000, "msg": "生成令牌失败"})
		return
	}

	refresh, err := GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5000, "msg": "生成刷新令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": gin.H{
			"token":         token,
			"refresh_token": refresh,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
			},
		},
	})
}

func RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "刷新成功", "data": gin.H{}})
}

func DemoLoginHandler(demoMode bool) gin.HandlerFunc {
	if demoMode {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "登录成功（Demo模式）",
				"data": gin.H{
					"token":         "demo-token",
					"refresh_token": "demo-refresh",
					"user":          gin.H{"id": 1, "username": "admin", "role": "admin"},
				},
			})
		}
	}
	return Login
}
