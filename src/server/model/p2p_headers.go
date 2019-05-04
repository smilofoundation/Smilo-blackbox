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

package model

var (
	//NETWORK_LIST p2p message type
	NETWORK_LIST = "NETWORK_LIST"
	//NETWORK_STATE p2p message type
	NETWORK_STATE = "NETWORK_STATE"
	//REQUEST_NET_STATE p2p message type
	REQUEST_NET_STATE = "REQUEST_NET_STATE"
	//PEER p2p message type
	PEER = "PEER"
	//GET_PEER p2p message type
	GET_PEER = "GET_PEER"
	//PEER_LIST p2p message type
	PEER_LIST = "PEER_LIST"
	//GET_PEER_LIST p2p message type
	GET_PEER_LIST = "GET_PEER_LIST"
	//MESSAGE p2p message type
	MESSAGE = "MESSAGE"

	//P2P BLOCK

	//COMMIT p2p message type
	COMMIT = "COMMIT"
	//APPROVE p2p message type
	APPROVE = "APPROVE"
	//DECLINE p2p message type
	DECLINE = "DECLINE"
	//BLOCK p2p message type
	BLOCK = "BLOCK"
	//GET_BLOCK p2p message type
	GET_BLOCK = "GET_BLOCK"
	//TRANSACTION p2p message type
	TRANSACTION = "TRANSACTION"
)
