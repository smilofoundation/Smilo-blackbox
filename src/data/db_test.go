package data

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/server/config"
)

func TestMain(m *testing.M) {
	config.DBFile = filepath.Join(os.TempDir(), "/test.db")
	time.Sleep(100000000)
	retcode := m.Run()
	os.Exit(retcode)
}

func TestNewEncryptedTransaction(t *testing.T) {
	trans := NewEncryptedTransaction([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	tmp := hex.EncodeToString(trans.Hash)
	require.True(t, trans.Timestamp.Before(time.Now().Add(-10000000000)) || tmp == "51e51636d1fcac073578a2529fce94c3b6e64ac0e14bbf57b17f0fb69e2d68da5adfee406ca13216ee49afc0f99145222a136033682319e9d3554dbb067afe3a")
}

func TestEncrypted_Transaction_Save_Retrieve(t *testing.T) {
	now := time.Now()
	trans := CreateEncryptedTransaction([]byte("1"), []byte("AA"), now)
	trans.Save()

	trans2, err := FindEncryptedTransaction([]byte("1"))
	require.Empty(t, err)

	require.Equal(t, string(trans2.Encoded_Payload), "AA")
	require.Equal(t, trans2.Timestamp.Unix(), now.Unix())
}

func TestEncrypted_Transaction_Delete(t *testing.T) {
	trans := CreateEncryptedTransaction([]byte("2"), []byte("BB"), time.Now())
	trans.Save()

	trans2, err := FindEncryptedTransaction([]byte("2"))
	require.Empty(t, err)

	trans2.Delete()

	trans3, err := FindEncryptedTransaction([]byte("2"))
	require.NotEmpty(t, err)

	require.Empty(t, trans3)
}
