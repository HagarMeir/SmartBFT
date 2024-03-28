// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package bft

import (
	"sync"

	"github.com/hyperledger-labs/SmartBFT/pkg/api"
	"github.com/hyperledger-labs/SmartBFT/pkg/types"
	protos "github.com/hyperledger-labs/SmartBFT/smartbftprotos"
	"github.com/pkg/errors"
)

type CensorProtector struct {
	SelfID uint64
	N      uint64
	f      int
	q      int

	Logger api.Logger

	incMsgs chan *incMsg

	pools *voteSet

	stopOnce sync.Once
	stopChan chan struct{}

	set []string
}

func (c *CensorProtector) Start() {
	c.q, c.f = computeQuorum(c.N)
	c.stopChan = make(chan struct{})
	c.stopOnce = sync.Once{}
	c.incMsgs = make(chan *incMsg, c.N)

	acceptPool := func(_ uint64, message *protos.Message) bool {
		return message.GetTxPoolBroadcast() != nil
	}
	c.pools = &voteSet{
		validVote: acceptPool,
	}
	c.pools.clear(c.N)
}

func (c *CensorProtector) close() {
	c.stopOnce.Do(
		func() {
			select {
			case <-c.stopChan:
				return
			default:
				close(c.stopChan)
			}
		},
	)
}

func (c *CensorProtector) Stop() {
	c.close()
}

func (c *CensorProtector) HandleMessage(sender uint64, m *protos.Message) {
	if m.GetTxPoolBroadcast() == nil {
		c.Logger.Panicf("Node %d handling a message which is not a pool", c.SelfID)
	}
	msg := &incMsg{sender: sender, Message: m}
	c.Logger.Debugf("Node %d handling pool messages: %v", c.SelfID, msg)
	select {
	case <-c.stopChan:
		return
	case c.incMsgs <- msg:
	default: // if incMsgs is full do nothing
		c.Logger.Debugf("Node %d reached default in handling pool messages: %v", c.SelfID, msg)
	}
}

func (c *CensorProtector) ClearCollected() {
	// drain message channel
	for len(c.incMsgs) > 0 {
		<-c.incMsgs
	}
	c.pools.clear(c.N)
	c.set = nil
}

func (c *CensorProtector) CollectPools() [][]byte {
	for {
		select {
		case <-c.stopChan:
			return nil
		case msg := <-c.incMsgs:
			c.pools.registerVote(msg.sender, msg.Message)
			c.Logger.Debugf("Node %d registered a pool: %+v; sender: %d", c.SelfID, msg.Message, msg.sender)
			if c.collectEnoughPools() {
				c.Logger.Debugf("Node %d collected enough pools", c.SelfID)
				return c.calculateSet()
			}
		}
	}
}

func (c *CensorProtector) collectEnoughPools() bool {
	c.Logger.Debugf("Node %d so far collected %d pools", c.SelfID, len(c.pools.voted))
	if len(c.pools.voted) < c.q {
		return false
	}
	return true
}

func (c *CensorProtector) calculateSet() [][]byte {
	counters := make(map[string]int, 0)
	num := len(c.pools.votes)
	var requests [][]byte
	for i := 0; i < num; i++ {
		vote := <-c.pools.votes
		pool := vote.GetTxPoolBroadcast()
		if pool == nil {
			c.Logger.Panicf("Node %d collected a message which is not a pool", c.SelfID)
			return nil
		}
		for _, tx := range pool.Txs {
			counters[tx.Id]++
			requests = append(requests, tx.Req)
		}
	}

	var set []string
	for tx, count := range counters {
		if count >= c.q {
			set = append(set, tx)
			c.Logger.Debugf("Node %d added tx %s to its set", c.SelfID, tx)
		}
	}

	c.set = set

	return requests
}

func (c *CensorProtector) VerifyProposed(requests []types.RequestInfo) error {
	for _, id := range c.set {
		found := false
		for _, req := range requests {
			if id == req.ID {
				found = true
			}
		}
		if !found {
			return errors.Errorf("Node %d did not find request %s in the proposal", c.SelfID, id)
		}
	}
	return nil
}
