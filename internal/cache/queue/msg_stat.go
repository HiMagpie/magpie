package queue

import (
	"encoding/json"
	"magpie/internal/cache"
)

const (
	MSG_STATUS_QUEUE = "queue_msg_status" // 在HiMagpie中的各种处理,都会相应地把消息的最新状态push到队列

	MSG_STAT_PUSH_ACK = 3 // 刚发送, 等待ACK
	MSG_STAT_PUSH_SUCC = 4 // 发送成功(收到消息ACK)
	MSG_STAT_OFFLINE = 5 // 客户端不在线, 进入离线状态
	MSG_STAT_OVERDUE = 6 // 过期失效
	MSG_STAT_FAIL = 7 // 失败
)

type MsgStatQueue struct {
	MsgId  int64 `json:"msg_id"`
	Cid    string `json:"cid"`
	Status int `json:"status"`
}

func NewMsgStatQueue(msgId int64, cid string) *MsgStatQueue {
	return &MsgStatQueue{
		MsgId: msgId,
		Cid: cid,
	}
}

func (this *MsgStatQueue) SetMsgStatPushAck() {
	this.Status = MSG_STAT_PUSH_ACK
	this.setMsgStat()
}

func (this *MsgStatQueue) SetMsgStatPushSucc() {
	this.Status = MSG_STAT_PUSH_SUCC
	this.setMsgStat()
}

func (this *MsgStatQueue) SetMsgStatOffline() {
	this.Status = MSG_STAT_OFFLINE
	this.setMsgStat()
}

func (this *MsgStatQueue) SetMsgStatOverdue() {
	this.Status = MSG_STAT_OVERDUE
	this.setMsgStat()
}

func (this *MsgStatQueue) SetMsgStatFail() {
	this.Status = MSG_STAT_FAIL
	this.setMsgStat()
}

func (this *MsgStatQueue) setMsgStat() {
	b, err := json.Marshal(this)
	if err != nil {
		return
	}
	cache.PushToQueue(MSG_STATUS_QUEUE, string(b))
}



