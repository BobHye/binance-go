package spot

import (
	"context"
	stdjson "encoding/json"
	"net/http"
)

// CreateOrderService create order
type CreateOrderService struct {
	c                *Client
	symbol           string
	side             SideType
	orderType        OrderType
	timeInForce      *TimeInForceType
	newOrderRespType *NewOrderRespType // 指定响应类型 ACK, RESULT, or FULL; MARKET 与 LIMIT 订单默认为FULL, 其他默认为ACK。
	quantity         *string
	quoteOrderQty    *string
	price            *string
	newClientOrderID *string // 用户自定义的orderid，如空缺系统会自动赋值。
	stopPrice        *string // 仅 STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, TAKE_PROFIT_LIMIT 需要此参数。
	trailingDelta    *string // 用于 STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, 和 TAKE_PROFIT_LIMIT 类型的订单。
	icebergQuantity  *string
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
	s.quantity = &quantity
	return s
}

// SetQuoteOrderQty set quoteOrderQty
func (s *CreateOrderService) SetQuoteOrderQty(quoteOrderQty string) *CreateOrderService {
	s.quoteOrderQty = &quoteOrderQty
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

// SetTrailingDelta set trailingDelta
func (s *CreateOrderService) SetTrailingDelta(trailingDelta string) *CreateOrderService {
	s.trailingDelta = &trailingDelta
	return s
}

// SetIcebergQuantity set icebergQuantity
func (s *CreateOrderService) SetIcebergQuantity(icebergQuantity string) *CreateOrderService {
	s.icebergQuantity = &icebergQuantity
	return s
}

// SetNewOrderRespType set icebergQuantity
func (s *CreateOrderService) SetNewOrderRespType(newOrderRespType NewOrderRespType) *CreateOrderService {
	s.newOrderRespType = &newOrderRespType
	return s
}

func (s *CreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: endpoint,
		secType:  secTypeSigned,
	}
	m := params{
		"symbol": s.symbol,
		"side":   s.side,
		"type":   s.orderType,
	}
	if s.quantity != nil {
		m["quantity"] = *s.quantity
	}
	if s.quoteOrderQty != nil {
		m["quoteOrderQty"] = *s.quoteOrderQty
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
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
	if s.trailingDelta != nil {
		m["trailingDelta"] = *s.trailingDelta
	}
	if s.icebergQuantity != nil {
		m["icebergQty"] = *s.icebergQuantity
	}
	if s.newOrderRespType != nil {
		m["newOrderRespType"] = *s.newOrderRespType
	}
	r.setFormParams(m)
	data, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateOrderResponse, err error) {
	// POST /api/v3/order | 下单 (TRADE)
	data, err := s.createOrder(ctx, "/api/v3/order", opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Test send test api to check if the request is valid
func (s *CreateOrderService) Test(ctx context.Context, opts ...RequestOption) (err error) {
	_, err = s.createOrder(ctx, "/api/v3/order/test", opts...)
	return err
}

// CreateOrderResponse define create order response
type CreateOrderResponse struct {
	Symbol                   string `json:"symbol"`        // 交易对
	OrderID                  int64  `json:"orderId"`       // 系统的订单ID
	ClientOrderID            string `json:"clientOrderId"` // 客户自己设置的ID
	TransactTime             int64  `json:"transactTime"`
	Price                    string `json:"price"`
	OrigQuantity             string `json:"origQty"`
	ExecutedQuantity         string `json:"executedQty"`
	CummulativeQuoteQuantity string `json:"cummulativeQuoteQty"`
	IsIsolated               bool   `json:"isIsolated"` // for isolated margin

	Status      OrderStatusType `json:"status"`
	TimeInForce TimeInForceType `json:"timeInForce"`
	Type        OrderType       `json:"type"`
	Side        SideType        `json:"side"`

	// for order response is set to FULL
	Fills                 []*Fill `json:"fills"`
	MarginBuyBorrowAmount string  `json:"marginBuyBorrowAmount"` // for margin
	MarginBuyBorrowAsset  string  `json:"marginBuyBorrowAsset"`
}

// Fill may be returned in an array of fills in a CreateOrderResponse.
type Fill struct {
	TradeID         int    `json:"tradeId"`
	Price           string `json:"price"`
	Quantity        string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}

// CreateOCOService create order
type CreateOCOService struct {
	c      *Client
	symbol string
	// 整个 OCO order list 的唯一ID。 如果未发送则自动生成。 仅当前一个订单已填满或完全过期时，才会接受具有相同的listClientOrderId。
	// listClientOrderId 与 aboveClientOrderId 和 belowCLientOrderId 不同。
	listClientOrderID    *string
	side                 SideType // 订单方向：BUY or SELL
	quantity             *string  // 两个订单的数量
	aboveType            *string  // 支持值：STOP_LOSS_LIMIT, STOP_LOSS, LIMIT_MAKER, TAKE_PROFIT, TAKE_PROFIT_LIMIT
	limitClientOrderID   *string
	price                *string
	limitIcebergQty      *string
	stopClientOrderID    *string
	stopPrice            *string
	stopLimitPrice       *string
	stopIcebergQty       *string
	stopLimitTimeInForce *TimeInForceType
	newOrderRespType     *NewOrderRespType
}

// SetSymbol set symbol
func (s *CreateOCOService) SetSymbol(symbol string) *CreateOCOService {
	s.symbol = symbol
	return s
}

// SetSide set side
func (s *CreateOCOService) SetSide(side SideType) *CreateOCOService {
	s.side = side
	return s
}

// SetQuantity set quantity
func (s *CreateOCOService) SetQuantity(quantity string) *CreateOCOService {
	s.quantity = &quantity
	return s
}

// SetListClientOrderID set listClientOrderID
func (s *CreateOCOService) SetListClientOrderID(listClientOrderID string) *CreateOCOService {
	s.listClientOrderID = &listClientOrderID
	return s
}

// SetLimitClientOrderID set limitClientOrderID
func (s *CreateOCOService) SetLimitClientOrderID(limitClientOrderID string) *CreateOCOService {
	s.limitClientOrderID = &limitClientOrderID
	return s
}

// SetPrice set price
func (s *CreateOCOService) SetPrice(price string) *CreateOCOService {
	s.price = &price
	return s
}

// SetLimitIcebergQuantity set limitIcebergQuantity
func (s *CreateOCOService) SetLimitIcebergQuantity(limitIcebergQty string) *CreateOCOService {
	s.limitIcebergQty = &limitIcebergQty
	return s
}

// SetStopClientOrderID set stopClientOrderID
func (s *CreateOCOService) SetStopClientOrderID(stopClientOrderID string) *CreateOCOService {
	s.stopClientOrderID = &stopClientOrderID
	return s
}

// SetStopPrice set stop price
func (s *CreateOCOService) SetStopPrice(stopPrice string) *CreateOCOService {
	s.stopPrice = &stopPrice
	return s
}

// SetStopLimitPrice set stop limit price
func (s *CreateOCOService) SetStopLimitPrice(stopLimitPrice string) *CreateOCOService {
	s.stopLimitPrice = &stopLimitPrice
	return s
}

// SetStopIcebergQty set stop limit price
func (s *CreateOCOService) SetStopIcebergQty(stopIcebergQty string) *CreateOCOService {
	s.stopIcebergQty = &stopIcebergQty
	return s
}

// SetStopLimitTimeInForce set stopLimitTimeInForce
func (s *CreateOCOService) SetStopLimitTimeInForce(stopLimitTimeInForce TimeInForceType) *CreateOCOService {
	s.stopLimitTimeInForce = &stopLimitTimeInForce
	return s
}

// SetNewOrderRespType set icebergQuantity
func (s *CreateOCOService) SetNewOrderRespType(newOrderRespType NewOrderRespType) *CreateOCOService {
	s.newOrderRespType = &newOrderRespType
	return s
}

func (s *CreateOCOService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: endpoint,
		secType:  secTypeSigned,
	}
	m := params{
		"symbol":    s.symbol,
		"side":      s.side,
		"quantity":  *s.quantity,
		"price":     *s.price,
		"stopPrice": *s.stopPrice,
	}

	if s.listClientOrderID != nil {
		m["listClientOrderId"] = *s.listClientOrderID
	}
	if s.limitClientOrderID != nil {
		m["limitClientOrderId"] = *s.limitClientOrderID
	}
	if s.limitIcebergQty != nil {
		m["limitIcebergQty"] = *s.limitIcebergQty
	}
	if s.stopClientOrderID != nil {
		m["stopClientOrderId"] = *s.stopClientOrderID
	}
	if s.stopLimitPrice != nil {
		m["stopLimitPrice"] = *s.stopLimitPrice
	}
	if s.stopIcebergQty != nil {
		m["stopIcebergQty"] = *s.stopIcebergQty
	}
	if s.stopLimitTimeInForce != nil {
		m["stopLimitTimeInForce"] = *s.stopLimitTimeInForce
	}
	if s.newOrderRespType != nil {
		m["newOrderRespType"] = *s.newOrderRespType
	}
	r.setFormParams(m)
	data, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// Do send request
func (s *CreateOCOService) Do(ctx context.Context, opts ...RequestOption) (res *CreateOCOResponse, err error) {
	// POST /api/v3/orderList/oco | New Order list - OCO (TRADE)
	// 发送新 one-cancels-the-other (OCO) 订单，激活其中一个订单会立即取消另一个订单。
	// OCO 包含了两个订单，分别被称为 上方订单 和 下方订单。
	// 其中一个订单必须是 LIMIT_MAKER/TAKE_PROFIT/TAKE_PROFIT_LIMIT 订单，另一个订单必须是 STOP_LOSS 或 STOP_LOSS_LIMIT 订单。
	// 针对价格限制：
	// 如果 OCO 订单方向是 SELL：
	// LIMIT_MAKER/TAKE_PROFIT_LIMIT price > 最后交易价格 > STOP_LOSS/STOP_LOSS_LIMIT stopPrice
	// TAKE_PROFIT stopPrice > 最后交易价格 > STOP_LOSS/STOP_LOSS_LIMIT stopPrice
	// 如果 OCO 订单方向是 BUY：
	// LIMIT_MAKER/TAKE_PROFIT_LIMIT price < 最后交易价格 < stopPrice
	// TAKE_PROFIT stopPrice < 最后交易价格 < STOP_LOSS/STOP_LOSS_LIMIT stopPrice
	data, err := s.createOrder(ctx, "/api/v3/order/oco", opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateOCOResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateOCOResponse define create order response
type CreateOCOResponse struct {
	OrderListID       int64             `json:"orderListId"`
	ContingencyType   string            `json:"contingencyType"`
	ListStatusType    string            `json:"listStatusType"`
	ListOrderStatus   string            `json:"listOrderStatus"`
	ListClientOrderID string            `json:"listClientOrderId"`
	TransactionTime   int64             `json:"transactionTime"`
	Symbol            string            `json:"symbol"`
	Orders            []*OCOOrder       `json:"orders"`
	OrderReports      []*OCOOrderReport `json:"orderReports"`
}

// OCOOrder may be returned in an array of OCOOrder in a CreateOCOResponse.
type OCOOrder struct {
	Symbol        string `json:"symbol"`
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
}

// OCOOrderReport may be returned in an array of OCOOrderReport in a CreateOCOResponse.
type OCOOrderReport struct {
	Symbol                   string          `json:"symbol"`
	OrderID                  int64           `json:"orderId"`
	OrderListID              int64           `json:"orderListId"`
	ClientOrderID            string          `json:"clientOrderId"`
	OrigClientOrderID        string          `json:"origClientOrderId"`
	TransactionTime          int64           `json:"transactionTime"`
	Price                    string          `json:"price"`
	OrigQuantity             string          `json:"origQty"`
	ExecutedQuantity         string          `json:"executedQty"`
	CummulativeQuoteQuantity string          `json:"cummulativeQuoteQty"`
	Status                   OrderStatusType `json:"status"`
	TimeInForce              TimeInForceType `json:"timeInForce"`
	Type                     OrderType       `json:"type"`
	Side                     SideType        `json:"side"`
	StopPrice                string          `json:"stopPrice"`
	IcebergQuantity          string          `json:"icebergQty"`
}

// ListOpenOcoService list opened oco
type ListOpenOcoService struct {
	c *Client
}

// oco define oco info
type Oco struct {
	Symbol            string   `json:"symbol"`
	OrderListId       int64    `json:"orderListId"`
	ContingencyType   string   `json:"contingencyType"`
	ListStatusType    string   `json:"listStatusType"`
	ListOrderStatus   string   `json:"listOrderStatus"`
	ListClientOrderID string   `json:"listClientOrderId"`
	TransactionTime   int64    `json:"transactionTime"`
	Orders            []*Order `json:"orders"`
}

// Do send request
func (s *ListOpenOcoService) Do(ctx context.Context, opts ...RequestOption) (res []*Oco, err error) {
	// GET /api/v3/openOrderList | 查询订单列表挂单 (USER_DATA)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/openOrderList",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*Oco{}, err
	}
	res = make([]*Oco, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*Oco{}, err
	}
	return res, nil
}

// ListOpenOrdersService list opened orders
type ListOpenOrdersService struct {
	c      *Client
	symbol string
}

// SetSymbol set symbol
func (s *ListOpenOrdersService) SetSymbol(symbol string) *ListOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *ListOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*Order, err error) {
	// GET /api/v3/openOrders | 查看账户当前挂单 (USER_DATA)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/openOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
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

// GetOrderService get an order
type GetOrderService struct {
	c                 *Client
	symbol            string
	orderID           *int64
	origClientOrderID *string
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
	// GET /api/v3/order | 查询订单 (USER_DATA)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/order",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setParam("origClientOrderId", *s.origClientOrderID)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
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

// Order define order info
type Order struct {
	Symbol                   string          `json:"symbol"`              // 交易对
	OrderID                  int64           `json:"orderId"`             // 系统的订单ID
	OrderListId              int64           `json:"orderListId"`         // 除非此单是订单列表的一部分, 否则此值为 -1
	ClientOrderID            string          `json:"clientOrderId"`       // 客户自己设置的ID
	Price                    string          `json:"price"`               // 订单价格
	OrigQuantity             string          `json:"origQty"`             // 用户设置的原始订单数量
	ExecutedQuantity         string          `json:"executedQty"`         // 用户设置的原始订单数量
	CummulativeQuoteQuantity string          `json:"cummulativeQuoteQty"` // 累计交易的金额
	Status                   OrderStatusType `json:"status"`              // 订单状态
	TimeInForce              TimeInForceType `json:"timeInForce"`         // 订单的时效方式
	Type                     OrderType       `json:"type"`                // 订单类型， 比如市价单，现价单等
	Side                     SideType        `json:"side"`                // 订单方向，买还是卖
	StopPrice                string          `json:"stopPrice"`           // 止损价格
	IcebergQuantity          string          `json:"icebergQty"`          // 冰山数量
	Time                     int64           `json:"time"`                // 订单时间
	UpdateTime               int64           `json:"updateTime"`          // 最后更新时间
	WorkingTime              int64           `json:"workingTime"`         // 订单添加到 order book 的时间
	IsWorking                bool            `json:"isWorking"`           // 订单是否出现在orderbook中
	//IsIsolated               bool            `json:"isIsolated"`
	OrigQuoteOrderQuantity  string `json:"origQuoteOrderQty"`       // 原始的交易金额
	SelfTradePreventionMode string `json:"selfTradePreventionMode"` // 如何处理自我交易模式
}

// ListOrdersService all account orders; active, canceled, or filled
type ListOrdersService struct {
	c         *Client
	symbol    string
	orderID   *int64 // 只返回此orderID之后的订单，缺省返回最近的订单
	startTime *int64
	endTime   *int64
	limit     *int // Default 500; max 1000.
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
	// GET /api/v3/allOrders | 查询所有订单（包括历史订单） (USER_DATA)
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/allOrders",
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
	data, err := s.c.callAPI(ctx, r, opts...)
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

// CancelOrderService cancel an order
type CancelOrderService struct {
	c                 *Client
	symbol            string
	orderID           *int64
	origClientOrderID *string
	newClientOrderID  *string // 用户自定义的本次撤销操作的ID(注意不是被撤销的订单的自定义ID)。如无指定会自动赋值。
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

// SetNewClientOrderID set newClientOrderID
func (s *CancelOrderService) SetNewClientOrderID(newClientOrderID string) *CancelOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// Do send request
func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOrderResponse, err error) {
	// DELETE /api/v3/order | 撤销订单 (TRADE)
	// orderId 与 origClientOrderId 必须至少发送一个.
	// 如果两个参数一起发送, orderId优先被考虑.
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/api/v3/order",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setFormParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setFormParam("origClientOrderId", *s.origClientOrderID)
	}
	if s.newClientOrderID != nil {
		r.setFormParam("newClientOrderId", *s.newClientOrderID)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
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

// CancelOCOService cancel all active orders on the list order.
type CancelOCOService struct {
	c                 *Client
	symbol            string
	listClientOrderID string // orderListId 或 listClientOrderId 必须被提供
	orderListID       int64  // orderListId 或 listClientOrderId 必须被提供
	newClientOrderID  string // 用户自定义的本次撤销操作的ID(注意不是被撤销的订单的自定义ID)。如无指定会自动赋值
}

// SetSymbol set symbol
func (s *CancelOCOService) SeSymbol(symbol string) *CancelOCOService {
	s.symbol = symbol
	return s
}

// SeListClientOrderID sets listClientOrderId
func (s *CancelOCOService) SeListClientOrderID(listClientOrderID string) *CancelOCOService {
	s.listClientOrderID = listClientOrderID
	return s
}

// SeOrderListID sets orderListId
func (s *CancelOCOService) SeOrderListID(orderListID int64) *CancelOCOService {
	s.orderListID = orderListID
	return s
}

// SeNewClientOrderID sets newClientOrderId
func (s *CancelOCOService) SeNewClientOrderID(newClientOrderID string) *CancelOCOService {
	s.newClientOrderID = newClientOrderID
	return s
}

// Do send request
func (s *CancelOCOService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOCOResponse, err error) {
	// DELETE /api/v3/orderList | 取消整个订单列表
	// 取消订单列表中的单个订单将取消整个订单列表.
	// 如果 orderListId 和 listClientOrderId 一起发送, orderListId 优先被考虑.
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/api/v3/orderList",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	if s.listClientOrderID != "" {
		r.setFormParam("listClientOrderId", s.listClientOrderID)
	}
	if s.orderListID != 0 {
		r.setFormParam("orderListId", s.orderListID)
	}
	if s.newClientOrderID != "" {
		r.setFormParam("newClientOrderId", s.newClientOrderID)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelOCOResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CancelOpenOrdersService cancel all active orders on a symbol.
type CancelOpenOrdersService struct {
	c      *Client
	symbol string
}

// SetSymbol set symbol
func (s *CancelOpenOrdersService) SetSymbol(symbol string) *CancelOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CancelOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOpenOrdersResponse, err error) {
	// DELETE /api/v3/openOrders | 撤销单一交易对下所有挂单。这也包括了来自订单列表的挂单
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/api/v3/openOrders",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return &CancelOpenOrdersResponse{}, err
	}
	rawMessages := make([]*stdjson.RawMessage, 0)
	err = json.Unmarshal(data, &rawMessages)
	if err != nil {
		return &CancelOpenOrdersResponse{}, err
	}
	cancelOpenOrdersResponse := new(CancelOpenOrdersResponse)
	for _, j := range rawMessages {
		o := new(CancelOrderResponse)
		if err := json.Unmarshal(*j, o); err != nil {
			return &CancelOpenOrdersResponse{}, err
		}
		// Non-OCO orders guaranteed to have order list ID of -1
		if o.OrderListID == -1 {
			cancelOpenOrdersResponse.Orders = append(cancelOpenOrdersResponse.Orders, o)
			continue
		}
		oco := new(CancelOCOResponse)
		if err := json.Unmarshal(*j, oco); err != nil {
			return &CancelOpenOrdersResponse{}, err
		}
		cancelOpenOrdersResponse.OCOOrders = append(cancelOpenOrdersResponse.OCOOrders, oco)
	}
	return cancelOpenOrdersResponse, nil
}

// CancelOpenOrdersResponse defines cancel open orders response.
type CancelOpenOrdersResponse struct {
	Orders    []*CancelOrderResponse
	OCOOrders []*CancelOCOResponse
}

// CancelOrderResponse may be returned included in a CancelOpenOrdersResponse.
type CancelOrderResponse struct {
	Symbol                   string          `json:"symbol"`
	OrigClientOrderID        string          `json:"origClientOrderId"`
	OrderID                  int64           `json:"orderId"`
	OrderListID              int64           `json:"orderListId"`
	ClientOrderID            string          `json:"clientOrderId"`
	TransactTime             int64           `json:"transactTime"`
	Price                    string          `json:"price"`
	OrigQuantity             string          `json:"origQty"`
	ExecutedQuantity         string          `json:"executedQty"`
	CummulativeQuoteQuantity string          `json:"cummulativeQuoteQty"`
	Status                   OrderStatusType `json:"status"`
	TimeInForce              TimeInForceType `json:"timeInForce"`
	Type                     OrderType       `json:"type"`
	Side                     SideType        `json:"side"`
}

// CancelOCOResponse may be returned included in a CancelOpenOrdersResponse.
type CancelOCOResponse struct {
	OrderListID       int64             `json:"orderListId"`
	ContingencyType   string            `json:"contingencyType"`
	ListStatusType    string            `json:"listStatusType"`
	ListOrderStatus   string            `json:"listOrderStatus"`
	ListClientOrderID string            `json:"listClientOrderId"`
	TransactionTime   int64             `json:"transactionTime"`
	Symbol            string            `json:"symbol"`
	Orders            []*OCOOrder       `json:"orders"`
	OrderReports      []*OCOOrderReport `json:"orderReports"`
}
