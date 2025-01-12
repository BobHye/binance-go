package margin

import (
	"context"
	"net/http"
)

// GetAccountSnapshotService all account orders; active, canceled, or filled
type GetAccountSnapshotService struct {
	c           *Client
	accountType string
	startTime   *int64
	endTime     *int64
	limit       *int
}

// SetType set account type ("SPOT", "MARGIN", "FUTURES")
func (s *GetAccountSnapshotService) SetType(accountType string) *GetAccountSnapshotService {
	s.accountType = accountType
	return s
}

// SetStartTime set starttime
func (s *GetAccountSnapshotService) SetStartTime(startTime int64) *GetAccountSnapshotService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endtime
func (s *GetAccountSnapshotService) SetEndTime(endTime int64) *GetAccountSnapshotService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *GetAccountSnapshotService) SetLimit(limit int) *GetAccountSnapshotService {
	s.limit = &limit
	return s
}

// Do send request
func (s *GetAccountSnapshotService) Do(ctx context.Context, opts ...RequestOption) (res *Snapshot, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/accountSnapshot",
		secType:  secTypeSigned,
	}
	r.setParam("type", s.accountType)

	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return &Snapshot{}, err
	}
	res = new(Snapshot)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return &Snapshot{}, err
	}
	return res, nil
}

// Snapshot define snapshot
type Snapshot struct {
	Code     int            `json:"code"`
	Msg      string         `json:"msg"`
	Snapshot []*SnapshotVos `json:"snapshotVos"`
}

// SnapshotVos define content of a snapshot
type SnapshotVos struct {
	Data       *SnapshotData `json:"data"`
	Type       string        `json:"type"`
	UpdateTime int64         `json:"updateTime"`
}

// SnapshotData define content of a snapshot
type SnapshotData struct {
	MarginLevel         string `json:"marginLevel"`
	TotalAssetOfBtc     string `json:"totalAssetOfBtc"`
	TotalLiabilityOfBtc string `json:"totalLiabilityOfBtc"`
	TotalNetAssetOfBtc  string `json:"totalNetAssetOfBtc"`

	Balances   []*SnapshotBalances   `json:"balances"`
	UserAssets []*SnapshotUserAssets `json:"userAssets"`
	Assets     []*SnapshotAssets     `json:"assets"`
	Positions  []*SnapshotPositions  `json:"position"`
}

// SnapshotBalances define snapshot balances
type SnapshotBalances struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// SnapshotUserAssets define snapshot user assets
type SnapshotUserAssets struct {
	Asset    string `json:"asset"`
	Borrowed string `json:"borrowed"`
	Free     string `json:"free"`
	Interest string `json:"interest"`
	Locked   string `json:"locked"`
	NetAsset string `json:"netAsset"`
}

// SnapshotAssets define snapshot assets
type SnapshotAssets struct {
	Asset         string `json:"asset"`
	MarginBalance string `json:"marginBalance"`
	WalletBalance string `json:"walletBalance"`
}

// SnapshotPositions define snapshot positions
type SnapshotPositions struct {
	EntryPrice       string `json:"entryPrice"`
	MarkPrice        string `json:"markPrice"`
	PositionAmt      string `json:"positionAmt"`
	Symbol           string `json:"symbol"`
	UnRealizedProfit string `json:"unRealizedProfit"`
}

// GetAPIKeyPermission get API Key permission info
type GetAPIKeyPermission struct {
	c *Client
}

// Do send request
func (s *GetAPIKeyPermission) Do(ctx context.Context, opts ...RequestOption) (res *APIKeyPermission, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/account/apiRestrictions",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(APIKeyPermission)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// APIKeyPermission define API key permission
type APIKeyPermission struct {
	IPRestrict                     bool   `json:"ipRestrict"`
	CreateTime                     uint64 `json:"createTime"`
	EnableWithdrawals              bool   `json:"enableWithdrawals"`
	EnableInternalTransfer         bool   `json:"enableInternalTransfer"`
	PermitsUniversalTransfer       bool   `json:"permitsUniversalTransfer"`
	EnableVanillaOptions           bool   `json:"enableVanillaOptions"`
	EnableReading                  bool   `json:"enableReading"`
	EnableFutures                  bool   `json:"enableFutures"`
	EnableMargin                   bool   `json:"enableMargin"`
	EnableSpotAndMarginTrading     bool   `json:"enableSpotAndMarginTrading"`
	TradingAuthorityExpirationTime uint64 `json:"tradingAuthorityExpirationTime"`
}
