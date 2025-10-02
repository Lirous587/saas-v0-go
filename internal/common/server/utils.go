package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func GetUserID(ctx *gin.Context) (int64, error) {
	uidStr, exist := ctx.Get(UserIDKey)
	if !exist {
		return 0, fmt.Errorf("user_id不存在")
	}

	userID, ok := uidStr.(int64)
	if !ok {
		return 0, fmt.Errorf("user_id类型错误")
	}

	return userID, nil
}
