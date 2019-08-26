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
	"fmt"

	"github.com/ethereum/go-ethereum/p2p"

	"Smilo-blackbox/src/server/model"
)

//Message holds header and body
type Message struct {
	Header string `json:"content"`
	Body   string `json:"body"`
}

//Peer is the main peer struct
type Peer struct {
	ID         string
	Dest       string
	SourcePort int
}

var (
	//protocols       []p2p.Protocol
	srv             *p2p.Server
//	maxPeersNetwork int
)

var (
	proto = p2p.Protocol{
		Name:    "blackbox",
		Version: 1,
		Length:  1,
		Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {

			go func() {

				for content := range msgC {

					switch content.Header {

					case model.GET_PEER_LIST:
						err := p2p.Send(rw, 0, content)
						if err != nil {
							log.
								WithField("peer", peer).WithError(err).Error("GET_PEER_LIST, p2p.Send, Could not send message")
							return
						}
						continue
					}

				}
				log.WithField("peer", peer).Error("p2p.Send, Could not get message from channel")
			}()

			for {
				msg, err := rw.ReadMsg()
				if err != nil {
					log.WithError(err).Error("ERROR: p2p ReadMsg")
					break
				}

				var p2pMessage Message
				err = msg.Decode(&p2pMessage)
				if err != nil {
					log.WithError(err).Error("ERROR: p2p Decode")
					continue
				}
				// make sure that the payload has been fully consumed
				err = msg.Discard()
				if err != nil {
					log.WithError(err).Error("Could not msg.Discard")
				}

				switch p2pMessage.Header {

				// ################# BEGIN P2P ###########################

				case model.GET_PEER:
					//TODO: return what we think about a peer random if none is required
					//TODO: eg: peer could have been blacklisted / banned temp
					continue

				case model.PEER:
					//TODO: process a peer into our list of peers / try to connect
					continue

				case model.PEER_LIST:
					PeerList(p2pMessage)
					continue

				case model.GET_PEER_LIST:
					GetPeerListSend(peer, rw)
					continue

				case model.MESSAGE:
					//TODO: implement
					continue
					// ################# END P2P ###########################

					// ################# BEGIN P2P BLOCK ###########################
				case model.NETWORK_STATE:
					//TODO: implement
					continue

				case model.REQUEST_NET_STATE:
					//TODO: implement
					continue

				case model.COMMIT:
					//TODO: implement
					continue

				case model.APPROVE:
					//TODO: implement
					continue

				case model.DECLINE:
					//TODO: implement
					continue

				case model.BLOCK:
					//TODO: implement
					continue

				case model.GET_BLOCK:
					//TODO: implement
					continue

				case model.TRANSACTION:
					//TODO: implement
					continue
					// ################# END P2P BLOCK ###########################

				default:
					fmt.Println("GOT UNKNOWN MSG!!!!!!!!!!!!!!!!!!!!:", p2pMessage)
					continue
				}
			}
			log.Info("terminate the protocol ??")

			return nil
		},
	}
)

// SendMsg will send a message
func SendMsg(peer *p2p.Peer, rw p2p.MsgWriter, err error, outmsg Message) {

	if outmsg.Header != "" {
		err = p2p.Send(rw, 0, outmsg)
	}
	if err != nil {
		log.
			WithError(err).
			WithField("Header", outmsg.Header).
			WithField("peer", peer).
			WithField("peerCount", srv.PeerCount()).
			Error("p2p.Send, Could not send second network")
	} else {
		log.
			WithField("Header", outmsg.Header).
			WithField("peer", peer).
			WithField("peerCount", srv.PeerCount()).
			Debug("p2p.Send, sent second network")
	}
}
