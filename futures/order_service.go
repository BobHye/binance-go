package futures

import (
	"context"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"github.com/BobHye/binance-go/api"
	"net/http"
	"strings"
)

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
	workingType      *WorkingType
	activationPrice  *string
	callbackRate     *string
	priceProtect     *bool
	newOrderRespType NewOrderRespType
	closePosition    *bool
}

// Symbol set symbol
func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CreateOrderService) Side(side SideType) *CreateOrderService {
	s.side = side
	return s
}

// PositionSide set side
func (s *CreateOrderService) PositionSide(positionSide PositionSideType) *CreateOrderService {
	s.positionSide = &positionSide
	return s
}

// Type set type
func (s *CreateOrderService) Type(orderType OrderType) *CreateOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *CreateOrderService) TimeInForce(timeInForce TimeInForceType) *CreateOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CreateOrderService) Quantity(quantity string) *CreateOrderService {
	s.quantity = quantity
	return s
}

// ReduceOnly set reduceOnly
func (s *CreateOrderService) ReduceOnly(reduceOnly bool) *CreateOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// Price set price
func (s *CreateOrderService) Price(price string) *CreateOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateOrderService) NewClientOrderID(newClientOrderID string) *CreateOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// StopPrice set stopPrice
func (s *CreateOrderService) StopPrice(stopPrice string) *CreateOrderService {
	s.stopPrice = &stopPrice
	return s
}

// WorkingType set workingType
func (s *CreateOrderService) WorkingType(workingType WorkingType) *CreateOrderService {
	s.workingType = &workingType
	return s
}

// ActivationPrice set activationPrice
func (s *CreateOrderService) ActivationPrice(activationPrice string) *CreateOrderService {
	s.activationPrice = &activationPrice
	return s
}

// CallbackRate set callbackRate
func (s *CreateOrderService) CallbackRate(callbackRate string) *CreateOrderService {
	s.callbackRate = &callbackRate
	return s
}

// PriceProtect set priceProtect
func (s *CreateOrderService) PriceProtect(priceProtect bool) *CreateOrderService {
	s.priceProtect = &priceProtect
	return s
}

// NewOrderResponseType set newOrderResponseType
func (s *CreateOrderService) NewOrderResponseType(newOrderResponseType NewOrderRespType) *CreateOrderService {
	s.newOrderRespType = newOrderResponseType
	return s
}

// ClosePosition set closePosition
func (s *CreateOrderService) ClosePosition(closePosition bool) *CreateOrderService {
	s.closePosition = &closePosition
	return s
}

type CreateOrderResponse struct {
	Symbol            string           `json:"symbol"`
	OrderID           int64            `json:"orderId"`
	ClientOrderID     string           `json:"clientOrderId"`
	Price             float64          `json:"price,string"`
	OrigQuantity      float64          `json:"origQty,string"`
	ExecutedQuantity  float64          `json:"executedQty,string"`
	CumQuote          float64          `json:"cumQuote,string"`
	ReduceOnly        bool             `json:"reduceOnly"`
	Status            OrderStatusType  `json:"status"`
	StopPrice         float64          `json:"stopPrice,string"`
	TimeInForce       TimeInForceType  `json:"timeInForce"`
	Type              OrderType        `json:"type"`
	Side              SideType         `json:"side"`
	UpdateTime        int64            `json:"updateTime"`
	WorkingType       WorkingType      `json:"workingType"`
	ActivatePrice     float64          `json:"activatePrice,string"`
	PriceRate         float64          `json:"priceRate,string"`
	AvgPrice          float64          `json:"avgPrice,string"`
	PositionSide      PositionSideType `json:"positionSide"`
	ClosePosition     bool             `json:"closePosition"`
	PriceProtect      bool             `json:"priceProtect"`
	RateLimitOrder10s string           `json:"rateLimitOrder10s,omitempty"`
	RateLimitOrder1m  string           `json:"rateLimitOrder1m,omitempty"`
}

// createOrder 新建订单
func (s *CreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...api.RequestOption) (data []byte, header *http.Header, err error) {
	r := &api.Request{
		Method:   http.MethodPost,
		Endpoint: endpoint,
		SecType:  api.SecTypeSigned,
	}
	m := api.Params{
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
	r.SetFormParams(m)
	data, header, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	return data, header, nil
}

// Do 发送请求
func (s *CreateOrderService) Do(ctx context.Context, opts ...api.RequestOption) (res *CreateOrderResponse, err error) {
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

// Order 定义订单信息
type Order struct {
	Symbol           string           `json:"symbol"`
	OrderID          int64            `json:"orderId"`
	ClientOrderID    string           `json:"clientOrderId"`
	Price            float64          `json:"price,string"`
	ReduceOnly       bool             `json:"reduceOnly"`
	OrigQuantity     float64          `json:"origQty,string"`
	ExecutedQuantity float64          `json:"executedQty,string"`
	CumQuantity      float64          `json:"cumQty,string"`
	CumQuote         float64          `json:"cumQuote,string"`
	Status           OrderStatusType  `json:"status"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	Side             SideType         `json:"side"`
	StopPrice        float64          `json:"stopPrice,string"`
	Time             int64            `json:"time"`
	UpdateTime       int64            `json:"updateTime"`
	WorkingType      WorkingType      `json:"workingType"`
	ActivatePrice    float64          `json:"activatePrice,string"`
	PriceRate        float64          `json:"priceRate,string"`
	AvgPrice         float64          `json:"avgPrice,string"`
	OrigType         string           `json:"origType"`
	PositionSide     PositionSideType `json:"positionSide"`
	PriceProtect     bool             `json:"priceProtect"`
	ClosePosition    bool             `json:"closePosition"`
}

// ListOpenOrdersService 列出当前全部挂单
type ListOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *ListOpenOrdersService) Symbol(symbol string) *ListOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *ListOpenOrdersService) Do(ctx context.Context, opts ...api.RequestOption) (res []*Order, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/openOrders",
		SecType:  api.SecTypeSigned,
	}
	if s.symbol != "" {
		r.SetParam("symbol", s.symbol)
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
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

func (s *GetOpenOrderService) Symbol(symbol string) *GetOpenOrderService {
	s.symbol = symbol
	return s
}

func (s *GetOpenOrderService) OrderID(orderID int64) *GetOpenOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID  查询当前挂单
func (s *GetOpenOrderService) OrigClientOrderID(origClientOrderID string) *GetOpenOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

func (s *GetOpenOrderService) Do(ctx context.Context, opts ...api.RequestOption) (res *Order, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/openOrder",
		SecType:  api.SecTypeSigned,
	}
	r.SetParam("symbol", s.symbol)
	if s.orderID == nil && s.origClientOrderID == nil {
		return nil, errors.New("either orderId or origClientOrderId must be sent")
	}
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetParam("origClientOrderId", *s.origClientOrderID)
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

// GetOrderService 查询订单service
type GetOrderService struct {
	c                 *Client
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Symbol set symbol
func (s *GetOrderService) Symbol(symbol string) *GetOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *GetOrderService) OrderID(orderID int64) *GetOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *GetOrderService) OrigClientOrderID(origClientOrderID string) *GetOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *GetOrderService) Do(ctx context.Context, opts ...api.RequestOption) (res *Order, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/order",
		SecType:  api.SecTypeSigned,
	}
	r.SetParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetParam("origClientOrderId", *s.origClientOrderID)
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

// ListOrdersService 查询所有订单(包括历史订单); active, canceled, or filled
type ListOrdersService struct {
	c         *Client
	symbol    string
	orderID   *int64
	startTime *int64
	endTime   *int64
	limit     *int
}

// Symbol set symbol
func (s *ListOrdersService) Symbol(symbol string) *ListOrdersService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *ListOrdersService) OrderID(orderID int64) *ListOrdersService {
	s.orderID = &orderID
	return s
}

// StartTime set starttime
func (s *ListOrdersService) StartTime(startTime int64) *ListOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endtime
func (s *ListOrdersService) EndTime(endTime int64) *ListOrdersService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ListOrdersService) Limit(limit int) *ListOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListOrdersService) Do(ctx context.Context, opts ...api.RequestOption) (res []*Order, err error) {
	r := &api.Request{
		Method:   http.MethodGet,
		Endpoint: "/fapi/v1/allOrders",
		SecType:  api.SecTypeSigned,
	}
	r.SetParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetParam("orderId", *s.orderID)
	}
	if s.startTime != nil {
		r.SetParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.SetParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.SetParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// CancelOrderResponse define response of canceling order
type CancelOrderResponse struct {
	ClientOrderID    string           `json:"clientOrderId"`
	CumQuantity      float64          `json:"cumQty,string"`
	CumQuote         float64          `json:"cumQuote,string"`
	ExecutedQuantity float64          `json:"executedQty,string"`
	OrderID          int64            `json:"orderId"`
	OrigQuantity     float64          `json:"origQty,string"`
	Price            float64          `json:"price,string"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             SideType         `json:"side"`
	Status           OrderStatusType  `json:"status"`
	StopPrice        float64          `json:"stopPrice,string"`
	Symbol           string           `json:"symbol"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	UpdateTime       int64            `json:"updateTime"`
	WorkingType      WorkingType      `json:"workingType"`
	ActivatePrice    float64          `json:"activatePrice,string"`
	PriceRate        float64          `json:"priceRate,string"`
	OrigType         string           `json:"origType"`
	PositionSide     PositionSideType `json:"positionSide"`
	PriceProtect     bool             `json:"priceProtect"`
}

// CancelOrderService 取消订单service
type CancelOrderService struct {
	c                 *Client
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Symbol set symbol
func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CancelOrderService) OrderID(orderID int64) *CancelOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelOrderService) OrigClientOrderID(origClientOrderID string) *CancelOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *CancelOrderService) Do(ctx context.Context, opts ...api.RequestOption) (res *CancelOrderResponse, err error) {
	r := &api.Request{
		Method:   http.MethodDelete,
		Endpoint: "/fapi/v1/order",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParam("symbol", s.symbol)
	if s.orderID != nil {
		r.SetFormParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.SetFormParam("origClientOrderId", *s.origClientOrderID)
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

// CancelAllOpenOrdersService 撤销所有订单
type CancelAllOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *CancelAllOpenOrdersService) Symbol(symbol string) *CancelAllOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CancelAllOpenOrdersService) Do(ctx context.Context, opts ...api.RequestOption) (err error) {
	r := &api.Request{
		Method:   http.MethodDelete,
		Endpoint: "/fapi/v1/allOpenOrders",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParam("symbol", s.symbol)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// CancelMultiplesOrdersService 批量撤销订单service
type CancelMultiplesOrdersService struct {
	c                     *Client
	symbol                string
	orderIDList           []int64
	origClientOrderIDList []string
}

// Symbol set symbol
func (s *CancelMultiplesOrdersService) Symbol(symbol string) *CancelMultiplesOrdersService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CancelMultiplesOrdersService) OrderIDList(orderIDList []int64) *CancelMultiplesOrdersService {
	s.orderIDList = orderIDList
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelMultiplesOrdersService) OrigClientOrderIDList(origClientOrderIDList []string) *CancelMultiplesOrdersService {
	s.origClientOrderIDList = origClientOrderIDList
	return s
}

// Do send request 批量撤销订单
func (s *CancelMultiplesOrdersService) Do(ctx context.Context, opts ...api.RequestOption) (res []*CancelOrderResponse, err error) {
	r := &api.Request{
		Method:   http.MethodDelete,
		Endpoint: "/fapi/v1/batchOrders",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParam("symbol", s.symbol)
	if s.orderIDList != nil {
		// convert a slice of integers to a string e.g. [1 2 3] => "[1,2,3]"
		orderIDListString := strings.Join(strings.Fields(fmt.Sprint(s.orderIDList)), ",")
		r.SetFormParam("orderIdList", orderIDListString)
	}
	if s.origClientOrderIDList != nil {
		r.SetFormParam("origClientOrderIdList", s.origClientOrderIDList)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}

// CreateBatchOrdersService 批量下单service
type CreateBatchOrdersService struct {
	c      *Client
	orders []*CreateOrderService
}

type CreateBatchOrdersResponse struct {
	Orders []*Order
}

func (s *CreateBatchOrdersService) OrderList(orders []*CreateOrderService) *CreateBatchOrdersService {
	s.orders = orders
	return s
}

func (s *CreateBatchOrdersService) Do(ctx context.Context, opts ...api.RequestOption) (res *CreateBatchOrdersResponse, err error) {
	r := &api.Request{
		Method:   http.MethodPost,
		Endpoint: "/fapi/v1/batchOrders",
		SecType:  api.SecTypeSigned,
	}

	var orders []api.Params
	for _, order := range s.orders {
		m := api.Params{
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
	m := api.Params{
		"batchOrders": string(b),
	}

	r.SetFormParams(m)

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
