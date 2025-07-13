package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// RESTClient REST客户端结构
type RESTClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
}

// RESTRequest REST请求结构
type RESTRequest struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

// RESTBatchResult 批量请求结果
type RESTBatchResult struct {
	Index    int
	Response *RESTResponse
	Error    error
}

// RESTResponse REST响应结构
type RESTResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

type REST interface {
	GET(ctx context.Context, path string, headers map[string]string) (*RESTResponse, error)
	POST(ctx context.Context, path string, body interface{}, headers map[string]string) (*RESTResponse, error)
	PUT(ctx context.Context, path string, body interface{}, headers map[string]string) (*RESTResponse, error)
	DELETE(ctx context.Context, path string, headers map[string]string) (*RESTResponse, error)
	PATCH(ctx context.Context, path string, body interface{}, headers map[string]string) (*RESTResponse, error)
	BatchRequest(ctx context.Context, requests []RESTRequest, maxConcurrency int) []RESTBatchResult
	BatchGET(ctx context.Context, paths []string, maxConcurrency int) []RESTBatchResult
}

// NewRESTClient 创建REST客户端
func NewRESTClient(baseURL string) REST {
	return &RESTClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Duration(3 * time.Second),
		},
		Headers: make(map[string]string),
	}
}

// GET 发送GET请求
func (c *RESTClient) GET(ctx context.Context, path string, headers map[string]string) (*RESTResponse, error) {
	return c.request(ctx, "GET", path, nil, headers)
}

// POST 发送POST请求
func (c *RESTClient) POST(ctx context.Context, path string, body interface{}, headers map[string]string) (*RESTResponse, error) {
	return c.request(ctx, "POST", path, body, headers)
}

// PUT 发送PUT请求
func (c *RESTClient) PUT(ctx context.Context, path string, body interface{}, headers map[string]string) (*RESTResponse, error) {
	return c.request(ctx, "PUT", path, body, headers)
}

// DELETE 发送DELETE请求
func (c *RESTClient) DELETE(ctx context.Context, path string, headers map[string]string) (*RESTResponse, error) {
	return c.request(ctx, "DELETE", path, nil, headers)
}

// PATCH 发送PATCH请求
func (c *RESTClient) PATCH(ctx context.Context, path string, body interface{}, headers map[string]string) (*RESTResponse, error) {
	return c.request(ctx, "PATCH", path, body, headers)
}

// BatchRequest 批量执行REST请求
func (c *RESTClient) BatchRequest(ctx context.Context, requests []RESTRequest, maxConcurrency int) []RESTBatchResult {
	if maxConcurrency <= 0 {
		maxConcurrency = 10 // 默认并发数
	}

	results := make([]RESTBatchResult, len(requests))
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request RESTRequest) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			resp, err := c.request(ctx, request.Method, request.Path, request.Body, request.Headers)
			results[index] = RESTBatchResult{
				Index:    index,
				Response: resp,
				Error:    err,
			}
		}(i, req)
	}

	wg.Wait()
	return results
}

// BatchGET 批量GET请求
func (c *RESTClient) BatchGET(ctx context.Context, paths []string, maxConcurrency int) []RESTBatchResult {
	requests := make([]RESTRequest, len(paths))
	for i, path := range paths {
		requests[i] = RESTRequest{
			Method: "GET",
			Path:   path,
		}
	}
	return c.BatchRequest(ctx, requests, maxConcurrency)
}

// BatchPOST 批量POST请求
func (c *RESTClient) BatchPOST(ctx context.Context, pathBodyPairs []struct {
	Path string
	Body interface{}
}, maxConcurrency int) []RESTBatchResult {
	requests := make([]RESTRequest, len(pathBodyPairs))
	for i, pair := range pathBodyPairs {
		requests[i] = RESTRequest{
			Method: "POST",
			Path:   pair.Path,
			Body:   pair.Body,
		}
	}
	return c.BatchRequest(ctx, requests, maxConcurrency)
}
func (c *RESTClient) request(ctx context.Context, method, path string, body interface{}, headers map[string]string) (*RESTResponse, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 设置默认头部
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	// 设置请求特定头部
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 如果有body，设置Content-Type
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return &RESTResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// SetHeader 设置默认请求头
func (c *RESTClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

// SetHeaders 批量设置请求头
func (c *RESTClient) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.Headers[k] = v
	}
}

// UnmarshalJSON 解析JSON响应
func (r *RESTResponse) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// String 返回响应体字符串
func (r *RESTResponse) String() string {
	return string(r.Body)
}

// IsSuccess 检查响应是否成功
func (r *RESTResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
