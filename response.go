package xhttp

import (
	"errors"
	"io"
	"net/http"

	"github.com/bitly/go-simplejson"
)

// BuildResponseFunc 响应函数
type BuildResponseFunc func(resp *http.Response, err error) IResponse

// IResponse 当客户端发送请求后，返回一个已实现此接口的对象
type IResponse interface {
	// StatusCode 返回响应的状态码
	StatusCode() int

	// Header 返回响应的Header信息
	Header() http.Header

	// Cookie 根据名称返回响应中的Cookie
	Cookie(name string) *http.Cookie

	// Content 返回响应的内容
	Content() interface{}

	// ContentLength 返回响应的内容长度
	ContentLength() int64

	// Error 返回响应的错误信息
	Error() error

	// Response 返回原始响应
	Response() *http.Response

	// Request 返回响应的请求信息
	Request() *http.Request
}

// NewXResponse 创建响应
func NewXResponse(rep *http.Response, err error) IResponse {
	response := new(XResponse)
	if err != nil {
		response.err = err
		return response
	}
	response.response = rep

	content, err := io.ReadAll(rep.Body)
	if err != nil {
		response.err = err
		return response
	}

	// 填充内容
	response.content = content

	// 关闭句柄
	_ = rep.Body.Close()

	return response
}

// XResponse 标准的响应
type XResponse struct {
	err      error
	content  []byte
	response *http.Response
}

// StatusCode 返回响应的状态码
func (x *XResponse) StatusCode() int {
	if x.response == nil {
		return 0
	}
	return x.response.StatusCode
}

// Header 返回响应的Header信息
func (x *XResponse) Header() http.Header {
	if x.response == nil {
		return nil
	}
	return x.response.Header
}

// Cookie 根据名称返回响应中的Cookie
func (x *XResponse) Cookie(name string) *http.Cookie {
	for _, cookie := range x.response.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// Content 返回响应的内容
func (x *XResponse) Content() interface{} {
	return x.content
}

// ContentLength 返回响应的内容长度
func (x *XResponse) ContentLength() int64 {
	if x.response == nil {
		return 0
	}
	return x.response.ContentLength
}

// Error 返回响应的错误信息
func (x *XResponse) Error() error {
	return x.err
}

// Response 返回原始响应
func (x *XResponse) Response() *http.Response {
	return x.response
}

// Request 返回响应的请求信息
func (x *XResponse) Request() *http.Request {
	if x.response == nil {
		return nil
	}
	return x.response.Request
}

// NewJResponse 创建响应
func NewJResponse(rep *http.Response, err error) IResponse {
	response := new(JResponse)
	if err != nil {
		response.err = err
		return response
	}
	response.response = rep

	content, err := io.ReadAll(rep.Body)
	if err != nil {
		response.err = err
		return response
	}

	contentJSON, err := simplejson.NewJson(content)
	if err != nil {
		response.err = errors.New("can't serialize data to JSON format")
		return response
	}

	// 填充内容
	response.content = contentJSON

	// 关闭句柄
	_ = rep.Body.Close()

	return response
}

// JResponse 基于SimpleJSON的响应
type JResponse struct {
	err      error
	content  *simplejson.Json
	response *http.Response
}

// StatusCode 返回响应的状态码
func (j *JResponse) StatusCode() int {
	if j.response == nil {
		return 0
	}
	return j.response.StatusCode
}

// Header 返回响应的Header信息
func (j *JResponse) Header() http.Header {
	if j.response == nil {
		return nil
	}
	return j.response.Header
}

// Cookie 根据名称返回响应中的Cookie
func (j *JResponse) Cookie(name string) *http.Cookie {
	for _, cookie := range j.response.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// Content 返回响应的内容
func (j *JResponse) Content() interface{} {
	return j.content
}

// ContentLength 返回响应的内容长度
func (j *JResponse) ContentLength() int64 {
	if j.response == nil {
		return 0
	}
	return j.response.ContentLength
}

// Error 返回响应的错误信息
func (j *JResponse) Error() error {
	return j.err
}

// Response 返回原始响应
func (j *JResponse) Response() *http.Response {
	return j.response
}

// Request 返回响应的请求信息
func (j *JResponse) Request() *http.Request {
	if j.response == nil {
		return nil
	}
	return j.response.Request
}
