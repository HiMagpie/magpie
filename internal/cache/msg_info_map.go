package cache

import (
	"fmt"
"magpie/internal/com/logger"
	"magpie/internal/com/utils"
)

/**
 * 通过msg_id获取消息体内容
 */
func GetMsgInfo(msgId int64) (string, error) {
	key := GetMsgInfoKey(msgId)
	c := rc.HGet(key, fmt.Sprintf("%d", msgId))
	if c.Err() != nil {
		logger.Error("cache.get.msg.info", logger.Format(
			"err", c.Err().Error(),
			"msg_id", msgId,
			"key", key,
		))
		return "", c.Err()
	}

	return c.Val(), nil
}

/**
 * 通过msg_id获取消息体所在的map的key
 */
func GetMsgInfoKey(msgId int64) string {
	return "msg_info_" + utils.Md5AndSub(fmt.Sprintf("%d", msgId), 0, 4)
}
