// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package bft_test

import (
	"testing"

	"github.com/hyperledger-labs/SmartBFT/internal/bft"
	"github.com/hyperledger-labs/SmartBFT/pkg/types"
	protos "github.com/hyperledger-labs/SmartBFT/smartbftprotos"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCollectPools(t *testing.T) {
	basicLog, err := zap.NewDevelopment()
	assert.NoError(t, err)
	log := basicLog.Sugar()

	protector := &bft.CensorProtector{
		SelfID: 1,
		N:      4,
		Logger: log,
	}

	protector.Start()

	txs2 := make([]*protos.TX, 0)
	txs2 = append(txs2, &protos.TX{
		Id:       "222",
		ClientId: "2",
		Req:      nil,
	})

	msg2 := &protos.Message{
		Content: &protos.Message_TxPoolBroadcast{
			TxPoolBroadcast: &protos.TXPoolBroadcast{
				Txs: txs2,
			},
		},
	}

	txs3 := make([]*protos.TX, 0)
	txs3 = append(txs3, &protos.TX{
		Id:       "333",
		ClientId: "3",
		Req:      nil,
	})

	msg3 := &protos.Message{
		Content: &protos.Message_TxPoolBroadcast{
			TxPoolBroadcast: &protos.TXPoolBroadcast{
				Txs: txs3,
			},
		},
	}

	txs4 := make([]*protos.TX, 0)
	txs4 = append(txs4, &protos.TX{
		Id:       "444",
		ClientId: "4",
		Req:      nil,
	})

	msg4 := &protos.Message{
		Content: &protos.Message_TxPoolBroadcast{
			TxPoolBroadcast: &protos.TXPoolBroadcast{
				Txs: txs4,
			},
		},
	}

	protector.ClearCollected()

	go func() {
		protector.HandleMessage(2, msg2)
		protector.HandleMessage(3, msg3)
		protector.HandleMessage(4, msg4)
	}()

	protector.CollectPools()

	protector.Stop()

}

func TestCalculateSetAndVerify(t *testing.T) {
	basicLog, err := zap.NewDevelopment()
	assert.NoError(t, err)
	log := basicLog.Sugar()

	protector := &bft.CensorProtector{
		SelfID: 1,
		N:      4,
		Logger: log,
	}

	protector.Start()

	tx1 := &protos.TX{Id: "tx1", Req: []byte{11}}
	tx2 := &protos.TX{Id: "tx2", Req: []byte{22}}
	tx3 := &protos.TX{Id: "tx3", Req: []byte{33}}
	txs := make([]*protos.TX, 0)
	txs = append(txs, tx1, tx2, tx3)
	reqs := make([]types.RequestInfo, 0)
	reqs = append(reqs, types.RequestInfo{ID: tx1.Id}, types.RequestInfo{ID: tx2.Id}, types.RequestInfo{ID: tx3.Id})

	msg := &protos.Message{
		Content: &protos.Message_TxPoolBroadcast{
			TxPoolBroadcast: &protos.TXPoolBroadcast{
				Txs: txs,
			},
		},
	}

	protector.ClearCollected()

	go func() {
		protector.HandleMessage(2, msg)
		protector.HandleMessage(3, msg)
		protector.HandleMessage(4, msg)
	}()

	protector.CollectPools()
	assert.NoError(t, protector.VerifyProposed(reqs))

	txs12 := make([]*protos.TX, 0)
	txs12 = append(txs12, tx1, tx2)
	msg12 := &protos.Message{
		Content: &protos.Message_TxPoolBroadcast{
			TxPoolBroadcast: &protos.TXPoolBroadcast{
				Txs: txs12,
			},
		},
	}

	protector.ClearCollected()

	go func() {
		protector.HandleMessage(2, msg)
		protector.HandleMessage(3, msg12)
		protector.HandleMessage(4, msg12)
	}()

	protector.CollectPools()
	assert.NoError(t, protector.VerifyProposed(reqs))
	assert.NoError(t, protector.VerifyProposed(reqs[0:2]))
	assert.Error(t, protector.VerifyProposed(reqs[0:1]))

	txs3 := make([]*protos.TX, 0)
	txs3 = append(txs3, tx3)
	msg3 := &protos.Message{
		Content: &protos.Message_TxPoolBroadcast{
			TxPoolBroadcast: &protos.TXPoolBroadcast{
				Txs: txs3,
			},
		},
	}

	protector.ClearCollected()

	go func() {
		protector.HandleMessage(2, msg)
		protector.HandleMessage(3, msg3)
		protector.HandleMessage(4, msg)
	}()

	protector.CollectPools()
	assert.NoError(t, protector.VerifyProposed(reqs))
	assert.Error(t, protector.VerifyProposed(reqs[0:2]))
	assert.NoError(t, protector.VerifyProposed(reqs[2:3]))

	protector.Stop()
}
