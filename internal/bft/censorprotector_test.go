// Copyright IBM Corp. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package bft_test

import (
	"github.com/hyperledger-labs/SmartBFT/internal/bft"
	protos "github.com/hyperledger-labs/SmartBFT/smartbftprotos"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
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
