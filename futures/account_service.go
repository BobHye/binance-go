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
	AccountAlias       string `json:"accountAlias"`              // 账户唯一识别码
	Asset              string `json:"asset"`                     // 资产
	Balance            string `json:"balance,string"`            // 总余额
	CrossWalletBalance string `json:"crossWalletBalance,string"` // 全仓余额
	CrossUnPnl         string `json:"crossUnPnl,string"`         // 全仓持仓未实现盈亏
	AvailableBalance   string `json:"availableBalance,string"`   // 下单可用余额
	MaxWithdrawAmount  string `json:"maxWithdrawAmount,string"`  // 最大可转出余额
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
	TotalInitialMargin          string             `json:"totalInitialMargin,string"`          // 以USD计价的所需起始保证金总额
	TotalMaintMargin            string             `json:"totalMaintMargin,string"`            // 以USD计价的维持保证金总额
	TotalWalletBalance          string             `json:"totalWalletBalance,string"`          // 以USD计价的账户总余额
	TotalUnrealizedProfit       string             `json:"totalUnrealizedProfit,string"`       // 以USD计价的持仓未实现盈亏总额
	TotalMarginBalance          string             `json:"totalMarginBalance,string"`          // 以USD计价的保证金总余额
	TotalPositionInitialMargin  string             `json:"totalPositionInitialMargin,string"`  // 以USD计价的持仓所需起始保证金(基于最新标记价格)
	TotalOpenOrderInitialMargin string             `json:"totalOpenOrderInitialMargin,string"` // 以USD计价的当前挂单所需起始保证金(基于最新标记价格)
	TotalCrossWalletBalance     string             `json:"totalCrossWalletBalance,string"`     // 以USD计价的全仓账户余额
	TotalCrossUnPnl             string             `json:"totalCrossUnPnl,string"`             // 以USD计价的全仓持仓未实现盈亏总额
	AvailableBalance            string             `json:"availableBalance,string"`            // 以USD计价的可用余额
	MaxWithdrawAmount           string             `json:"maxWithdrawAmount,string"`           // 以USD计价的最大可转出余额
	Positions                   []*AccountPosition `json:"positions"`
}

// AccountAsset define account asset | 账户资产
type AccountAsset struct {
	Asset                  string `json:"asset"`                         // 资产
	InitialMargin          string `json:"initialMargin,string"`          // 当前所需起始保证金
	MaintMargin            string `json:"maintMargin,string"`            // 维持保证金
	MarginBalance          string `json:"marginBalance,string"`          // 保证金余额
	MaxWithdrawAmount      string `json:"maxWithdrawAmount,string"`      // 最大可转出余额
	OpenOrderInitialMargin string `json:"openOrderInitialMargin,string"` // 当前挂单所需起始保证金(基于最新标记价格)
	PositionInitialMargin  string `json:"positionInitialMargin,string"`  // 持仓所需起始保证金(基于最新标记价格)
	UnrealizedProfit       string `json:"unrealizedProfit,string"`       // 未实现盈亏
	WalletBalance          string `json:"walletBalance,string"`          // 保证金余额
}

// AccountPosition define account position | 用户持仓风险V2
type AccountPosition struct {
	Isolated               bool             `json:"isolated"`                      // 是否是逐仓模式
	Leverage               string           `json:"leverage"`                      // 杠杆倍率
	InitialMargin          string           `json:"initialMargin,string"`          // 当前所需起始保证金(基于最新标记价格)
	MaintMargin            string           `json:"maintMargin,string"`            // 维持保证金
	OpenOrderInitialMargin string           `json:"openOrderInitialMargin,string"` // 当前挂单所需起始保证金(基于最新标记价格)
	PositionInitialMargin  string           `json:"positionInitialMargin,string"`  // 持仓所需起始保证金(基于最新标记价格)
	Symbol                 string           `json:"symbol"`                        // 交易对
	UnrealizedProfit       string           `json:"unrealizedProfit,string"`       // 持仓未实现盈亏
	EntryPrice             string           `json:"entryPrice,string"`             // 持仓成本价
	MaxNotional            string           `json:"maxNotional,string"`            // 当前杠杆下用户可用的最大名义价值
	PositionSide           PositionSideType `json:"positionSide"`                  // 持仓方向
	PositionAmt            string           `json:"positionAmt,string"`            // 持仓数量
	Notional               string           `json:"notional,string"`               // 名义价值
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
