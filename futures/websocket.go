package futures

import (
	"errors"
	"github.com/BobHye/binance-go/log"
	"github.com/BobHye/wsc"
)

// WsHandler handle raw websocket message | 处理原始 websocket 消息
type WsHandler func(message []byte)

// ErrHandler handles errors | 统一错误处理回调函数类型
type ErrHandler func(err error)

// WsConfig webservice configuration | webservice 配置
type WsConfig struct {
	Endpoint string // WebSocket服务端点地址
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

var wsServe = func(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	// 参数校验前置
	if cfg == nil || cfg.Endpoint == "" {
		return nil, nil, errors.New("invalid config")
	}

	done = make(chan struct{})

	ws = wsc.New(cfg.Endpoint)
	ws.OnConnected(func() {
		// 连接建立时
		if log.Default.OnConnected {
			log.Default.Log("websocket connected")
		}
	})
	ws.OnConnectError(errHandler) // 连接错误时
	ws.OnDisconnected(errHandler) // 连接断开时
	ws.OnClose(func(code int, text string) {
		// 连接关闭时
		if log.Default.OnClose {
			log.Default.Log("websocket closed, code: %d, message: %s", code, text)
		}
	})
	ws.OnSentError(errHandler) // 消息发送失败时
	ws.OnPingReceived(func(appData string) {
		// 接收PING帧
		if log.Default.OnPingReceived {
			log.Default.Log("ping received, data: %s", appData)
		}
	})
	ws.OnPongReceived(func(appData string) {
		// 接收PONG帧
		if log.Default.OnPongReceived {
			log.Default.Log("pong received, data: %s", appData)
		}
	})
	ws.OnTextMessageReceived(handler) // 收到文本消息时
	ws.OnKeepalive(func() {
		// 保活心跳
		if log.Default.OnKeepalive {
			log.Default.Log("keep alive")
		}
	})

	// 启动连接并在 goroutine 中运行
	go func() {
		ws.Connect()     // 建立 WebSocket 连接（可能阻塞，需在 goroutine 中执行）
		for range done { // 等待 done 通道的信号
			ws.Close() // 接收到信号时关闭连接
			return
		}
	}()
	return // 返回客户端实例、done 通道和可能的错误
}
