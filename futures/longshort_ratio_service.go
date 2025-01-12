package futures

import (
	"context"
	"net/http"
)

// LongShortRatioService list open history data of a symbol.
type LongShortRatioService struct {
	c         *Client
	symbol    string
	period    string // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	limit     *int   // default 30, max 500
	startTime *int64
	endTime   *int64
}

// LongShortRatio 多空比
type LongShortRatio struct {
	Symbol         string  `json:"symbol"`
	LongShortRatio float64 `json:"longShortRatio,string"` // 多空人数比值
	LongAccount    float64 `json:"longAccount,string"`    // 多仓人数比例
	ShortAccount   float64 `json:"shortAccount,string"`   // 空仓人数比例
	Timestamp      int64   `json:"timestamp"`
}

// SetSymbol set symbol
func (s *LongShortRatioService) SetSymbol(symbol string) *LongShortRatioService {
	s.symbol = symbol
	return s
}

// SetPeriod set period interval
func (s *LongShortRatioService) SetPeriod(period string) *LongShortRatioService {
	s.period = period
	return s
}

// SetLimit set limit
func (s *LongShortRatioService) SetLimit(limit int) *LongShortRatioService {
	s.limit = &limit
	return s
}

// SetStartTime set startTime
func (s *LongShortRatioService) SetStartTime(startTime int64) *LongShortRatioService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *LongShortRatioService) SetEndTime(endTime int64) *LongShortRatioService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *LongShortRatioService) Do(ctx context.Context, opts ...RequestOption) (res []*LongShortRatio, err error) {
	// GET /futures/data/globalLongShortAccountRatio | 多空持仓人数比
	// 若无 startime 和 endtime 限制， 则默认返回当前时间往前的limit值
	// 仅支持最近30天的数据
	// IP限频为1000次/5min
	r := &request{
		method:   http.MethodGet,
		endpoint: "/futures/data/globalLongShortAccountRatio",
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
		return []*LongShortRatio{}, err
	}

	res = make([]*LongShortRatio, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*LongShortRatio{}, err
	}

	return res, nil
}

// TopLongShortAccountRatioService list open history data of a symbol
type TopLongShortAccountRatioService struct {
	c         *Client
	symbol    string
	period    string // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	limit     *int   // default 30, max 500
	startTime *int64
	endTime   *int64
}

// SetSymbol set symbol
func (s *TopLongShortAccountRatioService) SetSymbol(symbol string) *TopLongShortAccountRatioService {
	s.symbol = symbol
	return s
}

// SetPeriod set period interval
func (s *TopLongShortAccountRatioService) SetPeriod(period string) *TopLongShortAccountRatioService {
	s.period = period
	return s
}

// SetLimit set limit
func (s *TopLongShortAccountRatioService) SetLimit(limit int) *TopLongShortAccountRatioService {
	s.limit = &limit
	return s
}

// SetStartTime set startTime
func (s *TopLongShortAccountRatioService) SetStartTime(startTime int64) *TopLongShortAccountRatioService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *TopLongShortAccountRatioService) SetEndTime(endTime int64) *TopLongShortAccountRatioService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *TopLongShortAccountRatioService) Do(ctx context.Context, opts ...RequestOption) (res []*LongShortRatio, err error) {
	// GET /futures/data/topLongShortAccountRatio | 大户账户数多空比
	// 持仓大户的净持仓多头和空头账户数占比，大户指保证金余额排名前20%的用户。一个账户记一次。 多仓账户数比例 = 持多仓大户数 / 总持仓大户数 空仓账户数比例 = 持空仓大户数 / 总持仓大户数 多空账户数比值 = 多仓账户数比例 / 空仓账户数比例
	// 若无 startime 和 endtime 限制， 则默认返回当前时间往前的limit值
	// 仅支持最近30天的数据
	// IP限频为1000次/5min
	r := &request{
		method:   http.MethodGet,
		endpoint: "/futures/data/topLongShortAccountRatio",
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
		return []*LongShortRatio{}, err
	}

	res = make([]*LongShortRatio, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*LongShortRatio{}, err
	}

	return res, nil
}

// TopLongShortPositionRatioService list open history data of a symbol
type TopLongShortPositionRatioService struct {
	c         *Client
	symbol    string
	period    string // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	limit     *int   // default 30, max 500
	startTime *int64
	endTime   *int64
}

// SetSymbol set symbol
func (s *TopLongShortPositionRatioService) SetSymbol(symbol string) *TopLongShortPositionRatioService {
	s.symbol = symbol
	return s
}

// SetPeriod set period interval
func (s *TopLongShortPositionRatioService) SetPeriod(period string) *TopLongShortPositionRatioService {
	s.period = period
	return s
}

// SetLimit set limit
func (s *TopLongShortPositionRatioService) SetLimit(limit int) *TopLongShortPositionRatioService {
	s.limit = &limit
	return s
}

// SetStartTime set startTime
func (s *TopLongShortPositionRatioService) SetStartTime(startTime int64) *TopLongShortPositionRatioService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *TopLongShortPositionRatioService) SetEndTime(endTime int64) *TopLongShortPositionRatioService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *TopLongShortPositionRatioService) Do(ctx context.Context, opts ...RequestOption) (res []*LongShortRatio, err error) {
	// GET /futures/data/topLongShortPositionRatio | 大户持仓量多空比
	// 大户的多头和空头总持仓量占比，大户指保证金余额排名前20%的用户。 多仓持仓量比例 = 大户多仓持仓量 / 大户总持仓量 空仓持仓量比例 = 大户空仓持仓量 / 大户总持仓量 多空持仓量比值 = 多仓持仓量比例 / 空仓持仓量比例
	r := &request{
		method:   http.MethodGet,
		endpoint: "/futures/data/topLongShortPositionRatio",
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
		return []*LongShortRatio{}, err
	}

	res = make([]*LongShortRatio, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*LongShortRatio{}, err
	}

	return res, nil
}
