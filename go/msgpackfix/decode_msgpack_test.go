package msgpackfix

import (
	"log"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

var testInputs []any = []any{
	map[string]interface{}{"foo": 1, "hello": "world"},
	[]map[string]interface{}{
		{"id": 1, "attrs": map[string]interface{}{"phone": 12345}},
		{"id": 2, "attrs": map[string]interface{}{"phone": 54321}},
	},
}

func TestDecodeMsgpack(t *testing.T) {
	for i := range testInputs {
		in := testInputs[i]
		b, err := msgpack.Marshal(in)
		ShouldSuccess(t, err)

		ret, err := Fix(b)
		log.Printf("ret=%v", ret)
		ShouldSuccess(t, err)

		// buf, err := json.MarshalIndent(ret, "", "  ")
		// ShouldSuccess(t, err)
		// log.Printf(string(buf))

		log.Println()
	}
}
