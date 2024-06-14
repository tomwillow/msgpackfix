package msgpackfix

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"unsafe"
)

const (
	FIXARRAY = iota
)

const (
	MSGPACK_U8  = 0xcc
	MSGPACK_U16 = 0xcd
	MSGPACK_U32 = 0xce
	MSGPACK_U64 = 0xcf
)

func FixHexString(hexString string) (ret any, err error) {
	buf, err := hex.DecodeString(hexString)
	if err != nil {
		return
	}
	return Fix(buf)
}

func pickUpValue[T uint8 | uint16 | uint32 | uint64](buf []byte, cvtFunc func([]byte) any) (remainBuf []byte, val any, err error) {
	dataSize := int(unsafe.Sizeof(T(0)))
	if len(buf)-1 < dataSize {
		err = fmt.Errorf("data not enough: want=%d, got=%d", dataSize, len(buf)-1)
		return
	}
	val = cvtFunc(buf[1 : 1+dataSize])
	remainBuf = buf[1+dataSize:]
	return
}

func pickUp(buf []byte) (remainBuf []byte, val any, err error) {
	remainBuf = buf
	if len(buf) == 0 {
		err = io.EOF
		return
	}

	// positive fixint
	if buf[0]>>7 == 0 {
		val = int8(buf[0])
		remainBuf = buf[1:]
		return
	}

	// fixmap
	// 0x8X(0x80 - 0x8F)
	if buf[0]>>4 == 0b1000 {
		return pickUpMap(buf)
	}

	// fixarray
	// 0x9X(0x90 - 0x9F)
	if buf[0]>>4 == 0b1001 {
		return pickUpArray(buf)
	}

	// fixstr
	// 0xA0 - 0xBF
	if buf[0]>>5 == 0b101 {
		return pickUpFixStr(buf)
	}

	if buf[0] == MSGPACK_U8 {
		return pickUpValue[uint8](buf, func(b []byte) any { return b[0] })
	}
	if buf[0] == MSGPACK_U16 {
		return pickUpValue[uint16](buf, func(b []byte) any { return binary.BigEndian.Uint16(b) })
	}
	if buf[0] == MSGPACK_U32 {
		return pickUpValue[uint32](buf, func(b []byte) any { return binary.BigEndian.Uint32(b) })
	}
	if buf[0] == MSGPACK_U64 {
		return pickUpValue[uint64](buf, func(b []byte) any { return binary.BigEndian.Uint64(b) })
	}
	err = fmt.Errorf("invalid first byte: 0x%X", buf[0])
	return
}

func Fix(buf []byte) (val any, err error) {
	remainBuf, val, err := pickUp(buf)
	if len(remainBuf) > 0 {
		appendError(&err, fmt.Errorf("unparsed: %v", remainBuf))
	}
	return
}
