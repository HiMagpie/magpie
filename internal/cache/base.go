package cache
import (
	"time"
	"errors"
)

func PushToQueue(queue, item string) (int64, error) {
	c := rc.LPush(queue, item)
	return c.Result()
}

func PopFromQueue(queue string) (string, error) {
	c := rc.RPop(queue)
	return c.Result()
}

func BPopFromQueue(queue string, timeout time.Duration) (string, error) {
	c := rc.BRPop(timeout, queue)
	res, err := c.Result()
	if err != nil {
		return "", err
	}

	if len(res) < 2 {
		return "", errors.New("Bpop from queue encounter invalid value.")
	}

	return res[1], nil
}

func GetQueueLen(queue string) (int64, error) {
	c := rc.LLen(queue)
	return c.Result()
}
