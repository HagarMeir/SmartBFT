// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package test

import (
	"os"
	"testing"

	"github.com/hyperledger-labs/SmartBFT/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestBasicCensorship(t *testing.T) {
	network := NewNetwork()
	defer network.Shutdown()

	testDir, err := os.MkdirTemp("", t.Name())
	assert.NoErrorf(t, err, "generate temporary test dir")
	defer os.RemoveAll(testDir)

	numberOfNodes := 5
	nodes := make([]*App, 0)
	for i := 1; i <= numberOfNodes; i++ {
		n := newNode(uint64(i), network, t.Name(), testDir, false, 0)
		n.Consensus.Config.CensorProtect = true
		n.Consensus.Config.RequestBatchMaxCount = 2
		n.censorTXLeaderID = 1
		n.censorTXInfo = types.RequestInfo{ID: "3", ClientID: "alice"}
		nodes = append(nodes, n)
	}

	startNodes(nodes, network)

	for i := 0; i < numberOfNodes; i++ {
		nodes[i].Submit(Request{ID: "1", ClientID: "alice"})
		nodes[i].Submit(Request{ID: "2", ClientID: "alice"})
		nodes[i].Submit(Request{ID: "3", ClientID: "alice"})
		nodes[i].Submit(Request{ID: "4", ClientID: "alice"})
	}

	data := make([]*AppRecord, 0)
	for i := 0; i < numberOfNodes; i++ {
		d := <-nodes[i].Delivered
		data = append(data, d)
	}
	for i := 0; i < numberOfNodes-1; i++ {
		assert.Equal(t, data[i], data[i+1])
	}

	data = make([]*AppRecord, 0)
	for i := 0; i < numberOfNodes; i++ {
		d := <-nodes[i].Delivered
		data = append(data, d)
	}
	for i := 0; i < numberOfNodes-1; i++ {
		assert.Equal(t, data[i], data[i+1])
	}
}
