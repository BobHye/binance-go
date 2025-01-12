package delivery

import (
	"context"
	"net/http"
)

// BookTicker define book ticker info
type BookTicker struct {
	Symbol      string `json:"symbol"`
	Pair        string `json:"pair"`
	BidPrice    string `json:"bidPrice"`
	BidQuantity string `json:"bidQty"`
	AskPrice    string `json:"askPrice"`
	AskQuantity string `json:"askQty"`
}

type ListBookTickersService struct {
	c      *Client
	symbol *string
	pair   *string
}

// SetSymbol set symbol.
func (s *ListBookTickersService) SetSymbol(symbol string) *ListBookTickersService {
	s.symbol = &symbol
	return s
}

// SetPair set pari
func (s *ListBookTickersService) SetPair(pair string) *ListBookTickersService {
	s.pair = &pair
	return s
}

// Do send request
func (s *ListBookTickersService) Do(ctx context.Context, opts ...RequestOption) (res []*BookTicker, err error) {
	// GET /dapi/v1/ticker/bookTicker | 返回当前最优的挂单(最高买单,最低卖单)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/ticker/bookTicker",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	if s.pair != nil {
		r.setParam("pair", *s.pair)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
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

// SymbolPrice define symbol, price and pair
type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Pair   string `json:"ps"`
	Price  string `json:"price"`
}

// ListPriceChangeStatsService show stats of price change in last 24 hours for single symbol, all symbols or pairs of symbols.
// 显示过去 24 小时内单个交易对、所有交易对或交易对的价格变化统计数据。
type ListPriceChangeStatsService struct {
	c      *Client
	symbol *string
	pair   *string
}

// SetSymbol set symbol
func (s *ListPriceChangeStatsService) SetSymbol(symbol string) *ListPriceChangeStatsService {
	s.symbol = &symbol
	return s
}

// SetPair set pair
func (s *ListPriceChangeStatsService) SetPair(pair string) *ListPriceChangeStatsService {
	s.pair = &pair
	return s
}

type PriceChangeStats struct {
	Symbol             string `json:"symbol"`
	Pair               string `json:"pair"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	LastPrice          string `json:"lastPrice"`
	LastQuantity       string `json:"lastQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	BaseVolume         string `json:"baseVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstID            int64  `json:"firstId"`
	LastID             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

// Do send request
func (s *ListPriceChangeStatsService) Do(ctx context.Context, opts ...RequestOption) (res []*PriceChangeStats, err error) {
	// GET /dapi/v1/ticker/24hr | 24hr价格变动情况(请注意,不携带symbol参数会返回全部交易对数据,不仅数据庞大,而且权重极高)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/ticker/24hr",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	if s.pair != nil {
		r.setParam("pair", *s.pair)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	res = make([]*PriceChangeStats, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type ListPricesService struct {
	c      *Client
	symbol *string
	pair   *string
}

// SetSymbol set symbol
func (s *ListPricesService) SetSymbol(symbol string) *ListPricesService {
	s.symbol = &symbol
	return s
}

// SetPair set pair
func (s *ListPricesService) SetPair(pair string) *ListPricesService {
	s.pair = &pair
	return s
}

// Do send request
func (s *ListPricesService) Do(ctx context.Context, opts ...RequestOption) (res []*SymbolPrice, err error) {
	// GET /dapi/v1/ticker/price | 返回最近价格
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/ticker/price",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	if s.pair != nil {
		r.setParam("pair", *s.pair)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*SymbolPrice{}, err
	}
	res = make([]*SymbolPrice, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*SymbolPrice{}, err
	}
	return res, nil
}
