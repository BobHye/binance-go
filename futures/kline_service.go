package futures

// 获取K线数据

import (
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

// KlinesService list klines
type KlinesService struct {
	c         *Client
	symbol    string // 交易对
	interval  string // 时间间隔
	startTime *int64 // 起始时间
	endTime   *int64 // 结束时间
	limit     *int   // 默认值:500 最大值:1500
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
	// GET /fapi/v1/klines | K线数据 (每根K线的开盘时间可视为唯一ID)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/klines",
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
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// Kline define kline info
type Kline struct {
	OpenTime                 int64   `json:"openTime"`                        // 开盘时间
	Open                     float64 `json:"open,string"`                     // 开盘价
	High                     float64 `json:"high,string"`                     // 最高价
	Low                      float64 `json:"low,string"`                      // 最低价
	Close                    float64 `json:"close,string"`                    // 收盘价(当前K线未结束的即为最新价)
	Volume                   float64 `json:"volume,string"`                   // 成交量
	CloseTime                int64   `json:"closeTime"`                       // 收盘时间
	QuoteAssetVolume         float64 `json:"quoteAssetVolume,string"`         // 成交额
	TradeNum                 int64   `json:"tradeNum"`                        // 成交笔数
	TakerBuyBaseAssetVolume  float64 `json:"takerBuyBaseAssetVolume,string"`  // 主动买入成交量
	TakerBuyQuoteAssetVolume float64 `json:"takerBuyQuoteAssetVolume,string"` // 主动买入成交额
}

func (kline *Kline) UnmarshalJSON(data []byte) error {
	iter := jsoniter.Get(data)
	if iter.Size() < 11 {
		return errors.New("invalid kline response")
	}

	openTime := iter.Get(0).ToInt64()
	open := iter.Get(1).ToFloat64()
	high := iter.Get(2).ToFloat64()
	low := iter.Get(3).ToFloat64()
	_close := iter.Get(4).ToFloat64()
	volume := iter.Get(5).ToFloat64()
	closeTime := iter.Get(6).ToInt64()
	quoteAssetVolume := iter.Get(7).ToFloat64()
	tradeNum := iter.Get(8).ToInt64()
	takerBuyBaseAssetVolume := iter.Get(9).ToFloat64()
	takerBuyQuoteAssetVolume := iter.Get(10).ToFloat64()

	kline.OpenTime = openTime
	kline.Open = open
	kline.High = high
	kline.Low = low
	kline.Close = _close
	kline.Volume = volume
	kline.CloseTime = closeTime
	kline.QuoteAssetVolume = quoteAssetVolume
	kline.TradeNum = tradeNum
	kline.TakerBuyBaseAssetVolume = takerBuyBaseAssetVolume
	kline.TakerBuyQuoteAssetVolume = takerBuyQuoteAssetVolume

	return nil
}
