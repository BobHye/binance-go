package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// WsConfig 客户端配置
type WsConfig struct {
	Endpoint          string        // WebSocket地址
	ReconnectBase     time.Duration // 基础重连间隔（默认 1s）
	MaxReconnectTime  time.Duration // 最大重连间隔（默认 1min）
	WriteTimeout      time.Duration // 写超时（默认 5s）
	HeartbeatInterval time.Duration // 心跳间隔（默认30秒）
	RequestHeader     http.Header   // 请求头
	ReadTimeout       time.Duration // 读超时（默认 10s）
	BufferSize        int           // 通道缓冲大小（默认 100）
}

func NewWsConfig(endpoint string) *WsConfig {
	// 高频交易场景推荐配置
	return &WsConfig{
		Endpoint:          endpoint,
		ReconnectBase:     500 * time.Millisecond,
		MaxReconnectTime:  10 * time.Second,
		WriteTimeout:      1 * time.Second,
		HeartbeatInterval: 5 * time.Second,
	}
}

// WebSocketClient WebSocket 客户端
type WebSocketClient struct {
	config        *WsConfig
	dialer        *websocket.Dialer
	conn          *websocket.Conn
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	mu            sync.RWMutex
	sendChan      chan []byte
	reconnectChan chan struct{}

	// 回调函数
	OnMessage func([]byte) // 消息回调
	OnError   func(error)  // 错误回调
	OnConnect func()       // 连接成功回调
}

// NewClient 创建客户端
func NewClient(cfg *WsConfig) *WebSocketClient {
	// 设置默认配置
	if cfg.ReconnectBase == 0 {
		cfg.ReconnectBase = 500 * time.Millisecond
	}
	if cfg.MaxReconnectTime == 0 {
		cfg.MaxReconnectTime = 10 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 1 * time.Second
	}
	if cfg.HeartbeatInterval == 0 {
		cfg.HeartbeatInterval = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WebSocketClient{
		config:        cfg,
		dialer:        websocket.DefaultDialer,
		ctx:           ctx,
		cancel:        cancel,
		sendChan:      make(chan []byte, 100),
		reconnectChan: make(chan struct{}, 1),
	}
}

// Connect 启动连接
func (c *WebSocketClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return errors.New("connection already exists")
	}
	go c.connectionManager()
	return nil
}

// Close 关闭连接
func (c *WebSocketClient) Close() {
	c.cancel()
	c.wg.Wait()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}
}

// Send 发送消息
func (c *WebSocketClient) Send(message []byte) error {
	select {
	case c.sendChan <- message:
		return nil
	case <-c.ctx.Done():
		return errors.New("connection closed")
	default:
		return errors.New("send buffer full")
	}
}

// Events 获取事件通道
//func (c *WebSocketClient) Events() <-chan ClientEvent {
//	return c.eventChan
//}

// 内部状态管理
//func (c *WebSocketClient) setState(state ClientState) {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	c.state = state
//}

// 连接管理主循环
func (c *WebSocketClient) connectionManager() {
	// 使用指数退避策略初始化重连间隔
	backoff := c.config.ReconnectBase
	// 主循环持续监控连接状态
	for {
		select {
		case <-c.ctx.Done(): // 监听上下文取消信号
			return // 彻底退出连接管理循环
		default:
			// 尝试建立新连接
			if err := c.connect(); err != nil { // 实际建立TCP/WS连接
				// 连接失败处理
				c.handleError(fmt.Errorf("connection failed: %w", err))

				// 调度下次重连（带退避策略）
				c.scheduleReconnect(backoff) // 这里会阻塞 d 时间

				// 指数增长重连间隔（不超过配置的最大值）
				backoff = minDuration(backoff*2, c.config.MaxReconnectTime) // 更新退避时间
				continue                                                    // 继续循环 -> 再次执行 c.connect()
			}

			// ---------- 连接成功后的处理 ----------
			// 重置退避时间为初始值（成功连接后恢复基准
			backoff = c.config.ReconnectBase

			// 启动三个核心协程（使用WaitGroup同步）
			c.wg.Add(3)
			go c.readPump()  // 消息接收循环
			go c.writePump() // 消息发送循环
			go c.heartbeat() // 心跳维持循环

			// ---------- 连接维持阶段 ----------
			// 阻塞等待重连信号（当收到信号时表示需要重新连接）
			<-c.reconnectChan // 可能来自：读写错误、心跳失败、服务器关闭等
			c.cleanup()       // 关闭现有连接，重置为nil

			// 循环将回到顶部，开始新的连接尝试
		}
	}
}

// connect 建立实际连接
func (c *WebSocketClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.config.ReconnectBase)
	defer cancel()

	conn, _, err := c.dialer.DialContext(ctx, c.config.Endpoint, c.config.RequestHeader)
	if err != nil {
		return err
	}

	c.conn = conn
	conn.SetCloseHandler(c.closeHandler)

	if c.OnConnect != nil {
		c.OnConnect()
	}
	return nil
}

// readPump 消息读取循环
func (c *WebSocketClient) readPump() {
	defer c.wg.Done()
	defer c.triggerReconnect()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.handleError(fmt.Errorf("read error: %w", err))
				return
			}

			if c.OnMessage != nil {
				c.OnMessage(message)
			}
		}
	}
}

// 消息写入循环
func (c *WebSocketClient) writePump() {
	defer c.wg.Done()
	defer c.triggerReconnect()

	for {
		select {
		case <-c.ctx.Done():
			return
		case message := <-c.sendChan:
			if err := c.writeWithTimeout(message); err != nil {
				c.handleError(fmt.Errorf("write error: %w", err))
				return
			}
		}
	}
}

// 带超时的写入
func (c *WebSocketClient) writeWithTimeout(message []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return errors.New("connection closed")
	}

	err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, message)
}

// 心跳检测
func (c *WebSocketClient) heartbeat() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.ping(); err != nil {
				c.handleError(fmt.Errorf("heartbeat failed: %w", err))
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// 发送Ping消息
func (c *WebSocketClient) ping() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return errors.New("connection closed")
	}
	err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

// 关闭处理器
func (c *WebSocketClient) closeHandler(code int, text string) error {
	c.handleError(fmt.Errorf("connection closed by server: code=%d, reason=%s", code, text))
	return nil
}

// 触发重连
func (c *WebSocketClient) triggerReconnect() {
	select {
	case c.reconnectChan <- struct{}{}:
	default:
	}
}

// 清理资源
func (c *WebSocketClient) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return
		}
		c.conn = nil
	}
}

// 错误处理
func (c *WebSocketClient) handleError(err error) {
	if c.OnError != nil {
		c.OnError(err)
	}
}

// 调度重连
func (c *WebSocketClient) scheduleReconnect(d time.Duration) {
	select {
	case <-c.ctx.Done(): // 1. 优先检查上下文是否已取消
		// 如果客户端已主动关闭，立即退出不执行重连
	case <-time.After(d): // 2. 阻塞等待指定时间间隔
		// 计时结束后不执行任何操作，控制权返回给调用方
	}
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
