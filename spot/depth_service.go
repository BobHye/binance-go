package spot

import (
	"context"
	"github.com/BobHye/binance-go/common"
	"net/http"
)

// DepthService show depth info
type DepthService struct {
	c      *Client
	symbol string
	limit  *int // 默认 100; 最大 5000. 可选值:[5, 10, 20, 50, 100, 500, 1000, 5000] 如果 limit > 5000, 最多返回5000条数据.
}

// SetSymbol set symbol
func (s *DepthService) SetSymbol(symbol string) *DepthService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *DepthService) Limit(limit int) *DepthService {
	s.limit = &limit
	return s
}

// Do send request
func (s *DepthService) Do(ctx context.Context, opts ...RequestOption) (res *DepthResponse, err error) {
	// GET /api/v3/depth | 深度信息
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/depth",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(DepthResponse)
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DepthResponse define depth info with bids and asks
type DepthResponse struct {
	LastUpdateID int64 `json:"lastUpdateId"`
	Bids         []Bid `json:"bids"`
	Asks         []Ask `json:"asks"`
}

// Ask is a type alias for PriceLevel.
type Ask = common.PriceLevel

// Bid is a type alias for PriceLevel.
type Bid = common.PriceLevel
