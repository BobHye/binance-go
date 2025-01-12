package futures

import (
	"context"
	"net/http"
)

// GetOpenInterestService get present open interest of a specific symbol.
type GetOpenInterestService struct {
	c      *Client
	symbol string
}

// SetSymbol set symbol
func (s *GetOpenInterestService) SetSymbol(symbol string) *GetOpenInterestService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *GetOpenInterestService) Do(ctx context.Context, opts ...RequestOption) (res *OpenInterest, err error) {
	// GET /fapi/v1/openInterest | 获取未平仓合约数
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/openInterest",
	}
	r.setParam("symbol", s.symbol)
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(OpenInterest)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type OpenInterest struct {
	OpenInterest string `json:"openInterest"` // 未平仓合约数量
	Symbol       string `json:"symbol"`       // 交易对
	Time         int64  `json:"time"`         // 撮合引擎时间
}

type OpenInterestStatisticsService struct {
	c         *Client
	symbol    string
	period    string
	limit     *int
	startTime *int64
	endTime   *int64
}

// SetSymbol set symbol
func (s *OpenInterestStatisticsService) SetSymbol(symbol string) *OpenInterestStatisticsService {
	s.symbol = symbol
	return s
}

// SetPeriod set period interval
func (s *OpenInterestStatisticsService) SetPeriod(period string) *OpenInterestStatisticsService {
	s.period = period
	return s
}

// SetLimit set limit
func (s *OpenInterestStatisticsService) SetLimit(limit int) *OpenInterestStatisticsService {
	s.limit = &limit
	return s
}

// SetStartTime set startTime
func (s *OpenInterestStatisticsService) SetStartTime(startTime int64) *OpenInterestStatisticsService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *OpenInterestStatisticsService) SetEndTime(endTime int64) *OpenInterestStatisticsService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *OpenInterestStatisticsService) Do(ctx context.Context, opts ...RequestOption) (res []*OpenInterestStatistic, err error) {
	// GET /futures/data/openInterestHist | 查询合约持仓量历史
	// 若无 startime 和 endtime 限制， 则默认返回当前时间往前的limit值
	// 仅支持最近30天的数据
	// IP限频为1000次/5min
	r := &request{
		method:   http.MethodGet,
		endpoint: "/futures/data/openInterestHist",
	}
	r.setParam("symbol", s.symbol)
	r.setParam("period", s.period)

	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*OpenInterestStatistic{}, err
	}
	res = make([]*OpenInterestStatistic, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*OpenInterestStatistic{}, err
	}

	return res, nil
}

type OpenInterestStatistic struct {
	Symbol               string `json:"symbol"`
	SumOpenInterest      string `json:"sumOpenInterest"`      // 持仓总数量
	SumOpenInterestValue string `json:"sumOpenInterestValue"` // 持仓总价值
	Timestamp            int64  `json:"timestamp"`
}
