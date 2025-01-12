package spot

import (
	"context"
	"net/http"
	"strings"
)

// ExchangeInfoService exchange info service
type ExchangeInfoService struct {
	c       *Client
	symbol  string
	symbols string
}

// SetSymbol set symbol
func (s *ExchangeInfoService) SetSymbol(symbol string) *ExchangeInfoService {
	s.symbols = symbol
	return s
}

// SetSymbols set symbols
func (s *ExchangeInfoService) SetSymbols(symbols ...string) *ExchangeInfoService {
	if len(symbols) == 0 {
		s.symbols = "[]"
	} else {
		s.symbols = "[\"" + strings.Join(symbols, "\",\"") + "\"]"
	}
	return s
}

// Do send request
func (s *ExchangeInfoService) Do(ctx context.Context, opts ...RequestOption) (res *ExchangeInfo, err error) {
	// GET /api/v3/exchangeInfo | 获取此时的交易规范信息
	// 如果参数 symbol 或者 symbols 提供的交易对不存在, 系统会返回错误并提示交易对不正确。
	// 所有的参数都是可选的。
	// permissions 可以支持单个或多个值（例如 SPOT, ["MARGIN","LEVERAGED"]）。此参数不能与 symbol 或 symbols 组合使用。
	// 如果未提供 permissions 参数，那么所有具有 SPOT、MARGIN 或 LEVERAGED 权限的交易对都将被返回。
	// 要显示具有任何权限的交易对，您需要在 permissions 中明确指定它们：（例如 ["SPOT","MARGIN",...])。有关完整列表，请参阅 可用权限
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/exchangeInfo",
		secType:  secTypeNone,
	}
	m := params{}
	if s.symbol != "" {
		m["symbol"] = s.symbol
	}
	if len(s.symbols) != 0 {
		m["symbols"] = s.symbol
	}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(ExchangeInfo)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ExchangeInfo exchange info
type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int64         `json:"serverTime"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []Symbol      `json:"symbols"`
}

// RateLimit struct
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"` // 每分钟调用的所有接口权重之和不得超过1200
}

// Symbol market symbol
type Symbol struct {
	Symbol                     string                   `json:"symbol"`
	Status                     string                   `json:"status"`
	BaseAsset                  string                   `json:"baseAsset"`
	BaseAssetPrecision         int                      `json:"baseAssetPrecision"`
	QuoteAsset                 string                   `json:"quoteAsset"`
	QuotePrecision             int                      `json:"quotePrecision"`
	QuoteAssetPrecision        int                      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int32                    `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int32                    `json:"quoteCommissionPrecision"`
	OrderTypes                 []string                 `json:"orderTypes"`
	IcebergAllowed             bool                     `json:"icebergAllowed"`
	OcoAllowed                 bool                     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool                     `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool                     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool                     `json:"isMarginTradingAllowed"`
	Filters                    []map[string]interface{} `json:"filters"`
	Permissions                []string                 `json:"permissions"`
}

// LotSizeFilter define lot size filter of symbol
type LotSizeFilter struct {
	MaxQuantity string `json:"maxQty"`
	MinQuantity string `json:"minQty"`
	StepSize    string `json:"stepSize"`
}

// PriceFilter define price filter of symbol
type PriceFilter struct {
	MaxPrice string `json:"maxPrice"`
	MinPrice string `json:"minPrice"`
	TickSize string `json:"tickSize"`
}

// PercentPriceFilter define percent price filter of symbol
type PercentPriceFilter struct {
	AveragePriceMins int    `json:"avgPriceMins"`
	MultiplierUp     string `json:"multiplierUp"`
	MultiplierDown   string `json:"multiplierDown"`
}

// MinNotionalFilter define min notional filter of symbol
type MinNotionalFilter struct {
	MinNotional      string `json:"minNotional"`
	AveragePriceMins int    `json:"avgPriceMins"`
	ApplyToMarket    bool   `json:"applyToMarket"`
}

// IcebergPartsFilter define iceberg part filter of symbol
type IcebergPartsFilter struct {
	Limit int `json:"limit"`
}

// MarketLotSizeFilter define market lot size filter of symbol
type MarketLotSizeFilter struct {
	MaxQuantity string `json:"maxQty"`
	MinQuantity string `json:"minQty"`
	StepSize    string `json:"stepSize"`
}

// MaxNumAlgoOrdersFilter define max num algo orders filter of symbol
type MaxNumAlgoOrdersFilter struct {
	MaxNumAlgoOrders int `json:"maxNumAlgoOrders"`
}

// LotSizeFilter return lot size filter of symbol
func (s *Symbol) LotSizeFilter() *LotSizeFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeLotSize) {
			f := &LotSizeFilter{}
			if i, ok := filter["maxQty"]; ok {
				f.MaxQuantity = i.(string)
			}
			if i, ok := filter["minQty"]; ok {
				f.MinQuantity = i.(string)
			}
			if i, ok := filter["stepSize"]; ok {
				f.StepSize = i.(string)
			}
			return f
		}
	}
	return nil
}

// PriceFilter return price filter of symbol
func (s *Symbol) PriceFilter() *PriceFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypePriceFilter) {
			f := &PriceFilter{}
			if i, ok := filter["maxPrice"]; ok {
				f.MaxPrice = i.(string)
			}
			if i, ok := filter["minPrice"]; ok {
				f.MinPrice = i.(string)
			}
			if i, ok := filter["tickSize"]; ok {
				f.TickSize = i.(string)
			}
			return f
		}
	}
	return nil
}

// PercentPriceFilter return percent price filter of symbol
func (s *Symbol) PercentPriceFilter() *PercentPriceFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypePercentPrice) {
			f := &PercentPriceFilter{}
			if i, ok := filter["avgPriceMins"]; ok {
				f.AveragePriceMins = int(i.(float64))
			}
			if i, ok := filter["multiplierUp"]; ok {
				f.MultiplierUp = i.(string)
			}
			if i, ok := filter["multiplierDown"]; ok {
				f.MultiplierDown = i.(string)
			}
			return f
		}
	}
	return nil
}

// MinNotionalFilter return min notional filter of symbol
func (s *Symbol) MinNotionalFilter() *MinNotionalFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeMinNotional) {
			f := &MinNotionalFilter{}
			if i, ok := filter["minNotional"]; ok {
				f.MinNotional = i.(string)
			}
			if i, ok := filter["avgPriceMins"]; ok {
				f.AveragePriceMins = int(i.(float64))
			}
			if i, ok := filter["applyToMarket"]; ok {
				f.ApplyToMarket = i.(bool)
			}
			return f
		}
	}
	return nil
}

// IcebergPartsFilter return iceberg part filter of symbol
func (s *Symbol) IcebergPartsFilter() *IcebergPartsFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeIcebergParts) {
			f := &IcebergPartsFilter{}
			if i, ok := filter["limit"]; ok {
				f.Limit = int(i.(float64))
			}
			return f
		}
	}
	return nil
}

// MarketLotSizeFilter return market lot size filter of symbol
func (s *Symbol) MarketLotSizeFilter() *MarketLotSizeFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeMarketLotSize) {
			f := &MarketLotSizeFilter{}
			if i, ok := filter["maxQty"]; ok {
				f.MaxQuantity = i.(string)
			}
			if i, ok := filter["minQty"]; ok {
				f.MinQuantity = i.(string)
			}
			if i, ok := filter["stepSize"]; ok {
				f.StepSize = i.(string)
			}
			return f
		}
	}
	return nil
}

// MaxNumAlgoOrdersFilter return max num algo orders filter of symbol
func (s *Symbol) MaxNumAlgoOrdersFilter() *MaxNumAlgoOrdersFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeMaxNumAlgoOrders) {
			f := &MaxNumAlgoOrdersFilter{}
			if i, ok := filter["maxNumAlgoOrders"]; ok {
				f.MaxNumAlgoOrders = int(i.(float64))
			}
			return f
		}
	}
	return nil
}
