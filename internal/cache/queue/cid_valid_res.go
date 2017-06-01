package queue

import (
	"time"
	"encoding/json"
	"os"
	"gopkg.in/redis.v3"
	"magpie/internal/cache"
	"magpie/internal/com/logger"
)

const (
	KEY_QUEUE_CID_VALID_RES_PREFIX = "queue_cid_valid_res_"
)

type CidValidResQueue struct {
	Cid   string `json:"cid"`
	Valid int32 `json:"valid"`
}

func NewCidValidResQueue() *CidValidResQueue {
	return new(CidValidResQueue)
}

func (this CidValidResQueue) getQueue() string {
	hostname, _ := os.Hostname()
	return KEY_QUEUE_CID_VALID_RES_PREFIX + hostname
}

func (this *CidValidResQueue) BPopFromQueue() error {
	item, err := cache.BPopFromQueue(this.getQueue(), time.Second * time.Duration(3))
	if err != nil {
		goto err_bpop_cid_valid_res
	}

	logger.Debug("bpop.from.cid.valid.res", logger.Format("item", item))

	err = json.Unmarshal([]byte(item), this)
	if err != nil {
		goto err_bpop_cid_valid_res
	}

	// 操作失败
	err_bpop_cid_valid_res:
	if err != nil {
		if err != redis.Nil {
			logger.Error("pop.from.cid.valid.res", logger.Format("err", err.Error()))
		}
		return err
	}

	return nil
}
