package futures

import (
	"context"
	"net/http"
)

// MarkPriceKlinesService list mark price klines
type MarkPriceKlinesService struct {
	c         *Client
	symbol    string
	interval  string
	limit     *int
	startTime *int64
	endTime   *int64
}

// SetInterval set interval
func (mpks *MarkPriceKlinesService) SetInterval(interval string) *MarkPriceKlinesService {
	mpks.interval = interval
	return mpks
}

// SetLimit set limit
func (mpks *MarkPriceKlinesService) SetLimit(limit int) *MarkPriceKlinesService {
	mpks.limit = &limit
	return mpks
}

// SetStartTime set startTime
func (mpks *MarkPriceKlinesService) SetStartTime(startTime int64) *MarkPriceKlinesService {
	mpks.startTime = &startTime
	return mpks
}

// SetEndTime set endTime
func (mpks *MarkPriceKlinesService) SetEndTime(endTime int64) *MarkPriceKlinesService {
	mpks.endTime = &endTime
	return mpks
}

// Do send request
func (mpks *MarkPriceKlinesService) Do(ctx context.Context, opts ...RequestOption) (res []*Kline, err error) {
	// GET /fapi/v1/markPriceKlines | 标记价格K线数据(每根K线的开盘时间可视为唯一ID)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/markPriceKlines",
	}
	r.setParam("symbol", mpks.symbol)
	r.setParam("interval", mpks.interval)
	if mpks.limit != nil {
		r.setParam("limit", *mpks.limit)
	}
	if mpks.startTime != nil {
		r.setParam("startTime", *mpks.startTime)
	}
	if mpks.endTime != nil {
		r.setParam("endTime", *mpks.endTime)
	}
	data, _, err := mpks.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}
