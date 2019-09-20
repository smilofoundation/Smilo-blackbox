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

//PublicKeyURL holds URL and pub for a peer
type PublicKeyURL struct {
	PublicKey []byte `key:"true"`
	URL       string
}

//NewPublicKeyURL create new peer based on pk and URL
func NewPublicKeyURL(pKey []byte, nodeURL string) *PublicKeyURL {
	p := PublicKeyURL{PublicKey: pKey, URL: nodeURL}
	return &p
}

//FindPublicKeyURL will find a peer
func FindPublicKeyURL(publicKey []byte) (*PublicKeyURL, error) {
	var p PublicKeyURL
	err := DBI.Find("PublicKey", publicKey, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

//Save save a PublicKeyURL into db
func (p *PublicKeyURL) Save() error {
	return DBI.Save(p)
}

//Delete delete a PublicKeyURL on db
func (p *PublicKeyURL) Delete() error {
	return DBI.Delete(p)
}
