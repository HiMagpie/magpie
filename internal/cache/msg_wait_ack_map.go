package cache

import (
	"magpie/internal/com/logger"
)

var (
	WAIT_MSG_ACK_MAP = "hmap_wait_msg_ack"
)

type MsgWaitAckCache struct {
	MsgId int64 `json:"msg_id"`
	Cid   string `json:"cid"`
	Time  int64 `json:"time"` // 推送的时间,用于判断是否超过Msg Ack的时间
}

func NewMsgWaitAckCache() *MsgWaitAckCache {
	return new(MsgWaitAckCache)
}

// 发送完消息之后,将消息设置进等ack的map

/**
 * 获取所有有正在等消息ACK的cid
 * cid => MsgWaitAckCache
 */
func GetAllWaitMsgAckCids() map[string]string {
	cidMap := make(map[string]string)
	c := rc.HGetAll(WAIT_MSG_ACK_MAP)
	cids, err := c.Result()
	if err != nil {
		return cidMap
	}

	for i, v := range cids {
		if i % 2 == 1 {
			cidMap[cids[i - 1]] = v
		}
	}
	return cidMap
}

/**
 * 是否cid是否存在等待ack的消息
 */
func IsExistWaitMsgAck(cid string) (bool, error) {
	c := rc.HExists(WAIT_MSG_ACK_MAP, cid)
	return c.Result()
}

func SetWaitMsgACK(cid, msg string) (bool, error) {
	c := rc.HSet(WAIT_MSG_ACK_MAP, cid, msg)
	if c.Err() != nil {
		return false, c.Err()
	}

	return c.Val(), nil
}

func GetWaitMsgAckByCid(cid string) (string, error) {
	c := rc.HGet(WAIT_MSG_ACK_MAP, cid)
	if c.Err() != nil {
		logger.Error("get.wait.msg.ack.by.cid", logger.Format("err", c.Err().Error()))
		return "", c.Err()
	}

	return c.Val(), nil
}

func DelWaitMsgACKByCid(cid string) (int64, error) {
	return rc.HDel(WAIT_MSG_ACK_MAP, cid).Result()
}
