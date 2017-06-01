package main

import (
	"runtime"
	"magpie/internal/rpc_client"
	"magpie/internal/connector"
	"magpie/internal/task"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Run TCP server and handle connections
	go connector.Start()

	// Register push server info to Housekeeper
	rpc_client.RegPushServer()

	// Start background tasks, eg: push, sync status...etc
	task.Start()

	select {}
}

