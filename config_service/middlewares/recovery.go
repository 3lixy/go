package middlewares

import (
	"azoya/nova"
	"config_service/common"
	"fmt"
	"net/http"
)

//RecoveryMiddleware 捕获输出异常
func RecoveryMiddleware() nova.HandlerFunc {
	return func(c *nova.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err, nil))
				errMessage := fmt.Sprintf("error info:%s", err)
				common.Log(errMessage, "Recovery", "error")
			}
		}()
		c.Next()
	}
}
