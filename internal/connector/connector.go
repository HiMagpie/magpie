package connector

import (
	"fmt"
	"magpie/internal/com/errs"
	"magpie/internal/com/cfg"
)

func Start() {
	// Listen on port
	server, err := Serve("tcp", fmt.Sprintf(":%d", cfg.C.Srv.Port))
	if errs.CheckError(err) {
		return
	}

	// Accept & deal connections
	for {
		session, err := server.Accept()
		if err != nil {
			continue
		}

		go handleClient(session)
	}
}

// 1. 考虑包长度超过缓冲区最大限制
// 2. 心跳/登录/无效包(舍弃), 注意包不完整的情况处理
func handleClient(session *Session) {
	session.Receive()
}
