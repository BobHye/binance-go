package futures

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"github.com/BobHye/binance-go/common"
	"github.com/bitly/go-simplejson"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// SideType define side type of order | 订单方向类型
type SideType string

// PositionSideType define position side type of order | 持仓方向类型
type PositionSideType string

// OrderType define order type | 订单类型
type OrderType string

// TimeInForceType define time in force type of order | 有效方式类型 (timeInForce)
type TimeInForceType string

// NewOrderRespType define response JSON verbosity | 响应类型 (newOrderRespType) [ACK/RESULT]
type NewOrderRespType string

// OrderExecutionType define order execution type | 订单执行类型
type OrderExecutionType string

// OrderStatusType define order status type | 本次事件的具体执行类型[NEW/CANCELED/CALCULATED/EXPIRED/TRADE/AMENDMENT]
type OrderStatusType string

// SymbolType define symbol type | 交易对类型
type SymbolType string

// SymbolStatusType define symbol status type | 交易对状态类型
type SymbolStatusType string //

// SymbolFilterType define symbol filter type | 交易对过滤器类型
type SymbolFilterType string

// SideEffectType define side effect type for orders
type SideEffectType string

// WorkingType define working type | 条件价格触发类型 (workingType) [MARK_PRICE/CONTRACT_PRICE]
type WorkingType string

// MarginType define margin type | 保证金模式
type MarginType string

// ContractType define contract type | 合约类型 PERPETUAL 永续合约/CURRENT_MONTH 当月交割合约/NEXT_MONTH 次月交割合约/CURRENT_QUARTER 当季交割合约/NEXT_QUARTER 次季交割合约
type ContractType string

// UserDataEventType define user data event type | 账户信息事件类型
type UserDataEventType string

// UserDataEventReasonType define reason type for user data event | 账户信息事件推出的原因类型
type UserDataEventReasonType string

// ForceOrderCloseType define reason type for force order | 强平订单类型
type ForceOrderCloseType string

// Redefining the standard package
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Endpoints
const (
	baseApiMainUrl    = "https://fapi.binance.com"
	baseApiTestnetUrl = "https://testnet.binancefuture.com"
)

// Global enums
const (
	SideTypeBuy  SideType = "BUY"  // 买入
	SideTypeSell SideType = "SELL" // 卖出

	PositionSideTypeBoth  PositionSideType = "BOTH"  // 单一持仓方向
	PositionSideTypeLong  PositionSideType = "LONG"  // 多头(双向持仓下)
	PositionSideTypeShort PositionSideType = "SHORT" // 空头(双向持仓下)

	OrderTypeLimit              OrderType = "LIMIT"                // 限价单
	OrderTypeMarket             OrderType = "MARKET"               // 市价单
	OrderTypeStop               OrderType = "STOP"                 // 止损限价单
	OrderTypeStopMarket         OrderType = "STOP_MARKET"          // 止损市价单
	OrderTypeTakeProfit         OrderType = "TAKE_PROFIT"          // 止盈限价单
	OrderTypeTakeProfitMarket   OrderType = "TAKE_PROFIT_MARKET"   // 止盈市价单
	OrderTypeTrailingStopMarket OrderType = "TRAILING_STOP_MARKET" // 跟踪止损单

	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel 成交为止
	TimeInForceTypeIOC TimeInForceType = "IOC" // Immediate or Cancel 无法立即成交(吃单)的部分就撤销
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill 无法全部立即成交就撤销
	TimeInForceTypeGTX TimeInForceType = "GTX" // Good Till Crossing 无法成为挂单方就撤销

	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"

	OrderExecutionTypeNew        OrderExecutionType = "NEW"
	OrderExecutionTypeTrade      OrderExecutionType = "TRADE"
	OrderExecutionTypeCanceled   OrderExecutionType = "CANCELED"
	OrderExecutionTypeCalculated OrderExecutionType = "CALCULATED"
	OrderExecutionTypeExpired    OrderExecutionType = "EXPIRED"
	OrderExecutionTypeAmendment  OrderExecutionType = "AMENDMENT"

	OrderStatusTypeNew             OrderStatusType = "NEW"              // 新建订单
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED" // 部分成交
	OrderStatusTypeFilled          OrderStatusType = "FILLED"           // 全部成交
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"         // 已撤销
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"          // 订单过期(根据timeInForce参数规则)

	SymbolTypeFuture SymbolType = "FUTURE"

	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"     // 标记价格
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE" // 合约价格

	SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING" //
	SymbolStatusTypeTrading      SymbolStatusType = "TRADING"     // 交易中
	SymbolStatusTypePostTrade    SymbolStatusType = "POST_TRADING"
	SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
	SymbolStatusTypeHalt         SymbolStatusType = "HALT"
	SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
	SymbolStatusTypeBreak        SymbolStatusType = "BREAK"

	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypePrice            SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumOrders     SymbolFilterType = "MAX_NUM_ORDERS"
	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"

	SideEffectTypeNoSideEffect SideEffectType = "NO_SIDE_EFFECT"
	SideEffectTypeMarginBuy    SideEffectType = "MARGIN_BUY"
	SideEffectTypeAutoRepay    SideEffectType = "AUTO_REPAY"

	MarginTypeIsolated MarginType = "ISOLATED" // 逐仓
	MarginTypeCrossed  MarginType = "CROSSED"  // 全仓

	ContractTypePerpetual ContractType = "PERPETUAL" // 永续合约

	UserDataEventTypeListenKeyExpired    UserDataEventType = "listenKeyExpired"      // listenKey 过期推送
	UserDataEventTypeMarginCall          UserDataEventType = "MARGIN_CALL"           // 追加保证金通知
	UserDataEventTypeAccountUpdate       UserDataEventType = "ACCOUNT_UPDATE"        // Balance和Position更新推送
	UserDataEventTypeOrderTradeUpdate    UserDataEventType = "ORDER_TRADE_UPDATE"    // 订单/交易 更新推送
	UserDataEventTypeAccountConfigUpdate UserDataEventType = "ACCOUNT_CONFIG_UPDATE" // 杠杆倍数等账户配置 更新推送

	UserDataEventReasonTypeDeposit             UserDataEventReasonType = "DEPOSIT"
	UserDataEventReasonTypeWithdraw            UserDataEventReasonType = "WITHDRAW"
	UserDataEventReasonTypeOrder               UserDataEventReasonType = "ORDER"
	UserDataEventReasonTypeFundingFee          UserDataEventReasonType = "FUNDING_FEE"
	UserDataEventReasonTypeWithdrawReject      UserDataEventReasonType = "WITHDRAW_REJECT"
	UserDataEventReasonTypeAdjustment          UserDataEventReasonType = "ADJUSTMENT"
	UserDataEventReasonTypeInsuranceClear      UserDataEventReasonType = "INSURANCE_CLEAR"
	UserDataEventReasonTypeAdminDeposit        UserDataEventReasonType = "ADMIN_DEPOSIT"
	UserDataEventReasonTypeAdminWithdraw       UserDataEventReasonType = "ADMIN_WITHDRAW"
	UserDataEventReasonTypeMarginTransfer      UserDataEventReasonType = "MARGIN_TRANSFER"
	UserDataEventReasonTypeMarginTypeChange    UserDataEventReasonType = "MARGIN_TYPE_CHANGE"
	UserDataEventReasonTypeAssetTransfer       UserDataEventReasonType = "ASSET_TRANSFER"
	UserDataEventReasonTypeOptionsPremiumFee   UserDataEventReasonType = "OPTIONS_PREMIUM_FEE"
	UserDataEventReasonTypeOptionsSettleProfit UserDataEventReasonType = "OPTIONS_SETTLE_PROFIT"

	ForceOrderCloseTypeLiquidation ForceOrderCloseType = "LIQUIDATION" // 强平单
	ForceOrderCloseTypeADL         ForceOrderCloseType = "ADL"         // ADL减仓单

	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"
)

// currentTimestamp 当前时间戳
func currentTimestamp() int64 {
	return int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
}

// newJSON 将byte数组转成json对象
func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// getApiEndpoint return the base endpoint of the WS according the UseTestnet flag | 根据UseTestnet变量，返回模拟盘/实盘 baseurl
func getApiEndpoint() string {
	if UseTestnet {
		return baseApiTestnetUrl
	}
	return baseApiMainUrl
}

func MustFloat64(s string) float64 {
	def := 0.0
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return f
}

// Client 定义API客户端
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
// 创建新的API Client
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    getApiEndpoint(),
		UserAgent:  "Binance/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
	}
}

// NewProxiedClient passing a proxy url
func NewProxiedClient(apiKey, secretKey, proxyUrl string) *Client {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   getApiEndpoint(),
		UserAgent: "Binance/goland",
		HTTPClient: &http.Client{
			Transport: tr,
		},
		Logger: log.New(os.Stderr, "Binance-goland", log.LstdFlags),
	}
}

type doFunc func(req *http.Request) (*http.Response, error)

// debug 输出调试信息
func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

// parseRequest 解释请求
func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint) // 拼接成完整url
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow) // 设置recvWindow参数
	}
	if r.secType == secTypeSigned {
		r.setParam(timestampKey, currentTimestamp()-c.TimeOffset) // 设置timestamp参数
	}
	queryString := r.query.Encode() // 将值编码为“URL" 编码
	body := &bytes.Buffer{}
	bodyString := r.form.Encode() // 将值编码为“URL" 编码
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
		// 如果需要 API-key，应当在HTTP头中以X-MBX-APIKEY字段传递
		header.Set("X-MBX-APIKEY", c.APIKey)
	}

	if r.secType == secTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		// 签名使用HMAC SHA256算法. API-KEY所对应的API-Secret作为 HMAC SHA256 的密钥，其他所有参数作为HMAC SHA256的操作对象，得到的输出即为签名
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		// 计算检验和
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		//checksum := mac.Sum(nil)
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", (mac.Sum(nil)))) // 设置签名
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
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

// callAPI 调用API请求
func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, header *http.Header, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}

	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v", req)
	fun := c.do
	if fun != nil {
		fun = c.HTTPClient.Do
	}
	res, err := fun(req)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)
	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, &http.Header{}, apiErr
	}
	return data, &res.Header, nil
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

// NewDepthService init depth service
func (c *Client) NewDepthService() *DepthService {
	return &DepthService{c: c}
}

// NewAggTradesService init aggregate trades service
func (c *Client) NewAggTradesService() *AggTradesService {
	return &AggTradesService{c: c}
}

// NewRecentTradesService init recent trades service
func (c *Client) NewRecentTradesService() *RecentTradesService {
	return &RecentTradesService{c: c}
}

// NewKlinesService init klines service
func (c *Client) NewKlinesService() *KlinesService {
	return &KlinesService{c: c}
}

// NewListPriceChangeStatsService init list prices change stats service
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

// NewCreateOrderService init creating order service
func (c *Client) NewCreateOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

// NewCreateBatchOrdersService init creating batch order service
func (c *Client) NewCreateBatchOrdersService() *CreateBatchOrdersService {
	return &CreateBatchOrdersService{c: c}
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

// NewCancelMultipleOrdersService init cancel multiple orders service
func (c *Client) NewCancelMultipleOrdersService() *CancelMultiplesOrdersService {
	return &CancelMultiplesOrdersService{c: c}
}

// NewGetOpenOrderService init get open order service
func (c *Client) NewGetOpenOrderService() *GetOpenOrderService {
	return &GetOpenOrderService{c: c}
}

// NewListOpenOrdersService init list open orders service
func (c *Client) NewListOpenOrdersService() *ListOpenOrdersService {
	return &ListOpenOrdersService{c: c}
}

// NewListOrdersService init listing orders service
func (c *Client) NewListOrdersService() *ListOrdersService {
	return &ListOrdersService{c: c}
}

// NewGetAccountService init getting account service
func (c *Client) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

// NewGetBalanceService init getting balance service
func (c *Client) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

// NewGetPositionRiskService init getting position risk service
func (c *Client) NewGetPositionRiskService() *GetPositionRiskService {
	return &GetPositionRiskService{c: c}
}

// NewHistoricalTradesService init listing trades service
func (c *Client) NewHistoricalTradesService() *HistoricalTradesService {
	return &HistoricalTradesService{c: c}
}

// NewListAccountTradeService init account trade list service
func (c *Client) NewListAccountTradeService() *ListAccountTradeService {
	return &ListAccountTradeService{c: c}
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

// NewCommissionRateService returns commission rate
func (c *Client) NewCommissionRateService() *CommissionRateService {
	return &CommissionRateService{c: c}
}
