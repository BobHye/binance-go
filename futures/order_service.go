package futures

import (
	"context"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// CreateOrderService create order
type CreateOrderService struct {
	c                *Client
	symbol           string            // 交易对
	side             SideType          // 买卖方向 SELL, BUY
	positionSide     *PositionSideType // 持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
	orderType        OrderType         // type 订单类型 LIMIT, MARKET, STOP, TAKE_PROFIT, STOP_MARKET, TAKE_PROFIT_MARKET, TRAILING_STOP_MARKET
	reduceOnly       *bool             // true, false; 非双开模式下默认false；双开模式下不接受此参数； 使用closePosition不支持此参数。
	quantity         string            // 下单数量,使用closePosition不支持此参数。
	price            *string           // 委托价格
	newClientOrderID *string           // 用户自定义的订单号，不可以重复出现在挂单中。如空缺系统会自动赋值。必须满足正则规则 ^[\.A-Z\:/a-z0-9_-]{1,36}$
	stopPrice        *string           // 触发价, 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
	closePosition    *bool             // true, false；触发后全部平仓，仅支持STOP_MARKET和TAKE_PROFIT_MARKET；不与quantity合用；自带只平仓效果，不与reduceOnly 合用
	activationPrice  *string           // 追踪止损激活价格，仅TRAILING_STOP_MARKET 需要此参数, 默认为下单当前市场价格(支持不同workingType)
	callbackRate     *string           // 追踪止损回调比例，可取值范围[0.1, 5],其中 1代表1% ,仅TRAILING_STOP_MARKET 需要此参数
	timeInForce      *TimeInForceType  // 有效方法
	workingType      *WorkingType      // stopPrice 触发类型: MARK_PRICE(标记价格), CONTRACT_PRICE(合约最新价). 默认 CONTRACT_PRICE
	priceProtect     *bool             // 条件单触发保护："TRUE","FALSE", 默认"FALSE". 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
	newOrderRespType NewOrderRespType  // "ACK", "RESULT", 默认 "ACK"
}

// SetSymbol set symbol
func (s *CreateOrderService) SetSymbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

// SetSide set side
func (s *CreateOrderService) SetSide(side SideType) *CreateOrderService {
	s.side = side
	return s
}

// SetPositionSide set side
func (s *CreateOrderService) SetPositionSide(positionSide PositionSideType) *CreateOrderService {
	s.positionSide = &positionSide
	return s
}

// SetType set type
func (s *CreateOrderService) SetType(orderType OrderType) *CreateOrderService {
	s.orderType = orderType
	return s
}

// SetTimeInForce set timeInForce
func (s *CreateOrderService) SetTimeInForce(timeInForce TimeInForceType) *CreateOrderService {
	s.timeInForce = &timeInForce
	return s
}

// SetQuantity set quantity
func (s *CreateOrderService) SetQuantity(quantity string) *CreateOrderService {
	s.quantity = quantity
	return s
}

// SetReduceOnly set reduceOnly
func (s *CreateOrderService) SetReduceOnly(reduceOnly bool) *CreateOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// SetPrice set price
func (s *CreateOrderService) SetPrice(price string) *CreateOrderService {
	s.price = &price
	return s
}

// SetNewClientOrderID set newClientOrderID
func (s *CreateOrderService) SetNewClientOrderID(newClientOrderID string) *CreateOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// SetStopPrice set stopPrice
func (s *CreateOrderService) SetStopPrice(stopPrice string) *CreateOrderService {
	s.stopPrice = &stopPrice
	return s
}

// SetWorkingType set workingType
func (s *CreateOrderService) SetWorkingType(workingType WorkingType) *CreateOrderService {
	s.workingType = &workingType
	return s
}

// SetActivationPrice set activationPrice
func (s *CreateOrderService) SetActivationPrice(activationPrice string) *CreateOrderService {
	s.activationPrice = &activationPrice
	return s
}

// SetCallbackRate set callbackRate
func (s *CreateOrderService) SetCallbackRate(callbackRate string) *CreateOrderService {
	s.callbackRate = &callbackRate
	return s
}

// SetPriceProtect set priceProtect
func (s *CreateOrderService) SetPriceProtect(priceProtect bool) *CreateOrderService {
	s.priceProtect = &priceProtect
	return s
}

// SetNewOrderResponseType set newOrderResponseType
func (s *CreateOrderService) SetNewOrderResponseType(newOrderResponseType NewOrderRespType) *CreateOrderService {
	s.newOrderRespType = newOrderResponseType
	return s
}

// SetClosePosition set closePosition
func (s *CreateOrderService) SetClosePosition(closePosition bool) *CreateOrderService {
	s.closePosition = &closePosition
	return s
}

// createOrder 新建订单
func (s *CreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, header *http.Header, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: endpoint,
		secType:  secTypeSigned,
	}
	m := params{
		"symbol":           s.symbol,
		"side":             s.side,
		"type":             s.orderType,
		"quantity":         s.quantity,
		"newOrderRespType": s.newOrderRespType,
	}
	if s.positionSide != nil {
		m["positionSide"] = *s.positionSide
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.reduceOnly != nil {
		m["reduceOnly"] = *s.reduceOnly
	}
	if s.price != nil {
		m["price"] = *s.price
	}
	if s.newClientOrderID != nil {
		m["newClientOrderId"] = *s.newClientOrderID
	}
	if s.stopPrice != nil {
		m["stopPrice"] = *s.stopPrice
	}
	if s.workingType != nil {
		m["workingType"] = *s.workingType
	}
	if s.priceProtect != nil {
		m["priceProtect"] = *s.priceProtect
	}
	if s.activationPrice != nil {
		m["activationPrice"] = *s.activationPrice
	}
	if s.callbackRate != nil {
		m["callbackRate"] = *s.callbackRate
	}
	if s.closePosition != nil {
		m["closePosition"] = *s.closePosition
	}
	r.setFormParams(m)
	data, header, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	return data, header, nil
}

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateOrderResponse, err error) {
	// POST /fapi/v1/order | 下单
	data, header, err := s.createOrder(ctx, "/fapi/v1/order", opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateOrderResponse)
	err = json.Unmarshal(data, res)
	res.RateLimitOrder10s = header.Get("X-Mbx-Order-Count-10s")
	res.RateLimitOrder1m = header.Get("X-Mbx-Order-Count-1m")
	if err != nil {
		return nil, err
	}
	return res, nil
}

type CreateOrderResponse struct {
	ClientOrderID           string           `json:"clientOrderId"` // 用户自定义的订单号
	CumQty                  float64          `json:"cumQty,string"`
	CumQuote                float64          `json:"cumQuote,string"`         // 成交金额
	ExecutedQuantity        float64          `json:"executedQty,string"`      // 成交量
	OrderID                 int64            `json:"orderId"`                 // 系统订单号
	AvgPrice                float64          `json:"avgPrice,string"`         // 平均成交价
	OrigQuantity            float64          `json:"origQty,string"`          // 原始委托数量
	Price                   float64          `json:"price,string"`            // 委托价格
	ReduceOnly              bool             `json:"reduceOnly"`              // 仅减仓
	Side                    SideType         `json:"side"`                    // 买卖方向
	PositionSide            PositionSideType `json:"positionSide"`            // 持仓方向
	Status                  OrderStatusType  `json:"status"`                  // 订单状态
	StopPrice               float64          `json:"stopPrice,string"`        // 触发价，对`TRAILING_STOP_MARKET`无效
	ClosePosition           bool             `json:"closePosition"`           // 是否条件全平仓
	Symbol                  string           `json:"symbol"`                  // 交易对
	TimeInForce             TimeInForceType  `json:"timeInForce"`             // 有效方法
	Type                    OrderType        `json:"type"`                    // 订单类型
	OrigType                string           `json:"origType"`                // 触发前订单类型
	ActivatePrice           float64          `json:"activatePrice,string"`    // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate               float64          `json:"priceRate,string"`        // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	UpdateTime              int64            `json:"updateTime"`              // 更新时间
	WorkingType             WorkingType      `json:"workingType"`             // 条件价格触发类型
	PriceProtect            bool             `json:"priceProtect"`            // 是否开启条件单触发保护
	PriceMatch              string           `json:"priceMatch"`              //盘口价格下单模式
	SelfTradePreventionMode string           `json:"selfTradePreventionMode"` //订单自成交保护模式
	GoodTillDate            int64            `json:"goodTillDate"`            //订单TIF为GTD时的自动取消时间
	RateLimitOrder10s       string           `json:"rateLimitOrder10s,omitempty"`
	RateLimitOrder1m        string           `json:"rateLimitOrder1m,omitempty"`
}

// ListOpenOrdersService 列出当前全部挂单
type ListOpenOrdersService struct {
	c      *Client
	symbol string // 交易对
}

// Symbol set symbol
func (s *ListOpenOrdersService) Symbol(symbol string) *ListOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *ListOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*Order, err error) {
	// GET /fapi/v1/openOrders | 查看当前全部挂单
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/openOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*Order{}, err
	}
	res = make([]*Order, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*Order{}, err
	}
	return res, nil
}

// GetOpenOrderService query current open order
type GetOpenOrderService struct {
	c                 *Client
	symbol            string  // 交易对
	orderID           *int64  // 系统订单号
	origClientOrderID *string // 用户自定义的订单号
}

// SetSymbol set symbol
func (s *GetOpenOrderService) SetSymbol(symbol string) *GetOpenOrderService {
	s.symbol = symbol
	return s
}

// SetOrderID set orderID
func (s *GetOpenOrderService) SetOrderID(orderID int64) *GetOpenOrderService {
	s.orderID = &orderID
	return s
}

// SetOrigClientOrderID  set origClientOrderID
func (s *GetOpenOrderService) SetOrigClientOrderID(origClientOrderID string) *GetOpenOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *GetOpenOrderService) Do(ctx context.Context, opts ...RequestOption) (res *Order, err error) {
	// GET /fapi/v1/openOrder | 查询当前挂单(请小心使用不带symbol参数的调用)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/openOrder",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.orderID == nil && s.origClientOrderID == nil {
		return nil, errors.New("either orderId or origClientOrderId must be sent")
	}
	if s.orderID != nil {
		r.setParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setParam("origClientOrderId", *s.origClientOrderID)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Order)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetOrderService get an order | 查询订单
type GetOrderService struct {
	c                 *Client
	symbol            string  // 交易对
	orderID           *int64  // 系统订单号
	origClientOrderID *string // 用户自定义的订单号
}

// SetSymbol set symbol
func (s *GetOrderService) SetSymbol(symbol string) *GetOrderService {
	s.symbol = symbol
	return s
}

// SetOrderID set orderID
func (s *GetOrderService) SetOrderID(orderID int64) *GetOrderService {
	s.orderID = &orderID
	return s
}

// SetOrigClientOrderID set origClientOrderID
func (s *GetOrderService) SetOrigClientOrderID(origClientOrderID string) *GetOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *GetOrderService) Do(ctx context.Context, opts ...RequestOption) (res *Order, err error) {
	// GET /fapi/v1/order | 查询订单
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/order",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setParam("origClientOrderId", *s.origClientOrderID)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Order)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Order 定义订单信息
type Order struct {
	AvgPrice                float64          `json:"avgPrice,string"`         // 平均成交价
	ClientOrderID           string           `json:"clientOrderId"`           // 用户自定义的订单号
	CumQuote                float64          `json:"cumQuote,string"`         // 成交金额
	ExecutedQuantity        float64          `json:"executedQty,string"`      // 成交量
	OrderID                 int64            `json:"orderId"`                 // 系统订单号
	OrigQuantity            float64          `json:"origQty,string"`          // 原始委托数量
	OrigType                string           `json:"origType"`                // 触发前订单类型
	Price                   float64          `json:"price,string"`            // 委托价格
	ReduceOnly              bool             `json:"reduceOnly"`              // 是否仅减仓
	Side                    SideType         `json:"side"`                    // 买卖方向
	PositionSide            PositionSideType `json:"positionSide"`            // 持仓方向
	Status                  OrderStatusType  `json:"status"`                  // 订单状态
	StopPrice               float64          `json:"stopPrice,string"`        // 触发价，对`TRAILING_STOP_MARKET`无效
	ClosePosition           bool             `json:"closePosition"`           // 是否条件全平仓
	Symbol                  string           `json:"symbol"`                  // 交易对
	Time                    int64            `json:"time"`                    // 订单时间
	TimeInForce             TimeInForceType  `json:"timeInForce"`             // 有效方法
	Type                    OrderType        `json:"type"`                    // 订单类型
	ActivatePrice           float64          `json:"activatePrice,string"`    // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate               float64          `json:"priceRate,string"`        // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	UpdateTime              int64            `json:"updateTime"`              // 更新时间
	WorkingType             WorkingType      `json:"workingType"`             // 条件价格触发类型
	PriceProtect            bool             `json:"priceProtect"`            // 是否开启条件单触发保护
	PriceMatch              string           `json:"priceMatch"`              //盘口价格下单模式
	SelfTradePreventionMode string           `json:"selfTradePreventionMode"` //订单自成交保护模式
	GoodTillDate            int64            `json:"goodTillDate"`            //订单TIF为GTD时的自动取消时间
}

// ListOrdersService all account orders; active, canceled, or filled | 查询所有订单(包括历史订单); active, canceled, or filled
type ListOrdersService struct {
	c         *Client
	symbol    string // 交易对
	orderID   *int64 // 只返回此orderID及之后的订单，缺省返回最近的订单
	startTime *int64 // 起始时间
	endTime   *int64 // 结束时间
	limit     *int   // 返回的结果集数量 默认值:500 最大值:1000
}

// SetSymbol set symbol
func (s *ListOrdersService) SetSymbol(symbol string) *ListOrdersService {
	s.symbol = symbol
	return s
}

// SetOrderID set orderID
func (s *ListOrdersService) SetOrderID(orderID int64) *ListOrdersService {
	s.orderID = &orderID
	return s
}

// SetStartTime set starttime
func (s *ListOrdersService) SetStartTime(startTime int64) *ListOrdersService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endtime
func (s *ListOrdersService) SetEndTime(endTime int64) *ListOrdersService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *ListOrdersService) SetLimit(limit int) *ListOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*Order, err error) {
	// GET /fapi/v1/allOrders | 查询所有订单(包括历史订单)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/allOrders",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setParam("orderId", *s.orderID)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// CancelOrderService 取消订单service
type CancelOrderService struct {
	c                 *Client
	symbol            string  // 交易对
	orderID           *int64  // 系统订单号
	origClientOrderID *string // 用户自定义的订单号
}

// SetSymbol set symbol
func (s *CancelOrderService) SetSymbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

// SetOrderID set orderID
func (s *CancelOrderService) SetOrderID(orderID int64) *CancelOrderService {
	s.orderID = &orderID
	return s
}

// SetOrigClientOrderID set origClientOrderID
func (s *CancelOrderService) SetOrigClientOrderID(origClientOrderID string) *CancelOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOrderResponse, err error) {
	// DELETE /fapi/v1/order | 撤销订单
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/fapi/v1/order",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setFormParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setFormParam("origClientOrderId", *s.origClientOrderID)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CancelOrderResponse define response of canceling order
type CancelOrderResponse struct {
	ClientOrderID           string           `json:"clientOrderId"` // 用户自定义的订单号
	CumQuantity             float64          `json:"cumQty,string"`
	CumQuote                float64          `json:"cumQuote,string"`         // 成交金额
	ExecutedQuantity        float64          `json:"executedQty,string"`      // 成交量
	OrderID                 int64            `json:"orderId"`                 // 系统订单号
	OrigQuantity            float64          `json:"origQty,string"`          // 原始委托数量
	Price                   float64          `json:"price,string"`            // 委托价格
	ReduceOnly              bool             `json:"reduceOnly"`              // 仅减仓
	Side                    SideType         `json:"side"`                    // 买卖方向
	PositionSide            PositionSideType `json:"positionSide"`            // 持仓方向
	Status                  OrderStatusType  `json:"status"`                  // 订单状态
	StopPrice               float64          `json:"stopPrice,string"`        // 触发价，对`TRAILING_STOP_MARKET`无效
	ClosePosition           bool             `json:"closePosition"`           // 是否条件全平仓
	Symbol                  string           `json:"symbol"`                  // 交易对
	TimeInForce             TimeInForceType  `json:"timeInForce"`             // 有效方法
	OrigType                string           `json:"origType"`                // 触发前订单类型
	Type                    OrderType        `json:"type"`                    // 订单类型
	ActivatePrice           float64          `json:"activatePrice,string"`    // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate               float64          `json:"priceRate,string"`        // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	UpdateTime              int64            `json:"updateTime"`              // 更新时间
	WorkingType             WorkingType      `json:"workingType"`             // 条件价格触发类型
	PriceProtect            bool             `json:"priceProtect"`            // 是否开启条件单触发保护
	PriceMatch              string           `json:"priceMatch"`              //盘口价格下单模式
	SelfTradePreventionMode string           `json:"selfTradePreventionMode"` //订单自成交保护模式
	GoodTillDate            int64            `json:"goodTillDate"`            //订单TIF为GTD时的自动取消时间
}

// CancelAllOpenOrdersService cancel all open orders | 撤销所有订单
type CancelAllOpenOrdersService struct {
	c      *Client
	symbol string // 交易对
}

// SetSymbol set symbol
func (s *CancelAllOpenOrdersService) SetSymbol(symbol string) *CancelAllOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CancelAllOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// DELETE /fapi/v1/allOpenOrders | 撤销全部订单 (TRADE)
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/fapi/v1/allOpenOrders",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// CancelMultiplesOrdersService cancel a list of orders | 批量撤销订单
type CancelMultiplesOrdersService struct {
	c                     *Client
	symbol                string   // 交易对
	orderIDList           []int64  // 系统订单号, 最多支持10个订单 比如[1234567,2345678]
	origClientOrderIDList []string // 用户自定义的订单号, 最多支持10个订单 比如["my_id_1","my_id_2"] 需要encode双引号。逗号后面没有空格。
}

// SetSymbol set symbol
func (s *CancelMultiplesOrdersService) SetSymbol(symbol string) *CancelMultiplesOrdersService {
	s.symbol = symbol
	return s
}

// SetOrderIDList set orderID
func (s *CancelMultiplesOrdersService) SetOrderIDList(orderIDList []int64) *CancelMultiplesOrdersService {
	s.orderIDList = orderIDList
	return s
}

// SetOrigClientOrderIDList set origClientOrderID
func (s *CancelMultiplesOrdersService) SetOrigClientOrderIDList(origClientOrderIDList []string) *CancelMultiplesOrdersService {
	s.origClientOrderIDList = origClientOrderIDList
	return s
}

// Do send request
func (s *CancelMultiplesOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*CancelOrderResponse, err error) {
	// DELETE /fapi/v1/batchOrders | 批量撤销订单 (TRADE)
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/fapi/v1/batchOrders",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	if s.orderIDList != nil {
		// convert a slice of integers to a string e.g. [1 2 3] => "[1,2,3]"
		orderIDListString := strings.Join(strings.Fields(fmt.Sprint(s.orderIDList)), ",")
		r.setFormParam("orderIdList", orderIDListString)
	}
	if s.origClientOrderIDList != nil {
		r.setFormParam("origClientOrderIdList", s.origClientOrderIDList)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// ListLiquidationOrdersService list liquidation orders
type ListLiquidationOrdersService struct {
	c         *Client
	symbol    *string
	startTime *int64
	endTime   *int64
	limit     *int
}

// SetSymbol set symbol
func (s *ListLiquidationOrdersService) SetSymbol(symbol string) *ListLiquidationOrdersService {
	s.symbol = &symbol
	return s
}

// SetStartTime set startTime
func (s *ListLiquidationOrdersService) SetStartTime(startTime int64) *ListLiquidationOrdersService {
	s.startTime = &startTime
	return s
}

// SetEndTime set startTime
func (s *ListLiquidationOrdersService) SetEndTime(endTime int64) *ListLiquidationOrdersService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *ListLiquidationOrdersService) SetLimit(limit int) *ListLiquidationOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListLiquidationOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*LiquidationOrder, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/allForceOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*LiquidationOrder{}, err
	}
	res = make([]*LiquidationOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*LiquidationOrder{}, err
	}
	return res, nil
}

// LiquidationOrder define liquidation order
type LiquidationOrder struct {
	Symbol           string          `json:"symbol"`
	Price            float64         `json:"price,string"`
	OrigQuantity     float64         `json:"origQty,string"`
	ExecutedQuantity float64         `json:"executedQty,string"`
	AveragePrice     float64         `json:"averagePrice,string"`
	Status           OrderStatusType `json:"status"`
	TimeInForce      TimeInForceType `json:"timeInForce"`
	Type             OrderType       `json:"type"`
	Side             SideType        `json:"side"`
	Time             int64           `json:"time"`
}

// ListUserLiquidationOrdersService lists user's liquidation orders
type ListUserLiquidationOrdersService struct {
	c             *Client
	symbol        *string
	autoCloseType ForceOrderCloseType // "LIQUIDATION": 强平单, "ADL": ADL 减仓单.
	startTime     *int64
	endTime       *int64
	limit         *int // Default 50; max 100
}

// SetSymbol set symbol
func (s *ListUserLiquidationOrdersService) SetSymbol(symbol string) *ListUserLiquidationOrdersService {
	s.symbol = &symbol
	return s
}

// SetAutoCloseType set autoCloseType
func (s *ListUserLiquidationOrdersService) SetAutoCloseType(autoCloseType ForceOrderCloseType) *ListUserLiquidationOrdersService {
	s.autoCloseType = autoCloseType
	return s
}

// SetStartTime set startTime
func (s *ListUserLiquidationOrdersService) SetStartTime(startTime int64) *ListUserLiquidationOrdersService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
func (s *ListUserLiquidationOrdersService) SetEndTime(endTime int64) *ListUserLiquidationOrdersService {
	s.endTime = &endTime
	return s
}

// SetLimit set limit
func (s *ListUserLiquidationOrdersService) SetLimit(limit int) *ListUserLiquidationOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListUserLiquidationOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*UserLiquidationOrder, err error) {
	// GET /fapi/v1/forceOrders | 查询用户强平单历史
	// 如果没有传 "autoCloseType", 强平单和 ADL 减仓单都会被返回
	//如果没有传"startTime", 只会返回"endTime"之前 7 天内的数据
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/forceOrders",
		secType:  secTypeSigned,
	}

	r.setParam("autoCloseType", s.autoCloseType)
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.symbol)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UserLiquidationOrder{}, err
	}
	res = make([]*UserLiquidationOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UserLiquidationOrder{}, err
	}
	return res, nil
}

// UserLiquidationOrder defines user's liquidation order
type UserLiquidationOrder struct {
	OrderId          int64            `json:"orderId"`
	Symbol           string           `json:"symbol"`
	Status           OrderStatusType  `json:"status"`
	ClientOrderId    string           `json:"clientOrderId"`
	Price            float64          `json:"price,string"`
	AveragePrice     float64          `json:"avgPrice,string"`
	OrigQuantity     float64          `json:"origQty,string"`
	ExecutedQuantity float64          `json:"executedQty,string"`
	CumQuote         float64          `json:"cumQuote,string"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	ReduceOnly       bool             `json:"reduceOnly"`
	ClosePosition    bool             `json:"closePosition"`
	Side             SideType         `json:"side"`
	PositionSide     PositionSideType `json:"positionSide"`
	StopPrice        float64          `json:"stopPrice,string"`
	WorkingType      WorkingType      `json:"workingType"`
	OrigType         float64          `json:"origType,string"`
	Time             int64            `json:"time"`
	UpdateTime       int64            `json:"updateTime"`
}

// CreateBatchOrdersService 批量下单service
type CreateBatchOrdersService struct {
	c      *Client
	orders []*CreateOrderService
}

type CreateBatchOrdersResponse struct {
	Orders []*Order
}

// SetOrderList set orderService
func (s *CreateBatchOrdersService) SetOrderList(orders []*CreateOrderService) *CreateBatchOrdersService {
	s.orders = orders
	return s
}

func (s *CreateBatchOrdersService) Do(ctx context.Context, opts ...RequestOption) (res *CreateBatchOrdersResponse, err error) {
	// POST /fapi/v1/batchOrders | 批量下单
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/batchOrders",
		secType:  secTypeSigned,
	}

	var orders []params
	for _, order := range s.orders {
		m := params{
			"symbol":           order.symbol,
			"side":             order.side,
			"type":             order.orderType,
			"quantity":         order.quantity,
			"newOrderRespType": order.newOrderRespType,
		}

		if order.positionSide != nil {
			m["positionSide"] = *order.positionSide
		}
		if order.timeInForce != nil {
			m["timeInForce"] = *order.timeInForce
		}
		if order.reduceOnly != nil {
			m["reduceOnly"] = *order.reduceOnly
		}
		if order.price != nil {
			m["price"] = *order.price
		}
		if order.newClientOrderID != nil {
			m["newClientOrderId"] = *order.newClientOrderID
		}
		if order.stopPrice != nil {
			m["stopPrice"] = *order.stopPrice
		}
		if order.workingType != nil {
			m["workingType"] = *order.workingType
		}
		if order.priceProtect != nil {
			m["priceProtect"] = *order.priceProtect
		}
		if order.activationPrice != nil {
			m["activationPrice"] = *order.activationPrice
		}
		if order.callbackRate != nil {
			m["callbackRate"] = *order.callbackRate
		}
		if order.closePosition != nil {
			m["closePosition"] = *order.closePosition
		}
		orders = append(orders, m)
	}
	b, err := json.Marshal(orders)
	if err != nil {
		return &CreateBatchOrdersResponse{}, err
	}
	m := params{
		"batchOrders": string(b), // 订单列表，最多支持5个订单
	}

	r.setFormParams(m)

	data, _, err := s.c.callAPI(ctx, r, opts...)

	if err != nil {
		return &CreateBatchOrdersResponse{}, err
	}

	rawMessages := make([]*stdjson.RawMessage, 0)

	err = json.Unmarshal(data, &rawMessages)

	if err != nil {
		return &CreateBatchOrdersResponse{}, err
	}

	batchCreateOrdersResponse := new(CreateBatchOrdersResponse)

	for _, j := range rawMessages {
		o := new(Order)
		if err := json.Unmarshal(*j, o); err != nil {
			return &CreateBatchOrdersResponse{}, err
		}

		if o.ClientOrderID != "" {
			batchCreateOrdersResponse.Orders = append(batchCreateOrdersResponse.Orders, o)
			continue
		}
	}
	return batchCreateOrdersResponse, nil
}
