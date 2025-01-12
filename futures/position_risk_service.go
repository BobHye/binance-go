package futures

import (
	"context"
	"net/http"
)

// 用户持仓风险V2

// GetPositionRiskService get account balance
type GetPositionRiskService struct {
	c      *Client
	symbol string
}

// SetSymbol set symbol
func (s *GetPositionRiskService) SetSymbol(symbol string) *GetPositionRiskService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *GetPositionRiskService) Do(ctx context.Context, opts ...RequestOption) (res []*PositionRisk, err error) {
	// GET /fapi/v2/positionRisk | 查询持仓风险
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/positionRisk",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// PositionRisk define position risk info
type PositionRisk struct {
	Symbol           string  `json:"symbol"`                  // 交易对
	PositionAmt      float64 `json:"positionAmt,string"`      // 头寸数量，符号代表多空方向, 正数为多，负数为空
	EntryPrice       float64 `json:"entryPrice,string"`       // 开仓均价
	BreakEvenPrice   float64 `json:"breakEvenPrice,string"`   // 盈亏平衡价
	MarkPrice        float64 `json:"markPrice,string"`        // 当前标记价格
	UnRealizedProfit float64 `json:"unRealizedProfit,string"` // 持仓未实现盈亏
	LiquidationPrice float64 `json:"liquidationPrice,string"` // 参考强平价格
	Leverage         int     `json:"leverage,string"`         // 当前杠杆倍数
	MaxNotionalValue float64 `json:"maxNotionalValue,string"` // 当前杠杆倍数允许的名义价值上限
	MarginType       string  `json:"marginType"`              // 逐仓模式或全仓模式
	IsolatedMargin   float64 `json:"isolatedMargin,string"`   // 逐仓保证金
	IsAutoAddMargin  bool    `json:"isAutoAddMargin,string"`
	PositionSide     string  `json:"positionSide"` // 持仓方向
	Notional         float64 `json:"notional,string"`
	IsolatedWallet   float64 `json:"isolatedWallet,string"`
	UpdateTime       int64   `json:"updateTime"` // 更新时间
}
