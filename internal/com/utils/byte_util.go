package utils
import (
	"bytes"
	"encoding/binary"
)

/**
 * 装换[]byte为uint8
 */
func BytesToUint8(b []byte) (uint8, error) {
	var tmp uint8
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, &tmp)
	return tmp, err
}

/**
 * 将uint8转为bytes
 */
func Uint8ToBytes(tmp uint8) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, tmp)
	return buf.Bytes(), err
}

/**
 * 将uint16转为bytes
 */
func Uint16ToBytes(tmp uint16) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, tmp)
	return buf.Bytes(), err
}
