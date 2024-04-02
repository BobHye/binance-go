package futures

import (
	"context"
	"github.com/BobHye/binance-go/api"
	"net/http"
)

// SymbolLeverage define leverage info of symbol
type SymbolLeverage struct {
	Leverage         int    `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"`
	Symbol           string `json:"symbol"`
}

// ChangeLeverageService 调整指定交易对的开仓杠杆
type ChangeLeverageService struct {
	c        *Client
	symbol   string
	leverage int
}

// Symbol set symbol
func (s *ChangeLeverageService) Symbol(symbol string) *ChangeLeverageService {
	s.symbol = symbol
	return s
}

// Leverage set leverage
func (s *ChangeLeverageService) Leverage(leverage int) *ChangeLeverageService {
	s.leverage = leverage
	return s
}

// Do send request
func (s *ChangeLeverageService) Do(ctx context.Context, opts ...api.RequestOption) (res *SymbolLeverage, err error) {
	r := &api.Request{
		Method:   http.MethodPost,
		Endpoint: "/fapi/v1/leverage",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParams(api.Params{
		"symbol":   s.symbol,
		"leverage": s.leverage,
	})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(SymbolLeverage)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ChangeMarginTypeService 变换逐全仓模式
type ChangeMarginTypeService struct {
	c          *Client
	symbol     string
	marginType MarginType
}

// Symbol set symbol
func (s *ChangeMarginTypeService) Symbol(symbol string) *ChangeMarginTypeService {
	s.symbol = symbol
	return s
}

// MarginType set margin type
func (s *ChangeMarginTypeService) MarginType(marginType MarginType) *ChangeMarginTypeService {
	s.marginType = marginType
	return s
}

// Do send request
func (s *ChangeMarginTypeService) Do(ctx context.Context, opts ...api.RequestOption) (err error) {
	r := &api.Request{
		Method:   http.MethodPost,
		Endpoint: "/fapi/v1/marginType",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParams(api.Params{
		"symbol":     s.symbol,
		"marginType": s.marginType,
	})
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePositionMarginService 调整逐仓保证金
type UpdatePositionMarginService struct {
	c            *Client
	symbol       string
	positionSide *PositionSideType
	amount       string
	actionType   int
}

// Symbol set symbol
func (s *UpdatePositionMarginService) Symbol(symbol string) *UpdatePositionMarginService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *UpdatePositionMarginService) PositionSide(positionSide PositionSideType) *UpdatePositionMarginService {
	s.positionSide = &positionSide
	return s
}

// Amount set position margin amount
func (s *UpdatePositionMarginService) Amount(amount string) *UpdatePositionMarginService {
	s.amount = amount
	return s
}

// Type set action type: 1: Add postion margin，2: Reduce postion margin
func (s *UpdatePositionMarginService) Type(actionType int) *UpdatePositionMarginService {
	s.actionType = actionType
	return s
}

// Do send request
func (s *UpdatePositionMarginService) Do(ctx context.Context, opts ...api.RequestOption) (err error) {
	r := &api.Request{
		Method:   http.MethodPost,
		Endpoint: "/fapi/v1/positionMargin",
		SecType:  api.SecTypeSigned,
	}
	m := api.Params{
		"symbol": s.symbol,
		"amount": s.amount,
		"type":   s.actionType,
	}
	if s.positionSide != nil {
		m["positionSide"] = *s.positionSide
	}
	r.SetFormParams(m)

	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// ChangePositionModeService 更改持仓模式
type ChangePositionModeService struct {
	c        *Client
	dualSide string
}

// Change user's position mode: true - Hedge Mode, false - One-way Mode
func (s *ChangePositionModeService) DualSide(dualSide bool) *ChangePositionModeService {
	if dualSide {
		s.dualSide = "true"
	} else {
		s.dualSide = "false"
	}
	return s
}

// Do send request
func (s *ChangePositionModeService) Do(ctx context.Context, opts ...api.RequestOption) (err error) {
	r := &api.Request{
		Method:   http.MethodPost,
		Endpoint: "/fapi/v1/positionSide/dual",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParams(api.Params{
		"dualSidePosition": s.dualSide,
	})
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// GetPositionModeService 查询持仓模式
type GetPositionModeService struct {
	c *Client
}

// Response of user's position mode
type PositionMode struct {
	DualSidePosition bool `json:"dualSidePosition"`
}

// Do send request
func (s *GetPositionModeService) Do(ctx context.Context, opts ...api.RequestOption) (res *PositionMode, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/positionSide/dual",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParams(api.Params{})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = &PositionMode{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
