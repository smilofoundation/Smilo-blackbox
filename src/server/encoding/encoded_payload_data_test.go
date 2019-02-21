// Copyright 2019 The Smilo-blackbox Authors
// This file is part of the Smilo-blackbox library.
//
// The Smilo-blackbox library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Smilo-blackbox library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Smilo-blackbox library. If not, see <http://www.gnu.org/licenses/>.

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
