package spot

import (
	"context"
	"fmt"
	"net/http"
)

// KlinesService list klines
type KlinesService struct {
	c         *Client
	symbol    string
	interval  string
	limit     *int // Default 500; max 1000
	startTime *int64
	endTime   *int64
}

// SetSymbol set symbol
func (s *KlinesService) SetSymbol(symbol string) *KlinesService {
	s.symbol = symbol
	return s
}

// SetInterval set interval
func (s *KlinesService) SetInterval(interval string) *KlinesService {
	s.interval = interval
	return s
}

// SetLimit set limit
func (s *KlinesService) SetLimit(limit int) *KlinesService {
	s.limit = &limit
	return s
}

// SetStartTime set startTime
func (s *KlinesService) SetStartTime(startTime int64) *KlinesService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *KlinesService) SetEndTime(endTime int64) *KlinesService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *KlinesService) Do(ctx context.Context, opts ...RequestOption) (res []*Kline, err error) {
	// GET /api/v3/klines | K线数据
	// 如果未发送startTime和endTime，将返回最近的K线数据。
	// timeZone支持的值包括：
	// 小时和分钟（例如 -1:00，05:45）
	// 仅小时（例如 0，8，4）
	// 接受的值范围严格为 [-12:00 到 +14:00]（包括边界）
	// 如果提供了timeZone，K线间隔将在该时区中解释，而不是在UTC中。
	// 请注意，无论timeZone如何，startTime和endTime始终以UTC时区解释
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/klines",
	}
	r.setParam("symbol", s.symbol)
	r.setParam("interval", s.interval)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*Kline{}, err
	}
	j, err := newJSON(data)
	if err != nil {
		return []*Kline{}, err
	}
	num := len(j.MustArray())
	res = make([]*Kline, num)
	for i := 0; i < num; i++ {
		item := j.GetIndex(i)
		if len(item.MustArray()) < 11 {
			err = fmt.Errorf("invalid kline response")
			return []*Kline{}, err
		}
		res[i] = &Kline{
			OpenTime:                 item.GetIndex(0).MustInt64(),   // 开盘时间
			Open:                     item.GetIndex(1).MustString(),  // 开盘价
			High:                     item.GetIndex(2).MustString(),  // 最高价
			Low:                      item.GetIndex(3).MustString(),  // 最低价
			Close:                    item.GetIndex(4).MustString(),  // 收盘价(当前K线未结束的即为最新价)
			Volume:                   item.GetIndex(5).MustString(),  // 成交量
			CloseTime:                item.GetIndex(6).MustInt64(),   // 收盘时间
			QuoteAssetVolume:         item.GetIndex(7).MustString(),  // 成交额
			TradeNum:                 item.GetIndex(8).MustInt64(),   // 成交笔数
			TakerBuyBaseAssetVolume:  item.GetIndex(9).MustString(),  // 主动买入成交量
			TakerBuyQuoteAssetVolume: item.GetIndex(10).MustString(), // 主动买入成交额
		}
	}
	return res, nil
}

// Kline define kline info
type Kline struct {
	OpenTime                 int64  `json:"openTime"`                 // 开盘时间
	Open                     string `json:"open"`                     // 开盘价
	High                     string `json:"high"`                     // 最高价
	Low                      string `json:"low"`                      // 最低价
	Close                    string `json:"close"`                    // 收盘价(当前K线未结束的即为最新价)
	Volume                   string `json:"volume"`                   // 成交量
	CloseTime                int64  `json:"closeTime"`                // 收盘时间
	QuoteAssetVolume         string `json:"quoteAssetVolume"`         // 成交额
	TradeNum                 int64  `json:"tradeNum"`                 // 成交笔数
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`  // 主动买入成交量
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"` // 主动买入成交额
}
