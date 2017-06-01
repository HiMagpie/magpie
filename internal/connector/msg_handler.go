package connector

import (
	"magpie/internal/com/logger"
	"github.com/golang/protobuf/proto"
	"magpie/internal/protocol"
	"magpie/internal/protocol/protos"
	"magpie/internal/cache/queue"
	"sync/atomic"
	"time"
	"magpie/internal/cache"
	"magpie/internal/model"
	"encoding/json"
)

type MsgHandler struct {
	encoder *PbEncoder
	decoder *PbDecoder
	session *Session
}

func NewMsgHandler(session *Session) *MsgHandler {
	return &MsgHandler{
		encoder:NewPbEncoder(session.Conn()),
		decoder:NewPbDecoder(session.Conn()),
		session: session,
	}
}

// Receive data from client
func (this *MsgHandler) ReceiveData() {
	go func() {
		for {
			p, err := this.decoder.Decode()
			if err != nil {
				this.session.Close()
				logger.Error("receive.data", logger.Format("err", err.Error()))
				return
			}

			logger.Info("pkg", "pkg", *p)
			switch p.Type{
			case protocol.TYP_LOGIN:
				this.handleLogin(p)

			case protocol.TYP_HEARTBEAT:
				this.handleHeartbeat(p)

			case protocol.TYP_MSG_ACK:
				this.handleMsgAck(p)

			default:
				logger.Error("decode.type", logger.Format("v", p.Body))
			}
		}
	}()
}

func (this *MsgHandler) handleHeartbeat(p *protocol.Pkg) error {
	// update last heartbeat time (magpie has a task to scan all sessions)
	atomic.StoreInt64(&this.session.HeartbeatTime, time.Now().Unix())
	return nil
}

func (this *MsgHandler) handleLogin(p *protocol.Pkg) error {
	login := new(protos.Login)
	err := proto.Unmarshal(p.Body, login)
	if err != nil {
		logger.Error("login", logger.Format("err", err.Error(), "pkg", p))
		return err
	}

	// set cid to session-validated queue
	cidValidQueue := queue.NewCidValidQueue()
	cidValidQueue.Cid = *login.Cid
	cidValidQueue.PushToQueue()
	logger.Debug("login", logger.Format("cid", *login.Cid))

	// bind relations between cid and session
	atomic.StoreInt32(&this.session.validFlag, VALID_FLAG_DEALING)
	AddCidSessionIdRel(*login.Cid, this.session.Id())

	return nil
}

func (this *MsgHandler) handleMsgAck(p *protocol.Pkg) error {
	ack := new(protos.MsgAck)
	err := proto.Unmarshal(p.Body, ack)
	if err != nil {
		logger.Error("msg.ack", logger.Format("err", err.Error(), "bytes", p.Body))
		return err
	}

	infoStr, err := cache.GetMsgInfo(int64(*ack.MsgId))
	msg := models.NewMsgEntity()
	err = json.Unmarshal([]byte(infoStr), msg)
	if err != nil {
		logger.Error("msg.ack", logger.Format("err", err.Error()))
		return err
	}

	cid, err := GetCidBySessionId(this.session.Id())
	if err != nil {
		logger.Error("msg.ack", logger.Format("err", err.Error()))
		return err
	}

	msgStr, err := cache.GetWaitMsgAckByCid(cid)
	if err != nil {
		logger.Error("push.ack", logger.Format("err", err.Error(), "cid", cid))
		return err
	}

	waitAck := cache.NewMsgWaitAckCache()
	err = json.Unmarshal([]byte(msgStr), waitAck)
	if err != nil {
		logger.Error("push.ack", logger.Format("err:", err.Error(), ))
		return err
	}

	// 当ack的msg_id < 当前等待ack的msg_id则忽略
	if int64(*ack.MsgId) < waitAck.MsgId {
		return nil
	}

	queue.NewMsgStatQueue(waitAck.MsgId, waitAck.Cid).SetMsgStatPushSucc()
	cache.DelWaitMsgACKByCid(cid)
	logger.Debug("msg.ack", logger.Format("cid", cid, "msg_id", *ack.MsgId))
	return nil
}

// Send data to client
func (this *MsgHandler) SendData(p *protocol.Pkg) error {
	pkgBytes, err := protocol.NewPkgBytes(int(p.Type), p.Body)
	if err != nil {
		logger.Error("msg.send.data", logger.Format("err", err.Error(), "pkg", *p))
		return err
	}

	return this.encoder.Encode(pkgBytes)
}

