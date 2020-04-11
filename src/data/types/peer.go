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

import "time"

//Peer holds peer synchronization information
type Peer struct {
	URL         string `key:"true"`
	PublicKeys  [][]byte
	Failures    int
	LastFailure time.Time
	Tries       int
	NextUpdate  time.Time
}

//NewPublicKeyURL create new peer based on pk and URL
func NewPeer(url string) *Peer {
	p := Peer{URL: url, PublicKeys: make([][]byte, 0, 128), Failures: 0, Tries: 0, NextUpdate: time.Now()}
	return &p
}

//FindPublicKeyURL will find a peer
func FindPeer(URL string) (*Peer, error) {
	var p Peer
	err := DBI.Find("URL", URL, &p)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
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

func GetAllPeers() (*[]Peer, error) {
	return DBI.AllPeers()
}

func FindNextUpdatablePeer(postpone time.Duration) (*Peer, error) {
	return DBI.GetNextPeer(postpone)
}

func UpdateNewPeers(peers []string, hostURL string) error {
	for _, peer := range peers {
		if peer == hostURL {
			continue
		}
		p, err := FindPeer(peer)
		if err != nil && err != ErrNotFound {
			return err
		}
		if p == nil {
			p := NewPeer(peer)
			err = p.Save()
			if err != nil {
				return err
			}
		}
		// If peer already exists just ignore.
	}
	return nil
}
