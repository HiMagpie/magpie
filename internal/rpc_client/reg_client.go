package rpc_client

import (
	"time"
	"magpie/internal/com/rpctool"
	"magpie/internal/com/logger"
)

// Register push server into to housekeeper
var regClient *rpctool.RpcClient

type RegArgs struct {
	Hostname string
	Ip       string
	Port     string
}

type RegResponse struct {
	PushServerHb int
}

func initRegClient() {
	regClient = &rpctool.RpcClient{
		Server:":8880",
		Timeout: time.Duration(5) * time.Second,
	}
}

// Register push server to housekeeper
func RegPushServer() (int, error) {
	args := RegArgs{
		Hostname: "localhost",
		Ip: "127.0.0.1",
		Port: "7777",
	}
	reply := RegResponse{}
	err := regClient.Call("Reg.RegPushServer", args, &reply)
	if err != nil {
		logger.Error("rpc.client.reg.push.server", logger.Format(
			"err", err.Error(),
		))
		return 0, err
	}
	logger.Info("rpc.client.reg.push.server", logger.Format("info", args, "reply", reply))
	return reply.PushServerHb, nil
}