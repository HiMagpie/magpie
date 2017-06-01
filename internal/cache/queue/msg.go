package queue

import (
	"os"
	"time"
	"magpie/internal/cache"
"magpie/internal/com/logger"
)

func PushBackToCidMq(cid, msg string) error {
	queue := GetCidQueueName(cid)
	_, err := cache.PushToQueue(queue, msg)
	return err
}

func CountCidMqNum(cid string) (int64, error) {
	return cache.GetQueueLen(GetCidQueueName(cid))
}

/**
 * 组装cid的队列名称
 */
func GetCidQueueName(cid string) string {
	return "mq_" + cid
}

func NotifyServerByCid(cid string) {
	msgNum, err := CountCidMqNum(cid)
	if err != nil {
		logger.Error("count.cid.mq.num", logger.Format(
			"err", err.Error(),
			"msg_num", msgNum))
	}

	for i := 0; i < int(msgNum); i++ {
		PushToServerNotifyQueue(cid)
	}
}

func PushToServerNotifyQueue(cid string) (int64, error) {
	return cache.PushToQueue(GetServerNotifyQueue(), cid)
}

func GetServerNotifyQueue() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "queue_server_notify_localhost"
	}

	return "queue_server_notify_" + hostname
}

func BPopFromServerNotifyQueue() (string, error) {
	queue := GetServerNotifyQueue()
	res, err := cache.BPopFromQueue(queue, time.Second * time.Duration(3))
	if err != nil {
		return "", err
	}

	return res, nil
}


