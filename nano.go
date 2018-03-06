package main

import (
	"github.com/frankh/nano/node"
	"github.com/frankh/nano/store"
)

func main() {
	store.Init(store.LiveConfig)

	for _, n := range(node.PeerList) {
		node.SendKeepAlive(n)
	}
	node.ListenForUdp()
}

