package data

import (
	"encoding/hex"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Start(os.TempDir() + "/test.db")
	time.Sleep(100000000)
	retcode := m.Run()
	os.Exit(retcode)
}

func TestNewEncryptedTransaction(t *testing.T) {
	trans := NewEncryptedTransaction([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	tmp := hex.EncodeToString(trans.Hash)
	if trans.Timestamp.Before(time.Now().Add(-10000000000)) || tmp != "51e51636d1fcac073578a2529fce94c3b6e64ac0e14bbf57b17f0fb69e2d68da5adfee406ca13216ee49afc0f99145222a136033682319e9d3554dbb067afe3a" {
		t.Fail()
	}
}

func TestEncrypted_Transaction_Save_Retrieve(t *testing.T) {
	now := time.Now()
	trans := CreateEncryptedTransaction([]byte("1"), []byte("AA"), now)
	trans.Save()

	trans2 := FindEncryptedTransaction([]byte("1"))
	if string(trans2.Encoded_Payload) != "AA" || !trans2.Timestamp.Equal(now) {
		t.Fail()
	}
}

func TestEncrypted_Transaction_Delete(t *testing.T) {
	trans := CreateEncryptedTransaction([]byte("2"), []byte("BB"), time.Now())
	trans.Save()

	trans2 := FindEncryptedTransaction([]byte("2"))
	trans2.Delete()

	trans3 := FindEncryptedTransaction([]byte("2"))
	if trans3 != nil {
		t.Fail()
	}
}
