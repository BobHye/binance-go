package futures

import (
	"context"
	"github.com/BobHye/binance-go/common"
	"net/http"
)

// 深度信息

// Ask is a type alias for PriceLevel.
type Ask = common.PriceLevel

// Bid is a type alias for PriceLevel.
type Bid = common.PriceLevel

// DepthResponse define depth info with bids and asks
type DepthResponse struct {
	LastUpdateID int64 `json:"lastUpdateId"`
	Time         int64 `json:"E"`    // 消息时间
	TradeTime    int64 `json:"T"`    // 撮合引擎时间
	Bids         []Bid `json:"bids"` // 买方 价格/数量
	Asks         []Ask `json:"asks"` // 卖方 价格/数量
}

// DepthService show depth info
type DepthService struct {
	c      *Client
	symbol string
	limit  *int
}

// SetSymbol set symbol
func (s *DepthService) SetSymbol(symbol string) *DepthService {
	s.symbol = symbol
	return s
}

// SetLimit set limit
func (s *DepthService) SetLimit(limit int) *DepthService {
	s.limit = &limit
	return s
}

// Do send request
func (s *DepthService) Do(ctx context.Context, opts ...RequestOption) (res *DepthResponse, err error) {
	// GET /fapi/v1/depth | 交易对深度信息
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/depth",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res = new(DepthResponse)
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
