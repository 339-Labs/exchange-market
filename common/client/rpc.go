package client

import (
	"context"
	"fmt"
	"net/rpc"
	"sync"
)

// RPCRequest RPC请求结构
type RPCRequest struct {
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
}

// RPCBatchResult 批量RPC调用结果
type RPCBatchResult struct {
	Index   int
	Request *RPCRequest
	Error   error
}

// RPCClient RPC客户端结构
type RPCClient struct {
	client *rpc.Client
	addr   string
}

type RPC interface {
	Close() error
	Call(serviceMethod string, args interface{}, reply interface{}) error
	CallWithContext(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
	BatchCall(ctx context.Context, requests []RPCRequest, maxConcurrency int) []RPCBatchResult
	AsyncCall(serviceMethod string, args interface{}, reply interface{}) <-chan *rpc.Call
	BatchAsyncCall(requests []RPCRequest, maxConcurrency int) <-chan RPCBatchResult
	BatchCallSameMethod(ctx context.Context, serviceMethod string, argsList []interface{}, maxConcurrency int) []RPCBatchResult
	GetAddress() string
}

// NewRPCClient 创建RPC客户端
func NewRPCClient(network, address string) (RPC, error) {
	client, err := rpc.Dial(network, address)
	if err != nil {
		return nil, fmt.Errorf("dial rpc server: %w", err)
	}
	return &RPCClient{
		client: client,
		addr:   address,
	}, nil
}

// Call 同步调用RPC方法
func (c *RPCClient) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return c.client.Call(serviceMethod, args, reply)
}

// CallWithContext 带上下文的同步调用
func (c *RPCClient) CallWithContext(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	// 创建一个通道来接收调用结果
	call := c.client.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1))

	select {
	case <-ctx.Done():
		return ctx.Err()
	case rpcCall := <-call.Done:
		return rpcCall.Error
	}
}

// AsyncCall 异步调用RPC方法
func (c *RPCClient) AsyncCall(serviceMethod string, args interface{}, reply interface{}) <-chan *rpc.Call {
	return c.client.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done
}

// BatchCall 批量RPC调用
func (c *RPCClient) BatchCall(ctx context.Context, requests []RPCRequest, maxConcurrency int) []RPCBatchResult {
	if maxConcurrency <= 0 {
		maxConcurrency = 10 // 默认并发数
	}

	results := make([]RPCBatchResult, len(requests))
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request RPCRequest) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := c.CallWithContext(ctx, request.ServiceMethod, request.Args, request.Reply)
			results[index] = RPCBatchResult{
				Index:   index,
				Request: &request,
				Error:   err,
			}
		}(i, req)
	}

	wg.Wait()
	return results
}

// BatchAsyncCall 批量异步RPC调用
func (c *RPCClient) BatchAsyncCall(requests []RPCRequest, maxConcurrency int) <-chan RPCBatchResult {
	if maxConcurrency <= 0 {
		maxConcurrency = 10
	}

	resultChan := make(chan RPCBatchResult, len(requests))
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request RPCRequest) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			callChan := c.AsyncCall(request.ServiceMethod, request.Args, request.Reply)
			call := <-callChan

			resultChan <- RPCBatchResult{
				Index:   index,
				Request: &request,
				Error:   call.Error,
			}
		}(i, req)
	}

	// 等待所有请求完成后关闭通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

// BatchCallSameMethod 批量调用相同方法
func (c *RPCClient) BatchCallSameMethod(ctx context.Context, serviceMethod string, argsList []interface{}, maxConcurrency int) []RPCBatchResult {
	requests := make([]RPCRequest, len(argsList))
	for i, args := range argsList {
		// 为每个请求创建独立的reply
		var reply interface{}
		requests[i] = RPCRequest{
			ServiceMethod: serviceMethod,
			Args:          args,
			Reply:         &reply,
		}
	}
	return c.BatchCall(ctx, requests, maxConcurrency)
}
func (c *RPCClient) Close() error {
	return c.client.Close()
}

// GetAddress 获取RPC服务器地址
func (c *RPCClient) GetAddress() string {
	return c.addr
}
