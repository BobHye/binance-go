package delivery

import (
	"context"
	"net/http"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// GET /dapi/v1/ping | 测试服务器连通性 PING(测试能否联通)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/ping",
	}
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

type ServiceTime struct {
	ServiceTime int64 `json:"serviceTime"`
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...RequestOption) (int64, error) {
	// GET /dapi/v1/time | 获取服务器时间
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/time",
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, nil
	}
	var res ServiceTime
	err = json.Unmarshal(data, &res)
	return res.ServiceTime, err
}

// SetServerTimeService set server time
type SetServerTimeService struct {
	c *Client
}

// Do send request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx)
	if err != nil {
		return 0, err
	}
	timeOffset = currentTimestamp() - serverTime
	s.c.TimeOffset = timeOffset
	return timeOffset, nil
}
