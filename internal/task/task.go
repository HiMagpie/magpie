package task

import (
	"magpie/internal/task/workers"
)

/**
 * 开始消费所有的session的消息队列
 */
func Start() {
	go workers.ConsumeNotify()
	go workers.MonitorWaitAckMsg()


	go workers.UpPushServerInfo()
}
