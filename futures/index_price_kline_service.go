package futures

import (
	"context"
	"net/http"
)

// IndexPriceKlinesService list klines
type IndexPriceKlinesService struct {
	c         *Client
	pair      string
	interval  string
	limit     *int
	startTime *int64
	endTime   *int64
}

// SetPair sets pair
func (ipks *IndexPriceKlinesService) SetPair(pair string) *IndexPriceKlinesService {
	ipks.pair = pair
	return ipks
}

// SetInterval set interval
func (ipks *IndexPriceKlinesService) SetInterval(interval string) *IndexPriceKlinesService {
	ipks.interval = interval
	return ipks
}

// SetLimit set limit
func (ipks *IndexPriceKlinesService) SetLimit(limit int) *IndexPriceKlinesService {
	ipks.limit = &limit
	return ipks
}

// SetStartTime set startTime
func (ipks *IndexPriceKlinesService) SetStartTime(startTime int64) *IndexPriceKlinesService {
	ipks.startTime = &startTime
	return ipks
}

// SetEndTime set endTime
func (ipks *IndexPriceKlinesService) SetEndTime(endTime int64) *IndexPriceKlinesService {
	ipks.endTime = &endTime
	return ipks
}

// Do send request
func (ipks *IndexPriceKlinesService) Do(ctx context.Context, opts ...RequestOption) (res []*Kline, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapiv1/indexPriceKlines",
	}
	r.setParam("pair", ipks.pair)
	r.setParam("interval", ipks.interval)
	if ipks.limit != nil {
		r.setParam("limit", *ipks.limit)
	}
	if ipks.startTime != nil {
		r.setParam("startTime", *ipks.startTime)
	}
	if ipks.endTime != nil {
		r.setParam("endTime", *ipks.endTime)
	}

	data, _, err := ipks.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}
