package spots

import (
	"context"
	"net/http"
)

type ServerTime struct {
	ServerTime int64 `json:"serverTime"`
}

type ServerTimeService struct {
	c *Client
}

// Do 发送请求
func (s *ServerTimeService) Do(ctx context.Context, opts ...RequestOption) (int64, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/time",
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	var res ServerTime
	err = json.Unmarshal(data, &res)
	return res.ServerTime, err
}

type SetServerTimeService struct {
	c *Client
}

// Do 发送请求
func (s *SetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx)
	if err != nil {
		return 0, err
	}
	timeOffset = currentTimestamp() - serverTime
	s.c.TimeOffset = timeOffset
	return timeOffset, nil
}
