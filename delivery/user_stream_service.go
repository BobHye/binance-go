package delivery

import (
	"context"
	"net/http"
)

// StartUserStreamService create listen key for user stream service
type StartUserStreamService struct {
	c *Client
}

type ListenKey struct {
	ListenKey string `json:"listenKey"`
}

// Do send request
func (s *StartUserStreamService) Do(ctx context.Context, opts ...RequestOption) (listenKey string, err error) {
	// POST /papi/v1/listenKey | 生成listenKey
	r := &request{
		method:   http.MethodPost,
		endpoint: "/dapi/v1/listenKey",
		secType:  SecTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return "", err
	}
	var res ListenKey
	err = json.Unmarshal(data, &res)
	return res.ListenKey, err
}

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
	// PUT /fapi/v1/listenKey | 延长listenKey有效期
	r := &request{
		method:   http.MethodPut,
		endpoint: "dapi/v1/listenKey",
		secType:  SecTypeSigned,
	}
	r.setFormParam("listenKey", s.listenKey)
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// CloseUserStreamService delete listenKey
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
	// DELETE /fapi/v1/listenKey | 关闭listenKey
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/dapi/v1/listenKey",
		secType:  SecTypeSigned,
	}
	r.setFormParam("listenKey", s.listenKey)
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}
