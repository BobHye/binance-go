package spot

import (
	"context"
	"github.com/BobHye/binance-go/common"
	"net/http"
)

// ListBookTickersService list best price/qty on the order book for a symbol or symbols
type ListBookTickersService struct {
	c      *Client
	symbol *string // 数 `symbol` 和 `symbols` 不可以一起使用 如果都不提供, 所有symbol的bookTicker数据都会返回.
}

// SetSymbol set symbol
func (s *ListBookTickersService) SetSymbol(symbol string) *ListBookTickersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListBookTickersService) Do(ctx context.Context, opts ...RequestOption) (res []*BookTicker, err error) {
	// GET /api/v3/ticker/bookTicker | 返回当前最优的挂单(最高买单，最低卖单)
	// 不发送交易对参数，则会返回所有交易对信息
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/ticker/bookTicker",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	data = common.ToJSONList(data)
	if err != nil {
		return []*BookTicker{}, err
	}
	res = make([]*BookTicker, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*BookTicker{}, err
	}
	return res, nil
}

// BookTicker define book ticker info
type BookTicker struct {
	Symbol      string `json:"symbol"`
	BidPrice    string `json:"bidPrice"` // 最优买单价
	BidQuantity string `json:"bidQty"`   // 挂单量
	AskPrice    string `json:"askPrice"` // 最优卖单价
	AskQuantity string `json:"askQty"`   // 挂单量
}

// ListPricesService list latest price for a symbol or symbols
type ListPricesService struct {
	c       *Client
	symbol  *string  // 参数 `symbol` 和 `symbols` 不可以一起使用 如果都不提供, 所有symbol的价格数据都会返回
	symbols []string // symbols参数可接受的格式： ["BTCUSDT","BNBUSDT"] 或 %5B%22BTCUSDT%22,%22BNBUSDT%22%5D
}

// SetSymbol set symbol
func (s *ListPricesService) SetSymbol(symbol string) *ListPricesService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListPricesService) Do(ctx context.Context, opts ...RequestOption) (res []*SymbolPrice, err error) {
	// GET /api/v3/ticker/price | 返回最近价格
	// 不发送交易对参数，则会返回所有交易对信息
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/ticker/price",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	} else if s.symbols != nil {
		s, _ := json.Marshal(s.symbols)
		r.setParam("symbols", string(s))
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*SymbolPrice{}, err
	}
	data = common.ToJSONList(data)
	res = make([]*SymbolPrice, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*SymbolPrice{}, err
	}
	return res, nil
}

// SymbolPrice define symbol and price pair
type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// ListPriceChangeStatsService show stats of price change in last 24 hours for all symbols
type ListPriceChangeStatsService struct {
	c      *Client
	symbol *string
}

// SetSymbol set symbol
func (s *ListPriceChangeStatsService) SetSymbol(symbol string) *ListPriceChangeStatsService {
	s.symbol = &symbol
	return s
}

// Symbols set symbols
func (s *ListPricesService) Symbols(symbols []string) *ListPricesService {
	s.symbols = symbols
	return s
}

// Do send request
func (s *ListPriceChangeStatsService) Do(ctx context.Context, opts ...RequestOption) (res []*PriceChangeStats, err error) {
	// GET /api/v3/ticker/24hr | 24hr价格变动情况
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/ticker/24hr",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	data = common.ToJSONList(data)
	res = make([]*PriceChangeStats, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PriceChangeStats define price change stats
type PriceChangeStats struct {
	Symbol             string `json:"symbol"` // 交易对
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"` // 间隔收盘价
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	AskPrice           string `json:"askPrice"`
	OpenPrice          string `json:"openPrice"`   // 间隔开盘价
	HighPrice          string `json:"highPrice"`   // 间隔最高价
	LowPrice           string `json:"lowPrice"`    // 间隔最低价
	Volume             string `json:"volume"`      // 总交易量 (base asset)
	QuoteVolume        string `json:"quoteVolume"` // 总交易量 (quote asset)
	OpenTime           int64  `json:"openTime"`    // ticker间隔的开始时间
	CloseTime          int64  `json:"closeTime"`   // ticker间隔的结束时间
	FristID            int64  `json:"firstId"`     // 首笔成交id
	LastID             int64  `json:"lastId"`      // 末笔成交id
	Count              int64  `json:"count"`       // 成交笔数
}

// AveragePriceService show current average price for a symbol
type AveragePriceService struct {
	c      *Client
	symbol string
}

// SetSymbol set symbol
func (s *AveragePriceService) SetSymbol(symbol string) *AveragePriceService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *AveragePriceService) Do(ctx context.Context, opts ...RequestOption) (res *AvgPrice, err error) {
	// GET /api/v3/avgPrice | 当前平均价格
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/avgPrice",
	}
	r.setParam("symbol", s.symbol)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	res = new(AvgPrice)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AvgPrice define average price
type AvgPrice struct {
	Mins      int64  `json:"mins"`
	Price     string `json:"price"`
	CloseTime int64  `json:"closeTime"`
}

type ListSymbolTickerService struct {
	c          *Client
	symbol     *string
	symbols    []string
	windowSize *string
}

type SymbolTicker struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`        // 价格变化
	PriceChangePercent string `json:"priceChangePercent"` // 价格变化百分比
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	LastPrice          string `json:"lastPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`  // ticker的开始时间
	CloseTime          int64  `json:"closeTime"` // ticker的结束时间
	FirstId            int64  `json:"firstId"`   // 统计时间内的第一笔trade id
	LastId             int64  `json:"lastId"`
	Count              int64  `json:"count"` // 统计时间内交易笔数
}

func (s *ListSymbolTickerService) SetSymbol(symbol string) *ListSymbolTickerService {
	s.symbol = &symbol
	return s
}

func (s *ListSymbolTickerService) SetSymbols(symbols []string) *ListSymbolTickerService {
	s.symbols = symbols
	return s
}

// SetWindowSize Defaults to 1d if no parameter provided
// Supported windowSize values:
// - 1m,2m....59m for minutes
// - 1h, 2h....23h - for hours
// - 1d...7d - for days
// Units cannot be combined (e.g. 1d2h is not allowed).
// Reference: https://binance-docs.github.io/apidocs/spot/en/#rolling-window-price-change-statistics
func (s *ListSymbolTickerService) SetWindowSize(windowSize string) *ListSymbolTickerService {
	s.windowSize = &windowSize
	return s
}

func (s *ListSymbolTickerService) Do(ctx context.Context, opts ...RequestOption) (res []*SymbolTicker, err error) {
	// GET /api/v3/ticker | 滚动窗口价格变动统计
	// 注意: 此接口和 GET /api/v3/ticker/24hr 有所不同.
	// 此接口统计的时间范围比请求的windowSize多不超过59999ms.
	// 接口的 openTime 是某一分钟的起始，而结束是当前的时间. 所以实际的统计区间会比请求的时间窗口多不超过59999ms.
	// 比如, 结束时间 closeTime 是 1641287867099 (January 04, 2022 09:17:47:099 UTC) , windowSize 为 1d. 那么开始时间 openTime 则为 1641201420000 (January 3, 2022, 09:17:00 UTC)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/ticker",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	} else if s.symbols != nil {
		s, _ := json.Marshal(s.symbols)
		r.setParam("symbols", string(s))
	}

	if s.windowSize != nil {
		r.setParam("windowSize", *s.windowSize)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	data = common.ToJSONList(data)
	if err != nil {
		return []*SymbolTicker{}, err
	}
	res = make([]*SymbolTicker, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*SymbolTicker{}, err
	}
	return res, nil
}
