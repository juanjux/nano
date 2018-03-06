package node

import (
	"bytes"
	"log"
	"math/rand"
	"net"
)

const packetSize = 512
const numberOfPeersToShare = 8

var DefaultPeer = Peer{
	// rai.raiblocks.net node
	net.ParseIP("::ffff:139.162.199.142"),
	7075,
}

var PeerList = []Peer{
	DefaultPeer,
	Peer{net.ParseIP("::ffff:217.111.215.36"), 7075},
	Peer{net.ParseIP("::ffff:103.6.12.90"), 7075},
	Peer{net.ParseIP("::ffff:181.199.69.145"), 7075},
}
var PeerSet = map[string]bool{
	DefaultPeer.String(): true,
	PeerList[1].String(): true,
	PeerList[2].String(): true,
	PeerList[3].String(): true,
}

func ListenForUdp() {
	log.Printf("Listening for udp packets on 7075")
	ln, err := net.ListenPacket("udp", ":7075")
	if err != nil {
		panic(err)
	}

	buf := make([]byte, packetSize)

	for {
		n, addr, err := ln.ReadFrom(buf)
		log.Printf("Received message from %v", addr)
		if err != nil {
			continue
		}
		if n > 0 {
			handleMessage(bytes.NewBuffer(buf[:n]))
		}
	}
}

func SendKeepAlive(peer Peer) error {
	addr := peer.Addr()
	randomPeers := make([]Peer, 0)

	outConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	randIndices := rand.Perm(len(PeerList))
	for n, i := range randIndices {
		if n == numberOfPeersToShare {
			break
		}
		randomPeers = append(randomPeers, PeerList[i])
	}

	m := CreateKeepAlive(randomPeers)
	buf := bytes.NewBuffer(nil)
	m.Write(buf)

	outConn.Write(buf.Bytes())
	log.Printf("Sent KeepAlive to %v", peer)
	return nil
}

