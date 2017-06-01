package main

import (
	"fmt"
	"net"
	"encoding/json"
	"math/rand"
	"magpie/internal/com/logger"
	"magpie/internal/com/utils/httptool"
	"magpie/internal/com/errs"
	"magpie/internal/protocol"
	"time"
	"magpie/internal/com/cfg"
	"magpie/internal/connector"
	"magpie/internal/protocol/protos"
	"github.com/golang/protobuf/proto"
)

// 用来控制每次只有一个协程在操作conn
var lock = make(chan net.Conn, 1)

type RegistrationData struct {
	Cid     string `json:"cid"`
	Servers []string `json:"servers"`
}

/**
 * Housekeeper checkin api 返回值结构
 */
type CheckInRet struct {
	Code int `json:"code"`
	Data *RegistrationData `json:"data"`
	Msg  string `json:"msg"`
}

func main() {
	fmt.Println("Hello client.")

	// 发送请求
	regUrl := "http://127.0.0.1:7780/registration/check-in"
	p := httptool.Params{
		"app_id": "8de4714b0454797ef0a403b1d50f0f20",
		"app_secret": "N2U2ZTYyMTVhNmExYmVhMjMzNTBhYzg5",
		"client_key": fmt.Sprintf("client_key_auto_generate_%d", rand.Intn(10000000)),
		"ver": "1.1.1",
		"os": "Android 6.0",
	}

	sbs, err := httptool.Get(regUrl, p)
	if err != nil {
		logger.Error("doregistration", logger.Format(
			"err", err.Error(),
			"url", regUrl,
			"params", p,
		))

		return
	}

	// json解析
	logger.Debug("doregistration", logger.Data{"response": string(sbs)})
	ret := new(CheckInRet)
	err = json.Unmarshal(sbs, ret)
	if err != nil {
		logger.Error("doregistration", logger.Format(
			"err", err.Error(),
			"resonse", string(sbs),
		))
		return
	}

	// @TODO cid
	//ret.Data.Cid = "9860391f747fe438952ff67ff5bc6d02"

	logger.Debug("doregistration", logger.Format(
		"params", p,
		"url", regUrl,
		"response", string(sbs),
	))
	if err != nil {
		logger.Debug("main", logger.Format(
			"msg", "client failed to register to housekeeper.",
			"err", err.Error()))
		return
	}

	server := ret.Data.Servers[len(ret.Data.Servers) - 1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if errs.CheckError(err) {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if errs.CheckError(err) {
		return
	}

	lock <- conn

	// 发送, 如果并发同时向一个conn中写入bytes, 那么会抢占式地使得每个包的字符顺序被打乱,
	// 所以必须保证每次只有一个进程或协程在读/写一个conn
	tmp := <-lock
	go login(tmp, ret.Data.Cid)
	tmp = <-lock

	// 心跳
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(cfg.C.Srv.HbInterval))
			p, err := protocol.NewHbPkg()
			if err != nil {
				logger.Error("hb", logger.Format("err", err.Error()))
				continue
			}
			conn.Write(p)
			fmt.Println("hb: ", p)
		}
	}()

	decoder := connector.NewPbDecoder(conn)

	// 读取返回结果
	fmt.Println("start to read.")
	defer func() {
		conn.Close()
	}()

	for {
		p, err := decoder.Decode()
		if err != nil {
			logger.Error("client.decode", logger.Format("err", err.Error()))
			conn.Close()
			break
		}

		switch p.Type{
		case protocol.TYP_MSG:
			msg := new(protos.Msg)
			err := proto.Unmarshal(p.Body, msg)
			if err != nil {
				logger.Error("decode.msg", logger.Format("err", err.Error()))
				continue
			}

			logger.Info("client.receiv", logger.Format("text", msg.GetText(), "msg_id", msg.GetMsgId()))
			ack := new(protos.MsgAck)
			ack.MsgId = msg.MsgId
			ackBytes, err := proto.Marshal(ack)
			if err != nil {
				logger.Error("encode.ack", logger.Format("err", err.Error(), "pkg", *p))
				continue
			}
			ap, err := protocol.NewPkgBytes(protocol.TYP_MSG_ACK, ackBytes)
			if err != nil {
				logger.Error("encode.ack.pkg", logger.Format("err", err.Error(), "ack_pkg", ap))
				continue
			}

			num, err := conn.Write(ap)
			logger.Info("ack", logger.Format("num", num, "err", err, "ack_bytes", ap))
		default:
			logger.Error("decode.type", logger.Format("v", p.Body))
		}
	}
}

// 发送
func login(conn net.Conn, cid string) {
	fmt.Println("start to login: ", conn)

	req, err := protocol.NewLoginPkg(cid)
	fmt.Println(err)
	for i := 0; i < len(req); i++ {
		_, err := conn.Write(req[i:i + 1])
		if errs.CheckError(err) {
			return
		}
	}

	fmt.Println("finish login.", conn)
	lock <- conn
}