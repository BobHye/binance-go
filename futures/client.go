package futures

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/BobHye/binance-go/api"
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

type SideType string           // 订单方向类型
type PositionSideType string   // 持仓方向类型
type OrderType string          // 订单类型
type TimeInForceType string    // 有效方式类型 (timeInForce)
type NewOrderRespType string   // 响应类型 (newOrderRespType) [ACK/RESULT]
type OrderExecutionType string // 订单执行类型
type OrderStatusType string    // 本次事件的具体执行类型[NEW/CANCELED/CALCULATED/EXPIRED/TRADE/AMENDMENT]
type SymbolType string         // 交易对类型
type SymbolStatusType string   // 交易对状态类型
type SymbolFilterType string   // 交易对过滤器类型
type SideEffectType string
type WorkingType string             // 条件价格触发类型 (workingType) [MARK_PRICE/CONTRACT_PRICE]
type MarginType string              // 保证金模式
type ContractType string            // 合约类型 PERPETUAL 永续合约/CURRENT_MONTH 当月交割合约/NEXT_MONTH 次月交割合约/CURRENT_QUARTER 当季交割合约/NEXT_QUARTER 次季交割合约
type UserDataEventType string       // 账户信息事件类型
type UserDataEventReasonType string // 账户信息事件推出的原因类型
type ForceOrderCloseType string     // 强平订单类型

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Endpoints
const (
	ApiUseTestnet     = false
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
	TimeInForceTypeGTD TimeInForceType = "GTD" // Good Till Date 在特定时间之前有效，到期自动撤销

	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"

	OrderExecutionTypeNew         OrderExecutionType = "NEW"
	OrderExecutionTypePartialFill OrderExecutionType = "PARTIAL_FILL"
	OrderExecutionTypeFill        OrderExecutionType = "FILL"
	OrderExecutionTypeCanceled    OrderExecutionType = "CANCELED"
	OrderExecutionTypeCalculated  OrderExecutionType = "CALCULATED"
	OrderExecutionTypeExpired     OrderExecutionType = "EXPIRED"
	OrderExecutionTypeTrade       OrderExecutionType = "TRADE"

	OrderStatusTypeNew             OrderStatusType = "NEW"              // 新建订单
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED" // 部分成交
	OrderStatusTypeFilled          OrderStatusType = "FILLED"           // 全部成交
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"         // 已撤销
	OrderStatusTypeRejected        OrderStatusType = "REJECTED"         // 订单被拒绝
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"          // 订单过期(根据timeInForce参数规则)
	OrderStatusTypeNewInsurance    OrderStatusType = "NEW_INSURANCE"    // 订单被STP过期

	SymbolTypeFuture SymbolType = "FUTURE"

	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"     // 标记价格
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE" // 合约价格

	SymbolStatusTypePendingTrading SymbolStatusType = "PENDING_TRADING" // 待上市
	SymbolStatusTypeTrading        SymbolStatusType = "TRADING"         // 交易中
	SymbolStatusTypePreDelivering  SymbolStatusType = "PRE_DELIVERING"  // 预交割
	SymbolStatusTypeDelivering     SymbolStatusType = "DELIVERING"      // 交割中
	SymbolStatusTypeDelivered      SymbolStatusType = "DELIVERED"       // 已交割
	SymbolStatusTypePreSettle      SymbolStatusType = "PRE_SETTLE"      // 预结算
	SymbolStatusTypeSettling       SymbolStatusType = "SETTLING"        // 结算中
	SymbolStatusTypeClose          SymbolStatusType = "CLOSE"           // 已下架

	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypePrice            SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumOrders     SymbolFilterType = "MAX_NUM_ORDERS"
	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"

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

// getApiEndpoint 根据UseTestnet变量，返回模拟盘/实盘 baseurl
func getApiEndpoint() string {
	if ApiUseTestnet {
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

type doFunc func(req *http.Request) (*http.Response, error)

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

// NewClient 创建新的API Client
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

// debug 输出调试信息
func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

// parseRequest 解释请求
func (c *Client) parseRequest(r *api.Request, opts ...api.RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.Validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint) // 拼接成完整url
	if r.RecvWindow > 0 {
		r.SetParam(recvWindowKey, r.RecvWindow) // 设置recvWindow参数
	}
	if r.SecType == api.SecTypeSigned {
		r.SetParam(timestampKey, currentTimestamp()-c.TimeOffset) // 设置timestamp参数
	}
	queryString := r.Query.Encode() // 将值编码为“URL" 编码
	body := &bytes.Buffer{}
	bodyString := r.Form.Encode() // 将值编码为“URL" 编码
	header := http.Header{}
	if r.Header != nil {
		header = r.Header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.SecType == api.SecTypeAPIKey || r.SecType == api.SecTypeSigned {
		// 如果需要 API-key，应当在HTTP头中以X-MBX-APIKEY字段传递
		header.Set("X-MBX-APIKEY", c.APIKey)
	}

	if r.SecType == api.SecTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		// 签名使用HMAC SHA256算法. API-KEY所对应的API-Secret作为 HMAC SHA256 的密钥，其他所有参数作为HMAC SHA256的操作对象，得到的输出即为签名
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		// 计算检验和
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		checksum := mac.Sum(nil)
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", checksum)) // 设置签名
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

	r.FullURL = fullURL
	r.Header = header
	r.Body = body
	return nil
}

// callAPI 调用API请求
func (c *Client) callAPI(ctx context.Context, r *api.Request, opts ...api.RequestOption) (data []byte, header *http.Header, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}

	req, err := http.NewRequest(r.Method, r.FullURL, r.Body)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.Header
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

// NewServerTimeService init server time service
func (c *Client) NewServerTimeService() *ServerTimeService {
	return &ServerTimeService{c: c}
}
