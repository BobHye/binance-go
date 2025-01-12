package spot

import (
	"context"
	"net/http"
)

// ListTradesService list trades
type ListTradesService struct {
	c         *Client
	symbol    string
	startTime *int64
	endTime   *int64
	limit     *int   // Default 500; max 1000.
	fromID    *int64 // 返回该fromId之后的成交，缺省返回最近的成交
	orderId   *int64 // 必须要和参数symbol一起使用.
}

// Symbol set symbol
func (s *ListTradesService) Symbol(symbol string) *ListTradesService {
	s.symbol = symbol
	return s
}

// SetStartTime set starttime
func (s *ListTradesService) SetStartTime(startTime int64) *ListTradesService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endtime
func (s *ListTradesService) SetEndTime(endTime int64) *ListTradesService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *ListTradesService) SetLimit(limit int) *ListTradesService {
	s.limit = &limit
	return s
}

// SetFromID set fromID
func (s *ListTradesService) SetFromID(fromID int64) *ListTradesService {
	s.fromID = &fromID
	return s
}

// SetOrderId set OrderId
func (s *ListTradesService) SetOrderId(OrderId int64) *ListTradesService {
	s.orderId = &OrderId
	return s
}

// Do send request
func (s *ListTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*TradeV3, err error) {
	// GET /api/v3/myTrades | 获取某交易对的成交历史
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/myTrades",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.fromID != nil {
		r.setParam("fromId", *s.fromID)
	}
	if s.orderId != nil {
		r.setParam("orderId", *s.orderId)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*TradeV3{}, err
	}
	res = make([]*TradeV3, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*TradeV3{}, err
	}
	return res, nil
}

// TradeV3 define v3 trade info
type TradeV3 struct {
	ID              int64  `json:"id"`
	Symbol          string `json:"symbol"`
	OrderID         int64  `json:"orderId"`
	OrderListId     int64  `json:"orderListId"`
	Price           string `json:"price"`
	Quantity        string `json:"qty"`
	QuoteQuantity   string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
	IsIsolated      bool   `json:"isIsolated"`
}

// HistoricalTradesService trades
type HistoricalTradesService struct {
	c      *Client
	symbol string
	limit  *int   // Default 500; max 1000.
	fromID *int64 // 从哪一条成交id开始返回. 缺省返回最近的成交记录
}

// SetSymbol set symbol
func (s *HistoricalTradesService) SetSymbol(symbol string) *HistoricalTradesService {
	s.symbol = symbol
	return s
}

// SetLimit set limit
func (s *HistoricalTradesService) SetLimit(limit int) *HistoricalTradesService {
	s.limit = &limit
	return s
}

// SetFromID set fromID
func (s *HistoricalTradesService) SetFromID(fromID int64) *HistoricalTradesService {
	s.fromID = &fromID
	return s
}

// Do send request
func (s *HistoricalTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*Trade, err error) {
	// GET /api/v3/historicalTrades | 查询历史成交
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/historicalTrades",
		secType:  secTypeAPIKey,
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	if s.fromID != nil {
		r.setParam("fromId", *s.fromID)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = make([]*Trade, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return
}

// Trade define trade info
type Trade struct {
	ID            int64  `json:"id"`
	Price         string `json:"price"`
	Quantity      string `json:"qty"`
	QuoteQuantity string `json:"quoteQty"`
	Time          int64  `json:"time"`
	IsBuyerMaker  bool   `json:"isBuyerMaker"`
	IsBestMatch   bool   `json:"isBestMatch"`
	IsIsolated    bool   `json:"isIsolated"`
}

// AggTradesService list aggregate trades
type AggTradesService struct {
	c         *Client
	symbol    string
	fromID    *int64 // 从包含fromID的成交开始返回结果
	startTime *int64 // 从该时刻之后的成交记录开始返回结果
	endTime   *int64 // 返回该时刻为止的成交记录
	limit     *int   // 默认 500; 最大 1000.
}

// SetSymbol set symbol
func (s *AggTradesService) SetSymbol(symbol string) *AggTradesService {
	s.symbol = symbol
	return s
}

// SetFromID set fromID
func (s *AggTradesService) SetFromID(fromID int64) *AggTradesService {
	s.fromID = &fromID
	return s
}

// SetStartTime set startTime
func (s *AggTradesService) SetStartTime(startTime int64) *AggTradesService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *AggTradesService) SetEndTime(endTime int64) *AggTradesService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *AggTradesService) SetLimit(limit int) *AggTradesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *AggTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*AggTrade, err error) {
	// GET /api/v3/aggTrades | 近期成交(归集)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/aggTrades",
	}
	r.setParam("symbol", s.symbol)
	if s.fromID != nil {
		r.setParam("fromId", *s.fromID)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*AggTrade{}, err
	}
	res = make([]*AggTrade, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*AggTrade{}, err
	}
	return res, nil
}

// AggTrade define aggregate trade info
type AggTrade struct {
	AggTradeID       int64  `json:"a"` // 归集成交ID
	Price            string `json:"p"` // 成交价
	Quantity         string `json:"q"` // 成交量
	FirstTradeID     int64  `json:"f"` // 被归集的首个成交ID
	LastTradeID      int64  `json:"l"` // 被归集的末个成交ID
	Timestamp        int64  `json:"T"` // 成交时间
	IsBuyerMaker     bool   `json:"m"` // 是否为主动卖出单
	IsBestPriceMatch bool   `json:"M"` // 是否为最优撮合单(可忽略，目前总为最优撮合)
}

// RecentTradesService list recent trades
type RecentTradesService struct {
	c      *Client
	symbol string
	limit  *int
}

// SetSymbol set symbol
func (s *RecentTradesService) SetSymbol(symbol string) *RecentTradesService {
	s.symbol = symbol
	return s
}

// SetLimit set limit
func (s *RecentTradesService) SetLimit(limit int) *RecentTradesService {
	s.limit = &limit
	return s
}

// Do send request
func (s *RecentTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*Trade, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v1/trades",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*Trade{}, err
	}
	res = make([]*Trade, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*Trade{}, err
	}
	return res, nil
}
