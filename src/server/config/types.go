package config

//Start of Private Key json file specification
type PrivateKeyBytes struct {
	Bytes string `json:"bytes"`
}

type PrivateKey struct {
	Data PrivateKeyBytes `json:"data"`
	Type string          `json:"type"`
}

//End of Private Key json file specification

//Start of Config json file specification
type Server struct {
	Port int `json:"port"`
}

type Peer struct {
	URL string `json:"url"`
}

type Key struct {
	PrivateKeyFile string `json:"config"`
	PublicKeyFile  string `json:"publicKey"`
}

type Keys struct {
	Passwords []string `json:"passwords"`
	KeyData   []Key    `json:"keyData"`
}

type Config struct {
	Server      Server `json:"server"`
	HostName    string `json:"hostName"`
	Peers       []Peer `json:"peer"`
	Keys        Keys   `json:"keys"`
	UnixSocket  string `json:"socket"`
	DBFile      string `json:"dbfile,omitempty"`
	PeersDBFile string `json:"peersdbfile,omitempty"`
}

//End of Config json file specification
