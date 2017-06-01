package connector

import (
	"io"
)

// pb encoder
type PbEncoder struct {
	w io.Writer
}

func NewPbEncoder(writer io.Writer) *PbEncoder {
	e := new(PbEncoder)
	e.w = writer
	return e
}

func (this *PbEncoder) Encode(d []byte) error {
	_, err := this.w.Write(d)
	return err
}