// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package bft

import (
	"github.com/hyperledger-labs/SmartBFT/pkg/api"
	protos "github.com/hyperledger-labs/SmartBFT/smartbftprotos"
	"sync"
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
}

func (c *CensorProtector) CollectPools() {
	c.pools.clear(c.N)
	for {
		select {
		case <-c.stopChan:
			return
		case msg := <-c.incMsgs:
			c.pools.registerVote(msg.sender, msg.Message)
			c.Logger.Debugf("Node %d registered a pool: %+v; sender: %d", c.SelfID, msg.Message, msg.sender)
			if c.collectEnoughPools() {
				c.Logger.Debugf("Node %d collected enough pools", c.SelfID)
				c.calculateSet()
				return
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

func (c *CensorProtector) calculateSet() {

}
