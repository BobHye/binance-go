package futures

import (
	"context"
	"github.com/BobHye/binance-go/api"
	"net/http"
)

// 用户持仓风险V2

// PositionRisk define position risk info
type PositionRisk struct {
	EntryPrice       float64 `json:"entryPrice,string"`
	MarginType       string  `json:"marginType"`
	IsAutoAddMargin  bool    `json:"isAutoAddMargin,string"`
	IsolatedMargin   float64 `json:"isolatedMargin,string"`
	Leverage         int     `json:"leverage,string"`
	LiquidationPrice float64 `json:"liquidationPrice,string"`
	MarkPrice        float64 `json:"markPrice,string"`
	MaxNotionalValue float64 `json:"maxNotionalValue,string"`
	PositionAmt      float64 `json:"positionAmt,string"`
	Symbol           string  `json:"symbol"`
	UnRealizedProfit float64 `json:"unRealizedProfit,string"`
	PositionSide     string  `json:"positionSide"`
	Notional         float64 `json:"notional,string"`
	IsolatedWallet   float64 `json:"isolatedWallet,string"`
}

type GetPositionRiskService struct {
	c      *Client
	symbol string
}

func (s *GetPositionRiskService) Symbol(symbol string) *GetPositionRiskService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *GetPositionRiskService) Do(ctx context.Context, opts ...api.RequestOption) (res []*PositionRisk, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v2/positionRisk",
		SecType:  api.SecTypeSigned,
	}
	if s.symbol != "" {
		r.SetParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}
