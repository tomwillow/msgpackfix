package msgpackfix

import (
	"fmt"
	"log"
	"testing"
)

func ShouldSuccess(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func ShouldFail(t *testing.T, err error) {
	if err == nil {
		t.Error("should fail")
		return
	}
	fmt.Printf("expected error: %v\n", err)
}

func TestDecode(t *testing.T) {
	hexStr := "91ce0108000E"

	ret, err := FixHexString(hexStr)
	log.Printf("ret=%v", ret)
	ShouldSuccess(t, err)
}

func TestDecodeNotEnough(t *testing.T) {
	hexStr := "95ce0108000E"

	ret, err := FixHexString(hexStr)
	log.Printf("ret=%v", ret)
	ShouldFail(t, err)
}
