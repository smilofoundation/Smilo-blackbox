package types

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewEncryptedTransaction(t *testing.T) {
	trans := NewEncryptedTransaction([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	tmp := hex.EncodeToString(trans.Hash)
	require.True(t, trans.Timestamp.Before(time.Now().Add(-10000000000)) || tmp == "51e51636d1fcac073578a2529fce94c3b6e64ac0e14bbf57b17f0fb69e2d68da5adfee406ca13216ee49afc0f99145222a136033682319e9d3554dbb067afe3a")
}
