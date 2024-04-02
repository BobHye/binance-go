package spots

import (
	"binance-go/common"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/bitly/go-simplejson"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// SideType define side type of order (BUY/SELL)
type SideType string

// OrderType define order type
type OrderType string

// TimeInForceType define time in force type of order
type TimeInForceType string

// NewOrderRespType define response JSON verbosity
type NewOrderRespType string

// OrderStatusType defind order status type
type OrderStatusType string

// SymbolType define symbol type
type SymbolType string

// SymbolStatusType define symbol status type
type SymbolStatusType string

// SymbolFilterType define symbol filter type
type SymbolFilterType string

// UserDataEventType define spot user data event type
type UserDataEventType string

// MarginTransferType define margin transfer type
type MarginTransferType int

// MarginLoanStatusType define margin loan status type
type MarginLoanStatusType string

// MarginRepayStatusType define margin repay status type
type MarginRepayStatusType string

// FuturesTransferStatusType define futures transfer status type
type FuturesTransferStatusType string

// SideEffectType define side effect type for orders
type SideEffectType string

// FuturesTransferType define futures transfer type
type FuturesTransferType int

// TransactionType define transaction type
type TransactionType string

// LendingType define the type of lending (flexible saving, activity, ...)
type LendingType string

// StakingProduct define the staking product (locked staking, flexible defi staking, locked defi staking, ...)
type StakingProduct string

// StakingTransactionType define the staking transaction type (subscription, redemption, interest)
type StakingTransactionType string

// LiquidityOperationType define the type of adding/removing liquidity to a liquidity pool(COMBINATION, SINGLE)
type LiquidityOperationType string

// SwappingStatus define the status of swap when querying the swap history
type SwappingStatus int

// LiquidityRewardType define the type of reward we'd claim
type LiquidityRewardType int

// RewardClaimStatus define the status of claiming a reward
type RewardClaimStatus int

// RateLimitType define the rate limitation types
// see https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#enum-definitions
type RateLimitType string

// RateLimitInterval define the rate limitation intervals
type RateLimitInterval string

// Endpoints
const (
	baseAPIMainURL    = "https://api.binance.com"
	baseAPITestnetURL = "https://testnet.binance.vision"
)

// UseTestnet 将所有 API 切换到测试网络
var UseTestnet = false

// 重新定义标准包
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 全局枚举类
const (
	// 订单方向 (side):
	SideBuy  SideType = "BUY"
	SideSell SideType = "SELL"

	// 订单种类 (orderTypes, type):
	OrderTypeLimit              OrderType = "LIMIT"                // 限价单
	OrderTypeMarket             OrderType = "MARKET"               // 市价单
	OrderTypeStop               OrderType = "STOP"                 // 止损限价单
	OrderTypeTakeProfit         OrderType = "TAKE_PROFIT"          // 止盈限价单
	OrderTypeStopMarket         OrderType = "STOP_MARKET"          // 止损市价单
	OrderTypeTakeProfitMarket   OrderType = "TAKE_PROFIT_MARKET"   // 止盈市价单
	OrderTypeTrailingStopMarket OrderType = "TRAILING_STOP_MARKET" // 跟踪止损单

	// 有效方式 (timeInForce):
	TimeInForceTypeGTC TimeInForceType = "GTC" // 成交为止
	TimeInForceTypeIOC TimeInForceType = "IOC" // 无法立即成交(吃单)的部分就撤销
	TimeInForceTypeFOK TimeInForceType = "FOK" // 无法全部立即成交就撤销
	TmeInForceTypeGTX  TimeInForceType = "GTX" // 无法成为挂单方就撤销
	TmeInForceTypeGTD  TimeInForceType = "GTD" // 在特定时间之前有效，到期自动撤销

	// 响应类型
	NewOrderRespTypeACK    NewOrderRespType = "ACK" // 默认值
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"

	// 订单状态 (status):
	OrderStatusTypeNew             OrderStatusType = "NEW"              // 新建订单
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED" // 部分成交
	OrderStatusTypeFilled          OrderStatusType = "FILLED"           // 全部成交
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"         // 已撤销
	OrderStatusTypeRejected        OrderStatusType = "REJECTED"         // 订单被拒绝
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"          // 订单过期(根据timeInForce参数规则)
	OrderStatusTypeExpiredInMatch  OrderStatusType = "EXPIRED_IN_MATCH" // 订单被STP过期

	SymbolTypeSpot SymbolType = "SPOT"

	// 合约状态 (contractStatus, status):
	SymbolStatusTypePendingTrading SymbolStatusType = "PENDING_TRADING" // 待上市
	SymbolStatusTypeTrading        SymbolStatusType = "TRADING"         // 交易中
	SymbolStatusTypePreDelivering  SymbolStatusType = "PRE_DELIVERING"  // 预交割
	SymbolStatusTypeDelivering     SymbolStatusType = "DELIVERING"      // 交割中
	SymbolStatusTypeDelivered      SymbolStatusType = "DELIVERED"       // 已交割
	SymbolStatusTypePreSettle      SymbolStatusType = "PRE_SETTLE"      // 预结算
	SymbolStatusTypeSettling       SymbolStatusType = "SETTLING"        // 预结算
	SymbolStatusTypeClose          SymbolStatusType = "CLOSE"           // 已下架

	// 交易对过滤器类型
	SymbolFilterTypePriceFilter      SymbolFilterType = "PRICE_FILTER"        // 价格过滤器
	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"            // 订单尺寸
	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"     // 市价订单尺寸
	SymbolFilterTypeMaxNumOrders     SymbolFilterType = "MAX_NUM_ORDERS"      // 最多订单数
	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS" // 最多条件订单数
	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"       // 价格振幅过滤器
	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"        // 最小名义价值

	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"
)

// currentTimestamp 获取当前时间戳(毫秒)
func currentTimestamp() int64 {
	return FormatTimestamp(time.Now())
}

// FormatTimestamp 按照 Binance 的要求，将时间格式化为 Unix 时间戳（以毫秒为单位）
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// newJSON 将数据转成json
func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// getAPIEndpoint 根据 UseTestnet 标志返回 Rest API 的baseurl
func getAPIEndpoint() string {
	if UseTestnet {
		return baseAPITestnetURL
	}
	return baseAPIMainURL
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client 定义API客户端
type Client struct {
	APIKey     string
	SecretKey  string
	BaseUrl    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// 设置用户的请求选项
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseUrl, r.endpoint)
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
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
	if r.secType == secTypeAPIKey {
		header.Set("X-MBX-APIKEY", c.APIKey)
	}
	if r.secType == secTypeSigned {
		header.Set("X-MBX-APIKEY", c.APIKey)
		r.setParam(timestampKey, currentTimestamp()-c.TimeOffset)
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", mac.Sum(nil)))
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s%s", fullURL, queryString)
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
	c.debug("request: %#v", req)
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
		cerr := res.Body.Close()
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
		return nil, apiErr
	}
	return data, nil
}

func (c *Client) NewPingService() *PingService {
	return &PingService{c: c}
}

func (c *Client) NewServerTimeService() *ServerTimeService {
	return &ServerTimeService{c: c}
}
