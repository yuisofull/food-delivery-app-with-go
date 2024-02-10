package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yuisofull/food-delivery-app-with-go/common"
	"github.com/yuisofull/food-delivery-app-with-go/component/appctx"
)

func CheckRole(appCtx appctx.AppContext, allowRoles ...string) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		for _, role := range allowRoles {
			if requester.GetRole() == role {
				c.Next()
				return
			}
		}
		panic(common.ErrNoPermission(errors.New("invalid role")))
	}

}
