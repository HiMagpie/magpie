package connector

import (
	"time"
	"gopkg.in/redis.v3"
	"magpie/internal/cache/queue"
	"magpie/internal/cache"
"magpie/internal/com/logger"
)

func init() {
	go startDealCidValidQueue()
}

/**
 * 处理验证cid的结果队列
 */
func startDealCidValidQueue() {
	cidValidResQueue := queue.NewCidValidResQueue()

	for {
		err := cidValidResQueue.BPopFromQueue()
		if err != nil {
			if err != redis.Nil {
				logger.Error("deal.ci.valid.queue", logger.Format("err", err.Error()))
			}
			time.Sleep(time.Second)
			continue
		}
		logger.Info("cid.valid", logger.Format("res", cidValidResQueue))

		// 1.1 Cid验证通过
		if cidValidResQueue.Valid == VALID_FLAG_SUCC {
			cache.SetCidServerQueue(cidValidResQueue.Cid, queue.GetServerNotifyQueue())
			queue.NotifyServerByCid(cidValidResQueue.Cid)

			// 设置登录验证完成ACK
			loginAckQueue := new(queue.LoginVerifyAckQueue)
			loginAckQueue.Cid = cidValidResQueue.Cid
			loginAckQueue.Ok = true
			loginAckQueue.PushToQueue()
			continue
		}

		// 2.1 Cid验证失败
		sess, err := GetSessionByCid(cidValidResQueue.Cid)
		if err != nil {
			logger.Error("deal.cid.valid.queue", logger.Format(
				"err", err.Error(),
				"cid", cidValidResQueue.Cid,
				"valid", cidValidResQueue.Valid,
			))
			continue
		}

		// 2.2 关闭连接 & 删除cid和session的映射关系
		sess.Close()
		RemoveCidRel(cidValidResQueue.Cid)
	}
}
