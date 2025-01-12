package delivery

import (
	"context"
	"net/http"
)

// CreateOrderService create order
type CreateOrderService struct {
	c                *Client
	symbol           string
	side             SideType
	positionSide     *PositionSideType
	orderType        OrderType
	timeInForce      *TimeInForceType
	quantity         string
	reduceOnly       *bool
	price            *string
	newClientOrderID *string
	stopPrice        *string
	closePosition    *bool
	activationPrice  *string
	callbackRate     *string
	workingType      *WorkingType
	priceProtect     *bool
	newOrderRespType NewOrderRespType
}

// SetSymbol set symbol
func (s *CreateOrderService) setSymbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

// SetSide set side
func (s *CreateOrderService) SetSide(side SideType) *CreateOrderService {
	s.side = side
	return s
}

// SetPositionSide set positionSide
func (s *CreateOrderService) SetPositionSide(positionSide PositionSideType) *CreateOrderService {
	s.positionSide = &positionSide
	return s
}

// SetType set order type
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

// SetNewClientOrderID set NewClientOrderID
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

// createOrder
func (s *CreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, err error) {
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
		m["priceProduct"] = *s.priceProtect
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
	data, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateOrderResponse, err error) {
	data, err := s.createOrder(ctx, "/dapi/v1/order", opts...)
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

// CreateOrderResponse define create order response
type CreateOrderResponse struct {
	ClientOrderID    string           `json:"clientOrderId"`
	CumQuantity      string           `json:"cumQty"`
	CumBase          string           `json:"cumBase"`
	ExecutedQuantity string           `json:"executedQty"`
	OrderID          int64            `json:"orderId"`
	AvgPrice         string           `json:"avgPrice"`
	OrigQuantity     string           `json:"origQty"`
	Price            string           `json:"price"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             SideType         `json:"side"`
	PositionSide     PositionSideType `json:"positionSide"`
	Status           OrderStatusType  `json:"status"`
	StopPrice        string           `json:"stopPrice"`
	ClosePosition    bool             `json:"closePosition"`
	Symbol           string           `json:"symbol"`
	Pair             string           `json:"pair"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	OrigType         OrderType        `json:"origType"`
	ActivatePrice    string           `json:"activatePrice"`
	PriceRate        string           `json:"priceRate"`
	UpdateTime       int64            `json:"updateTime"`
	WorkingType      WorkingType      `json:"workingType"`
	PriceProtect     bool             `json:"priceProtect"`
}

// ListOpenOrdersService list opened orders
type ListOpenOrdersService struct {
	c      *Client
	symbol string
	pair   string
}

// SetSymbol set symbol
func (s *ListOpenOrdersService) SetSymbol(symbol string) *ListOpenOrdersService {
	s.symbol = symbol
	return s
}

// SetPair set pair
func (s *ListOpenOrdersService) SetPair(pair string) *ListOpenOrdersService {
	s.pair = pair
	return s
}

// Do send request
func (s *ListOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*Order, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/openOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	if s.pair != "" {
		r.setParam("pair", s.symbol)
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

type Order struct {
	AvgPrice         string           `json:"avgPrice"`
	ClientOrderID    string           `json:"clientOrderId"`
	CumBase          string           `json:"cumBase"`
	ExecutedQuantity string           `json:"executedQty"`
	OrderID          int64            `json:"orderId"`
	OrigQuantity     string           `json:"origQty"`
	OrigType         OrderType        `json:"origType"`
	Price            string           `json:"price"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             SideType         `json:"side"`
	PositionSide     PositionSideType `json:"positionSide"`
	Status           OrderStatusType  `json:"status"`
	StopPrice        string           `json:"stopPrice"`
	ClosePosition    bool             `json:"closePosition"`
	Symbol           string           `json:"symbol"`
	Pair             string           `json:"pair"`
	Time             int64            `json:"time"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	ActivatePrice    string           `json:"activatePrice"`
	PriceRate        string           `json:"priceRate"`
	UpdateTime       int64            `json:"updateTime"`
	WorkingType      WorkingType      `json:"workingType"`
	PriceProtect     bool             `json:"priceProtect"`
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
func (s GetOrderService) Do(ctx context.Context, opts ...RequestOption) (res *Order, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/order",
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
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListOrdersService all account orders; active, canceled, or filled | 账户所有订单；有效、已取消或已完成
type ListOrdersService struct {
	c         *Client
	symbol    string
	pair      string
	OrderID   *int64
	startTime *int64
	endTime   *int64
	limit     *int
}

// SetSymbol set symbol
func (s *ListOrdersService) SetSymbol(symbol string) *ListOrdersService {
	s.symbol = symbol
	return s
}

// SetPair set pair
func (s *ListOrdersService) SetPair(pair string) *ListOrdersService {
	s.pair = pair
	return s
}

// SetOrderID set orderID
func (s *ListOrdersService) SetOrderID(orderID int64) *ListOrdersService {
	s.OrderID = &orderID
	return s
}

// SetStartTime set startTime
func (s *ListOrdersService) SetStartTime(startTime int64) *ListOrdersService {
	s.startTime = &startTime
	return s
}

// SetEndTime set endTime
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
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/allOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	if s.pair != "" {
		r.setParam("pair", s.pair)
	}
	if s.OrderID != nil {
		r.setParam("orderId", *s.OrderID)
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
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/dapi/v1/order",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setFormParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setFormParam("origClientOrderId", *s.origClientOrderID)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelOrderResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CancelOrderResponse define response of canceling order
type CancelOrderResponse struct {
	AvgPrice         string           `json:"avgPrice"`
	ClientOrderID    string           `json:"clientOrderID"`
	CumQuantity      string           `json:"cumQty"`
	CumBase          string           `json:"cumBase"`
	ExecutedQuantity string           `json:"executedQty"`
	OrderID          int64            `json:"orderId"`
	OrigQuantity     string           `json:"origQty"`
	OrigType         OrderType        `json:"origType"`
	Price            string           `json:"price"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             SideType         `json:"side"`
	PositionSide     PositionSideType `json:"positionSide"`
	Status           OrderStatusType  `json:"status"`
	StopPrice        string           `json:"stopPrice"`
	ClosePosition    bool             `json:"closePosition"`
	Symbol           string           `json:"symbol"`
	Pair             string           `json:"pair"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	ActivatePrice    string           `json:"activatePrice"`
	PriceRate        string           `json:"priceRate"`
	UpdateTime       int64            `json:"updateTime"`
	WorkingType      WorkingType      `json:"workingType"`
	PriceProtect     bool             `json:"priceProtect"`
}

type CancelAllOpenOrdersService struct {
	c      *Client
	symbol string
}

// SetSymbol set symbol
func (s *CancelAllOpenOrdersService) SetSymbol(symbol string) *CancelAllOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CancelAllOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/dapi/v1/allOpenOrders",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	_, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// ListLiquidationOrdersService list liquidation orders
type ListLiquidationOrdersService struct {
	c         *Client
	symbol    *string
	pair      *string
	startTime *int64
	endTime   *int64
	limit     *int
}

// SetSymbol set symbol
func (s *ListLiquidationOrdersService) SetSymbol(symbol string) *ListLiquidationOrdersService {
	s.symbol = &symbol
	return s
}

// SetPair set pair
func (s *ListLiquidationOrdersService) SetPair(pair string) *ListLiquidationOrdersService {
	s.pair = &pair
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
		endpoint: "/dapi/v1/allForceOrder",
		secType:  secTypeSigned,
	}
	if s.pair != nil {
		r.setParam("pair", *s.pair)
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
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
		return []*LiquidationOrder{}, err
	}
	res = make([]*LiquidationOrder, 0)
	err = json.Unmarshal(data, res)
	if err != nil {
		return []*LiquidationOrder{}, err
	}
	return res, nil
}

// LiquidationOrder define liquidation order
type LiquidationOrder struct {
	Symbol           string          `json:"symbol"`
	Price            string          `json:"price"`
	OrigQuantity     string          `json:"origQty"`
	ExecutedQuantity string          `json:"executedQty"`
	AveragePrice     string          `json:"averagePrice"`
	Status           OrderStatusType `json:"status"`
	TimeInForce      TimeInForceType `json:"timeInForce"`
	Type             OrderType       `json:"type"`
	Side             SideType        `json:"side"`
	Time             int64           `json:"time"`
}
