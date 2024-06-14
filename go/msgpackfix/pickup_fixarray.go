package msgpackfix

import (
	"fmt"
	"io"
	"log"
	"reflect"
)

func pickUpArray(buf []byte) (remainBuf []byte, val any, err error) {
	remainBuf = buf
	length := int(remainBuf[0] & 0b00001111)
	log.Printf("fixarray: len=%d", length)
	fieldIndex := 0
	defer func(gotLength *int, curErr *error) {
		if *gotLength == length {
			return
		}
		var innerErr error
		if *gotLength < length {
			innerErr = fmt.Errorf("element of array is not enough: want=%d, got=%d", length, *gotLength)
		} else {
			innerErr = fmt.Errorf("element of array is too much: want=%d, got=%d", length, *gotLength)
		}
		appendError(curErr, innerErr)
	}(&fieldIndex, &err)

	remainBuf = remainBuf[1:]
	arr := make([]any, 0)
	for {
		remainBuf_, val_, err_ := pickUp(remainBuf)
		remainBuf = remainBuf_
		err = err_
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		log.Printf("index=%v, type=%v, value=%v(0x%X)", fieldIndex, reflect.TypeOf(val_), val_, val_)

		arr = append(arr, val_)
		val = arr
		fieldIndex++
	}
	return
}
