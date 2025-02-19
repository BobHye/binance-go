package futures

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BobHye/wsc"
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
	// WebsocketTimeout 如果启用了 WebsocketKeepalive，则发送 ping/pong 消息的间隔 | 如果启用了 WebsocketKeepalive，则发送 ping/pong 消息的时间间隔
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability | 允许发送 ping/pong 消息以检查连接稳定性
	WebsocketKeepalive = true
	// UseTestnet switch all the WS streams from production to the testnet | 将所有 WS 流从生产环境切换到测试网络
	UseTestnet = false
)

// getWsEndpoint return the base endpoint of the WS according the UseTestnet flag | 根据 UseTestnet 标志返回 WS 的基本端点
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

// WsAggTradeEvent define websocket aggTrade event | 定义websocket @aggTrade 事件
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

// WsAggTradeHandler handle websocket that push trade information that is aggregated for a single taker order. | 处理 websocket，推送单个交易对的的聚合交易信息
type WsAggTradeHandler func(event *WsAggTradeEvent)

// WsAggTradeServe serve websocket that push trade information that is aggregated for a single taker order | 提供 websocket 服务，推送单个交易对的订单聚合交易信息。
func WsAggTradeServe(symbol string, handler WsAggTradeHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@aggTrade", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsAggTradeEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedAggTradeServe is similar to WsAggTradeServe, but it handles multiple symbols | 与 WsAggTradeServe 类似，但它处理多个交易对
func WsCombinedAggTradeServe(symbols []string, handler WsAggTradeHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@aggTrade", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}
		stream := j.Get("stream").MustString()
		data := j.Get("data").MustMap()

		symbols := strings.Split(stream, "@")[0]

		jsonData, _ := json.Marshal(data)

		event := new(WsAggTradeEvent)
		err = json.Unmarshal(jsonData, event)
		if err != nil {
			errHandler(err)
			return
		}
		event.Symbol = strings.ToUpper(symbols)
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
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

// wsMarkPriceServe
func wsMarkPriceServe(endpoint string, handler WsMarkPriceHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMarkPriceEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsMarkPriceServe serve websocket that pushes price and funding rate for a single symbol | 提供 websocket 服务，推送单个交易对的价格和资金费率信息。
func WsMarkPriceServe(symbol string, handler WsMarkPriceHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	// <symbol>@markPrice 或 <symbol>@markPrice@1s | 最新MarkPrice推送
	endpoint := fmt.Sprintf("%s/%s@markPrice", getWsEndpoint(), strings.ToLower(symbol))
	return wsMarkPriceServe(endpoint, handler, errHandler)
}

// WsMarkPriceServeWithRate serve websocket that pushes price and funding rate for a single symbol and rate
func WsMarkPriceServeWithRate(symbol string, rate time.Duration, handler WsMarkPriceHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	var rateStr string
	switch rate {
	case 3 * time.Second:
		rateStr = ""
	case 1 * time.Second:
		rateStr = "@1s"
	default:
		return nil, nil, errors.New("invalid rate")
	}
	endpoint := fmt.Sprintf("%s/%s@markPrice%s", getWsEndpoint(), strings.ToLower(symbol), rateStr)
	return wsMarkPriceServe(endpoint, handler, errHandler)
}

// WsAllMarkPriceEvent 定义 websocket markPriceUpdate 事件数组
type WsAllMarkPriceEvent []*WsMarkPriceEvent

// WsAllMarkPriceHandler handle websocket that pushes price and funding rate for all symbol | 处理推送所有交易对的价格和资金费率
type WsAllMarkPriceHandler func(event WsAllMarkPriceEvent)

// wsAllMarkPriceServe 提供可推送所有交易对价格和资金费率的 websocket
func wsAllMarkPriceServe(endpoint string, handler WsAllMarkPriceHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMarkPriceEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// wsAllMarkPriceServeWithRate 提供可推送所有交易对价格和资金费率的websocket
func wsAllMarkPriceServeWithRate(rate time.Duration, handler WsAllMarkPriceHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	var rateStr string
	switch rate {
	case 3 * time.Second:
		rateStr = ""
	case 1 * time.Second:
		rateStr = "@1s"
	default:
		return nil, nil, errors.New("invalid rate")
	}
	endpoint := fmt.Sprintf("%s/!markPrice@arr%s", getWsEndpoint(), rateStr)
	return wsAllMarkPriceServe(endpoint, handler, errHandler)
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
func WsKlineServe(symbol string, interval string, handler WsKlineHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@kline_%s", getWsEndpoint(), strings.ToLower(symbol), interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(data []byte) {
		var event WsKlineEvent
		if err := json.Unmarshal(data, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return wsServe(cfg, wsHandler, errHandler)
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

func WsContractInfoServe(handler WsContractInfoHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	// !contractInfo || 交易对信息信息流,Symbol状态更改时推送（上架/下架/bracket调整）; bks仅在bracket调整时推出。
	endpoint := fmt.Sprintf("%s/!contractInfo", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(data []byte) {
		var event WsContractInfoEvent
		if err := json.Unmarshal(data, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedKlineServe is similar to WsKlineServe, but it handles multiple symbols with it interval
func WsCombinedKlineServe(symbolIntervalPair map[string]string, handler WsKlineHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for symbol, interval := range symbolIntervalPair {
		endpoint += fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}

		stream := j.Get("stream").MustString()
		data := j.Get("data").MustMap()

		symbol := strings.Split(stream, "@")[0]

		jsonData, _ := json.Marshal(data)

		event := new(WsKlineEvent)
		err = json.Unmarshal(jsonData, event)
		if err != nil {
			errHandler(err)
			return
		}
		event.Symbol = strings.ToUpper(symbol)

		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsMiniMarketTickerEvent 定义 精简Ticker websocket事件
type WsMiniMarketTickerEvent struct {
	Event       string `json:"e"` // 事件类型
	Time        int64  `json:"E"` // 事件时间(毫秒)
	Symbol      string `json:"s"` // 交易对
	ClosePrice  string `json:"c"` // 最新成交价格
	OpenPrice   string `json:"o"` // 24小时前开始第一笔成交价格
	HighPrice   string `json:"h"` // 24小时内最高成交价
	LowPrice    string `json:"l"` // 24小时内最低成交价
	Volume      string `json:"v"` // 成交量
	QuoteVolume string `json:"q"` // 成交额
}

// WsMiniMarketTickerHandler 处理 websocket，推送单个交易对的精简ticker数据
type WsMiniMarketTickerHandler func(event *WsMiniMarketTickerEvent)

func WsMiniMarketTickerServe(symbol string, handler WsMiniMarketTickerHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@miniTicker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMiniMarketTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllMiniMarketTickerEvent 全市场的精简Ticker事件
type WsAllMiniMarketTickerEvent []*WsMiniMarketTickerEvent

// WsAllMiniMarketTickerHandler 处理全市场的精简Ticker
type WsAllMiniMarketTickerHandler func(event WsAllMiniMarketTickerEvent)

// WsAllMiniMarketTickerServe 全市场的精简Ticker
func WsAllMiniMarketTickerServe(handler WsAllMiniMarketTickerHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!miniTicker@arr", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMiniMarketTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsMarketTickerEvent 定义完整Ticker websocket事件
type WsMarketTickerEvent struct {
	Event              string `json:"e"` // 事件类型
	Time               int64  `json:"E"` // 事件时间
	Symbol             string `json:"s"` // 交易对
	PriceChange        string `json:"p"` // 24小时价格变化
	PriceChangePercent string `json:"P"` // 24小时价格变化(百分比)
	WeightedAvgPrice   string `json:"w"` // 平均价格
	ClosePrice         string `json:"c"` // 最新成交价格
	CloseQty           string `json:"Q"` // 最新成交价格上的成交量
	OpenPrice          string `json:"o"` // 24小时内第一比成交的价格
	HighPrice          string `json:"h"` // 24小时内最高成交价
	LowPrice           string `json:"l"` // 24小时内最低成交价
	BaseVolume         string `json:"v"` // 24小时内成交量
	QuoteVolume        string `json:"q"` // 24小时内成交额
	OpenTime           int64  `json:"O"` // 统计开始时间
	CloseTime          int64  `json:"C"` // 统计关闭时间
	FirstID            int64  `json:"F"` // 24小时内第一笔成交交易ID
	LastID             int64  `json:"L"` // 24小时内最后一笔成交交易ID
	TradeCount         int64  `json:"n"` // 24小时内成交数
}

// WsMarketTickerHandler 处理websocket推送的单个交易对 完整Ticker信息
type WsMarketTickerHandler func(event *WsMarketTickerEvent)

func WsMarketTickerServe(symbol string, handler WsMarketTickerHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@ticker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMarketTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllMarketTickerEvent 定义全市场完整Ticker websocket事件数组
type WsAllMarketTickerEvent []*WsMarketTickerEvent

type WsAllMarketTickerHandler func(event WsAllMarketTickerEvent)

func WsAllMarketTickerServe(handler WsAllMarketTickerHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!ticker@arr", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMarketTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsBookTickerEvent 定义 最优挂单信息websocket事件
type WsBookTickerEvent struct {
	Event           string  `json:"e"`        // 事件类型
	UpdateID        int64   `json:"u"`        // 更新ID
	Time            int64   `json:"E"`        // 事件推送时间
	TransactionTime int64   `json:"T"`        // 撮合时间
	Symbol          string  `json:"s"`        // 交易对
	BestBidPrice    float64 `json:"b,string"` // 买单最优挂单价格
	BestBidQty      float64 `json:"B,string"` // 买单最优挂单数量
	BestAskPrice    float64 `json:"a,string"` // 卖单最优挂单价格
	BestAskQty      float64 `json:"A,string"` // 卖单最优挂单数量
}

// WsBookTickerHandler handle websocket that pushes updates to the best bid or ask price or quantity in real-time for a specified symbol.
type WsBookTickerHandler func(event *WsBookTickerEvent)

// WsBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for a specified symbol.
func WsBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@bookTicker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for all symbols.
func WsAllBookTickerServe(handler WsBookTickerHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!bookTicker", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsLiquidationOrderEvent define websocket liquidation order event
type WsLiquidationOrderEvent struct {
	Event            string             `json:"e"`
	Time             int64              `json:"E"`
	LiquidationOrder WsLiquidationOrder `json:"o"`
}

// WsLiquidationOrder define websocket liquidation order
type WsLiquidationOrder struct {
	Symbol               string          `json:"s"`  // 交易对
	Side                 SideType        `json:"S"`  // 订单方向
	OrderType            OrderType       `json:"o"`  // 订单类型
	TimeInForce          TimeInForceType `json:"f"`  // 有效方式
	OrigQuantity         string          `json:"q"`  // 订单数量
	Price                string          `json:"p"`  // 订单价格
	AvgPrice             string          `json:"ap"` // 平均价格
	OrderStatus          OrderStatusType `json:"X"`  // 订单状态
	LastFilledQty        string          `json:"l"`  // 订单最近成交量
	AccumulatedFilledQty string          `json:"z"`  // 订单累计成交量
	TradeTime            int64           `json:"T"`  // 交易时间
}

// WsLiquidationOrderHandler handle websocket that pushes force liquidation order information for specific symbol.
type WsLiquidationOrderHandler func(event *WsLiquidationOrderEvent)

// WsLiquidationOrderServe serve websocket that pushes force liquidation order information for specific symbol.
func WsLiquidationOrderServe(symbol string, handler WsLiquidationOrderHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	// <symbol>@forceOrder || 推送特定symbol的强平订单快照信息。 1000ms内至多仅推送一条最近的强平订单作为快照
	endpoint := fmt.Sprintf("%s/%s@forceOrder", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(data []byte) {
		event := new(WsLiquidationOrderEvent)
		err := json.Unmarshal(data, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllLiquidationOrderServe serve websocket that pushes force liquidation order information for all symbols.
func WsAllLiquidationOrderServe(handler WsLiquidationOrderHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	// !forceOrder@arr || 推送全市场强平订单快照信息 每个symbol，1000ms内至多仅推送一条最近的强平订单作为快照
	endpoint := fmt.Sprintf("%s/!forceOrder@arr", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsLiquidationOrderEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsDepthEvent define websocket depth book event
type WsDepthEvent struct {
	Event            string `json:"e"`
	Time             int64  `json:"E"`
	TransactionTime  int64  `json:"T"`
	Symbol           string `json:"s"`
	FirstUpdateID    int64  `json:"U"`
	LastUpdateID     int64  `json:"u"`
	PrevLastUpdateID int64  `json:"pu"`
	Bids             []Bid  `json:"b"`
	Asks             []Ask  `json:"a"`
}

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(event *WsDepthEvent)

func wsPartialDepthServe(symbol string, levels int, rate *time.Duration, handler WsDepthHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	if levels != 5 && levels != 10 && levels != 20 {
		return nil, nil, errors.New("invalid levels")
	}
	levelsStr := fmt.Sprintf("%d", levels)
	return wsDepthServe(symbol, levelsStr, rate, handler, errHandler)
}

// WsPartialDepthServe serve websocket partial depth handler
func WsPartialDepthServe(symbol string, levels int, handler WsDepthHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	return wsPartialDepthServe(symbol, levels, nil, handler, errHandler)
}

// WsPartialDepthServeWithRate serve websocket partial depth handler with rate.
func WsPartialDepthServeWithRate(symbol string, levels int, rate time.Duration, handler WsDepthHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	return wsPartialDepthServe(symbol, levels, &rate, handler, errHandler)
}

// WsDiffDepthServe serve websocket diff. depth handler.
func WsDiffDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	return wsDepthServe(symbol, "", nil, handler, errHandler)
}

// WsCombinedDepthServe is similar to WsPartialDepthServe, but it for multiple symbols
func WsCombinedDepthServe(symbolLevels map[string]string, handler WsDepthHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for s, l := range symbolLevels {
		endpoint += fmt.Sprintf("%s@depth%s", strings.ToLower(s), l) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedDiffDepthServe is similar to WsDiffDepthServe, but it for multiple symbols
func WsCombinedDiffDepthServe(symbols []string, handler WsDepthHandler, errHandler ErrHandler) (wsc *wsc.Wsc, done chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@depth", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsDiffDepthServeWithRate serve websocket diff. depth handler with rate.
func WsDiffDepthServeWithRate(symbol string, rate time.Duration, handler WsDepthHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	return wsDepthServe(symbol, "", &rate, handler, errHandler)
}

func wsDepthServe(symbol string, levels string, rate *time.Duration, handler WsDepthHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	var rateStr string
	if rate != nil {
		switch *rate {
		case 250 * time.Millisecond:
			rateStr = ""
		case 500 * time.Millisecond:
			rateStr = "@500ms"
		case 100 * time.Millisecond:
			rateStr = "@100ms"
		default:
			return nil, nil, errors.New("invalid rate")
		}
	}
	endpoint := fmt.Sprintf("%s/%s@depth%s%s", getWsEndpoint(), strings.ToLower(symbol), levels, rateStr)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsDepthEvent)
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsBLVTInfoEvent define websocket BLVT info event 杠杆代币信息更新
type WsBLVTInfoEvent struct {
	Event          string         `json:"e"`
	Time           int64          `json:"E"`
	Symbol         string         `json:"s"`
	Issued         float64        `json:"m"`
	Baskets        []WsBLVTBasket `json:"b"`
	Nav            float64        `json:"n"`
	Leverage       float64        `json:"l"`
	TargetLeverage int64          `json:"t"`
	FundingRate    float64        `json:"f"`
}

// WsBLVTBasket define websocket BLVT basket
type WsBLVTBasket struct {
	Symbol   string `json:"s"`
	Position int64  `json:"n"`
}

// WsBLVTLogger handle websocket BLVT event
type WsBLVTLogger func(event *WsBLVTInfoEvent)

// WsBLVTInfoServe serve BLVT info stream
func WsBLVTInfoServe(name string, handler WsBLVTLogger, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@tokenNav", getWsEndpoint(), strings.ToLower(name))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBLVTInfoEvent)
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsBLVTKLineEvent define BLVT kline event || 净值 K 线更新事件
type WsBLVTKLineEvent struct {
	Event  string      `json:"e"`
	Time   int64       `json:"E"`
	Symbol string      `json:"s"`
	Kline  WsBLVTKLine `json:"k"`
}

// WsBLVTKLine BLVT Kline
type WsBLVTKLine struct {
	StartTime       int64  `json:"t"`
	CloseTime       int64  `json:"T"`
	Symbol          string `json:"s"`
	Interval        string `json:"i"`
	FirstUpdateTime int64  `json:"f"`
	LastUpdateTime  int64  `json:"L"`
	OpenPrice       string `json:"o"`
	ClosePrice      string `json:"c"`
	HighPrice       string `json:"h"`
	LowPrice        string `json:"l"`
	Leverage        string `json:"v"`
	Count           int64  `json:"n"`
}

// WsBLVTKLineHandler BLVT Kline handler
type WsBLVTKLineHandler func(event *WsBLVTKLineEvent)

// WsBLVTKLineServe serve BLVT Kline stream
func WsBLVTKLineServe(name string, interval string, handler WsBLVTKLineHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@nav_Kline_%s", getWsEndpoint(), strings.ToUpper(name), interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBLVTKLineEvent)
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCompositeIndexEvent websocket composite index event || 综合指数交易对信息流
type WsCompositeIndexEvent struct {
	Event       string          `json:"e"` // 事件类型
	Time        int64           `json:"E"` // 事件时间
	Symbol      string          `json:"s"` // 交易对
	Price       string          `json:"p"` // 价格
	Composition []WsComposition `json:"c"` // 成分信息
}

// WsComposition websocket composite index event composition
type WsComposition struct {
	BaseAsset    string `json:"b"` // 基础资产
	QuoteAsset   string `json:"q"` // 报价资产
	WeightQty    string `json:"w"` // 权重(数量)
	WeighPercent string `json:"W"` // 权重(比例)
	IndexPrice   string `json:"i"` // 指数价格
}

// WsCompositeIndexHandler websocket composite index handler
type WsCompositeIndexHandler func(event *WsCompositeIndexEvent)

// WsCompositiveIndexServe serve composite index information for index symbols
func WsCompositiveIndexServe(symbol string, handler WsCompositeIndexHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@compositeIndex", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsCompositeIndexEvent)
		if err := json.Unmarshal(message, &event); err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsUserDataEvent define user data event
type WsUserDataEvent struct {
	Event               UserDataEventType     `json:"e"`
	Time                int64                 `json:"E"`
	CrossWalletBalance  string                `json:"cw"`
	MarginCallPositions []WsPosition          `json:"p"`
	TransactionTime     int64                 `json:"T"`
	AccountUpdate       WsAccountUpdate       `json:"a"`
	OrderTradeUpdate    WsOrderTradeUpdate    `json:"o"`
	AccountConfigUpdate WsAccountConfigUpdate `json:"ac"`
}

// WsAccountUpdate define account update
type WsAccountUpdate struct {
	Reason    UserDataEventReasonType `json:"m"` // 事件推出原因
	Balances  []WsBalance             `json:"B"` // 余额信息
	Positions []WsPosition            `json:"P"`
}

type WsBalance struct {
	Asset              string  `json:"a"`         // 资产名称
	Balance            float64 `json:"wb,string"` // 钱包余额
	CrossWalletBalance float64 `json:"cw,string"` // 除去逐仓仓位保证金的钱包余额
	ChangeBalance      float64 `json:"bc,string"` // 除去盈亏与交易手续费以外的钱包余额改变量
}

// WsPosition define position
type WsPosition struct {
	Symbol                    string           `json:"s"`          // 交易对
	Side                      PositionSideType `json:"ps"`         // 持仓方向
	Amount                    float64          `json:"pa,string"`  // 仓位
	MarginType                MarginType       `json:"mt"`         // 保证金模式
	IsolatedWallet            float64          `json:"iw,string"`  // 若为逐仓，仓位保证金
	MarkPrice                 float64          `json:"mp,string"`  // 标记价格
	UnrealizedPnL             float64          `json:"up,string"`  // 持仓未实现盈亏
	MaintenanceMarginRequired string           `json:"mm"`         // 持仓需要的维持保证金
	EntryPrice                float64          `json:"ep,string"`  // 入仓价格
	AccumulatedRealized       float64          `json:"cr,string"`  // (费前)累计实现损益
	BreakEvenPrice            float64          `json:"bep,string"` // 盈亏平衡价
}

type WsOrderTradeUpdate struct {
	Symbol               string             `json:"s"`         // 交易对
	ClientOrderID        string             `json:"c"`         // 客户端自定订单ID
	Side                 SideType           `json:"S"`         // 订单方向
	Type                 OrderType          `json:"o"`         // 订单类型
	TimeInForce          TimeInForceType    `json:"f"`         // 有效方式
	OriginalQty          float64            `json:"q,string"`  // 订单原始数量
	OriginalPrice        float64            `json:"p,string"`  // 订单原始价格
	AveragePrice         float64            `json:"ap,string"` // 订单平均价格
	StopPrice            float64            `json:"sp,string"` // 条件订单触发价格，对追踪止损单无效
	ExecutionType        OrderExecutionType `json:"x"`         // 本次事件的具体执行类型
	Status               OrderStatusType    `json:"X"`         // 订单的当前状态
	ID                   int64              `json:"i"`         // 订单ID
	LastFilledQty        float64            `json:"l,string"`  // 订单末次成交量
	AccumulatedFilledQty float64            `json:"z,string"`  // 订单累计已成交量
	LastFilledPrice      float64            `json:"L,string"`  // 订单末次成交价格
	CommissionAsset      string             `json:"N"`         // 手续费资产类型
	Commission           float64            `json:"n,string"`  // 手续费数量
	TradeTime            int64              `json:"T"`         // 成交时间
	TradeID              int64              `json:"t"`         // 成交ID
	BidsNotional         float64            `json:"b,string"`  // 买单净值
	AsksNotional         float64            `json:"a,string"`  // 卖单净值
	IsMaker              bool               `json:"m"`         // 该成交是作为挂单成交吗？
	IsReduceOnly         bool               `json:"R"`         // 是否是只减仓单
	WorkingType          WorkingType        `json:"wt"`        // 触发价类型
	OriginalType         OrderType          `json:"ot"`        // 原始订单类型
	PositionSide         PositionSideType   `json:"ps"`        // 持仓方向
	IsClosingPosition    bool               `json:"cp"`        // 是否为触发平仓单; 仅在条件订单情况下会推送此字段
	ActivationPrice      float64            `json:"AP,string"` // 追踪止损激活价格, 仅在追踪止损单时会推送此字段
	CallbackRate         float64            `json:"cr,string"` // 追踪止损回调比例, 仅在追踪止损单时会推送此字段
	RealizedPnL          float64            `json:"rp,string"` // 该交易实现盈亏
}

// WsAccountConfigUpdate define account config update
type WsAccountConfigUpdate struct {
	Symbol   string `json:"s"`
	Leverage int    `json:"l"`
}

// WsUserDataHandler handle WsUserDataEvent
type WsUserDataHandler func(event *WsUserDataEvent)

func WsUserDataServe(listenKey string, handler WsUserDataHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		if bytes.Contains(message, []byte("\"e\":\"TRADE_LITE\"")) {
			return
		}
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
