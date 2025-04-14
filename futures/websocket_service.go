package common

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Endpoints
const (
	baseWsMainUrl          = "wss://fstream.binance.com/ws"
	baseWsTestnetUrl       = "wss://stream.binancefuture.com/ws"
	baseCombinedMainURL    = "wss://fstream.binance.com/stream?streams="
	baseCombinedTestnetURL = "wss://stream.binancefuture.com/stream?streams="
)

var (
	// HeartbeatInterval 如果启用了 WebsocketKeepalive，则发送 ping/pong 消息的间隔
	//HeartbeatInterval = 5 * time.Second
	// WebsocketKeepalive 允许发送 ping/pong 消息以检查连接稳定性
	//WebsocketKeepalive = true
	// 将所有 WS 流从生产环境切换到测试网络
	UseTestnet = false
)

// getWsEndpoint | 根据 UseTestnet 标志返回 WS 的基本端点
func getWsEndpoint() string {
	if UseTestnet {
		return baseWsTestnetUrl
	}
	return baseWsMainUrl
}

func getCombinedEndpoint() string {
	if UseTestnet {
		return baseCombinedTestnetURL
	}
	return baseCombinedMainURL
}

// WsAggTradeEvent 定义websocket @aggTrade 事件
type WsAggTradeEvent struct {
	Event            string  `json:"e"`        // 事件类型
	Time             string  `json:"E"`        // 事件时间
	Symbol           string  `json:"s"`        // 交易对
	AggregateTradeID string  `json:"a"`        // 归集成交 ID
	Price            string  `json:"p,string"` // 成交价格
	Quantity         float64 `json:"q,string"` // 成交量
	FirstTradeID     int64   `json:"f"`        // 被归集的首个交易ID
	LastTradeID      int64   `json:"l"`        // 被归集的末次交易ID
	TradeTime        int64   `json:"T"`        // 成交时间
	Maker            bool    `json:"m"`        // 买方是否是做市方。如true，则此次成交是一个主动卖出单，否则是一个主动买入单。
}

// WsAggTradeHandler 处理 websocket，推送单个交易对的的聚合交易信息
type WsAggTradeHandler func(event *WsAggTradeEvent)
type WsAggTradeConnect func()

// WsAggTradeServe 提供 websocket 服务，推送单个交易对的订单聚合交易信息。
func WsAggTradeServe(symbol string, handler WsAggTradeHandler, onConnect WsAggTradeConnect) (client *WebSocketClient, err error) {
	endpoint := fmt.Sprintf("%s/%s@aggTrade", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	client = NewClient(cfg)
	client.OnMessage = func(message []byte) {
		event := new(WsAggTradeEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			//errHandler(err)
			return
		}
		handler(event)
	}
	client.OnConnect = func() {
		onConnect()
	}
	err = client.Connect()
	return client, err
}

// WsMarkPriceEvent 定义 websocket markPriceUpdate 事件。
type WsMarkPriceEvent struct {
	Event                string `json:"e"` // 事件类型
	Time                 int64  `json:"E"` // 事件时间
	Symbol               string `json:"s"` // 交易对
	MarkPrice            string `json:"p"` // 标记价格
	EstimatedSettlePrice string `json:"P"` // 预估结算价,仅在结算前最后一小时有参考价值
	IndexPrice           string `json:"i"` // 指数价格
	FundingRate          string `json:"r"` // 资金费率，对非永续合约显示""
	NextFundingTime      int64  `json:"T"` // 下个资金时间,对非永续合约显示0
}

// WsMarkPriceHandler 处理单个交易对价格和资金费率
type WsMarkPriceHandler func(event *WsMarkPriceEvent)
type WsMarkPriceConnect func()

func wsMarkPriceServe(endpoint string, handler WsMarkPriceHandler, onConnect WsMarkPriceConnect) (client *WebSocketClient, err error) {
	cfg := newWsConfig(endpoint)
	client = NewClient(cfg)
	client.OnMessage = func(message []byte) {
		event := new(WsMarkPriceEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			return
		}
		handler(event)
	}
	client.OnConnect = func() {
		onConnect()
	}
	err = client.Connect()
	return client, err
}

// WsMarkPriceServeWithRate serve websocket that pushes price and funding rate for a single symbol and rate
func WsMarkPriceServeWithRate(symbol string, rate time.Duration, handler WsMarkPriceHandler, onConnect WsMarkPriceConnect) (client *WebSocketClient, err error) {
	var rateStr string
	switch rate {
	case 3 * time.Second:
		rateStr = ""
	case 1 * time.Second:
		rateStr = "@1s"
	default:
		return nil, errors.New("invalid rate")
	}
	endpoint := fmt.Sprintf("%s/%s@markPrice%s", getWsEndpoint(), strings.ToLower(symbol), rateStr)
	return wsMarkPriceServe(endpoint, handler, onConnect)
}

// WsAllMarkPriceEvent 定义 websocket markPriceUpdate 事件数组
type WsAllMarkPriceEvent []*WsMarkPriceEvent

// WsAllMarkPriceHandler 处理推送所有交易对的价格和资金费率
type WsAllMarkPriceHandler func(event WsAllMarkPriceEvent)

// wsAllMarkPriceServe 提供可推送所有交易对价格和资金费率的 websocket
func wsAllMarkPriceServe(endpoint string, handler WsAllMarkPriceHandler) (client *WebSocketClient, err error) {
	cfg := newWsConfig(endpoint)
	client = NewClient(cfg)
	client.OnMessage = func(message []byte) {
		var event WsAllMarkPriceEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			return
		}
		handler(event)
	}
	err = client.Connect()
	return client, err
}

// wsAllMarkPriceServeWithRate 提供可推送所有交易对价格和资金费率的websocket
func wsAllMarkPriceServeWithRate(rate time.Duration, handler WsAllMarkPriceHandler) (client *WebSocketClient, err error) {
	var rateStr string
	switch rate {
	case 3 * time.Second:
		rateStr = ""
	case 1 * time.Second:
		rateStr = "@1s"
	default:
		return nil, errors.New("invalid rate")
	}
	endpoint := fmt.Sprintf("%s/!markPrice@arr%s", getWsEndpoint(), rateStr)
	return wsAllMarkPriceServe(endpoint, handler)
}

// WsKline K线数据
type WsKline struct {
	StartTime            int64   `json:"t"`        // 这根K线的起始时间
	EndTime              int64   `json:"T"`        // 这根K线的结束时间
	Symbol               string  `json:"s"`        // 交易对
	Interval             string  `json:"i"`        // K线间隔
	FirstTradeID         int64   `json:"F"`        // 这根K线期间第一笔成交ID
	LastTradeID          int64   `json:"L"`        // 这根K线期间末一笔成交ID
	Open                 float64 `json:"o,string"` // 这根K线期间第一笔成交价
	Close                float64 `json:"c,string"` // 这根K线期间末一笔成交价
	High                 float64 `json:"h,string"` // 这根K线期间最高成交价
	Low                  float64 `json:"l,string"` // 这根K线期间最低成交价
	Volume               float64 `json:"v,string"` // 这根K线期间成交量
	TradeNum             int64   `json:"n"`        // 这根K线期间成交笔数
	IsFinal              bool    `json:"x"`        // 这根K线是否完结(是否已经开始下一根K线)
	QuoteVolume          float64 `json:"q,string"` // 这根K线期间成交额
	ActiveBuyVolume      float64 `json:"V,string"` // 主动买入的成交额
	ActiveBuyQuoteVolume float64 `json:"Q,string"` // 主动买入的成交量
}

// WsKlineEvent 定义websocket kline线事件
type WsKlineEvent struct {
	Event  string  `json:"e"` // 事件类型
	Time   int64   `json:"E"` // 事件时间
	Symbol string  `json:"s"` // 交易对
	Kline  WsKline `json:"k"` // k线数据
}

// WsKlineHandler 处理 websocket kline 事件
type WsKlineHandler func(event *WsKlineEvent)

// WsKlineServe 为 websocket kline 处理程序提供符号和间隔，如 15m、30s
func WsKlineServe(symbol string, interval string, handler WsKlineHandler) (client *WebSocketClient, err error) {
	endpoint := fmt.Sprintf("%s/%s@kline_%s", getWsEndpoint(), strings.ToLower(symbol), interval)
	cfg := newWsConfig(endpoint)
	client = NewClient(cfg)
	client.OnMessage = func(message []byte) {
		var event WsKlineEvent
		if err := json.Unmarshal(message, &event); err != nil {
			return
		}
		handler(&event)
	}
	err = client.Connect()
	return client, err
}

// WsContractInfoEvent 交易对信息信息流
type WsContractInfoEvent struct {
	Type           string `json:"e"`  // 事件类型
	Time           int64  `json:"E"`  // 事件时间
	Symbol         string `json:"s"`  // 交易对
	Pair           string `json:"ps"` // 交易对标的
	ContractType   string `json:"ct"` // 合约类型
	DeliveryTime   int64  `json:"dt"` // 结算时间
	OnboardTime    int64  `json:"ot"` // 上架时间
	ContractStatus string `json:"cs"` // 交易对状态
}

type WsContractInfoHandler func(event *WsContractInfoEvent)

func WsContractInfoServe(handler WsContractInfoHandler) (client *WebSocketClient, err error) {
	// !contractInfo || 交易对信息信息流,Symbol状态更改时推送（上架/下架/bracket调整）; bks仅在bracket调整时推出。
	endpoint := fmt.Sprintf("%s/!contractInfo", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	client = NewClient(cfg)
	client.OnMessage = func(message []byte) {
		var event WsContractInfoEvent
		if err := json.Unmarshal(message, &event); err != nil {
			return
		}
		handler(&event)
	}
	err = client.Connect()
	return client, err
}

// WsCombinedKlineServe 与 WsKlineServe 类似，但它处理多个交易对
func WsCombinedKlineServe(symbolIntervalPair map[string]string) (client *WebSocketClient, err error) {
	endpoint := getCombinedEndpoint()
	for symbol, interval := range symbolIntervalPair {
		endpoint += fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	client = NewClient(cfg)
	client.OnMessage = func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			return
		}

		stream := j.
	}
}
