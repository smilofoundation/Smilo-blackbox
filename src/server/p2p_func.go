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

package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tidwall/buntdb"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"

	"errors"
	"net"

	"Smilo-blackbox/src/server/config"
	"Smilo-blackbox/src/server/model"
)

func InitP2PServer(bootstrapNodes []*discover.Node) (*p2p.Server, error) {

	host := config.P2PDestination.Value

	if host == "" {
		var err error
		host, err = GetExternalIP()
		if err != nil {
			log.WithError(err).Fatal("Could not get external IP ADDR")
		}
	}

	port := config.P2PPort.Value

	strCon := fmt.Sprintf("%s:%s", host, port)
	log.WithField("strCon", strCon).Info("Starting P2P ")

	nodekey, err := crypto.LoadECDSA("key.txt")
	if err != nil {
		nodekey, err = crypto.GenerateKey()
		if err != nil {
			log.WithError(err).Error("Generate private key failed")
			return nil, err
		}

		err = crypto.SaveECDSA("key.txt", nodekey)
		if err != nil {
			log.WithError(err).Error("Generated private key failed to be saved")
		}
	} else {
		log.WithField("nodekey", nodekey.Public()).Info("ecdsa.PrivateKey loaded ok ")
	}

	srv = &p2p.Server{
		Config: p2p.Config{
			MaxPeers:        100,
			PrivateKey:      nodekey,
			Name:            common.MakeName("go-smilo node "+strCon, "1"),
			ListenAddr:      strCon,
			Protocols:       []p2p.Protocol{proto},
			EnableMsgEvents: true,
			DiscoveryV5:     true,
			BootstrapNodes:  bootstrapNodes,
		},
	}

	// the Err() on the Subscription object returns when subscription is closed

	if err := srv.Start(); err != nil {
		log.WithError(err).Error("Could not start p2p Server")
		return srv, err
	}

	return srv, nil

}

func InitP2p() {
	maxPeersNetwork = 20
	var err error

	var peerNodes []model.PeerNode
	//get all peers except yourself
	query := StormDBPeers.Select().OrderBy("LastSeen").Limit(10)

	// exec query
	err = query.Find(&peerNodes)

	if err != nil {
		log.WithError(err).
			Error("InitP2p, StormDBPeers, Could not get any peers")
	}

	var bootstrapNodes []*discover.Node
	for _, thispeer := range peerNodes {

		var urlstr = thispeer.ID
		if strings.Contains(thispeer.ID, "@") == false && thispeer.RemoteAddr != "" {
			urlstr = thispeer.ID + "@" + thispeer.RemoteAddr
		}

		p, err := discover.ParseNode(urlstr)

		if err != nil {
			log.
				WithField("urlstr", urlstr).
				Error("Incomplete Peer")
		} else if p.Incomplete() {
			log.
				WithField("urlstr", urlstr).
				Error("Incomplete Peer")
		} else {
			bootstrapNodes = append(bootstrapNodes, p)
		}
	}

	srv, err := InitP2PServer(bootstrapNodes)

	if err != nil {
		log.WithError(err).Fatal("InitP2p Could not start P2P Server")
	}

	srv_node := srv.Self()

	//save to db

	myPeerNode := model.PeerNode{
		ID:       srv_node.String(),
		LastSeen: time.Now(),
	}

	err = StormDBPeers.Save(&myPeerNode)
	if err != nil {
		log.WithError(err).Error("InitP2p, Could NOT save peer")
		return
	}

	log.
		WithField("myPeerNode", myPeerNode).
		WithField("len(peerNodes)", len(peerNodes)).
		Info("InitP2p, Will start p2p bootstrap")

	InitP2PPeers(peerNodes)

	SubscribeP2P()

}

func SubscribeP2P() {
	eventOneC := make(chan *p2p.PeerEvent)
	sub_one := srv.SubscribeEvents(eventOneC)
	go func() {

		for {
			select {
			case peerevent := <-eventOneC:

				switch peerevent.Type {
				case "add":
					log.WithField("Type", peerevent.Type).
						WithField("Peer", peerevent.Peer).
						Warn("InitP2PPeers, Received peer add notification")
					//save to peers db

					continue
				case "msgsend":
					//log.WithField("Type", peerevent.Type).
					//	WithField("Event", peerevent).
					//	WithField("node", i).
					//	Debug("Received message send notification")
					continue
				case "drop":
					log.WithField("Type", peerevent.Type).
						WithField("Event", peerevent).
						Error("Received DROP message")
					continue
				default:
					//log.WithField("Type", peerevent.Type).
					//	WithField("Event", peerevent).
					//	WithField("node", i).
					//	Debug("Received message")
					continue
				}
			case err := <-sub_one.Err():
				log.
					WithError(err).
					Error("**********************  Received sub_one.Err(): message")
				break
			}
		}
	}()

	// wait for the connections to complete
	time.Sleep(time.Millisecond * 1000)
}

func AddPeer(node *discover.Node) error {
	log.
		WithField("node.ID.String()", node.ID.TerminalString()).
		WithField("peerCount", srv.PeerCount()).
		Debug("Going to connect to peer ", srv.PeerCount())

	srv.RemovePeer(node)

	// add it as a peer
	// the connection and crypto handshake will be performed automatically
	srv.AddPeer(node)

	// inspect the results

	return nil
}

func InitP2PPeers(peers []model.PeerNode) {
	//log.Infof("InitP2PPeers, total peers %d", len(peers))
	for _, peer := range peers {

		if peer.ID != "" {

			var urlstr = peer.ID
			if strings.Contains(peer.ID, "@") == false && peer.RemoteAddr != "" {
				urlstr = peer.ID + "@" + peer.RemoteAddr
			}

			thispeer, err := discover.ParseNode(urlstr)

			if err != nil {
				log.
					Debug("Incomplete Peer")
				continue
			} else if thispeer.Incomplete() {
				log.
					WithField("peer", thispeer.ID.TerminalString()).
					Debug("Incomplete Peer")
				continue
			}

			targetObject := model.PeerNode{
				ID:         peer.ID,
				LastSeen:   time.Now(),
				RemoteAddr: peer.RemoteAddr,
			}

			alreadyAdded := IsPeerAlreadyAdded(thispeer)

			if alreadyAdded == false {
				log.
					WithField("parsedPeer", thispeer.ID.TerminalString()).
					Warn("Peer not found on the srv.Peers() and database, will run AddPeer. ")
				err = AddPeer(thispeer)

				StormDBPeers.Update(func(tx *buntdb.Tx) error {
					targetObjectArr, err := json.Marshal(&targetObject)
					if err != nil {
						log.
							WithError(err).
							Error("ERROR: CreateNewNetwork, Could not json.Marshal targetObject ")
					}
					myPeerbjectStr := string(targetObjectArr)

					_, _, err = tx.Set(targetObject.ID, myPeerbjectStr, DefaultExpirationTime)

					return err
				})

			} else {
				log.WithField("method", "InitP2PPeers").
					WithField("peer", thispeer.ID.TerminalString()).
					WithError(err).
					Debug("Peer already on the list!")
				continue
			}
			if err != nil {
				log.WithField("method", "InitP2PPeers").
					WithField("peer", thispeer.ID.TerminalString()).
					WithError(err).
					Error("Failed to connect to peer")
				continue
			} else {

				var peerNodes []model.PeerNode
				query := StormDBPeers.Select(q.Eq("ID", thispeer.String())).OrderBy("LastSeen").Limit(10)

				// exec query
				err = query.Find(&peerNodes)
				if err != nil {
					log.WithError(err).
						WithField("peer", thispeer.ID.TerminalString()).
						WithField("targetObject", targetObject).
						WithField("peerNodes len", len(peerNodes)).
						WithField("peerNodes", peerNodes).
						Error("InitP2p, StormDBPeers, Could not get targetObject peer")
					continue
				}

				var found bool
				for _, v := range peerNodes {
					if GetPeerNodeID(targetObject.ID) == GetPeerNodeID(v.ID) {
						found = true
					}
				}
				if found == false {
					log.
						WithField("peer", thispeer.ID.TerminalString()).
						WithField("targetObject", targetObject).
						WithField("peerNodes len", len(peerNodes)).
						WithField("peerNodes", peerNodes).
						Error("InitP2p, StormDBPeers, Could not query.Find the valid peer")
				}

				continue
			}
		} else {
			log.WithField("method", "InitP2PPeers").
				WithField("peer.ID", peer.ID).Error("Could not get a Peer instance from this peer")
		}

	}

}

//TODO: process a list of peers into our database / try to connect
func PeerList(p2pMessage Message) {
	var peerList []model.PeerNode
	err := json.Unmarshal([]byte(p2pMessage.Body), &peerList)
	if err != nil {
		log.WithError(err).Error("PEER_LIST, Could not json.Unmarshal peerList")
		return
	}

	log.
		WithField("peerCount", len(peerList)).
		Debug("********** PROCESS PEER_LIST **********")

	InitP2PPeers(peerList)

}

//TODO: return 10 peers from our database order by last seen, make sure I'm not in it
func GetPeerListSend(peer *p2p.Peer, rw p2p.MsgReadWriter) {
	var peerNodes []model.PeerNode

	msgCmutex.Lock()
	for _, thispeer := range srv.Peers() {
		p := model.PeerNode{
			ID:         thispeer.String(),
			LastSeen:   time.Now(),
			RemoteAddr: thispeer.RemoteAddr().String(),
		}

		peerNodes = append(peerNodes, p)
	}
	msgCmutex.Unlock()

	body, err := json.Marshal(peerNodes)
	if err != nil {
		log.
			WithError(err).Errorf("ERROR: GET_PEER_LIST, p2p Marshal, peerNodes")
	}

	outmsg := &Message{
		Header: model.PEER_LIST,
		Body:   string(body),
	}

	err = p2p.Send(rw, 0, outmsg)
	if err != nil {
		log.WithField("peer", peer).WithError(err).Errorf("ERROR: GET_PEER_LIST, p2p.Send")
		return
	}

	return
}

func GetPeerNodeID(id string) string {
	newStr := strings.Split(id, "@")[0]
	return newStr
}

func IsPeerAlreadyConnected(parsedPeer *discover.Node) bool {
	msgCmutex.Lock()
	var found bool
	for _, thispeer := range srv.Peers() {

		p, err := discover.ParseNode(thispeer.String())
		if err == nil && parsedPeer.String() == p.String() {
			found = true
		}
	}

	msgCmutex.Unlock()

	return found
}

func IsPeerAlreadyAdded(parsedPeer *discover.Node) bool {
	var found bool
	msgCmutex.Lock()

	log.
		WithField("parsedPeer", parsedPeer.ID.TerminalString()).
		Debug("Peer not found on the srv.Peers(), will check database")

	var peerNodes []model.PeerNode
	query := StormDBPeers.Select(q.Eq("ID", parsedPeer.String())).OrderBy("LastSeen").Limit(1)

	// exec query
	err := query.Find(&peerNodes)

	if err != nil {
		log.WithError(err).Errorf("ERROR: p2p.IsPeerAlreadyAdded")
	}

	if len(peerNodes) > 0 {
		found = true
	}

	msgCmutex.Unlock()

	return found
}

func GetExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
