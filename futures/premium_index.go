package futures

import (
	"context"
	"github.com/BobHye/binance-go/common"
	"net/http"
)

// PremiumIndexService get premium index
type PremiumIndexService struct {
	c      *Client
	symbol *string
}

// SetSymbol set symbol
func (s *PremiumIndexService) SetSymbol(symbol string) *PremiumIndexService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *PremiumIndexService) Do(ctx context.Context, opts ...RequestOption) (res []*PremiumIndex, err error) {
	// GET /fapi/v1/premiumIndex | 采集各大交易所数据加权平均
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/premiumIndex",
		secType:  secTypeSigned,
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	data = common.ToJSONList(data)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// PremiumIndex define premium index of mark price
type PremiumIndex struct {
	Symbol               string  `json:"symbol"`                      // 交易对
	MarkPrice            float64 `json:"markPrice,string"`            // 标记价格
	IndexPrice           float64 `json:"indexPrice,string"`           // 指数价格
	EstimatedSettlePrice float64 `json:"estimatedSettlePrice,string"` // 预估结算价,仅在交割开始前最后一小时有意义
	LastFundingRate      float64 `json:"lastFundingRate,string"`      // 最近更新的资金费率
	NextFundingTime      int64   `json:"nextFundingTime"`             // 下次资金费时间
	InterestRate         float64 `json:"interestRate,string"`         // 标的资产基础利率
	Time                 int64   `json:"time"`                        // 更新时间
}

// FundingRateService get funding rate
type FundingRateService struct {
	c         *Client
	symbol    string
	startTime *int64
	endTime   *int64
	limit     *int // 默认值:100 最大值:1000
}

// SetSymbol set symbol
func (s *FundingRateService) SetSymbol(symbol string) *FundingRateService {
	s.symbol = symbol
	return s
}

// SetStartTime set startTime
func (s *FundingRateService) SetStartTime(startTime int64) *FundingRateService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *FundingRateService) SetEndTime(endTime int64) *FundingRateService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *FundingRateService) SetLimit(limit int) *FundingRateService {
	s.limit = &limit
	return s
}

// Do send request
func (s *FundingRateService) Do(ctx context.Context, opts ...RequestOption) (res []*FundingRate, err error) {
	// GET /fapi/v1/fundingRate | 查询资金费率历史
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/fundingRate",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
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
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// FundingRate define funding rate of mark price
type FundingRate struct {
	Symbol      string  `json:"symbol"`             // 交易对
	FundingRate float64 `json:"fundingRate,string"` // 资金费率
	FundingTime int64   `json:"fundingTime"`        // 资金费时间
	Time        int64   `json:"time"`               // 资金费对应标记价格
}

// GetLeverageBracketService get funding rate
type GetLeverageBracketService struct {
	c      *Client
	symbol string
}

// SetSymbol set symbol
func (s *GetLeverageBracketService) SetSymbol(symbol string) *GetLeverageBracketService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *GetLeverageBracketService) Do(ctx context.Context, opts ...RequestOption) (res []*LeverageBracket, err error) {
	// GET /fapi/v1/leverageBracket | 查询账户特定交易对的杠杆分层标准
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/leverageBracket",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	if s.symbol != "" {
		data = common.ToJSONList(data)
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// LeverageBracket define the leverage bracket
type LeverageBracket struct {
	Symbol       string    `json:"symbol"`
	NotionalCoef float64   `json:"notionalCoef"` // 用户bracket相对默认bracket的倍数，仅在和交易对默认不一样时显示
	Brackets     []Bracket `json:"brackets"`
}

// Bracket define the bracket
type Bracket struct {
	Bracket          int     `json:"bracket"`          // 层级
	InitialLeverage  int     `json:"initialLeverage"`  // 该层允许的最高初始杠杆倍数
	NotionalCap      float64 `json:"notionalCap"`      // 该层对应的名义价值上限
	NotionalFloor    float64 `json:"notionalFloor"`    // 该层对应的名义价值下限
	MaintMarginRatio float64 `json:"mainMarginRation"` // 该层对应的维持保证金率
	Cum              float64 `json:"cum"`              // 速算数
}
