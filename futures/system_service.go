package futures

import (
	"context"
	"github.com/BobHye/binance-go/api"
	"net/http"
)

type ServerTime struct {
	ServerTime int64 `json:"serverTime"`
}

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...api.RequestOption) (err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/ping",
	}
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...api.RequestOption) (int64, error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/time",
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	var res ServerTime
	err = json.Unmarshal(data, &res)
	return res.ServerTime, err
}

// SetServerTimeService set server time
type SetServerTimeService struct {
	c *Client
}

// Do send request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...api.RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx)
	if err != nil {
		return 0, err
	}
	timeOffset = currentTimestamp() - serverTime
	s.c.TimeOffset = timeOffset
	return timeOffset, nil
}
