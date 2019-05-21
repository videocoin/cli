package listener

import (
	"fmt"
	"math/big"

	"github.com/VideoCoin/go-videocoin/common"
)

var (
	EventStreamRequested       = "EventStreamRequested"
	EventStreamCreated         = "EventStreamCreated"
	EventStreamApproved        = "EventStreamApproved"
	EventStreamInputChunkAdded = "EventStreamInputChunkAdded"
)

type event struct {
	Name string `json:"name"`

	StreamID *big.Int `json:"streamId"`
	ChunkID  *big.Int `json:"chunkId"`

	Address       common.Address `json:"address"`
	StreamAddress common.Address `json:"streamAddress"`

	BlockNumber uint64      `json:"blockNumber"`
	BlockHash   common.Hash `json:"blockHash"`

	TxHash  common.Hash `json:"transactionHash" gencodec:"required"`
	TxIndex uint        `json:"transactionIndex" gencodec:"required"`
}

func (e *event) String() string {
	switch e.Name {
	case EventStreamRequested, EventStreamApproved:
		return fmt.Sprintf(
			"%s\n\tAddress: %s\n\tStreamId: %s\n\tBlockNumber: %d\n\tBlockHash: %s\n\tTxHash: %s\n\tTxIndex: %d\n",
			e.Name, e.Address.String(), e.StreamID, e.BlockNumber, e.BlockHash.String(), e.TxHash.String(), e.TxIndex)
	case EventStreamCreated:
		return fmt.Sprintf(
			"%s\n\tAddress: %s\n\tStreamId: %s\n\tStreamAddress: %s\n\tBlockNumber: %d\n\tBlockHash: %s\n\tTxHash: %s\n\tTxIndex: %d\n",
			e.Name, e.Address.String(), e.StreamID, e.StreamAddress.String(), e.BlockNumber, e.BlockHash.String(), e.TxHash.String(), e.TxIndex)
	default:
		return ""
	}
}
