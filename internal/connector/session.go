package connector

import (
	"net"
	"sync"
	"container/list"
	"sync/atomic"
	"errors"
	"time"
	"magpie/internal/com/logger"
	"magpie/internal/protocol"
)

var (
	ErrClosed = errors.New("link.Session closed")
	globalSessionId uint64

// Session是否已验证
	VALID_FLAG_DEALING int32 = 1 // 验证中
	VALID_FLAG_SUCC int32 = 2 // 验证成功
	VALID_FLAG_FAIL int32 = 3 // 验证失败
	VALID_FLAT_NOT_START int32 = 0 // 未验证

// Session是否已关闭
	CLOSE_FLAG_TRUE int32 = 1
	CLOSE_FLAT_FALSE int32 = 2
)

type Session struct {
	id              uint64
	conn            net.Conn

	closeChan       chan int

	closeEventMutex sync.Mutex
	closeCallbacks  *list.List // 关闭conn前的回调

    // 消息的处理器(解码得到的消息的处理 & 发送前消息的处理 & 异常处理)
	handler         *MsgHandler
	validFlag       int32      // 是否已验证
	closeFlag       int32      // 关闭的标志

	HeartbeatTime   int64      // 最后收到心跳的时间戳
	State           interface{}
}

type closeCallback struct {
	Handler interface{}
	Func    func()
}

func NewSession(conn net.Conn, handler *MsgHandler) *Session {
	session := &Session{
		id:             atomic.AddUint64(&globalSessionId, 1),
		conn:           conn,
		closeCallbacks: list.New(),
		handler: handler,
		closeFlag: CLOSE_FLAT_FALSE,
		HeartbeatTime: time.Now().Unix(),
	}
	return session
}

func (this *Session) Id() uint64 {
	return this.id
}

func (this *Session) Conn() net.Conn {
	return this.conn
}

func (this *Session) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == CLOSE_FLAG_TRUE
}

func (this *Session ) IsValid() bool {
	return atomic.LoadInt32(&this.validFlag) == VALID_FLAG_SUCC
}

/**
 * 关闭session
 */
func (this *Session) Close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, CLOSE_FLAT_FALSE, CLOSE_FLAG_TRUE) {
		this.invokeCloseCallbacks()
		if this.closeChan != nil {
			close(this.closeChan)
		}
		this.conn.Close()
		logger.Error("sess.close", logger.Format("sess", this))
	}
}

/**
 * 接收内容
 */
func (this *Session) Receive() (error) {
	if this.IsClosed() {
		return ErrClosed
	}

	go this.handler.ReceiveData()
	return nil
}

/**
 * 发送内容
 */
func (this *Session) Send(p *protocol.Pkg) (err error) {
	if this.IsClosed() {
		return ErrClosed
	}

	return this.handler.SendData(p)
}

/**
 * 添加关闭session的回调
 */
func (this *Session) AddCloseCallback(handler interface{}, callback func()) {
	if this.IsClosed() {
		return
	}

	this.closeEventMutex.Lock()
	defer this.closeEventMutex.Unlock()

	this.closeCallbacks.PushBack(closeCallback{handler, callback})
}

/**
 * 移除关闭session的回调
 */
func (this *Session) RemoveCloseCallback(handler interface{}) {
	if this.IsClosed() {
		return
	}

	this.closeEventMutex.Lock()
	defer this.closeEventMutex.Unlock()

	for i := this.closeCallbacks.Front(); i != nil; i = i.Next() {
		if i.Value.(closeCallback).Handler == handler {
			this.closeCallbacks.Remove(i)
			return
		}
	}
}

/**
 * 调用关闭session的回调
 */
func (this *Session) invokeCloseCallbacks() {
	this.closeEventMutex.Lock()
	defer this.closeEventMutex.Unlock()

	for i := this.closeCallbacks.Front(); i != nil; i = i.Next() {
		callback := i.Value.(closeCallback)
		callback.Func()
	}
}