package config

type PrivateKeyBytes struct {
	Bytes string `json:"bytes"`
}

// PrivateKey is a container for a private key.
type PrivateKey struct {
	Data PrivateKeyBytes `json:"data"`
	Type string          `json:"type"`
}
