package delivery

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/BobHye/binance-go/common"
	"github.com/bitly/go-simplejson"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// SideType define side type of order | 定义订单类型
type SideType string

// PositionSideType define position side type of order | 定义仓位订单类型
type PositionSideType string

// OrderType define order type | 定义订单类型
type OrderType string

// TimeInForceType define time in force type of order | 定义订单有效时间类型
type TimeInForceType string

// NewOrderRespType define response JSON verbosity | 定义新订单响应的 JSON 详细程度
type NewOrderRespType string

// OrderExecutionType define order execution type | 定义订单执行类型
type OrderExecutionType string

// OrderStatusType define order status type | 定义订单状态类型
type OrderStatusType string

// SymbolType define symbol type | 定义交易对类型
type SymbolType string

// SymbolStatusType define symbol status type | 定义交易对状态类型
type SymbolStatusType string

// SymbolFilterType define symbol filter type | 定义交易对过滤器类型
type SymbolFilterType string

// SideEffectType define side effect type for orders
type SideEffectType string

// WorkingType define working type
type WorkingType string

// MarginType define margin type
type MarginType string

// UserDataEventType define user data event type | 定义用户数据事件类型
type UserDataEventType string

// UserDataEventReasonType define reason type for user data event | 定义用户数据事件的原因类型
type UserDataEventReasonType string

// Redefining the standard package
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Endpoints
const (
	baseApiMainUrl    = "https://dapi.binance.com"
	baseApiTestnetUrl = "https://testnet.binancefuture.com"
)

// Global enums | 全局常量
const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"

	PositionSideTypeBoth  PositionSideType = "BOTH"
	PositionSideTypeLong  PositionSideType = "LONG"
	PositionSideTypeShort PositionSideType = "SHORT"

	OrderTypeLimit              OrderType = "LIMIT"
	OrderTypeMarket             OrderType = "MARKET"
	OrderTypeStop               OrderType = "STOP"
	OrderTypeStopMarket         OrderType = "STOP_MARKET"
	OrderTypeTakeProfit         OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket   OrderType = "TAKE_PROFIT_MARKET"
	OrderTypeTrailingStopMarket OrderType = "TRAILING_STOP_MARKET"

	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceTYPEIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceTypeGTX TimeInForceType = "GTX" // Good Till Crossing (Post Only)

	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
	NewOrderRespTypeFULL   NewOrderRespType = "FULL"

	OrderExecutionTypeNew         OrderExecutionType = "NEW"
	OrderExecutionTypePartialFill OrderExecutionType = "PARTIAL_FILL"
	OrderExecutionTypeFill        OrderExecutionType = "FILL"
	OrderExecutionTypeCanceled    OrderExecutionType = "CANCELED"
	OrderExecutionTypeCalculated  OrderExecutionType = "CALCULATED"
	OrderExecutionTypeExpired     OrderExecutionType = "EXPIRED"
	OrderExecutionTypeTrade       OrderExecutionType = "TRADE"

	OrderStatusTypeNew             OrderStatusType = "NEW"
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
	OrderStatusTypeFilled          OrderStatusType = "FILLED"
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"
	OrderStatusTypeTrade           OrderStatusType = "TRADE"

	SymbolTypeFuture SymbolType = "FUTURE"

	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE"

	SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING"
	SymbolStatusTypeTrading      SymbolStatusType = "TRADING"
	SymbolStatusTypePostTrading  SymbolStatusType = "POST_TRADING"
	SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
	SymbolStatusTypeHalt         SymbolStatusType = "HALT"
	SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
	SymbolStatusTypeBreak        SymbolStatusType = "BREAK"

	SymbolFilterTypeLotSize       SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypePrice         SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypePercentPrice  SymbolFilterType = "PERCENT_PRICE"
	SymbolFilterTypeMarketLotSize SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumOrders  SymbolFilterType = "MAX_NUM_ORDERS"

	SideEffectTypeNoSideEffect    SideEffectType = "NO_SIDE_EFFECT"
	SideEffectTypeMarginBuy       SideEffectType = "MARGIN_BUY"
	SideEffectTypeAutoBorrowRepay SideEffectType = "AUTO_BORROW_REPAY"
	SideEffectTypeAutoRepay       SideEffectType = "AUTO_REPAY"

	MarginTypeIsolated MarginType = "ISOLATED" // 逐仓
	MarginTypeCrossed  MarginType = "CROSSED"  // 全仓

	UserDataEventTypeListenKeyExpired    UserDataEventType = "listenKeyExpired"      // listenKey过期推送
	UserDataEventTypeMarginCall          UserDataEventType = "MARGIN_CALL"           // 追加保证金通知
	UserDataEventTypeAccountUpdate       UserDataEventType = "ACCOUNT_UPDATE"        // 账户更新事件(Balance 和 Position 更新推送)
	UserDataEventTypeOrderTradeUpdate    UserDataEventType = "ORDER_TRADE_UPDATE"    //
	UserDataEventTypeAccountConfigUpdate UserDataEventType = "ACCOUNT_CONFIG_UPDATE" // 杠杆倍数等账户配置 更新推送

	UserDataEventReasonTypeDeposit          UserDataEventReasonType = "DEPOSIT"
	UserDataEventReasonTypeWithdraw         UserDataEventReasonType = "WITHDRAW"
	UserDataEventReasonTypeOrder            UserDataEventReasonType = "ORDER"
	UserDataEventReasonTypeFundingFee       UserDataEventReasonType = "FUNDING_FEE"
	UserDataEventReasonTypeAdjustment       UserDataEventReasonType = "ADJUSTMENT"
	UserDataEventReasonTypeInsuranceClear   UserDataEventReasonType = "INSURANCE_CLEAR"
	UserDataEventReasonTypeAdminDeposit     UserDataEventReasonType = "ADMIN_DEPOSIT"
	UserDataEventReasonTypeAdminWithdraw    UserDataEventReasonType = "ADMIN_WITHDRAW"
	UserDataEventReasonTypeMarginTransfer   UserDataEventReasonType = "MARGIN_TRANSFER"
	UserDataEventReasonTypeMarginTypeChange UserDataEventReasonType = "MARGIN_TYPE_CHANGE"
	UserDataEventReasonTypeAssetTransfer    UserDataEventReasonType = "ASSET_TRANSFER"
	UserDataEventReasonTypeCoinSwapDeposit  UserDataEventReasonType = "COIN_SWAP_DEPOSIT"
	UserDataEventReasonTypeCoinSwapWithdraw UserDataEventReasonType = "COIN_SWAP_WITHDRAW"

	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"
)

func currentTimestamp() int64 {
	return int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, err
}

func getApiEndpoint() string {
	if Uset {
		return baseApiTestnetUrl
	}
	return baseApiMainUrl
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    getApiEndpoint(),
		UserAgent:  "Binance/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "Binance-golang", log.LstdFlags),
	}
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options form user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint) // 拼接完整url
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow) // 设置recvWindow参数
	}
	if r.secType == secTypeSigned {
		r.setParam(timestampKey, currentTimestamp()) // 设置timestamp参数
	}
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
		header.Set("X-MBX-APIKEY", c.APIKey)
	}

	if r.secType == secTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", (mac.Sum(nil))))
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.debug("full url: %s, body: %s", fullURL, bodyString)

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request:%#v", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cErr := res.Body.Close()
		if err == nil && cErr != nil {
			err = cErr
		}
	}()
	c.debug("response:%#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, apiErr
	}
	return data, nil
}

// NewPingService init ping service
func (c *Client) NewPingService() *PingService {
	return &PingService{c: c}
}

// NewServerTimeService init server time service
func (c *Client) NewServerTimeService() *ServerTimeService {
	return &ServerTimeService{c: c}
}

// NewSetServerTimeService init set server time service
func (c *Client) NewSetServerTimeService() *SetServerTimeService {
	return &SetServerTimeService{c: c}
}

// NewKlinesService init klines service
func (c *Client) NewKlinesService() *KlinesService {
	return &KlinesService{c: c}
}

// NewListPriceChangeStatsService init list prices change stats service | 初始化列表价格变更统计服务
func (c *Client) NewListPriceChangeStatsService() *ListPriceChangeStatsService {
	return &ListPriceChangeStatsService{c: c}
}

// NewListPricesService init listing prices service
func (c *Client) NewListPricesService() *ListPricesService {
	return &ListPricesService{c: c}
}

// NewListBookTickersService init listing booking tickers service
func (c *Client) NewListBookTickersService() *ListBookTickersService {
	return &ListBookTickersService{c: c}
}

// NewStartUserStreamService init starting user stream service
func (c *Client) NewStartUserStreamService() *StartUserStreamService {
	return &StartUserStreamService{c: c}
}

// NewKeepaliveUserStreamService init keep alive user stream service
func (c *Client) NewKeepaliveUserStreamService() *KeepaliveUserStreamService {
	return &KeepaliveUserStreamService{c: c}
}

// NewCloseUserStreamService init closing user stream service
func (c *Client) NewCloseUserStreamService() *CloseUserStreamService {
	return &CloseUserStreamService{c: c}
}

// NewExchangeInfoService init exchange info service
func (c *Client) NewExchangeInfoService() *ExchangeInfoService {
	return &ExchangeInfoService{c: c}
}

// NewCreateOrderService init creating order service
func (c *Client) NewCreateOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

// NewGetOrderService init get order service
func (c *Client) NewGetOrderService() *GetOrderService {
	return &GetOrderService{c: c}
}

// NewCancelOrderService init cancel order service
func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

// NewCancelAllOpenOrdersService init cancel all open orders service
func (c *Client) NewCancelAllOpenOrdersService() *CancelAllOpenOrdersService {
	return &CancelAllOpenOrdersService{c: c}
}

// NewListOpenOrdersService init list open orders service
func (c *Client) NewListOpenOrdersService() *ListOpenOrdersService {
	return &ListOpenOrdersService{c: c}
}

// NewListLiquidationOrdersService init funding rate service
func (c *Client) NewListLiquidationOrdersService() *ListLiquidationOrdersService {
	return &ListLiquidationOrdersService{c: c}
}

// NewGetAccountService init account service
func (c *Client) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

// NewGetBalanceService init balance service
func (c *Client) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

// NewGetPositionRiskService init getting position risk service
func (c *Client) NewGetPositionRiskService() *GetPositionRiskService {
	return &GetPositionRiskService{c: c}
}

// NewChangeLeverageService init change leverage service
func (c *Client) NewChangeLeverageService() *ChangeLeverageService {
	return &ChangeLeverageService{c: c}
}

// NewChangeMarginTypeService init change margin type service
func (c *Client) NewChangeMarginTypeService() *ChangeMarginTypeService {
	return &ChangeMarginTypeService{c: c}
}

// NewUpdatePositionMarginService init update position margin
func (c *Client) NewUpdatePositionMarginService() *UpdatePositionMarginService {
	return &UpdatePositionMarginService{c: c}
}

// NewChangePositionModeService init change position mode service
func (c *Client) NewChangePositionModeService() *ChangePositionModeService {
	return &ChangePositionModeService{c: c}
}

// NewGetPositionModeService init get position mode service
func (c *Client) NewGetPositionModeService() *GetPositionModeService {
	return &GetPositionModeService{c: c}
}
