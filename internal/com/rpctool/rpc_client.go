package rpctool

import (
	"sync"
	"net/rpc"
	"time"
	"github.com/toolkits/net"
	"math"
)

type RpcClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	Server    string
	Timeout   time.Duration
}

// Close deals
func (this *RpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

// Ensure rpc client's connection
// Retry to connect rpc server where client is broken
func (this *RpcClient) ensureConn() {
	var err error
	retry := 1

	for {
		if this.rpcClient != nil {
			return
		}

		this.rpcClient, err = net.JsonRpcClient("tcp", this.Server, this.Timeout)
		if err == nil {
			return
		}

		// Sleep when init client fail
		if retry > 8 {
			retry = 1
		}
		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
		retry++
	}
}

// Invoke rpc call
func (this *RpcClient) Call(method string, args interface{}, reply interface{}) error {
	this.Lock()
	defer this.Unlock()

	this.ensureConn()

	// Add basic max timeout limitation
	timeout := time.Duration(50 * time.Second)
	done := make(chan error)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		this.close()

	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}

