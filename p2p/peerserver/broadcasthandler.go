package peerserver

import "github.com/HalalChain/qitmeer-lib/log"

// handleBroadcastMsg deals with broadcasting messages to peers.  It is invoked
// from the peerHandler goroutine.
func (s *PeerServer) handleBroadcastMsg(state *peerState, bmsg *broadcastMsg) {
	log.Error("TODO handleBroadcastMsg()")
}

