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
	"Smilo-blackbox/src/utils"
	"os"
	"strconv"
	"testing"
	"time"

	"Smilo-blackbox/src/data/types"

	"github.com/stretchr/testify/require"
)

type testEngine struct {
	Filename string
	Engine   string
	CleanUp  func()
}

func TestMain(m *testing.M) {
	// TODO: Tests to dynamodb and redis rely on services running to accept requests, for now they are just commented out.
	//       To run all tests we need to start services using docker instances and configure environment for aws.
	engines := []testEngine{
		{Filename: utils.BuildFilename("blackbox.db"), Engine: "boltdb", CleanUp: func() { os.Remove(utils.BuildFilename("blackbox.db")) }},
		//{Filename: "", Engine: "dynamodb", CleanUp: func() {}},
		//{Filename: "redis/test.conf", Engine: "redis", CleanUp: func() {}},
	}
	for _, eng := range engines {
		eng.CleanUp()
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
	testValues := []string{"teste1", "teste2", "teste3", "teste4"}
	err := types.UpdateNewPeers(testValues, "")
	require.NoError(t, err)
	peers, err := types.GetAllPeers()
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peers")
	}
	require.Equal(t, len(*peers), 4)
	require.Contains(t, testValues, (*peers)[0].URL)
	require.Contains(t, testValues, (*peers)[3].URL)
	require.NotEqual(t, (*peers)[0].URL, (*peers)[3].URL)

	for _, peer := range *peers {
		err = peer.Delete()
		require.NoError(t, err)
	}
}

func TestGetAll(t *testing.T) {
	testValues := []string{"teste1", "teste2", "teste3", "teste4"}
	err := types.UpdateNewPeers(testValues, "")
	require.NoError(t, err)
	peers := make([]types.Peer, 0)
	err = types.GetAll(&peers)
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peers")
	}
	require.Equal(t, len(peers), 4)
	require.Contains(t, testValues, (peers)[0].URL)
	require.Contains(t, testValues, (peers)[3].URL)
	require.NotEqual(t, (peers)[0].URL, (peers)[3].URL)

	for _, peer := range peers {
		err = peer.Delete()
		require.NoError(t, err)
	}
}

func TestGetNextPeer(t *testing.T) {
	testValues := []string{"teste1", "teste2"}
	err := types.UpdateNewPeers(testValues, "")
	require.NoError(t, err)
	peer1, err := types.FindNextUpdatablePeer(10 * time.Second)
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peer")
	}
	require.Contains(t, testValues, peer1.URL)
	peer2, err := types.FindNextUpdatablePeer(10 * time.Second)
	if err != nil {
		require.Fail(t, "Unexpected error retrieving peer")
	}
	require.Contains(t, testValues, peer2.URL)
	require.NotEqual(t, peer1.URL, peer2.URL)
	err = peer1.Delete()
	require.NoError(t, err)
	err = peer2.Delete()
	require.NoError(t, err)
}

func TestMigrateBoltDB(t *testing.T) {
	var peers []types.Peer
	var transactions []types.EncryptedTransaction
	var rawTransactions []types.EncryptedRawTransaction
	var publicKeys []types.PublicKeyURL

	for i := 0; i < 100; i++ {
		now := time.Now()
		trans := types.CreateEncryptedTransaction([]byte(strconv.Itoa(i)), []byte("Payload: "+strconv.Itoa(i)), now)
		err := trans.Save()
		require.NoError(t, err)
	}
	for i := 0; i < 100; i++ {
		trans := types.NewEncryptedRawTransaction([]byte("Payload: "+strconv.Itoa(i)), []byte(""))
		err := trans.Save()
		require.NoError(t, err)
	}

	for i := 0; i < 200; i++ {
		peer := types.NewPeer("teste " + strconv.Itoa(i))
		for j := 0; j < 2; j++ {
			peer.PublicKeys = append(peer.PublicKeys, []byte("pk_"+strconv.Itoa(i)+"_"+strconv.Itoa(j)))
		}
		err := peer.Save()
		require.NoError(t, err)
	}
	err := types.GetAll(&peers)
	require.NoError(t, err)
	err = types.DBI.Close()
	require.NoError(t, err)
	_ = os.Remove("blackbox2.db")
	err = Migrate(dbEngine, dbFile, BOLTDBENGINE, "blackbox2.db")
	require.NoError(t, err)

	err = types.GetAll(&peers)
	require.NoError(t, err)
	err = types.GetAll(&transactions)
	require.NoError(t, err)
	err = types.GetAll(&rawTransactions)
	require.NoError(t, err)
	err = types.GetAll(&publicKeys)
	require.NoError(t, err)

	require.Equal(t, 100, len(transactions))
	require.Equal(t, 100, len(rawTransactions))
	require.Equal(t, 200, len(peers))
	require.Equal(t, 400, len(publicKeys))
}
