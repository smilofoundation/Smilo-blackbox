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

package config

//PrivateKeyBytes Start of Private Key json file specification
type PrivateKeyBytes struct {
	Bytes string `json:"bytes"`
}

//PrivateKey holds data and type
type PrivateKey struct {
	Data PrivateKeyBytes `json:"data"`
	Type string          `json:"type"`
}

//End of Private Key json file specification

//Server Start of Config json file specification
type Server struct {
	Port     int    `json:"port"`
	Hostaddr string `json:"hostaddr,omitempty"`
	TLSCert  string `json:"tlscert,omitempty"`
	TLSKey   string `json:"tlskey,omitempty"`
}

//Peer json file specification
type Peer struct {
	URL string `json:"url"`
}

//Key json file specification
type Key struct {
	PrivateKeyFile string `json:"config"`
	PublicKeyFile  string `json:"publicKey"`
}

//Keys json file specification
type Keys struct {
	Passwords []string `json:"passwords"`
	KeyData   []Key    `json:"keyData"`
}

//Config json file specification
type Config struct {
	Server      Server   `json:"server"`
	HostName    string   `json:"hostName"`
	RootCA      []string `json:"rootCA,omitempty"`
	Peers       []Peer   `json:"peer"`
	Keys        Keys     `json:"keys"`
	UnixSocket  string   `json:"socket"`
	DBEngine    string   `json:"dbengine,omitempty"`
	DBFile      string   `json:"dbfile,omitempty"`
	PeersDBFile string   `json:"peersdbfile,omitempty"`
}

//End of Config json file specification
