package futures

import (
	"context"
	"github.com/BobHye/binance-go/api"
	"net/http"
)

// 查询用户交易手续费率

type CommissionRate struct {
	Symbol              string  `json:"symbol"`
	MakerCommissionRate float64 `json:"makerCommissionRate,string"`
	TakerCommissionRate float64 `json:"takerCommissionRate,string"`
}

type CommissionRateService struct {
	c      *Client
	symbol string
}

// Symbol 设置交易对
func (s *CommissionRateService) Symbol(symbol string) *CommissionRateService {
	s.symbol = symbol
	return s
}

func (s *CommissionRateService) Do(ctx context.Context, opts ...api.RequestOption) (res *CommissionRate, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/commissionRate",
		SecType:  api.SecTypeSigned,
	}
	if s.symbol != "" {
		r.SetParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CommissionRate)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
