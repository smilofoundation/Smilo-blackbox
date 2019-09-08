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

package types

//PublicKeyUrl holds URL and pub for a peer
type PublicKeyUrl struct {
	PublicKey []byte `key:"true"`
	URL       string
}

//NewPublicKeyUrl create new peer based on pk and URL
func NewPublicKeyUrl(pKey []byte, nodeURL string) *PublicKeyUrl {
	p := PublicKeyUrl{PublicKey: pKey, URL: nodeURL}
	return &p
}

//FindPublicKeyUrl will find a peer
func FindPublicKeyUrl(publicKey []byte) (*PublicKeyUrl, error) {
	var p PublicKeyUrl
	err := DBI.Find("PublicKey", publicKey, &p)
	if err != nil {
		//data.log.Error("Unable to find Peer.")
		return nil, err
	}
	return &p, nil
}

//Save save a PublicKeyUrl into db
func (p *PublicKeyUrl) Save() error {
	return DBI.Save(p)
}

//Delete delete a PublicKeyUrl on db
func (p *PublicKeyUrl) Delete() error {
	return DBI.Delete(p)
}

