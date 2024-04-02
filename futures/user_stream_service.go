package futures

import (
	"context"
	"github.com/BobHye/binance-go/api"
	"net/http"
)

// StartUserStreamService 生成listenKey
type StartUserStreamService struct {
	c *Client
}

// Do send request
func (s *StartUserStreamService) Do(ctx context.Context, opts ...api.RequestOption) (listenKey string, err error) {
	r := &api.Request{
		Method:   http.MethodPost,
		Endpoint: "/fapi/v1/listenKey",
		SecType:  api.SecTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return "", err
	}
	var res ListenKey
	err = json.Unmarshal(data, &res)
	return res.ListenKey, err
}

type ListenKey struct {
	ListenKey string `json:"listenKey"`
}

// KeepaliveUserStreamService 延长listenKey有效期
type KeepaliveUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *KeepaliveUserStreamService) ListenKey(listenKey string) *KeepaliveUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *KeepaliveUserStreamService) Do(ctx context.Context, opts ...api.RequestOption) (err error) {
	r := &api.Request{
		Method:   http.MethodPut,
		Endpoint: "/fapi/v1/listenKey",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParam("listenKey", s.listenKey)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// CloseUserStreamService 关闭某账户数据流
type CloseUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *CloseUserStreamService) ListenKey(listenKey string) *CloseUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *CloseUserStreamService) Do(ctx context.Context, opts ...api.RequestOption) (err error) {
	r := &api.Request{
		Method:   http.MethodDelete,
		Endpoint: "/fapi/v1/listenKey",
		SecType:  api.SecTypeSigned,
	}
	r.SetFormParam("listenKey", s.listenKey)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}
