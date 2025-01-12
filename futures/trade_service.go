package futures

import (
	"context"
	"net/http"
)

// HistoricalTradesService trades | 查询历史成交记录
type HistoricalTradesService struct {
	c      *Client
	symbol string // 交易对
	limit  *int   // 默认值:500 最大值:1000.
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
	// GET /fapi/v1/historicalTrades | 查询订单簿历史成交
	// 仅返回订单簿成交，即不会返回保险基金和自动减仓(ADL)成交
	// 仅支持返回最近3个月的数据
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/historicalTrades",
		secType:  secTypeAPIKey,
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	if s.fromID != nil {
		r.setParam("fromId", *s.fromID)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
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
	ID            int64  `json:"id"`           // 成交ID
	Price         string `json:"price"`        // 成交价格
	Quantity      string `json:"qty"`          // 成交量
	QuoteQuantity string `json:"quoteQty"`     // 成交额
	Time          int64  `json:"time"`         // 时间
	IsBuyerMaker  bool   `json:"isBuyerMaker"` // 买方是否为挂单方
}

// TradeV3 define v3 trade info
type TradeV3 struct {
	ID              int64  `json:"id"`
	Symbol          string `json:"symbol"`
	OrderID         int64  `json:"orderId"`
	Price           string `json:"price"`
	Quantity        string `json:"qty"`
	QuoteQuantity   string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}

// AggTradesService list aggregate trades | 获取近期成交(归集)
type AggTradesService struct {
	c         *Client
	symbol    string // 交易对
	fromID    *int64 // 从包含fromID的成交开始返回结果
	startTime *int64 // 从该时刻之后的成交记录开始返回结果
	endTime   *int64 // 返回该时刻为止的成交记录
	limit     *int   // 默认 500; 最大 1000
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
	// GET /fapi/v1/aggTrades | 近期成交(归集)，归集交易与逐笔交易的区别在于，同一价格、同一方向、同一时间(100ms计算)的订单簿trade会被聚合为一条
	// 接口仅支持查询最近1年的交易数据
	// 如果同时发送startTime和endTime，间隔必须小于一小时
	// 如果没有发送任何筛选参数(fromId, startTime, endTime)，默认返回最近的成交记录
	// 保险基金和自动减仓(ADL)成交不属于订单簿成交，故不会被归并聚合
	// 同时发送startTime/endTime和fromId可能导致请求超时，建议仅发送fromId或仅发送startTime和endTime
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/aggTrades",
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
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// AggTrade aggregate trades | 定义近期成交(归集)数据结构
type AggTrade struct {
	AggTradeID   int64  `json:"a"` // 归集成交ID
	Price        string `json:"p"` // 成交价
	Quantity     string `json:"q"` // 成交量
	FirstTradeID int64  `json:"f"` // 被归集的首个成交ID
	LastTradeID  int64  `json:"l"` // 被归集的末个成交ID
	Timestamp    int64  `json:"T"` // 成交时间
	IsBuyerMaker bool   `json:"m"` // 是否为主动卖出单
}

// RecentTradesService list recent trades | 获取近期订单簿成交
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
		endpoint: "/fapi/v1/trades",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// ListAccountTradeService define account trade list service | 获取某交易对的成交历史
type ListAccountTradeService struct {
	c         *Client
	symbol    string // 交易对
	startTime *int64 // 起始时间
	endTime   *int64 // 结束时间
	fromID    *int64 // 返回该fromId及之后的成交，缺省返回最近的成交
	limit     *int   // 返回的结果集数量 默认值:500 最大值:1000
}

// SetSymbol set symbol
func (s *ListAccountTradeService) SetSymbol(symbol string) *ListAccountTradeService {
	s.symbol = symbol
	return s
}

// SetStartTime set startTime
func (s *ListAccountTradeService) SetStartTime(startTime int64) *ListAccountTradeService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *ListAccountTradeService) SetEndTime(endTime int64) *ListAccountTradeService {
	s.endTime = &endTime
	return s
}

// SetFromID set fromID
func (s *ListAccountTradeService) SetFromID(fromID int64) *ListAccountTradeService {
	s.fromID = &fromID
	return s
}

// SetLimit set limit
func (s *ListAccountTradeService) SetLimit(limit int) *ListAccountTradeService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListAccountTradeService) Do(ctx context.Context, opts ...RequestOption) (res []*AccountTrade, err error) {
	// GET /fapi/v1/userTrades | 获取某交易对的成交历史
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/userTrades",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.fromID != nil {
		r.setParam("fromID", *s.fromID)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*AccountTrade{}, err
	}
	res = make([]*AccountTrade, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*AccountTrade{}, err
	}
	return res, nil
}

// AccountTrade 定义成交历史数据结构
type AccountTrade struct {
	Buyer           bool             `json:"buyer"`           // 是否是买方
	Commission      string           `json:"commission"`      // 手续费
	CommissionAsset string           `json:"commissionAsset"` // 手续费计价单位
	ID              int64            `json:"id"`              // 交易ID
	Maker           bool             `json:"maker"`           // 是否是挂单方
	OrderID         int64            `json:"orderId"`         // 订单编号
	Price           string           `json:"price"`           // 成交价
	Quantity        string           `json:"qty"`             // 成交量
	QuoteQuantity   string           `json:"quoteQty"`        // 成交额
	RealizedPnl     string           `json:"realizedPnl"`     // 实现盈亏
	Side            SideType         `json:"side"`            // 买卖方向
	PositionSide    PositionSideType `json:"positionSide"`    // 持仓方向
	Symbol          string           `json:"symbol"`          // 交易对
	Time            int64            `json:"time"`            // 时间
}
