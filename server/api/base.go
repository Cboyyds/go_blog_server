package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type BaseApi struct {
}

var store = base64Captcha.DefaultMemStore

// Captcha 生成数字验证码
func (baseApi *BaseApi) Captcha(c *gin.Context) {
	// 创建数字验证码的驱动
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)

	// 创建验证码对象
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, _, err := captcha.Generate()

	if err != nil {
		global.Log.Error("Failed to generate captcha:", zap.Error(err))
		response.FailWithMessage("Failed to generate captcha", c)
		return
	}
	response.OkWithData(response.Captcha{
		CaptchaID: id,
		PicPath:   b64s,
	}, c)
}

// SendEmailVerificationCode 处理发送邮件验证码的请求。
// 该函数首先从请求中解析出验证码和邮箱信息，然后验证验证码是否正确。
// 如果验证码正确，则调用基础服务发送邮件验证码，并根据发送结果返回相应的响应。
// 如果验证码错误或发送邮件失败，则返回错误信息。
//
// 参数:
//   - c: *gin.Context, Gin框架的上下文对象，用于处理HTTP请求和响应。
//
// 返回值:
//
//	无返回值，但通过Gin的上下文对象返回HTTP响应。
func (baseApi *BaseApi) SendEmailVerificationCode(c *gin.Context) {
	// 解析请求体中的JSON数据到SendEmailVerificationCode结构体
	var req request.SendEmailVerificationCode
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// 如果解析失败，返回错误信息
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证验证码是否正确
	if store.Verify(req.CaptchaID, req.Captcha, true) { //
		// 如果验证码正确，调用基础服务发送邮件验证码
		err = baseService.SendEmailVerificationCode(c, req.Email)
		if err != nil {
			// 如果发送邮件失败，记录错误日志并返回错误信息
			global.Log.Error("Failed to send email:", zap.Error(err))
			response.FailWithMessage("Failed to send email", c)
			return
		}
		// 发送成功，返回成功信息
		response.OkWithMessage("Successfully sent email", c)
		return
	}

	// 验证码错误，返回错误信息
	response.FailWithMessage("Incorrect verification code", c)
}

// QQLoginURL 返回 QQ 登录链接
func (baseApi *BaseApi) QQLoginURL(c *gin.Context) {
	url := global.Config.QQ.QQLoginURL()
	response.OkWithData(url, c)
}
