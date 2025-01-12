package futures

import (
	"context"
	"net/http"
)

// ChangeLeverageService change user's initial leverage of specific symbol market | 调整指定交易对的开仓杠杆
type ChangeLeverageService struct {
	c        *Client
	symbol   string
	leverage int
}

// SetSymbol set symbol
func (s *ChangeLeverageService) SetSymbol(symbol string) *ChangeLeverageService {
	s.symbol = symbol
	return s
}

// SetLeverage set leverage
func (s *ChangeLeverageService) SetLeverage(leverage int) *ChangeLeverageService {
	s.leverage = leverage
	return s
}

// Do send request
func (s *ChangeLeverageService) Do(ctx context.Context, opts ...RequestOption) (res *SymbolLeverage, err error) {
	// POST /fapi/v1/leverage | 调整用户在指定symbol合约的开仓杠杆
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/leverage",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
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

// SymbolLeverage define leverage info of symbol
type SymbolLeverage struct {
	Leverage         int    `json:"leverage"`         // 杠杆倍数
	MaxNotionalValue string `json:"maxNotionalValue"` // 当前杠杆倍数下允许的最大名义价值
	Symbol           string `json:"symbol"`           // 交易对
}

// ChangeMarginTypeService change user's margin type of specific symbol market | 变换逐全仓模式
type ChangeMarginTypeService struct {
	c          *Client
	symbol     string
	marginType MarginType // 保证金模式 ISOLATED(逐仓), CROSSED(全仓)
}

// SetSymbol set symbol
func (s *ChangeMarginTypeService) SetSymbol(symbol string) *ChangeMarginTypeService {
	s.symbol = symbol
	return s
}

// SetMarginType set margin type
func (s *ChangeMarginTypeService) SetMarginType(marginType MarginType) *ChangeMarginTypeService {
	s.marginType = marginType
	return s
}

// Do send request
func (s *ChangeMarginTypeService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// POST /fapi/v1/marginType | 变换用户在指定symbol合约上的保证金模式：逐仓或全仓。
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/marginType",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"symbol":     s.symbol,
		"marginType": s.marginType,
	})
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePositionMarginService update isolated position margin | 调整逐仓保证金
type UpdatePositionMarginService struct {
	c            *Client
	symbol       string            // 交易对
	positionSide *PositionSideType // 持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
	amount       string            // 保证金资金
	actionType   int               // 调整方向 1: 增加逐仓保证金，2: 减少逐仓保证金
}

// SetSymbol set symbol
func (s *UpdatePositionMarginService) SetSymbol(symbol string) *UpdatePositionMarginService {
	s.symbol = symbol
	return s
}

// SetPositionSide set positionSide
func (s *UpdatePositionMarginService) SetPositionSide(positionSide PositionSideType) *UpdatePositionMarginService {
	s.positionSide = &positionSide
	return s
}

// SetAmount set position margin amount
func (s *UpdatePositionMarginService) SetAmount(amount string) *UpdatePositionMarginService {
	s.amount = amount
	return s
}

// SetType set action type: 1: Add postion margin，2: Reduce postion margin
func (s *UpdatePositionMarginService) SetType(actionType int) *UpdatePositionMarginService {
	s.actionType = actionType
	return s
}

// Do send request
func (s *UpdatePositionMarginService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// POST /fapi/v1/positionMargin | 针对逐仓模式下的仓位，调整其逐仓保证金资金。
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/positionMargin",
		secType:  secTypeSigned,
	}
	m := params{
		"symbol": s.symbol,
		"amount": s.amount,
		"type":   s.actionType,
	}
	if s.positionSide != nil {
		m["positionSide"] = *s.positionSide
	}
	r.setFormParams(m)

	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// ChangePositionModeService change user's position mode | 更改持仓模式
type ChangePositionModeService struct {
	c                *Client
	dualSidePosition string // "true": 双向持仓模式；"false": 单向持仓模式
}

// SetDualSide Change user's position mode: true - Hedge Mode, false - One-way Mode
func (s *ChangePositionModeService) SetDualSide(dualSide bool) *ChangePositionModeService {
	if dualSide {
		s.dualSidePosition = "true"
	} else {
		s.dualSidePosition = "false"
	}
	return s
}

// Do send request
func (s *ChangePositionModeService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// POST /fapi/v1/positionSide/dual | 变换用户在 所有symbol 合约上的持仓模式：双向持仓或单向持仓。
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/positionSide/dual",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"dualSidePosition": s.dualSidePosition,
	})
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// GetPositionModeService get user's position mode | 查询持仓模式
type GetPositionModeService struct {
	c *Client
}

// Do send request
func (s *GetPositionModeService) Do(ctx context.Context, opts ...RequestOption) (res *PositionMode, err error) {
	// GET /fapi/v1/positionSide/dual | 查询用户目前在 所有symbol 合约上的持仓模式：双向持仓或单向持仓。
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/positionSide/dual",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{})
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

// PositionMode Response of user's position mode
type PositionMode struct {
	DualSidePosition bool `json:"dualSidePosition"` // "true": 双向持仓模式；"false": 单向持仓模式
}
