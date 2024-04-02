package futures

import (
	"errors"
	"fmt"
	"github.com/BobHye/binance-go/api"
	"strings"
	"time"
)

// Endpoints
const (
	baseWsMainUrl       = "wss://fstream.binance.com/ws"
	baseCombinedMainURL = "wss://fstream.binance.com/stream?streams="

	baseWsTestnetUrl       = "wss://stream.binancefuture.com/ws"
	baseCombinedTestnetURL = "wss://stream.binancefuture.com/stream?streams="

	// UseTestnet 是否将所有 WS 流从生产网切换到测试网
	WsUseTestnet = false
)

var (
	// WebsocketTimeout 如果启用了 WebsocketKeepalive，则发送 ping/pong 消息的时间间隔
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive 允许发送 ping/pong 消息以检查连接稳定性
	WebsocketKeepalive = true
)

// getWsEndpoint 返回测试url或正式url
func getWsEndpoint() string {
	if WsUseTestnet {
		return baseWsTestnetUrl
	}
	return baseWsMainUrl
}

func getCombinedEndpoint() string {
	if WsUseTestnet {
		return baseCombinedTestnetURL
	}
	return baseCombinedMainURL
}

// WsAggTradeEvent 定义websocket @aggTrade 事件
type WsAggTradeEvent struct {
	Event            string  `json:"e"`
	Time             string  `json:"E"`
	Symbol           string  `json:"s"`
	AggregateTradeID string  `json:"a"`
	Price            string  `json:"p,string"`
	Quantity         float64 `json:"q,string"`
	FirstTradeID     int64   `json:"f"`
	LastTradeID      int64   `json:"l"`
	TradeTime        int64   `json:"T"`
	Maker            bool    `json:"m"`
}

// WsAggTradeHandler 处理 websocket，推送单个交易对的的聚合交易信息
type WsAggTradeHandler func(event *WsAggTradeEvent)

// WsAggTradeServe 提供 websocket 服务，推送单个交易对的订单聚合交易信息。
func WsAggTradeServe(symbol string, handler WsAggTradeHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@aggTrade", getWsEndpoint(), strings.ToLower(symbol))
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsAggTradeEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsCombinedAggTradeServe 与 WsAggTradeServe 类似，但它处理多个交易对
func WsCombinedAggTradeServe(symbols []string, handler WsAggTradeHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@aggTrade", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := api.NewWsConfig(endpoint)
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
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsMarkPriceEvent 定义 websocket markPriceUpdate 事件。
type WsMarkPriceEvent struct {
	Event                string `json:"e"`
	Time                 int64  `json:"E"`
	Symbol               string `json:"s"`
	MarkPrice            string `json:"p"`
	IndexPrice           string `json:"i"`
	EstimatedSettlePrice string `json:"P"`
	FundingRate          string `json:"r"`
	NextFundingTime      int64  `json:"T"`
}

// WsMarkPriceHandler 处理单个交易对价格和资金费率
type WsMarkPriceHandler func(event *WsMarkPriceEvent)

// wsMarkPriceServe
func wsMarkPriceServe(endpoint string, handler WsMarkPriceHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMarkPriceEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsMarkPriceServe 提供 websocket 服务，推送单个交易对的价格和资金费率信息。
func WsMarkPriceServe(symbol string, handler WsMarkPriceHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@markPrice", getWsEndpoint(), strings.ToLower(symbol))
	return wsMarkPriceServe(endpoint, handler, errHandler)
}

// WsMarkPriceServeWithRate
func WsMarkPriceServeWithRate(symbol string, rate time.Duration, handler WsMarkPriceHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
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
	return wsMarkPriceServe(endpoint, handler, errHandler)
}

// WsAllMarkPriceEvent 定义 websocket markPriceUpdate 事件数组
type WsAllMarkPriceEvent []*WsMarkPriceEvent

// WsAllMarkPriceHandler 处理推送所有交易对的价格和资金费率
type WsAllMarkPriceHandler func(event WsAllMarkPriceEvent)

// wsAllMarkPriceServe 提供可推送所有交易对价格和资金费率的 websocket
func wsAllMarkPriceServe(endpoint string, handler WsAllMarkPriceHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMarkPriceEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// wsAllMarkPriceServeWithRate 提供可推送所有交易对价格和资金费率的websocket
func wsAllMarkPriceServeWithRate(rate time.Duration, handler WsAllMarkPriceHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
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
	return wsAllMarkPriceServe(endpoint, handler, errHandler)
}

// WsKline K线数据
type WsKline struct {
	StartTime            int64   `json:"t"`
	EndTime              int64   `json:"T"`
	Symbol               string  `json:"s"`
	Interval             string  `json:"i"`
	FirstTradeID         int64   `json:"f"`
	LastTradeID          int64   `json:"L"`
	Open                 float64 `json:"o,string"`
	Close                float64 `json:"c,string"`
	High                 float64 `json:"h,string"`
	Low                  float64 `json:"l,string"`
	Volume               float64 `json:"v,string"`
	TradeNum             int64   `json:"n"`
	IsFinal              bool    `json:"x"`
	QuoteVolume          float64 `json:"q,string"`
	ActiveBuyVolume      float64 `json:"V,string"`
	ActiveBuyQuoteVolume float64 `json:"Q,string"`
}

// WsKlineEvent 定义websocket kline线事件
type WsKlineEvent struct {
	Event  string  `json:"e"`
	Time   int64   `json:"E"`
	Symbol string  `json:"s"`
	Kline  WsKline `json:"k"`
}

// WsKlineHandler 处理 websocket kline 事件
type WsKlineHandler func(event *WsKlineEvent)

// WsKlineServe 为 websocket kline 处理程序提供符号和间隔，如 15m、30s
func WsKlineServe(symbol string, interval string, handler WsKlineHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@kline_%s", getWsEndpoint(), strings.ToLower(symbol), interval)
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(data []byte) {
		var event WsKlineEvent
		if err := json.Unmarshal(data, &event); err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}
func WsCombinedKlineServe(symbolIntervalPair map[string]string, handler WsKlineHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for symbol, interval := range symbolIntervalPair {
		endpoint += fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := api.NewWsConfig(endpoint)
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
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsMiniMarketTickerEvent 定义 精简Ticker websocket事件
type WsMiniMarketTickerEvent struct {
	Event       string `json:"e"`
	Time        int64  `json:"E"`
	Symbol      string `json:"s"`
	ClosePrice  string `json:"c"`
	OpenPrice   string `json:"o"`
	HighPrice   string `json:"h"`
	LowPrice    string `json:"l"`
	Volume      string `json:"v"`
	QuoteVolume string `json:"q"`
}

// WsMiniMarketTickerHandler 处理 websocket，推送单个交易对的精简ticker数据
type WsMiniMarketTickerHandler func(event *WsMiniMarketTickerEvent)

func WsMiniMarketTickerServe(symbol string, handler WsMiniMarketTickerHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@miniTicker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMiniMarketTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsAllMiniMarketTickerEvent 全市场的精简Ticker事件
type WsAllMiniMarketTickerEvent []*WsMiniMarketTickerEvent

// WsAllMiniMarketTickerHandler 处理全市场的精简Ticker
type WsAllMiniMarketTickerHandler func(event WsAllMiniMarketTickerEvent)

// WsAllMiniMarketTickerServe 全市场的精简Ticker
func WsAllMiniMarketTickerServe(handler WsAllMiniMarketTickerHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!miniTicker@arr", getWsEndpoint())
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMiniMarketTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsMarketTickerEvent 定义完整Ticker websocket事件
type WsMarketTickerEvent struct {
	Event              string `json:"e"`
	Time               int64  `json:"E"`
	Symbol             string `json:"s"`
	PriceChange        string `json:"p"`
	PriceChangePercent string `json:"P"`
	WeightedAvgPrice   string `json:"w"`
	ClosePrice         string `json:"c"`
	CloseQty           string `json:"Q"`
	OpenPrice          string `json:"o"`
	HighPrice          string `json:"h"`
	LowPrice           string `json:"l"`
	BaseVolume         string `json:"v"`
	QuoteVolume        string `json:"q"`
	OpenTime           int64  `json:"O"`
	CloseTime          int64  `json:"C"`
	FirstID            int64  `json:"F"`
	LastID             int64  `json:"L"`
	TradeCount         int64  `json:"n"`
}

// WsMarketTickerHandler 处理websocket推送的单个交易对 完整Ticker信息
type WsMarketTickerHandler func(event *WsMarketTickerEvent)

func WsMarketTickerServe(symbol string, handler WsMarketTickerHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@ticker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMarketTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsAllMarketTickerEvent 定义全市场完整Ticker websocket事件数组
type WsAllMarketTickerEvent []*WsMarketTickerEvent

type WsAllMarketTickerHandler func(event WsAllMarketTickerEvent)

func WsAllMarketTickerServe(handler WsAllMarketTickerHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!ticker@arr", getWsEndpoint())
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMarketTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsBookTickerEvent 定义 最优挂单信息websocket事件
type WsBookTickerEvent struct {
	Event           string  `json:"e"`
	UpdateID        int64   `json:"u"`
	Time            int64   `json:"E"`
	TransactionTime int64   `json:"T"`
	Symbol          string  `json:"s"`
	BestBidPrice    float64 `json:"b,string"`
	BestBidQty      float64 `json:"B,string"`
	BestAskPrice    float64 `json:"a,string"`
	BestAskQty      float64 `json:"A,string"`
}
type WsBookTickerHandler func(event *WsBookTickerEvent)

func WsBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@bookTicker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

type WsAllMarketBookTickerEvent []*WsBookTickerEvent

type WsAllMarketBookTickerHandler func(event WsAllMarketBookTickerEvent)

// WsAllBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for all symbols.
func WsAllBookTickerServe(handler WsAllMarketBookTickerHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!bookTicker", getWsEndpoint())
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMarketBookTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}

// WsPosition define position
type WsPosition struct {
	Symbol                    string           `json:"s"`
	Side                      PositionSideType `json:"ps"`
	Amount                    float64          `json:"pa,string"`
	MarginType                MarginType       `json:"mt"`
	IsolatedWallet            float64          `json:"iw,string"`
	EntryPrice                float64          `json:"ep,string"`
	MarkPrice                 float64          `json:"mp,string"`
	UnrealizedPnL             float64          `json:"up,string"`
	AccumulatedRealized       float64          `json:"cr,string"`
	MaintenanceMarginRequired string           `json:"mm"`
}

type WsBalance struct {
	Asset              string  `json:"a"`
	Balance            float64 `json:"wb,string"`
	CrossWalletBalance float64 `json:"cw,string"`
	ChangeBalance      float64 `json:"bc,string"`
}

type WsAccountUpdate struct {
	Reason    UserDataEventReasonType `json:"m"`
	Balances  []WsBalance             `json:"B"`
	Positions []WsPosition            `json:"P"`
}

type WsOrderTradeUpdate struct {
	Symbol               string             `json:"s"`
	ClientOrderID        string             `json:"c"`
	Side                 SideType           `json:"S"`
	Type                 OrderType          `json:"o"`
	TimeInForce          TimeInForceType    `json:"f"`
	OriginalQty          float64            `json:"q,string"`
	OriginalPrice        float64            `json:"p,string"`
	AveragePrice         float64            `json:"ap,string"`
	StopPrice            float64            `json:"sp,string"`
	ExecutionType        OrderExecutionType `json:"x"`
	Status               OrderStatusType    `json:"X"`
	ID                   int64              `json:"i"`
	LastFilledQty        float64            `json:"l,string"`
	AccumulatedFilledQty float64            `json:"z,string"`
	LastFilledPrice      float64            `json:"L,string"`
	CommissionAsset      string             `json:"N"`
	Commission           float64            `json:"n,string"`
	TradeTime            int64              `json:"T"`
	TradeID              int64              `json:"t"`
	BidsNotional         float64            `json:"b,string"`
	AsksNotional         float64            `json:"a,string"`
	IsMaker              bool               `json:"m"`
	IsReduceOnly         bool               `json:"R"`
	WorkingType          WorkingType        `json:"wt"`
	OriginalType         OrderType          `json:"ot"`
	PositionSide         PositionSideType   `json:"ps"`
	IsClosingPosition    bool               `json:"cp"`
	ActivationPrice      float64            `json:"AP,string"`
	CallbackRate         float64            `json:"cr,string"`
	RealizedPnL          float64            `json:"rp,string"`
}

type WsAccountConfigUpdate struct {
	Symbol   string `json:"s"`
	Leverage int64  `json:"l"`
}

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

// WsUserDataHandler
type WsUserDataHandler func(event *WsUserDataEvent)

func WsUserDataServe(listenKey string, handler WsUserDataHandler, errHandler api.ErrHandler) (done chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := api.NewWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return api.WsServe(cfg, wsHandler, errHandler)
}
