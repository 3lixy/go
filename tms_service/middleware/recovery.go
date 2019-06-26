package middleware

import (
	"azoya/lib/log"
	"azoya/nova"
	"net/http"
	"tms_service/common"
	"fmt"
)

//捕获输出异常
func RecoveryMiddleware() nova.HandlerFunc {
	return func(c *nova.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusOK, common.ErrorMessage(http.StatusBadRequest, err))
				common.GetLogger().Error("system error", log.String("error msg", fmt.Sprintf("%v",err)))
			}
		}()
		c.Next()
	}
}
