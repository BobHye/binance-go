package futures

import (
	"context"
	"github.com/BobHye/binance-go/api"
	"github.com/BobHye/binance-go/common"
	"net/http"
)

// BookTicker 定义最优挂单信息
type BookTicker struct {
	Symbol      string `json:"symbol"`
	BidPrice    string `json:"bidPrice"`
	BidQuantity string `json:"bidQty"`
	AskPrice    string `json:"askPrice"`
	AskQuantity string `json:"askQty"`
}

// ListBookTickersService 获取当前最优挂单，不发送交易对参数，则会返回所有交易对信息
type ListBookTickersService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *ListBookTickersService) Symbol(symbol string) *ListBookTickersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListBookTickersService) Do(ctx context.Context, opts ...api.RequestOption) (res []*BookTicker, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/ticker/bookTicker",
	}
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// SymbolPrice 定义最新价格结构
type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// ListPricesService 查询最新价格，不发送交易对参数，则会返回所有交易对信息
type ListPricesService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *ListPricesService) Symbol(symbol string) *ListPricesService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListPricesService) Do(ctx context.Context, opts ...api.RequestOption) (res []*SymbolPrice, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/ticker/price",
	}
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// PriceChangeStats 定义价格变动情况数据结构
type PriceChangeStats struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQuantity       string `json:"lastQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstID            int64  `json:"firstId"`
	LastID             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

// ListPriceChangeStatsService 24hr价格变动情况，不发送交易对参数，则会返回所有交易对信息
type ListPriceChangeStatsService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *ListPriceChangeStatsService) Symbol(symbol string) *ListPriceChangeStatsService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListPriceChangeStatsService) Do(ctx context.Context, opts ...api.RequestOption) (res []*PriceChangeStats, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/ticker/24hr",
	}
	if s.symbol != nil {
		r.SetParam("symbol", *s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
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
