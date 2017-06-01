package cache

import (
	"magpie/internal/com/logger"
)

const (
	KEY_MAP_CID_SERVER_NOTIFY_QUEUE = "hmap_cid_server_notify_queue"
)

func GetCidPushServerQueueKey(cid string) string {
	return KEY_MAP_CID_SERVER_NOTIFY_QUEUE + "_" + cid[0:3]
}

func SetCidServerQueue(cid string, queue string) error {
	c := rc.HSet(GetCidPushServerQueueKey(cid), cid, queue)
	return c.Err()
}

func RemCidServerQueue(cids []string) error {
	for _, c := range cids {
		c := rc.HDel(GetCidPushServerQueueKey(c), c)
		if c.Err() != nil {
			logger.Error("cid.srv.rm", logger.Format("err", c.Err().Error(), "cid", c))
		}
	}

	return nil
}

func DealClosedCids(cids []string) error {
	if len(cids) == 0 {
		return nil
	}

	err := RemCidServerQueue(cids)
	if err != nil {
		logger.Error("deal.closed.cids", logger.Format(
			"err", err.Error(),
			"cids", cids))
		return err
	}

	return nil
}