package workers

import (
	"time"
	"magpie/internal/cache/queue"
	"magpie/internal/cache"
	"magpie/internal/com/logger"
	"magpie/internal/com/utils"
	"magpie/internal/model"
	"encoding/json"
	"magpie/internal/protocol/protos"
	"github.com/golang/protobuf/proto"
	"magpie/internal/connector"
	"errors"
	"magpie/internal/protocol"
)

/**
 * 接收有新消息的提醒,将对应cid的queue进行消息推送
 */
func ConsumeNotify() {

	go func() {
		for {
			cid, err := queue.BPopFromServerNotifyQueue()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			if cid == "" {
				logger.Error("bpop.from.server.notify.empty.cid", logger.Format(
					"queue", queue.GetServerNotifyQueue(),
				))
				time.Sleep(time.Second)
				continue
			}

			logger.Debug("bpop.from.server.notify.empty.cid", logger.Format(
				"cid", cid,
			))

			waitAck, err := cache.IsExistWaitMsgAck(cid);
			if err != nil {
				logger.Error("is.exist.wait.msg.ack", logger.Format(
					"err", err.Error(),
				))
				continue
			}

			if waitAck {
				dealWhenWaitingAck(cid)
			}else {
				dealWhenNotWaitingAck(cid)
			}
		}
	}()
}

func dealWhenWaitingAck(cid string) {
	queue.PushToServerNotifyQueue(cid)
}

func dealWhenNotWaitingAck(cid string) {
	queueName := queue.GetCidQueueName(cid)
	msgIdStr, err := cache.PopFromQueue(queueName)
	msgId, err := utils.AtoInt64(msgIdStr)
	err = sendMsg(cid, msgId)
	if err != nil {
		logger.Error("send.msg.by.cid.and.msgid", logger.Format(
			"err", err.Error(),
			"cid", cid,
			"msg_id", msgId,
		))

		return
	}

	// 消息状态更新
	queue.NewMsgStatQueue(msgId, cid).SetMsgStatPushAck()
}

func sendMsg(cid string, msgId int64) error {
	// Get msg detail  by msg_id
	infoStr, err := cache.GetMsgInfo(msgId)
	logger.Debug("send.msg", logger.Format("cid", cid, "msg_id", msgId, "info", infoStr))

	msgEntity := models.NewMsgEntity()
	err = json.Unmarshal([]byte(infoStr), msgEntity)
	if err != nil {
		logger.Error("trans.to.payload", logger.Format(
			"err", err.Error(),
			"msg_str", infoStr,
		))

		// push msg_id notification back to queue
		queue.PushToServerNotifyQueue(cid)
		return err
	}

	// Organize msg & pkg to send
	msg := new(protos.Msg)
	msg.MsgId = proto.Uint64(uint64(msgEntity.MsgId))
	msg.Ring = proto.Bool(msgEntity.Ring)
	msg.Vibrate = proto.Bool(msgEntity.Vibrate)
	msg.Cleanable = proto.Bool(msgEntity.Cleanable)
	msg.Trans = proto.Int32(int32(msgEntity.Trans))
	msg.Title = proto.String(msgEntity.Title)
	msg.Text = proto.String(msgEntity.Text)
	msg.Logo = proto.String(msgEntity.Logo)
	msg.Url = proto.String(msgEntity.Url)
	//msg.Ctime = proto.Int32(msgEntity.Ctime)
	msgbytes, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("send.msg", logger.Format("err", err.Error(), "j_info", infoStr))
		return errors.New("fail to proto marshal msg.")
	}

	session, err := connector.GetSessionByCid(cid)
	if session == nil {
		//queue.PushBackToCidMq(cid, fmt.Sprintf("%d", msgId))
		cache.DelWaitMsgACKByCid(cid)
		return errors.New("fail to get session by cid.")
	}
	p := protocol.GenPkg(protocol.TYP_MSG, msgbytes)
	err = session.Send(p)
	if err != nil {
		logger.Error("send.msg", logger.Format("err", err.Error()))
		return err
	}

	// Set waiting msg ack
	msgWaitAck := new(cache.MsgWaitAckCache)
	msgWaitAck.MsgId = msgId
	msgWaitAck.Cid = cid
	msgWaitAck.Time = time.Now().Unix()
	ackBytes, err := json.Marshal(msgWaitAck)
	if err != nil {
		logger.Error("msg.consumer.mash.ack.byte", logger.Format(
			"err", err.Error(),
			"msg", msg,
		))

		queue.PushToServerNotifyQueue(cid)
		return err
	}
	cache.SetWaitMsgACK(cid, string(ackBytes))
	return nil
}
