package services

import (
	"contentsystem/internal/dao"
	"contentsystem/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginReq struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRsp struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Nickname  string `json:"nickname"`
}

func (c *CmsApp) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var (
		userID   = req.UserID
		password = req.Password
	)
	accountDao := dao.NewAccountDao(c.db)
	account, err := accountDao.FirstByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的账号ID"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的密码"})
		return
	}

	sessionID, err := c.generateSessionID(context.Background(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后重试"})
		return
	}
	// 回包
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &LoginRsp{
			SessionID: sessionID,
			UserID:    account.UserID,
			Nickname:  account.Nickname,
		},
	})
	return

}

func (c *CmsApp) generateSessionID(ctx context.Context, userID string) (string, error) {
	sessionID := uuid.New().String()

	// key : session_id:{user_id} val : session_id

	// 鉴权方式一：对当前用户校验是否有效
	// userID -> sessionID -> 再查 sessionID 是否过期
	sessionKey := utils.GetSessionKey(userID)
	err := c.rdb.Set(ctx, sessionKey, sessionID, 8*time.Hour).Err()
	if err != nil {
		fmt.Printf("rdb set error = %v \n", err)
		return "", err
	}

	// 鉴权方式二：对 sessionID 进行鉴权，即可保护接口
	// 取出客户端携带的 sessionID 拼接session_auth: -> redis 中查询是否有效
	// 若为空值，则说明该会话已过期
	authKey := utils.GetAuthKey(sessionID)
	err = c.rdb.Set(ctx, authKey, time.Now().Unix(), 1*time.Minute).Err()
	if err != nil {
		fmt.Printf("rdb set error = %v \n", err)
		return "", err
	}
	fmt.Println("sessionKey:", sessionKey)
	fmt.Println("authKey:", authKey)
	return sessionID, nil
}
