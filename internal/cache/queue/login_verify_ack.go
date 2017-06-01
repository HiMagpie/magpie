package queue

import (
	"fmt"
	"encoding/json"
	"os"
	"magpie/internal/cache"
)

const (
	QUEUE_LOGIN_VERIFY_ACK_PREFIX string = "queue_login_verify_ack_"
)

type LoginVerifyAckQueue struct {
	Cid string `json:"cid"`
	Ok  bool `json:"ok"`
}

func (this *LoginVerifyAckQueue) getQueue() string {
	hostname, _ := os.Hostname()
	return QUEUE_LOGIN_VERIFY_ACK_PREFIX + hostname
}

func (this *LoginVerifyAckQueue) PushToQueue() error {
	tmpB, err := json.Marshal(this)
	if err != nil {
		fmt.Println("Json mashal login verify ack fail, ", err.Error())
		return err
	}
	cache.PushToQueue(this.getQueue(), string(tmpB))
	return nil
}

func (this *LoginVerifyAckQueue)PopFromQueue() (error) {
	tmp, _ := cache.PopFromQueue(this.getQueue())
	err := json.Unmarshal([]byte(tmp), this)
	if err != nil {
		return err
	}

	return nil
}
