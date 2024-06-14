package msgpackfix

import (
	"io"
	"log"
	"reflect"
)

func pickUpMap(buf []byte) (remainBuf []byte, val any, err error) {
	remainBuf = buf
	length := remainBuf[0] & 0b00001111
	log.Printf("fixmap: size of obj=%d", length)

	remainBuf = remainBuf[1:]
	fieldIndex := 0
	mp := make(map[any]any, 0)
	var key any
	for {
		remainBuf_, val_, err_ := pickUp(remainBuf)
		remainBuf = remainBuf_
		err = err_
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		log.Printf("index=%v, type=%v, value=%v(0x%X)", fieldIndex, reflect.TypeOf(val_), val_, val_)

		if fieldIndex%2 == 0 {
			// key
			key = val_
		} else {
			// value
			mp[key] = val_
			val = mp
		}

		fieldIndex++
	}
	return
}
