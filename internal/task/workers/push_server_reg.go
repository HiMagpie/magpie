package workers

import (
	"time"
	"magpie/internal/rpc_client"
)

func UpPushServerInfo() {
	go func() {
		for {
			hb, err := rpc_client.RegPushServer()
			if err != nil || hb < 30 {
				// limit hb heartbeat interval (second)
				hb = 30
			}
			time.Sleep(time.Second * time.Duration(hb))
		}
	}()
}
