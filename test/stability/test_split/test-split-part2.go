// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
	"github.com/intel-go/yanff/rules"
)

var (
	L3Rules    *rules.L3Rules
	options    = `{"cores": {"Value": 16, "Locked": false}}`
	RECV_PORT  uint
	SEND_PORT1 uint
	SEND_PORT2 uint
	SEND_PORT3 uint
)

// Main function for constructing packet processing graph.
func main() {
	// If you modify port numbers with cmd line, provide modified test-split.conf accordingly
	filename := flag.String("FILE", "test-split.conf", "file with split rules in .conf format. If you change default port numbers, please, provide modified rules file too")
	flag.UintVar(&RECV_PORT, "RECV_PORT", 0, "port for receiver")
	flag.UintVar(&SEND_PORT1, "SEND_PORT1", 1, "port for 1st sender")
	flag.UintVar(&SEND_PORT2, "SEND_PORT2", 2, "port for 2nd sender")
	flag.UintVar(&SEND_PORT3, "SEND_PORT3", 3, "port for 3rd sender")

	// Init YANFF system at requested number of cores.
	flow.SystemInit(options)

	// Get splitting rules from access control file.
	L3Rules = rules.GetL3RulesFromORIG(*filename)

	// Receive packets from 0 port
	inputFlow := flow.SetReceiver(uint8(RECV_PORT))

	// Split packet flow based on ACL.
	flowsNumber := 4
	outputFlows := flow.SetSplitter(inputFlow, L3Splitter, uint(flowsNumber))

	// "0" flow is used for dropping packets without sending them.
	flow.SetStopper(outputFlows[0])

	// Send each flow to corresponding port. Send queues will be added automatically.
	flow.SetSender(outputFlows[1], uint8(SEND_PORT1))
	flow.SetSender(outputFlows[2], uint8(SEND_PORT2))
	flow.SetSender(outputFlows[3], uint8(SEND_PORT3))

	// Begin to process packets.
	flow.SystemStart()
}

func L3Splitter(currentPacket *packet.Packet) uint {
	// Firstly set up all fields at packet: MAC, IPv4 or IPv6, TCP or UDP.
	currentPacket.ParseL4()

	// Return number of flow to which put this packet. Based on ACL rules.
	return rules.L3_ACL_port(currentPacket, L3Rules)
}
