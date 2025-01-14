package futures

import (
	"context"
	"net/http"
)

// GetBalanceService get account balance | 获取账户余额
type GetBalanceService struct {
	c *Client
}

func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (res []*Balance, err error) {
	// GET /fapi/v2/balance | 账户余额V2
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/balance",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// Balance define user balance of your account | 账户余额
type Balance struct {
	AccountAlias       string  `json:"accountAlias"`              // 账户唯一识别码
	Asset              string  `json:"asset"`                     // 资产
	Balance            float64 `json:"balance,string"`            // 总余额
	CrossWalletBalance float64 `json:"crossWalletBalance,string"` // 全仓余额
	CrossUnPnl         float64 `json:"crossUnPnl,string"`         // 全仓持仓未实现盈亏
	AvailableBalance   float64 `json:"availableBalance,string"`   // 下单可用余额
	MaxWithdrawAmount  float64 `json:"maxWithdrawAmount,string"`  // 最大可转出余额
}

// GetAccountService get account info
type GetAccountService struct {
	c *Client
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *Account, err error) {
	// GET /fapi/v3/account | 现有账户信息。 用户在单资产模式和多资产模式下会看到不同结果，响应部分的注释解释了两种模式下的不同。
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/account",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Account)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Account define account info | 账户 (多资产模式)
type Account struct {
	Assets                      []*AccountAsset    `json:"assets"`                             // 资产
	FeeTier                     int                `json:"feeTier"`                            // 手续费等级
	CanTrade                    bool               `json:"canTrade"`                           // 是否可以交易
	CanDeposit                  bool               `json:"canDeposit"`                         // 是否可以入金
	CanWithdraw                 bool               `json:"canWithdraw"`                        // 是否可以出金
	UpdateTime                  int64              `json:"updateTime"`                         // 保留字段，请忽略
	TotalInitialMargin          float64            `json:"totalInitialMargin,string"`          // 以USD计价的所需起始保证金总额
	TotalMaintMargin            float64            `json:"totalMaintMargin,string"`            // 以USD计价的维持保证金总额
	TotalWalletBalance          float64            `json:"totalWalletBalance,string"`          // 以USD计价的账户总余额
	TotalUnrealizedProfit       float64            `json:"totalUnrealizedProfit,string"`       // 以USD计价的持仓未实现盈亏总额
	TotalMarginBalance          float64            `json:"totalMarginBalance,string"`          // 以USD计价的保证金总余额
	TotalPositionInitialMargin  float64            `json:"totalPositionInitialMargin,string"`  // 以USD计价的持仓所需起始保证金(基于最新标记价格)
	TotalOpenOrderInitialMargin float64            `json:"totalOpenOrderInitialMargin,string"` // 以USD计价的当前挂单所需起始保证金(基于最新标记价格)
	TotalCrossWalletBalance     float64            `json:"totalCrossWalletBalance,string"`     // 以USD计价的全仓账户余额
	TotalCrossUnPnl             float64            `json:"totalCrossUnPnl,string"`             // 以USD计价的全仓持仓未实现盈亏总额
	AvailableBalance            float64            `json:"availableBalance,string"`            // 以USD计价的可用余额
	MaxWithdrawAmount           float64            `json:"maxWithdrawAmount,string"`           // 以USD计价的最大可转出余额
	Positions                   []*AccountPosition `json:"positions"`
}

// AccountAsset define account asset | 账户资产
type AccountAsset struct {
	Asset                  string  `json:"asset"`                         // 资产
	InitialMargin          float64 `json:"initialMargin,string"`          // 当前所需起始保证金
	MaintMargin            float64 `json:"maintMargin,string"`            // 维持保证金
	MarginBalance          float64 `json:"marginBalance,string"`          // 保证金余额
	MaxWithdrawAmount      float64 `json:"maxWithdrawAmount,string"`      // 最大可转出余额
	OpenOrderInitialMargin float64 `json:"openOrderInitialMargin,string"` // 当前挂单所需起始保证金(基于最新标记价格)
	PositionInitialMargin  float64 `json:"positionInitialMargin,string"`  // 持仓所需起始保证金(基于最新标记价格)
	UnrealizedProfit       float64 `json:"unrealizedProfit,string"`       // 未实现盈亏
	WalletBalance          float64 `json:"walletBalance,string"`          // 保证金余额
}

// AccountPosition define account position | 用户持仓风险V2
type AccountPosition struct {
	Isolated               bool             `json:"isolated"`                      // 是否是逐仓模式
	Leverage               string           `json:"leverage"`                      // 杠杆倍率
	InitialMargin          float64          `json:"initialMargin,string"`          // 当前所需起始保证金(基于最新标记价格)
	MaintMargin            float64          `json:"maintMargin,string"`            // 维持保证金
	OpenOrderInitialMargin float64          `json:"openOrderInitialMargin,string"` // 当前挂单所需起始保证金(基于最新标记价格)
	PositionInitialMargin  float64          `json:"positionInitialMargin,string"`  // 持仓所需起始保证金(基于最新标记价格)
	Symbol                 string           `json:"symbol"`                        // 交易对
	UnrealizedProfit       float64          `json:"unrealizedProfit,string"`       // 持仓未实现盈亏
	EntryPrice             float64          `json:"entryPrice,string"`             // 持仓成本价
	MaxNotional            float64          `json:"maxNotional,string"`            // 当前杠杆下用户可用的最大名义价值
	PositionSide           PositionSideType `json:"positionSide"`                  // 持仓方向
	PositionAmt            float64          `json:"positionAmt,string"`            // 持仓数量
	Notional               float64          `json:"notional,string"`               // 名义价值
	IsolatedWallet         string           `json:"isolatedWallet"`                //
	UpdateTime             int64            `json:"updateTime"`                    // 更新时间
}

// GetSymbolConfigService get account info
type GetSymbolConfigService struct {
	c *Client
}

// Do send request
func (s *GetSymbolConfigService) Do(ctx context.Context, opts ...RequestOption) (res []*SymbolConfig, err error) {
	// GET /fapi/v1/symbolConfig | 查询交易对上的基础配置
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/symbolConfig",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SymbolConfig define symbol config
type SymbolConfig struct {
	Symbol           string     `json:"symbol"`
	MarginType       MarginType `json:"marginType"` // 保证金模式
	IsAutoAddMargin  bool       `json:"isAutoAddMargin"`
	Leverage         int        `json:"leverage"`
	MaxNotionalValue float64    `json:"maxNotionalValue,string"`
}
