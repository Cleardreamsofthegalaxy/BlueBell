package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求
func SignUpHandler(c *gin.Context) {
	//获取参数，参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil { //该函数仅能进行简单的数据判断：字段类型和请求数据的类型
		//参数有误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断error是不是validator类型
		validatorError, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(validatorError.Translate(trans)))
		return
	}
	//业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	//获取参数，参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//参数有误
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断error是不是validator类型
		validatorError, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(validatorError.Translate(trans)))
		return
	}
	//业务逻辑处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//返回响应
	ResponseSuccess(c, token)
}
