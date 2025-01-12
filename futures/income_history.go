package futures

import (
	"context"
	"net/http"
)

// GetIncomeHistoryService get position margin history service
type GetIncomeHistoryService struct {
	c      *Client
	symbol string
	// 收益类型： TRANSFER 转账, WELCOME_BONUS 欢迎奖金, REALIZED_PNL 已实现盈亏, FUNDING_FEE 资金费用, COMMISSION 佣金, INSURANCE_CLEAR 强平, REFERRAL_KICKBACK 推荐人返佣, COMMISSION_REBATE 被推荐人返佣, API_REBATE API佣金回扣, CONTEST_REWARD 交易大赛奖金, CROSS_COLLATERAL_TRANSFER cc转账, OPTIONS_PREMIUM_FEE 期权购置手续费, OPTIONS_SETTLE_PROFIT 期权行权收益, INTERNAL_TRANSFER 内部账户，给普通用户划转, AUTO_EXCHANGE 自动兑换, DELIVERED_SETTELMENT 下架结算, COIN_SWAP_DEPOSIT 闪兑转入, COIN_SWAP_WITHDRAW 闪兑转出, POSITION_LIMIT_INCREASE_FEE 仓位限制上调费用，STRATEGY_UMFUTURES_TRANSFER UM策略子账户划转，FEE_RETURN 策略交易手续费退还，BFUSD_REWARD BFUSD每日奖励
	incomeType string
	startTime  *int64 // 起始时间
	endTime    *int64 // 结束时间
	limit      *int64 // 返回的结果集数量 默认值:100 最大值:1000
}

// SetSymbol set symbol
func (s *GetIncomeHistoryService) SetSymbol(symbol string) *GetIncomeHistoryService {
	s.symbol = symbol
	return s
}

// SetIncomeType set income type
func (s *GetIncomeHistoryService) SetIncomeType(incomeType string) *GetIncomeHistoryService {
	s.incomeType = incomeType
	return s
}

// SetStartTime set startTime
func (s *GetIncomeHistoryService) SetStartTime(startTime int64) *GetIncomeHistoryService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *GetIncomeHistoryService) SetEndTime(endTime int64) *GetIncomeHistoryService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *GetIncomeHistoryService) SetLimit(limit int64) *GetIncomeHistoryService {
	s.limit = &limit
	return s
}

// Do send request
func (s *GetIncomeHistoryService) Do(ctx context.Context, opts ...RequestOption) (res []*IncomeHistory, err error) {
	// GET /fapi/v1/income | 获取账户损益资金流水
	// 如果startTime 和 endTime 均未发送, 只会返回最近7天的数据。
	// 如果incomeType没有发送，返回所有类型账户损益资金流水。
	// "trandId" 在相同用户的同一种收益流水类型中是唯一的。
	// 仅保留最近3个月的数据。
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/income",
	}
	r.setParam("symbol", s.symbol)
	if s.incomeType != "" {
		r.setParam("incomeType", s.incomeType)
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
	err = json.Unmarshal(data, &res)
	return res, err
}

// IncomeHistory define position margin history info
type IncomeHistory struct {
	Asset      string `json:"asset"`      // 资产内容
	Income     string `json:"income"`     // 资金流数量，正数代表流入，负数代表流出
	IncomeType string `json:"incomeType"` // 资金流类型
	Info       string `json:"info"`       // 备注信息，取决于流水类型
	Symbol     string `json:"symbol"`     // 交易对，仅针对涉及交易对的资金流
	Time       int64  `json:"time"`       // 时间
	TranID     int64  `json:"tranId"`     // 划转ID
	TradeID    string `json:"tradeId"`    // 引起流水产生的原始交易ID
}
