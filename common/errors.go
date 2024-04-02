package common

import (
	"errors"
	"fmt"
)

// APIError - default API error when response status is 4xx or 5xx
type APIError struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`
}

// Error APIError 类型实现error接口
func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%d, msg=%s", e.Code, e.Message)
}

func IsAPIError(e error) bool {
	var APIError *APIError
	// As - 判断err是否为target 类型
	ok := errors.As(e, &APIError)
	return ok
}
