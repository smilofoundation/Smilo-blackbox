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
	"os"
	"testing"
	"time"

	"Smilo-blackbox/src/data/types"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/utils"
)

type testEngine struct {
	Filename string
	Engine   string
}

func TestMain(m *testing.M) {
	// TODO: Tests to dynamodb and redis rely on services running to accept requests, for now they are just commented out.
	//       To run all tests we need to start services using docker instances and configure environment for aws.
	engines := []testEngine{
		{Filename: utils.BuildFilename("blackbox.db"), Engine: "boltdb"},
		//{Filename:"", Engine:"dynamodb"},
		//{Filename:"redis/test.conf", Engine:"redis"},
	}
	for _, eng := range engines {
		SetFilename(eng.Filename)
		SetEngine(eng.Engine)
		Start()
		time.Sleep(100000000)
		retcode := m.Run()
		if retcode != 0 {
			os.Exit(retcode)
		}
	}
	os.Exit(0)
}

func TestEncryptedTransaction_Save_Retrieve(t *testing.T) {
	now := time.Now()
	trans := types.CreateEncryptedTransaction([]byte("1"), []byte("AA"), now)
	err := trans.Save()
	require.NoError(t, err)

	trans2, err := types.FindEncryptedTransaction([]byte("1"))
	require.Empty(t, err)

	require.Equal(t, string(trans2.EncodedPayload), "AA")
	require.Equal(t, trans2.Timestamp.Unix(), now.Unix())
}

func TestEncryptedTransaction_Delete(t *testing.T) {
	trans := types.CreateEncryptedTransaction([]byte("2"), []byte("BB"), time.Now())
	err := trans.Save()
	require.NoError(t, err)

	trans2, err := types.FindEncryptedTransaction([]byte("2"))
	require.Empty(t, err)

	err = trans2.Delete()
	require.NoError(t, err)

	trans3, err := types.FindEncryptedTransaction([]byte("2"))
	require.NotEmpty(t, err)

	require.Empty(t, trans3)
}

func TestGetAllPeersEmpty(t *testing.T) {
	peers, err := types.GetAllPeers()
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peers")
	}
	require.Equal(t, peers, &[]types.Peer{})
}

func TestGetAllPeers(t *testing.T) {
	err := types.UpdateNewPeers([]string{"teste1", "teste2", "teste3", "teste4"}, "")
	require.NoError(t, err)
	peers, err := types.GetAllPeers()
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peers")
	}
	require.Equal(t, len(*peers), 4)
	require.Equal(t, (*peers)[0].URL, "teste1")
	require.Equal(t, (*peers)[3].URL, "teste4")

	for _, peer := range *peers {
		err = peer.Delete()
		require.NoError(t, err)
	}
}

func TestGetNextPeer(t *testing.T) {
	err := types.UpdateNewPeers([]string{"teste1", "teste2"}, "")
	require.NoError(t, err)
	peer1, err := types.FindNextUpdatablePeer(10 * time.Second)
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peer")
	}
	require.Equal(t, peer1.URL, "teste1")
	peer2, err := types.FindNextUpdatablePeer(10 * time.Second)
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peer")
	}
	require.Equal(t, peer2.URL, "teste2")
	err = peer1.Delete()
	require.NoError(t, err)
	err = peer2.Delete()
	require.NoError(t, err)
}
