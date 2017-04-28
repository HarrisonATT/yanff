// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/intel-go/yanff/flow"
import "github.com/intel-go/yanff/packet"

import "flag"

var (
	mode    uint
	cores   uint
	options = `{"cores": {"Value": 15, "Locked": false}}`

	RECV_PORT1 uint
	RECV_PORT2 uint
	SEND_PORT1 uint
	SEND_PORT2 uint
)

func main() {
	flag.UintVar(&mode, "mode", 0, "Benching mode: 0 - empty, 1 - parsing, 2 - parsing, reading, writing")
	flag.UintVar(&SEND_PORT1, "SEND_PORT1", 1, "port for 1st sender")
	flag.UintVar(&SEND_PORT2, "SEND_PORT2", 1, "port for 2nd sender")
	flag.UintVar(&RECV_PORT1, "RECV_PORT1", 0, "port for 1st receiver")
	flag.UintVar(&RECV_PORT2, "RECV_PORT2", 0, "port for 2nd receiver")

	// Initialize YANFF library at requested number of cores
	flow.SystemInit(options)

	// Receive packets from zero port. One queue per receive will be added automatically.
	firstFlow0 := flow.SetReceiver(uint8(RECV_PORT1))
	firstFlow1 := flow.SetReceiver(uint8(RECV_PORT2))

	firstFlow := flow.SetMerger(firstFlow0, firstFlow1)

	// Handle second flow via some heavy function
	if mode == 0 {
		flow.SetHandler(firstFlow, heavyFunc0)
	} else if mode == 1 {
		flow.SetHandler(firstFlow, heavyFunc1)
	} else {
		flow.SetHandler(firstFlow, heavyFunc2)
	}

	// Split for two senders and send
	secondFlow := flow.SetPartitioner(firstFlow, 150, 150)
	flow.SetSender(firstFlow, uint8(SEND_PORT1))
	flow.SetSender(secondFlow, uint8(SEND_PORT2))

	flow.SystemStart()
}

func heavyFunc0(currentPacket *packet.Packet) {
}

func heavyFunc1(currentPacket *packet.Packet) {
	currentPacket.ParseEtherIPv4()
}

func heavyFunc2(currentPacket *packet.Packet) {
	currentPacket.ParseEtherIPv4()
	T := (currentPacket.IPv4.DstAddr)
	currentPacket.IPv4.SrcAddr = 263 + (T)
}
