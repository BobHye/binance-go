package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type secType int

const (
	SecTypeNone secType = iota
	SecTypeAPIKey
	SecTypeSigned
)

type Params map[string]interface{} // 参数

// request define an API request
type Request struct {
	Method     string
	Endpoint   string
	Query      url.Values
	Form       url.Values
	RecvWindow int64
	SecType    secType
	Header     http.Header
	Body       io.Reader
	FullURL    string
}

// setParam 将键值对参数设置为查询参数
func (r *Request) SetParam(key string, value interface{}) *Request {
	if r.Query == nil {
		r.Query = url.Values{}
	}
	r.Query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// setFormParam 将键值对参数设置为请求表单body
func (r *Request) SetFormParam(key string, value interface{}) *Request {
	if r.Form == nil {
		r.Form = url.Values{}
	}
	r.Form.Set(key, fmt.Sprintf("%v", value))
	return r
}

// setFormParams 将键值对参数设置为请求表单body
func (r *Request) SetFormParams(m Params) *Request {
	for k, v := range m {
		r.SetFormParam(k, v)
	}
	return r
}

func (r *Request) Validate() (err error) {
	if r.Query == nil {
		r.Query = url.Values{}
	}
	if r.Form == nil {
		r.Form = url.Values{}
	}
	return nil
}

// RequestOption 定义请求的选项类型
type RequestOption func(*Request)

// WithRecvWindow 设置recvWindow参数
func WithRecvWindow(recvWindow int64) RequestOption {
	return func(r *Request) {
		r.RecvWindow = recvWindow
	}
}

// WithHeader 添加或设置请求头
func WithHeader(key, value string, replace bool) RequestOption {
	return func(r *Request) {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		if replace {
			r.Header.Set(key, value)
		} else {
			r.Header.Add(key, value)
		}
	}
}

// WithHeaders 设置或替换请求头
func WithHeaders(header http.Header) RequestOption {
	return func(r *Request) {
		r.Header = header.Clone()
	}
}
