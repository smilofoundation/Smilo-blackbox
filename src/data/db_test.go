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

package data

import (
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/utils"
)

func TestMain(m *testing.M) {
	SetFilename(utils.BuildFilename("blackbox.db"))
	Start()
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
	err := trans.Save()
	require.NoError(t, err)

	trans2, err := FindEncryptedTransaction([]byte("1"))
	require.Empty(t, err)

	require.Equal(t, string(trans2.Encoded_Payload), "AA")
	require.Equal(t, trans2.Timestamp.Unix(), now.Unix())
}

func TestEncrypted_Transaction_Delete(t *testing.T) {
	trans := CreateEncryptedTransaction([]byte("2"), []byte("BB"), time.Now())
	err := trans.Save()
	require.NoError(t, err)

	trans2, err := FindEncryptedTransaction([]byte("2"))
	require.Empty(t, err)

	err = trans2.Delete()
	require.NoError(t, err)

	trans3, err := FindEncryptedTransaction([]byte("2"))
	require.NotEmpty(t, err)

	require.Empty(t, trans3)
}
