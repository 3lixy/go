package common

import (
	"azoya/nova"
	"github.com/jinzhu/gorm"
	"net/http"
)

//ResponseResult 根据返回结果和错误集返回结果
func ResponseResult(c *nova.Context, result interface{}, err error) {
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, ErrorMessage(http.StatusNotFound, http.StatusText(http.StatusNotFound)))
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMessage(http.StatusInternalServerError, err.Error()))
	} else {
		c.JSON(http.StatusOK, Message(http.StatusOK, "success", result))
	}
}
