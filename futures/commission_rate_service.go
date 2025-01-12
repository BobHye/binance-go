package futures

import (
	"context"
	"net/http"
)

// 查询用户交易手续费率

type CommissionRateService struct {
	c      *Client
	symbol string
}

// SetSymbol 设置交易对
func (s *CommissionRateService) SetSymbol(symbol string) *CommissionRateService {
	s.symbol = symbol
	return s
}

func (s *CommissionRateService) Do(ctx context.Context, opts ...RequestOption) (res *CommissionRate, err error) {
	// GET /fapi/v1/commissionRate | 查询用户手续费率
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/commissionRate",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
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

type CommissionRate struct {
	Symbol              string  `json:"symbol"`
	MakerCommissionRate float64 `json:"makerCommissionRate,string"`
	TakerCommissionRate float64 `json:"takerCommissionRate,string"`
}
