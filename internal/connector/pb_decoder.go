package connector

import (
	"io"
	"magpie/internal/protocol"
	"magpie/internal/com/errs"
	"magpie/internal/com/utils"
	"encoding/binary"
"magpie/internal/com/logger"
)

/**
 * pb decoder
 */
type PbDecoder struct {
	r        io.Reader
	cache    []byte //解码时使用的缓冲区
	cacheLen int    //缓冲区中有效的字节长度
}

func NewPbDecoder(reader io.Reader) *PbDecoder {
	d := new(PbDecoder)
	d.r = reader
	d.cacheLen = 0
	d.cache = make([]byte, protocol.PACKAGE_LIMIT)
	return d
}

func (this *PbDecoder) Decode() (*protocol.Pkg, error) {
	req := make([]byte, protocol.PACKAGE_LIMIT)
	for {
		num, err := this.r.Read(req)
		if err != nil {
			logger.Error("decode", err.Error())
			break
		}

		// cope data to cache when data read.
		copy(this.cache[this.cacheLen:this.cacheLen + num], req[0:num])
		this.cacheLen += num

		// whether enough to parse pkg header to get all pkg info
		for this.cacheLen >= protocol.HEADER_LEN {

			// parse pkg TYPE and body len
			offset := 0
			pkgTyp, err := utils.BytesToUint8(this.cache[offset: offset + 1])
			offset += 1
			if err != nil {
				logger.Error("decode", err.Error())
				break
			}
			bodyLen := binary.BigEndian.Uint16(this.cache[offset:offset + 2])

			// whether enough to parse the whole pkg
			pkgLen := protocol.HEADER_LEN + int(bodyLen)
			if this.cacheLen < pkgLen {
				break
			}

			// parse pkg and reset counters
			p := new(protocol.Pkg)
			p.Type = pkgTyp
			p.BodyLen = bodyLen
			p.Body = append(p.Body, this.cache[protocol.HEADER_LEN:pkgLen]...)

			this.cacheLen -= pkgLen
			copy(this.cache[0:], this.cache[pkgLen:])

			return p, nil
		}
	}

	return nil, errs.ERR_OPERATION
}