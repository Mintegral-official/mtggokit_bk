package streamer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/vmihailenco/msgpack"
	"testing"
)

type MsgpackCodec struct{}

// Encode encodes an object into slice of bytes.
func (c MsgpackCodec) Encode(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	//enc.UseJSONTag(true)
	err := enc.Encode(i)
	return buf.Bytes(), err
}

// Decode decodes an object from slice of bytes.
func (c MsgpackCodec) Decode(data []byte, i interface{}) error {
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	//dec.UseJSONTag(true)
	err := dec.Decode(i)
	return err
}

//func TestMsgCodec(t *testing.T) {
//	base := &BaseInfo{
//		Name:     "bifrost_streamer_example",
//		Progress: 5,
//		Data: map[container.MapKey]interface{}{
//			container.StrKey("1"): 1,
//			container.StrKey("2"): 4,
//		},
//	}
//
//	codec := MsgpackCodec{}
//
//	data, err := codec.Encode(base)
//	if err != nil {
//		fmt.Println("Encode error, ", err.Error())
//		return
//	}
//	fmt.Println("Encode succ: ", len(data))
//	n := &BaseInfo{}
//	err = codec.Decode(data, n)
//	if err != nil {
//		fmt.Println("Decode error, ", err.Error())
//		return
//	}
//
//	fmt.Println("Decode succ, ", n)
//}

func init() {
	gob.Register(&container.Int64Key{})
	gob.Register(&container.StringKey{})
}

func TestGobCodec(t *testing.T) {
	base := &BaseInfo{
		Name:     "bifrost_streamer_example",
		Progress: 5,
		Data: map[container.MapKey]interface{}{
			container.StrKey("1"): 1,
			container.StrKey("2"): 4,
		},
	}

	codec := GobCodec{}

	data, err := codec.Encode(base)
	if err != nil {
		fmt.Println("Encode error, ", err.Error())
		return
	}
	fmt.Println("Encode succ: ", len(data))
	n := &BaseInfo{}
	err = codec.Decode(data, n)
	if err != nil {
		fmt.Println("Decode error, ", err.Error())
		return
	}

	fmt.Println("Decode succ, ", n)
}
