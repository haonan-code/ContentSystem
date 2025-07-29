package api

import (
	"contentsystem/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

const SessionKey = "session_id"

type SessionAuth struct {
	rdb *redis.Client
}

func NewSessionAuth() *SessionAuth {
	s := &SessionAuth{}
	connRdb(s)
	return s
}

func (s *SessionAuth) Auth(ctx *gin.Context) {
	sessionID := ctx.GetHeader(SessionKey)
	// TODO: imp auth
	if sessionID == "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "session id is null")
	}

	authKey := utils.GetAuthKey(sessionID)
	loginTime, err := s.rdb.Get(ctx, authKey).Result()
	if err != nil && err != redis.Nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "session auth error")
	}
	if loginTime == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "session auth failed")
	}

	fmt.Println("session id ", sessionID)
	ctx.Next()
	fmt.Println("end = ", sessionID)
}

func connRdb(s *SessionAuth) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.31.43:6379",
		Password: "redis123456", // no password set
		DB:       0,             // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	s.rdb = rdb
}
