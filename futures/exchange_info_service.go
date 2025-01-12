package futures

import (
	"context"
	"net/http"
	"strconv"
)

// 获取统一账户交易规则

// ExchangeInfoService exchange info service
type ExchangeInfoService struct {
	c *Client
}

// Do send request
func (s *ExchangeInfoService) Do(ctx context.Context, opts ...RequestOption) (res *ExchangeInfo, err error) {
	// GET /fapi/v1/exchangeInfo | 获取交易规则和交易对
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/exchangeInfo",
		secType:  secTypeNone,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(ExchangeInfo)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ExchangeInfo exchange info
type ExchangeInfo struct {
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	RateLimits      []RateLimit   `json:"rateLimits"` // API访问的限制
	ServerTime      int64         `json:"serverTime"` // 请忽略。如果需要获取当前系统时间，请查询接口 “GET /fapi/v1/time”
	Assets          []Asset       `json:"assets"`     // 资产信息
	Symbols         []Symbol      `json:"symbols"`    // 交易对信息
	Timezone        string        `json:"timezone"`   // 服务器所用的时间区域
}

// RateLimit API访问的限制 struct
type RateLimit struct {
	Interval      string `json:"interval"`      // 按照分钟计算
	IntervalNum   int64  `json:"intervalNum"`   // 按照1分钟计算
	Limit         int64  `json:"limit"`         // 上限次数
	RateLimitType string `json:"rateLimitType"` // 按照访问权重来计算
}

// Asset 资产信息 struct
type Asset struct {
	Asset             string `json:"asset"`
	MarginAvailable   bool   `json:"marginAvailable"`   // 是否可用作保证金
	AutoAssetExchange int64  `json:"autoAssetExchange"` // 保证金资产自动兑换阈值
}

// Symbol 交易对信息 struct
type Symbol struct {
	Symbol                string                   `json:"symbol"`                // 交易对
	Pair                  string                   `json:"pair"`                  // 标的交易对
	ContractType          ContractType             `json:"contractType"`          // 合约类型
	DeliveryDate          int64                    `json:"deliveryDate"`          // 交割日期
	OnboardDate           int64                    `json:"onboardDate"`           // 上线日期
	Status                string                   `json:"status"`                // 交易对状态
	MaintMarginPercent    string                   `json:"maintMarginPercent"`    // 请忽略
	RequiredMarginPercent string                   `json:"requiredMarginPercent"` // 请忽略
	BaseAsset             string                   `json:"baseAsset"`             // 标的资产
	QuoteAsset            string                   `json:"quoteAsset"`            // 报价资产
	MarginAsset           string                   `json:"marginAsset"`           // 保证金资产
	PricePrecision        int                      `json:"pricePrecision"`        // 价格小数点位数(仅作为系统精度使用，注意同tickSize 区分）
	QuantityPrecision     int                      `json:"quantityPrecision"`     // 数量小数点位数(仅作为系统精度使用，注意同stepSize 区分）
	BaseAssetPrecision    int                      `json:"baseAssetPrecision"`    // 标的资产精度
	QuotePrecision        int                      `json:"quotePrecision"`        // 报价资产精度
	UnderlyingType        string                   `json:"underlyingType"`
	UnderlyingSubType     []string                 `json:"underlyingSubType"`
	SettlePlan            int64                    `json:"settlePlan"`
	TriggerProtect        string                   `json:"triggerProtect"` // 开启"priceProtect"的条件订单的触发阈值
	Filters               []map[string]interface{} `json:"filters"`
	OrderType             []OrderType              `json:"OrderType"`              // 订单类型
	TimeInForce           []TimeInForceType        `json:"timeInForce"`            // 有效方式
	LiquidationFee        float64                  `json:"liquidationFee,string"`  // 强平费率
	MarketTakeBound       float64                  `json:"marketTakeBound,string"` // 市价吃单(相对于标记价格)允许可造成的最大价格偏离比例
}

// LotSizeFilter define lot size filter of symbol | 数量限制
type LotSizeFilter struct {
	MaxQuantity string `json:"maxQty"`   // 数量上限, 最大数量
	MinQuantity string `json:"minQty"`   // 数量下限, 最小数量
	StepSize    string `json:"stepSize"` // 订单最小数量间隔
}

// PriceFilter define price filter of symbol | 价格限制
type PriceFilter struct {
	MaxPrice string `json:"maxPrice"` // 价格上限, 最大价格
	MinPrice string `json:"minPrice"` // 价格下限, 最小价格
	TickSize string `json:"tickSize"` // 订单最小价格间隔
}

// PercentPriceFilter define percent price filter of symbol | // 价格比限制
type PercentPriceFilter struct {
	MultiplierDecimal int    `json:"multiplierDecimal"`
	MultiplierUp      string `json:"multiplierUp"`   // 价格上限百分比
	MultiplierDown    string `json:"multiplierDown"` // 价格下限百分比
}

// MarketLotSizeFilter define market lot size filter of symbol | 市价订单数量限制
type MarketLotSizeFilter struct {
	MaxQuantity string `json:"maxQty"`   // 数量上限, 最大数量
	MinQuantity string `json:"minQty"`   // 数量下限, 最小数量
	StepSize    string `json:"stepSize"` // 允许的步进值
}

// MaxNumOrdersFilter define max num orders filter of symbol | 最多订单数限制
type MaxNumOrdersFilter struct {
	Limit int64 `json:"limit"`
}

// MaxNumAlgoOrdersFilter define max num algo orders filter of symbol | 最多条件订单数限制
type MaxNumAlgoOrdersFilter struct {
	Limit int64 `json:"limit"`
}

// MinNotionalFilter define min notional filter of symbol | 最小名义价值
type MinNotionalFilter struct {
	Notional string `json:"notional"`
}

// LotSizeFilter return lot size filter of symbol | 数量限制
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

// PriceFilter return price filter of symbol | 价格限制
func (s *Symbol) PriceFilter() *PriceFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypePrice) {
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

// PercentPriceFilter return percent price filter of symbol | 价格比限制
func (s *Symbol) PercentPriceFilter() *PercentPriceFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypePercentPrice) {
			f := &PercentPriceFilter{}
			if i, ok := filter["multiplierDecimal"]; ok {
				smd, is := i.(string)
				if is {
					md, _ := strconv.Atoi(smd)
					f.MultiplierDecimal = md
				} else {
					f.MultiplierDecimal = int(i.(float64))
				}
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

// MarketLotSizeFilter return market lot size filter of symbol | 市价订单数量限制
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

// MaxNumOrdersFilter return max num orders filter of symbol | 最多订单数限制
func (s *Symbol) MaxNumOrdersFilter() *MaxNumOrdersFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeMaxNumOrders) {
			f := &MaxNumOrdersFilter{}
			if i, ok := filter["limit"]; ok {
				f.Limit = int64(i.(float64))
			}
			return f
		}
	}
	return nil
}

// MaxNumAlgoOrdersFilter return max num orders filter of symbol | 最多条件订单数限制
func (s *Symbol) MaxNumAlgoOrdersFilter() *MaxNumAlgoOrdersFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeMaxNumAlgoOrders) {
			f := &MaxNumAlgoOrdersFilter{}
			if i, ok := filter["limit"]; ok {
				f.Limit = int64(i.(float64))
			}
			return f
		}
	}
	return nil
}

// MinNotionalFilter return min notional filter of symbol | 最小名义价值
func (s *Symbol) MinNotionalFilter() *MinNotionalFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string(SymbolFilterTypeMinNotional) {
			f := &MinNotionalFilter{}
			if i, ok := filter["notional"]; ok {
				f.Notional = i.(string)
			}
			return f
		}
	}
	return nil
}
