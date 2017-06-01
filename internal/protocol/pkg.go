package protocol

import (
	"time"
	"magpie/internal/com/utils"
	"encoding/binary"
	"magpie/internal/protocol/protos"
	"github.com/golang/protobuf/proto"
	"magpie/internal/com/logger"
)

const (
// 每个包最长限制1000byte(参考极光推送)
	PACKAGE_LIMIT int = 1000

	HEART_BEAT_TIME = 5 * 60 // 心跳包5mins

	HEADER_LEN int = 3 // 3 bytes
)

// | 0000 0001 | 0000 0000 0000 0001 | 0000 0000 ...
type Pkg struct {
	Type    uint8
	BodyLen uint16
	Body    []byte
}

func NewPkg() *Pkg {
	return new(Pkg)
}

func GenPkg(typ uint8, body []byte) *Pkg {
	return &Pkg{
		Type:typ,
		BodyLen: uint16(len(body)),
		Body:body,
	}
}

func NewHeaderBytes(typ, bodyLen int) ([]byte, error) {
	header := make([]byte, HEADER_LEN)

	// type
	tmpBytes, err := utils.Uint8ToBytes(uint8(typ))
	if err != nil {
		return nil, err
	}
	copy(header[0:], tmpBytes) // 1 byte - type

	// body_len
	tmpBytes = make([]byte, 2) // 2 bytes - bodyLen
	binary.BigEndian.PutUint16(tmpBytes, uint16(bodyLen))
	copy(header[1:], tmpBytes)

	return header, nil
}

// typ - pkg type
// body - pkg body (bytes, eg: pb content)
func NewPkgBytes(typ int, body []byte) ([]byte, error) {
	// generate pkg header bytes
	header, err := NewHeaderBytes(typ, len(body))
	if err != nil {
		logger.Error("header.new", logger.Format("err", err.Error()))
		return nil, err
	}

	return append(header, body...), nil
}

// generate a login pkg
func NewLoginPkg(cid string) ([]byte, error) {
	// login pkg's body bytes
	login := new(protos.Login)
	login.Cid = proto.String(cid)
	body, err := proto.Marshal(login)
	if err != nil {
		logger.Error("pkg.new.login", logger.Format("err", err.Error()))
		return nil, err
	}

	return NewPkgBytes(TYP_LOGIN, body)
}

func NewHbPkg() ([]byte, error) {
	return NewPkgBytes(TYP_HEARTBEAT, nil)
}

/**
 * 通过时间判断一个Session是否可用
 */
func IsSessionAvailable(lastTime int64) bool {
	return lastTime >= time.Now().Unix() - HEART_BEAT_TIME
}

