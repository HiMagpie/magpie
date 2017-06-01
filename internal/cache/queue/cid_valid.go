package queue

import (
	"os"
	"encoding/json"
	"time"
	"magpie/internal/cache"
	"magpie/internal/com/logger"
)

const (
	KEY_QUEUE_CID_VALID = "queue_valid_cid"
)

type CidValidQueue struct {
	Cid      string `json:"cid"`
	Hostname string `json:"hostname"`
}

func NewCidValidQueue() *CidValidQueue {
	return new(CidValidQueue)
}

func (this *CidValidQueue) GetCidValidJson() (string, error) {
	cidValidJson, err := json.Marshal(this)
	if err != nil {
		return "", err
	}

	return string(cidValidJson), nil
}

func (this *CidValidQueue)PushToQueue() error {
	hostname, _ := os.Hostname()
	this.Hostname = hostname
	cidValid, err := this.GetCidValidJson()
	if err != nil {
		logger.Error("push.to.cid.valid.queue", logger.Format("err", err.Error()))
		return err
	}
	_, err = cache.PushToQueue(KEY_QUEUE_CID_VALID, cidValid)
	return err
}

func (this *CidValidQueue) PopFromQueue() error {
	item, err := cache.PopFromQueue(KEY_QUEUE_CID_VALID)
	if err != nil {
		goto err_pop_cid_valid
	}

	err = json.Unmarshal([]byte(item), this)
	if err != nil {
		goto err_pop_cid_valid
	}

	// 操作失败
	err_pop_cid_valid:
	if err != nil {
		logger.Error("pop.from.cid.valid", logger.Format("err", err.Error()))
		return err
	}

	return nil
}

func (this *CidValidQueue) BPopFromQueue() error {
	item, err := cache.BPopFromQueue(KEY_QUEUE_CID_VALID, time.Second * time.Duration(3))
	if err != nil {
		goto err_bpop_cid_valid
	}

	logger.Debug("bpop.from.cid.valid.queue", logger.Format("item", item))

	err = json.Unmarshal([]byte(item), this)
	if err != nil {
		goto err_bpop_cid_valid
	}

	// 操作失败
	err_bpop_cid_valid:
	if err != nil {
		logger.Error("pop.from.cid.valid", logger.Format("err", err.Error()))
		return err
	}

	return nil
}