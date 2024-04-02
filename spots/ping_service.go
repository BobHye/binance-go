package spots

import (
	"context"
	"net/http"
)

// 测试服务器连通性 PING

type PingService struct {
	c *Client
}

// Do 发送请求
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ping",
	}
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}
