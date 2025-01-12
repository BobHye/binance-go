package futures

import (
	"context"
	"net/http"
)

// GetPositionMarginHistoryService get position margin history service
type GetPositionMarginHistoryService struct {
	c         *Client
	symbol    string
	_type     *int
	startTime *int64
	endTime   *int64
	limit     *int64
}

// SetSymbol set symbol
func (s *GetPositionMarginHistoryService) SetSymbol(symbol string) *GetPositionMarginHistoryService {
	s.symbol = symbol
	return s
}

// SetType set type
func (s *GetPositionMarginHistoryService) SetType(_type int) *GetPositionMarginHistoryService {
	s._type = &_type
	return s
}

// SetStartTime set startTime
func (s *GetPositionMarginHistoryService) SetStartTime(startTime int64) *GetPositionMarginHistoryService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *GetPositionMarginHistoryService) SetEndTime(endTime int64) *GetPositionMarginHistoryService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *GetPositionMarginHistoryService) SetLimit(limit int64) *GetPositionMarginHistoryService {
	s.limit = &limit
	return s
}

// Do send request
func (s *GetPositionMarginHistoryService) Do(ctx context.Context, opts ...RequestOption) (res []*PositionMarginHistory, err error) {
	// /fapi/v1/positionMargin/history | 查询逐仓保证金变动历史
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/positionMargin/history",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s._type != nil {
		r.setParam("type", *s._type)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*PositionMarginHistory, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PositionMarginHistory define position margin history info
type PositionMarginHistory struct {
	Amount       string           `json:"amount"`       // 数量
	Asset        string           `json:"asset"`        // 资产
	DeltaType    string           `json:"deltaType"`    // 划转类型
	Symbol       string           `json:"symbol"`       // 交易对
	Time         int64            `json:"time"`         // 时间
	Type         int              `json:"type"`         // 调整方向
	PositionSide PositionSideType `json:"positionSide"` // 持仓方向
}
