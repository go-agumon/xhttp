package xhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Client 客户端
type Client struct {
	// 发起请求的Client
	client *http.Client
	// 请求超时时间
	timeout time.Duration
	// Header信息，整个生命周期
	header map[string]string
	// Cookie信息，整个生命周期
	cookies []*http.Cookie
	// Header信息，每次请求后会重置
	onceHeader map[string]string
	// Header信息，每次请求后会重置
	onceCookie []*http.Cookie
	// 处理响应的函数
	buildResponse BuildResponseFunc
}

// OnceHeader 设置单次Header信息
func (c *Client) OnceHeader(header map[string]string) *Client {
	c.onceHeader = header
	return c
}

// OnceCookie 设置单次Cookie信息
func (c *Client) OnceCookie(cookies []*http.Cookie) *Client {
	c.onceCookie = cookies
	return c
}

// AddHeader 添加Header信息，整个生命周期有效，若存在则更新，若不存在则添加
func (c *Client) AddHeader(header map[string]string) {
	for k, v := range header {
		c.header[k] = v
	}
}

// SetHeader 设置全新的Header信息，整个生命周期有效，删除原有的，添加新增的
func (c *Client) SetHeader(header map[string]string) {
	c.header = header
}

// AddCookies 添加Cookie信息，整个生命周期有效，若存在则更新，若不存在则添加
func (c *Client) AddCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		c.cookies = append(c.cookies, cookie)
	}
}

// SetCookies 设置全新的Cookie信息，整个生命周期有效，删除原有的，添加新增的
func (c *Client) SetCookies(cookies []*http.Cookie) {
	c.cookies = cookies
}

// Do 发起请求
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	return c.client.Do(request)
}

// buildRequest 构建请求
func (c *Client) buildRequest(method, url string, params map[string]string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range c.header {
		request.Header.Set(k, v)
	}
	for k, v := range c.onceHeader {
		request.Header.Set(k, v)
	}
	// 清空单次请求头信息
	c.onceHeader = nil

	// 设置请求参数
	if params != nil {
		query := request.URL.Query()
		for k, v := range params {
			query.Add(k, v)
		}
		request.URL.RawQuery = query.Encode()
	}

	if _, ok := request.Header["User-Agent"]; !ok {
		request.Header.Set("User-Agent", HttpUserAgentChromePC)
	}

	// 设置Cookie
	for _, v := range c.cookies {
		request.AddCookie(v)
	}
	for _, v := range c.onceCookie {
		request.AddCookie(v)
	}
	// 清空单次Cookie
	c.onceCookie = nil

	return request, nil
}

// SendWithMethod 以指定方法发送请求
func (c *Client) SendWithMethod(method, url string, params map[string]string, body map[string]interface{}) IResponse {
	// IO Reader
	var ioReader bytes.Reader

	if body != nil {
		// 序列化请求体
		buffer, err := json.Marshal(body)
		if err != nil {
			return nil
		}

		ioReader = *bytes.NewReader(buffer)
	}

	request, err := c.buildRequest(method, url, params, &ioReader)
	if err != nil {
		return c.buildResponse(nil, err)
	}

	// 格式化响应
	return c.buildResponse(c.Do(request))
}

// Get 发送GET请求
func (c *Client) Get(url string, params map[string]string) IResponse {
	return c.SendWithMethod(http.MethodGet, url, params, nil)
}

// Post 发送POST请求
func (c *Client) Post(url string, params map[string]string, body map[string]interface{}) IResponse {
	return c.SendWithMethod(http.MethodPost, url, params, body)
}
