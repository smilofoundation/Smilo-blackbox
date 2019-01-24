package model

import (
	"time"
)

// Peers contains the peer info
type PeerNode struct {
	ID            string    `json:"id" storm:"id"`
	LastSeen      time.Time `json:"last_seen"`
	NetworkStatus string    `json:"network_status"`
	RemoteAddr    string    `json:"remote_addr"`
}
