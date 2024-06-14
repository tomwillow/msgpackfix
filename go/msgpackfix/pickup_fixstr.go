package msgpackfix

import (
	"fmt"
	"log"
)

func pickUpFixStr(buf []byte) (remainBuf []byte, val any, err error) {
	remainBuf = buf
	length := int(remainBuf[0] & 0b00011111)
	log.Printf("fixstr: length=%d", length)
	remainBuf = remainBuf[1:]
	if len(remainBuf) < length {
		val = string(remainBuf)
		err = fmt.Errorf("data not enough: want=%d, got=%d", length, len(remainBuf))
		return
	}

	val = string(remainBuf[:length])
	remainBuf = remainBuf[length:]
	return
}
