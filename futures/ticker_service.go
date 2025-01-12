package futures

import (
	"context"
	"github.com/BobHye/binance-go/common"
	"net/http"
)

// ListBookTickersService list best price/qty on the order book for a symbol or symbols | 获取当前最优挂单，不发送交易对参数，则会返回所有交易对信息
type ListBookTickersService struct {
	c      *Client
	symbol *string // 交易对
}

// SetSymbol set symbol
func (s *ListBookTickersService) SetSymbol(symbol string) *ListBookTickersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListBookTickersService) Do(ctx context.Context, opts ...RequestOption) (res []*BookTicker, err error) {
	// GET /fapi/v1/ticker/bookTicker | 返回当前最优的挂单(最高买单，最低卖单)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ticker/bookTicker",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
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

// BookTicker 定义最优挂单信息
type BookTicker struct {
	Symbol      string  `json:"symbol"`          // 交易对
	BidPrice    float64 `json:"bidPrice,string"` //最优买单价
	BidQuantity float64 `json:"bidQty,string"`   //挂单量
	AskPrice    float64 `json:"askPrice,string"` //最优卖单价
	AskQuantity float64 `json:"askQty,string"`   //挂单量
	Time        int64   `json:"time"`            // 撮合引擎时间
}

// ListPricesService list latest price for a symbol or symbols | 查询最新价格，不发送交易对参数，则会返回所有交易对信息
type ListPricesService struct {
	c      *Client
	symbol *string // 交易对
}

// SetSymbol set symbol
func (s *ListPricesService) SetSymbol(symbol string) *ListPricesService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListPricesService) Do(ctx context.Context, opts ...RequestOption) (res []*SymbolPrice, err error) {
	// GET /fapi/v1/ticker/price | 返回最近价格
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ticker/price",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
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

// SymbolPrice 定义最新价格结构
type SymbolPrice struct {
	Symbol string  `json:"symbol"`       // 交易对
	Price  float64 `json:"price,string"` // 价格
	Time   int64   `json:"time"`         // 撮合引擎时间
}

// ListPriceChangeStatsService 24hr价格变动情况，不发送交易对参数，则会返回所有交易对信息
type ListPriceChangeStatsService struct {
	c      *Client
	symbol *string // 交易对
}

// Symbol set symbol
func (s *ListPriceChangeStatsService) Symbol(symbol string) *ListPriceChangeStatsService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListPriceChangeStatsService) Do(ctx context.Context, opts ...RequestOption) (res []*PriceChangeStats, err error) {
	// GET /fapi/v1/ticker/24hr | 24hr价格变动情况(请注意，不携带symbol参数会返回全部交易对数据，不仅数据庞大，而且权重极高)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ticker/24hr",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
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

// PriceChangeStats define price change stats | 定义价格变动情况数据结构
type PriceChangeStats struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`        //24小时价格变动
	PriceChangePercent string `json:"priceChangePercent"` //24小时价格变动百分比
	WeightedAvgPrice   string `json:"weightedAvgPrice"`   //加权平均价
	//PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice    string `json:"lastPrice"`   //最近一次成交价
	LastQuantity string `json:"lastQty"`     //最近一次成交额
	OpenPrice    string `json:"openPrice"`   //24小时内第一次成交的价格
	HighPrice    string `json:"highPrice"`   //24小时最高价
	LowPrice     string `json:"lowPrice"`    //24小时最低价
	Volume       string `json:"volume"`      //24小时成交量
	QuoteVolume  string `json:"quoteVolume"` //24小时成交额
	OpenTime     int64  `json:"openTime"`    //24小时内，第一笔交易的发生时间
	CloseTime    int64  `json:"closeTime"`   //24小时内，最后一笔交易的发生时间
	FirstID      int64  `json:"firstId"`     // 首笔成交id
	LastID       int64  `json:"lastId"`      // 末笔成交id
	Count        int64  `json:"count"`       // 成交笔数
}
