package encoding

import (
	"reflect"
	"testing"
)

func TestEncoded_Payload_Data_Serialize_Deserialize(t *testing.T) {
	a := Encoded_Payload_Data{
		Sender: []byte("abcde"),
		Cipher: []byte("asdfgh"),
		Nonce:  []byte("123456"),
		RecipientList: [][]byte{[]byte("123"),
			[]byte("456"),
			[]byte("789"),
		},
		RecipientNonce: []byte("qwerty"),
	}
	encodedString := a.Serialize()
	b := Deserialize(*encodedString)
	if !reflect.DeepEqual(a, *b) {
		t.Fail()
	}
}
