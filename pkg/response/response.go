package response

import (
	"net/http"

	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/gin-gonic/gin"
)

// Result 统一返回结果
type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func Success[T any](data T) Result[T] {
	return Result[T]{
		Code: errors.SuccessCode,
		Msg:  "success",
		Data: data,
	}
}

func Error(code int, msg string) Result[any] {
	return Result[any]{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// WriteSuccess 写入成功响应
func WriteSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Success(data))
}

// WriteError 写入错误响应
func WriteError(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Error(code, msg))
}

// WriteBizError 写入业务异常响应
func WriteBizError(c *gin.Context, err error) {
	if e, ok := err.(*errors.BizError); ok {
		c.JSON(http.StatusOK, Error(e.Code, e.Msg))
		return
	}
	c.JSON(http.StatusOK, Error(errors.ServerErrCode, err.Error()))
}

// WritePageData 写入分页响应 (Custom wrapper)
func WritePageData[T any](c *gin.Context, total int64, list []T) {
	c.JSON(http.StatusOK, Success(pagination.PageResult[T]{
		List:  list,
		Total: total,
	}))
}

// WritePage 写入分页响应 (Helper for pre-converted list)
func WritePage[T any](c *gin.Context, total int64, list []T) {
	c.JSON(http.StatusOK, Success(pagination.PageResult[T]{
		List:  list,
		Total: total,
	}))
}

// NormalizeIntPtr normalizes *int query parameter to handle empty string case
// Gin binds "?status=" (empty string) to &0 instead of nil, causing unintended filtering
// This function checks if the raw query param is empty and returns nil in that case
func NormalizeIntPtr(c *gin.Context, paramName string, value *int) *int {
	if value != nil && c.Query(paramName) == "" {
		return nil
	}
	return value
}
