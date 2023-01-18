package xhttp

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

// CheckRedirectFunc 请求重定向的回调函数
type CheckRedirectFunc func(req *http.Request, via []*http.Request) error

// NewBuilder 创建构建器
func NewBuilder() *Builder {
	return &Builder{
		skipVerify: true,
	}
}

// Builder 构建器
type Builder struct {
	// 请求超时时间
	timeout time.Duration
	// Header信息，默认会自动填充User-Agent
	header map[string]string
	// Cookie信息
	cookies []*http.Cookie
	// 是否跳过验证证书，默认跳过
	skipVerify bool
	// 重定向函数
	checkRedirect CheckRedirectFunc
	// 处理响应的函数
	buildResponse BuildResponseFunc
}

// Build 构建客户端
func (builder *Builder) Build() (*Client, error) {
	// 检查是否有处理响应的函数
	if builder.buildResponse == nil {
		return nil, errors.New("clint not set BuildResponse")
	}

	// 检查是否跳过验证证书
	tlsConfig := &tls.Config{
		InsecureSkipVerify: builder.skipVerify,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &Client{
		client: &http.Client{
			Transport:     transport,
			Timeout:       builder.timeout,
			CheckRedirect: builder.checkRedirect,
		},
		header:        builder.header,
		cookies:       builder.cookies,
		buildResponse: builder.buildResponse,
	}

	return client, nil
}

// Timeout 设置请求超时时间
func (builder *Builder) Timeout(timeout time.Duration) *Builder {
	builder.timeout = timeout
	return builder
}

// SkipVerify 设置是否跳过验证证书
func (builder *Builder) SkipVerify(skip bool) *Builder {
	builder.skipVerify = skip
	return builder
}

// Header 设置Header信息
func (builder *Builder) Header(header map[string]string) *Builder {
	builder.header = header
	return builder
}

// Cookies 设置Cookie信息
func (builder *Builder) Cookies(cookies []*http.Cookie) *Builder {
	builder.cookies = cookies
	return builder
}

// CheckRedirect 设置重定向函数
func (builder *Builder) CheckRedirect(checkRedirect CheckRedirectFunc) *Builder {
	builder.checkRedirect = checkRedirect
	return builder
}

// BuildResponse 设置处理响应函数
func (builder *Builder) BuildResponse(build BuildResponseFunc) *Builder {
	builder.buildResponse = build
	return builder
}

// Option 构建器的选项
type Option func(b *Builder)

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(b *Builder) {
		b.timeout = timeout
	}
}

// WithSkipVerify 设置是否跳过验证证书
func WithSkipVerify(skip bool) Option {
	return func(b *Builder) {
		b.skipVerify = skip
	}
}

// WithHeader 设置Header信息
func WithHeader(header map[string]string) Option {
	return func(b *Builder) {
		b.header = header
	}
}

// WithCookies 设置Cookie信息
func WithCookies(cookies []*http.Cookie) Option {
	return func(b *Builder) {
		b.cookies = cookies
	}
}

// WithCheckRedirect 设置重定向函数
func WithCheckRedirect(checkRedirect CheckRedirectFunc) Option {
	return func(b *Builder) {
		b.checkRedirect = checkRedirect
	}
}

// WithBuildResponse 设置重定向函数
func WithBuildResponse(build BuildResponseFunc) Option {
	return func(b *Builder) {
		b.buildResponse = build

	}
}

func Do(method, url string, params map[string]string, body map[string]interface{}, options ...Option) IResponse {
	// 创建一个构建器
	builder := NewBuilder()

	// 设置处理响应的函数（默认）
	builder.BuildResponse(NewXResponse)

	// 应用构建器的选项
	for _, option := range options {
		option(builder)
	}

	// 检查是否跳过验证证书
	tlsConfig := &tls.Config{
		InsecureSkipVerify: builder.skipVerify,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &Client{
		client: &http.Client{
			Transport:     transport,
			Timeout:       builder.timeout,
			CheckRedirect: builder.checkRedirect,
		},
		header:        builder.header,
		cookies:       builder.cookies,
		buildResponse: builder.buildResponse,
	}

	return client.SendWithMethod(method, url, params, body)
}
