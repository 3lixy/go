package controllers

import (
	"azoya/nova"
	"net/http"
)

//SampleController controller 都继承Monitor
type TestController struct {

}

//NewSampleController 初始化，把初始化的Monitor对象放入
func NewTestController() *TestController {
	return &TestController{
	}
}

//LoggerInfo 普通日志记录的示例，在span下面追加log
func (controller *TestController) LoggerInfo(c *nova.Context) {
	// controller.Logger().Bg().Info(typeof(cancelCtx.Context))
	c.JSON(http.StatusOK, nova.H{"code": "aaa"})
}