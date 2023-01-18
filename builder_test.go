package xhttp

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
)

func TestBuilder(t *testing.T) {
	// 创建一个构建器
	builder := NewBuilder()

	// 是否跳过服务器证书校验
	builder.SkipVerify(true)

	// 设置请求超时时间
	builder.Timeout(time.Second * 5)

	// 设置Header信息
	header := make(map[string]string)
	header["Content-Type"] = HttpContentTypeJson
	header["Accept-Language"] = "Accept-Language: en,zh"
	builder.Header(header)

	// 设置处理响应的函数
	builder.BuildResponse(NewXResponse)

	// 构建客户端
	client, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 发起请求
	response := client.Get("https://baidu.com", nil)

	// 输出结果
	fmt.Printf("Error: %v\n", response.Error())
	fmt.Printf("StatusCode: %v\n", response.StatusCode())
	fmt.Printf("ContentLength: %v\n", response.ContentLength())
}

func TestDo(t *testing.T) {
	// 选项
	options := make([]Option, 0)

	// 设置Header信息
	header := make(map[string]string)
	header["Content-Type"] = HttpContentTypeJson
	options = append(options, WithHeader(header))

	// 设置处理响应的函数
	options = append(options, WithBuildResponse(NewJResponse))

	response := Do(http.MethodGet, "http://myip.ipip.net/json", nil, nil, options...)
	fmt.Printf("Local IP: %v\n", response.Content().(*simplejson.Json).GetPath("data", "ip").MustString())
}
