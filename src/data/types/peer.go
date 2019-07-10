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

//Peer holds url and pub for a peer
type Peer struct {
	publicKey []byte `storm:"id"`
	url       string
}

//NewPeer create new peer based on pk and url
func NewPeer(pKey []byte, nodeURL string) *Peer {
	p := Peer{publicKey: pKey, url: nodeURL}
	return &p
}

//Update will update a peer
func Update(pKey []byte, nodeURL string) *Peer {
	p, err := FindPeer(pKey)
	if err != nil {
		p = NewPeer(pKey, nodeURL)
	} else {
		p.url = nodeURL
	}
	err = p.Save()
	if err != nil {
		//data.log.WithError(err).Error("Could not Update, failed to execute Save method")
	}
	return p
}

//FindPeer will find a peer
func FindPeer(publicKey []byte) (*Peer, error) {
	var p Peer
	err := DBI.Find("publicKey", publicKey, &p)
	if err != nil {
		//data.log.Error("Unable to find Peer.")
		return nil, err
	}
	return &p, nil
}

//Save save a peer into db
func (p *Peer) Save() error {
	return DBI.Save(p)
}

//Delete delete a peer on db
func (p *Peer) Delete() error {
	return DBI.Delete(p)
}
