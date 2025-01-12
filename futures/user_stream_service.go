package futures

import (
	"context"
	"net/http"
)

// StartUserStreamService create listen key for user stream service | 生成listenKey
type StartUserStreamService struct {
	c *Client
}

// Do send request
func (s *StartUserStreamService) Do(ctx context.Context, opts ...RequestOption) (listenKey string, err error) {
	// POST /fapi/v1/listenKey | 生成listenKey，创建一个新的user data stream，返回值为一个listenKey，即websocket订阅的stream名称。如果该帐户具有有效的listenKey，则将返回该listenKey并将其有效期延长60分钟。
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/listenKey",
		secType:  secTypeSigned,
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

// KeepaliveUserStreamService update listen key | 延长listenKey有效期
type KeepaliveUserStreamService struct {
	c         *Client
	listenKey string
}

// SetListenKey set listen key
func (s *KeepaliveUserStreamService) SetListenKey(listenKey string) *KeepaliveUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *KeepaliveUserStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// PUT /fapi/v1/listenKey | 有效期延长至本次调用后60分钟
	r := &request{
		method:   http.MethodPut,
		endpoint: "/fapi/v1/listenKey",
		secType:  secTypeSigned,
	}
	r.setFormParam("listenKey", s.listenKey) // 被延长的listenKey
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// CloseUserStreamService delete listen key | 关闭某账户数据流
type CloseUserStreamService struct {
	c         *Client
	listenKey string
}

// SetListenKey set listen key
func (s *CloseUserStreamService) SetListenKey(listenKey string) *CloseUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *CloseUserStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// DELETE /fapi/v1/listenKey | 关闭某账户数据流
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/fapi/v1/listenKey",
		secType:  secTypeSigned,
	}
	r.setFormParam("listenKey", s.listenKey)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}
