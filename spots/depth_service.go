package spots

import (
	"binance-go/common"
	"context"
	"net/http"
)

// 获取深度信息

type DepthService struct {
	c      *Client
	symbol string
	limit  *int
}

// Symbol 设置交易对
func (s *DepthService) Symbol(symbol string) *DepthService {
	s.symbol = symbol
	return s
}

// Limit 设置limit 默认 500; 可选值:[5, 10, 20, 50, 100, 500, 1000]
func (s *DepthService) Limit(limit int) *DepthService {
	s.limit = &limit
	return s
}

// Do 发送请求
func (s *DepthService) Do(ctx context.Context, opts ...RequestOption) (res *DepthService, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/depth",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(DepthService)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

type Ask = common.PriceLevel
type Bid = common.PriceLevel

type DepthResponse struct {
	LastUpdateID int64 `json:"lastUpdateId"`
	Bids         []Bid `json:"bids"`
	Asks         []Ask `json:"asks"`
}
