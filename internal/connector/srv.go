package connector

import (
	"net"
	"time"
	"sync"
	"sync/atomic"
	"magpie/internal/cache"
	"magpie/internal/com/logger"
	"magpie/internal/com/cfg"
	"magpie/internal/com/utils"
)

var (
	server *Server // 推送服务器实例
)

const sessionMapNum = 100

type sessionMap struct {
	sync.RWMutex
	sessions map[uint64]*Session
}

/**
 * 服务器结构体
 */
type Server struct {
	listener     net.Listener

	// 存储会话
	maxSessionId uint64
	sessionMaps  [sessionMapNum]sessionMap

	// 服务器开启和关闭相关
	stopOnce     sync.Once
	stopWait     sync.WaitGroup
}

func NewServer(listener net.Listener) *Server {
	server := &Server{
		listener:  listener,
	}

	for i := 0; i < sessionMapNum; i++ {
		server.sessionMaps[i].sessions = make(map[uint64]*Session)
	}

	// 启动定时清理失效Session的goroutine
	server.scavengeInvalidSessions()
	return server
}


func (this *Server) Listener() net.Listener {
	return this.listener
}

func (this *Server) Accept() (*Session, error) {
	var tempDelay time.Duration
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return nil, err
		}
		tempDelay = 0
		return this.newSession(conn), nil
	}
}

func (this *Server) Stop() {
	this.stopOnce.Do(func() {
		this.listener.Close()
		this.closeSessions()
		this.stopWait.Wait()
	})
}

/**
 * 通过sessionId获取session实例
 */
func (this *Server) GetSession(sessionId uint64) *Session {
	smap := &this.sessionMaps[sessionId % sessionMapNum]
	smap.RLock()
	defer smap.RUnlock()

	session, _ := smap.sessions[sessionId]
	return session
}

/**
 * 通过conn创建一个session
 */
func (this *Server) newSession(conn net.Conn) *Session {
	session := NewSession(conn, nil)
	session.handler = NewMsgHandler(session)
	this.putSession(session)

	// ValidTimeout秒之后, 如果仍然为发送登录验证, 则关闭连接
	this.setValidClock(session)
	return session
}

func (this *Server) setValidClock(sess *Session) {
	go func() {
		// ValidTimeout秒之后, 如果仍然为发送登录验证, 则关闭连接
		time.Sleep(time.Second * time.Duration(cfg.C.Srv.ValidTimeout))

		if atomic.LoadInt32(&sess.validFlag) == VALID_FLAT_NOT_START {
			sess.Close()
		}
	}()
}

/**
 * 将session设置进总得sessionMaps里面
 */
func (this *Server) putSession(session *Session) {
	smap := &this.sessionMaps[session.id % sessionMapNum]
	smap.Lock()
	defer smap.Unlock()

	// 移除Session实例
	session.AddCloseCallback(this, func() {
		this.delSession(session)
	})

	smap.sessions[session.id] = session
	this.stopWait.Add(1)
}

/**
 * 删除sessionMaps里面的某个session 以及 cid和Session的映射关系
 */
func (this *Server) delSession(session *Session) {
	smap := &this.sessionMaps[session.id % sessionMapNum]
	smap.Lock()
	defer smap.Unlock()

	session.RemoveCloseCallback(this)
	delete(smap.sessions, session.id)
	this.stopWait.Done()
}

/**
 * 复制某个sessionMap的所有session
 */
func (this *Server) copySessions(n int) []*Session {
	smap := &this.sessionMaps[n]
	smap.Lock()
	defer smap.Unlock()

	sessions := make([]*Session, 0, len(smap.sessions))
	for _, session := range smap.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

/**
 * 关闭所有session
 */
func (this *Server) closeSessions() {
	// 复制session防止死锁
	for i := 0; i < sessionMapNum; i++ {
		sessions := this.copySessions(i)
		for _, session := range sessions {
			session.Close()
		}
	}
}

/**
 * 清理double time没有心跳包到达的Session, 认为其已失效
 * Server起来后,开始定时执行清理
 */
func (this *Server ) scavengeInvalidSessions() {

	go func() {
		for {
			cids := GetAllCids()
			invalidCids := make([]string, 0)
			for _, cid := range cids {
				session, err := GetSessionByCid(cid)
				if err != nil {
					invalidCids = append(invalidCids, cid)
					continue
				}

				if session == nil {
					RemoveCidRel(cid)
					logger.Debug("scave.session.nil", logger.Format(
						"scave_cid", cid,
					))
					invalidCids = append(invalidCids, cid)
					continue
				}

				// 处理失效的Session
				if atomic.LoadInt64(&session.HeartbeatTime) < (time.Now().Unix() - int64(cfg.C.Srv.HbInterval) * 3) {
					logger.Info("scavenge.invalid.session", logger.Format(
						"cid", cid,
						"last_hb_time", utils.FormatTimeStamp(session.HeartbeatTime),
					))

					invalidCids = append(invalidCids, cid)
					this.RemoveCidRelAndCloseSession(cid)
				}
			}

			cache.DealClosedCids(invalidCids)
			time.Sleep(time.Second * time.Duration(cfg.C.Srv.ScavengeInterval))
		}
	}()
}

/**
 * 移除Cid和SessionId映射关系 & 关闭Session
 */
func (this *Server) RemoveCidRelAndCloseSession(cid string) {
	logger.Debug("remove.cid.rel.and.close", logger.Format(
		"cid", cid,
	))
	session, err := GetSessionByCid(cid)
	if err != nil {
		return
	}

	RemoveCidRel(cid)
	session.Close()
}

/**
 * 开启服务监听客户端连接
 */
func Serve(network, address string) (*Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	server = NewServer(listener)
	return server, nil
}

/**
 * 获取服务器实例
 */
func GetServer() *Server {
	return server
}