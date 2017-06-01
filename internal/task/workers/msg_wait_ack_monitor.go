package workers

import (
	"time"
	"encoding/json"
	"magpie/internal/cache"
	"magpie/internal/com/cfg"
)

/**
 * 定时将所有等待消息ACK的队列,判断是否需要重发消息
 */
func MonitorWaitAckMsg() {
	ackItem := new(cache.MsgWaitAckCache)

	for {
		waitAckCids := cache.GetAllWaitMsgAckCids()
		for cid, msgWaitAckStr := range waitAckCids {
			err := json.Unmarshal([]byte(msgWaitAckStr), ackItem)
			if err != nil {
				continue
			}

			// ACK未超时, 继续等
			if ackItem.Time > time.Now().Unix() - int64(cfg.C.Srv.MsgWaitAckSeconds){
				continue
			}

			// 超时则重新推送
			sendMsg(cid, ackItem.MsgId)
		}

		time.Sleep(time.Second * 2)
	}
}
