package listener

import (
	"fmt"
	"math/big"
	"time"

	sm "github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/sirupsen/logrus"
)

type EventListenerConfig struct {
	SmartContractManager *sm.Manager
	Timeout              time.Duration
	Logger               *logrus.Entry
}

type EventListener struct {
	smartContractManager *sm.Manager
	timeout              time.Duration
	logger               *logrus.Entry
}

func NewEventListener(c *EventListenerConfig) *EventListener {
	return &EventListener{
		smartContractManager: c.SmartContractManager,
		timeout:              c.Timeout,
		logger:               c.Logger.WithField("component", "event-listener"),
	}
}

func (e *EventListener) LogStreamRequestEvent(streamID *big.Int, address common.Address) (chan *event, chan error) {
	addresses := []common.Address{address}
	streamIDs := []*big.Int{streamID}

	resultCh := make(chan *event, 1)
	errCh := make(chan error, 1)

	go func() {
		for timeout := time.After(e.timeout * time.Second); ; {
			select {
			case <-timeout:
				errCh <- fmt.Errorf("failed to log stream request event and exit on timeout")
				break
			default:
				iterator, err := e.smartContractManager.FilterStreamRequested(
					new(bind.FilterOpts), addresses, streamIDs)
				if err != nil {
					errCh <- fmt.Errorf("failed to log stream request event: %s", err.Error())
				}

				for {
					if iterator.Error() != nil {
						errCh <- fmt.Errorf("failed to retrieve or parse log: %s", err.Error())
					}
					if iterator.Event != nil {
						e := iterator.Event
						resultCh <- &event{
							Name:        EventStreamRequested,
							StreamID:    e.StreamId,
							Address:     e.Raw.Address,
							BlockNumber: e.Raw.BlockNumber,
							BlockHash:   e.Raw.BlockHash,
							TxHash:      e.Raw.TxHash,
							TxIndex:     e.Raw.TxIndex,
						}
						return
					}

					if !iterator.Next() {
						break
					}
				}

				time.Sleep(5 * time.Second)
			}
		}
	}()

	return resultCh, errCh
}

func (e *EventListener) LogStreamCreateEvent(streamID *big.Int) (chan *event, chan error) {
	streamAddresses := []common.Address{}
	streamIDs := []*big.Int{streamID}

	resultCh := make(chan *event, 1)
	errCh := make(chan error, 1)

	go func() {
		for timeout := time.After(e.timeout * time.Second); ; {
			select {
			case <-timeout:
				errCh <- fmt.Errorf("failed to log stream created event and exit on timeout")
				break
			default:
				iterator, err := e.smartContractManager.FilterStreamCreated(
					new(bind.FilterOpts), streamAddresses, streamIDs)
				if err != nil {
					errCh <- fmt.Errorf("failed to log stream created event: %s", err.Error())
				}

				for {
					if iterator.Error() != nil {
						errCh <- fmt.Errorf("failed to retrieve or parse log: %s", err.Error())
					}
					if iterator.Event != nil {
						e := iterator.Event
						resultCh <- &event{
							Name:          EventStreamCreated,
							StreamID:      e.StreamId,
							Address:       e.Raw.Address,
							StreamAddress: e.StreamAddress,
							BlockNumber:   e.Raw.BlockNumber,
							BlockHash:     e.Raw.BlockHash,
							TxHash:        e.Raw.TxHash,
							TxIndex:       e.Raw.TxIndex,
						}
						return
					}

					if !iterator.Next() {
						break
					}
				}

				time.Sleep(5 * time.Second)
			}
		}
	}()

	return resultCh, errCh
}

func (e *EventListener) LogStreamApproveEvent(streamID *big.Int) (chan *event, chan error) {
	streamIDs := []*big.Int{streamID}

	resultCh := make(chan *event, 1)
	errCh := make(chan error, 1)

	go func() {
		for timeout := time.After(e.timeout * time.Second); ; {
			select {
			case <-timeout:
				errCh <- fmt.Errorf("failed to log stream approved event and exit on timeout")
				break
			default:
				iterator, err := e.smartContractManager.FilterStreamApproved(
					new(bind.FilterOpts), streamIDs)
				if err != nil {
					errCh <- fmt.Errorf("failed to log stream approved event: %s", err.Error())
				}

				for {
					if iterator.Error() != nil {
						errCh <- fmt.Errorf("failed to retrieve or parse log: %s", err.Error())
					}
					if iterator.Event != nil {
						e := iterator.Event
						resultCh <- &event{
							Name:        EventStreamApproved,
							StreamID:    e.StreamId,
							Address:     e.Raw.Address,
							BlockNumber: e.Raw.BlockNumber,
							BlockHash:   e.Raw.BlockHash,
							TxHash:      e.Raw.TxHash,
							TxIndex:     e.Raw.TxIndex,
						}
						return
					}

					if !iterator.Next() {
						break
					}
				}

				time.Sleep(5 * time.Second)
			}
		}
	}()

	return resultCh, errCh
}

func (e *EventListener) LogInputChunkAddEvent(streamID *big.Int, chunkID *big.Int) (chan *event, chan error) {
	streamIDs := []*big.Int{streamID}
	chunkIDs := []*big.Int{chunkID}

	resultCh := make(chan *event, 1)
	errCh := make(chan error, 1)

	go func() {
		for timeout := time.After(e.timeout * time.Second); ; {
			select {
			case <-timeout:
				errCh <- fmt.Errorf("failed to log input chunk added event and exit on timeout")
				break
			default:
				iterator, err := e.smartContractManager.FilterInputChunkAdded(
					new(bind.FilterOpts), streamIDs, chunkIDs)
				if err != nil {
					errCh <- fmt.Errorf("failed to log input chunk added event: %s", err.Error())
				}

				for {
					if iterator.Error() != nil {
						errCh <- fmt.Errorf("failed to retrieve or parse log: %s", err.Error())
					}
					if iterator.Event != nil {
						e := iterator.Event
						resultCh <- &event{
							Name:        EventStreamInputChunkAdded,
							StreamID:    e.StreamId,
							ChunkID:     e.ChunkId,
							Address:     e.Raw.Address,
							BlockNumber: e.Raw.BlockNumber,
							BlockHash:   e.Raw.BlockHash,
							TxHash:      e.Raw.TxHash,
							TxIndex:     e.Raw.TxIndex,
						}
						return
					}

					if !iterator.Next() {
						break
					}
				}

				time.Sleep(5 * time.Second)
			}
		}
	}()

	return resultCh, errCh
}
