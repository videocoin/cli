package emitter

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/VideoCoin/common/bcops"
	sm "github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/accounts/keystore"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/VideoCoin/cli/internal/listener"
	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type EmitterManagerConfig struct {
	NodeRPCAddr     string
	ContractAddress string
	Key             *keystore.Key
	Logger          *logrus.Entry
}

type emitterManager struct {
	ethClient     *ethclient.Client
	smManager     *sm.Manager
	eventListener *listener.EventListener
	transactOpts  *bind.TransactOpts
	key           *keystore.Key
	logger        *logrus.Entry
}

func NewEmitterManager(c EmitterManagerConfig) (*emitterManager, error) {
	client, err := ethclient.Dial(c.NodeRPCAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial eth client: %s", err.Error())
	}

	managerAddress := common.HexToAddress(c.ContractAddress)
	manager, err := sm.NewManager(managerAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create smart contract stream manager: %s", err.Error())
	}

	eventListenerConfig := &listener.EventListenerConfig{
		SmartContractManager: manager,
		Timeout:              60,
		Logger:               c.Logger,
	}
	eventListener := listener.NewEventListener(eventListenerConfig)

	transactOpts, err := bcops.GetBCAuth(client, c.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to init blockchain auth: %s", err.Error())
	}

	return &emitterManager{
		ethClient:     client,
		smManager:     manager,
		eventListener: eventListener,
		transactOpts:  transactOpts,
		key:           c.Key,
		logger:        c.Logger.WithField("component", "emitter"),
	}, nil
}

func (s *emitterManager) RequestStream() (*big.Int, error) {
	streamID := big.NewInt(int64(rand.Intn(math.MaxInt64)))

	_, err := s.smManager.RequestStream(
		s.transactOpts,
		streamID,
		"videocoin",
		[]*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(2)},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to request stream: %s", err.Error())
	}

	go func() {
		resultCh, errCh := s.eventListener.LogStreamRequestEvent(streamID, s.key.Address)

		select {
		case err := <-errCh:
			s.logger.WithError(err).Errorf("failed to watch stream requested")
		case e := <-resultCh:
			s.logger.Infof("received an event:%s\n", e.String())
		}
	}()

	go func() {
		resultCh, errCh := s.eventListener.LogStreamApproveEvent(streamID)

		select {
		case err := <-errCh:
			s.logger.WithError(err).Errorf("failed to watch stream approved")
		case e := <-resultCh:
			s.logger.Infof("received an event:%s\n", e.String())
		}
	}()

	return streamID, nil
}

func (s *emitterManager) CreateStream(streamID *big.Int) (string, error) {
	var i, e = big.NewInt(10), big.NewInt(19)
	s.transactOpts.Value = i.Exp(i, e, nil)
	s.transactOpts.From = s.key.Address

	_, err := s.smManager.CreateStream(
		s.transactOpts,
		streamID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create stream: %s", err.Error())
	}

	resultCh, errCh := s.eventListener.LogStreamCreateEvent(streamID)

	select {
	case err := <-errCh:
		return "", fmt.Errorf("failed to watch stream created: %s", err.Error())
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("failed to watch stream created: timeout")
	case e := <-resultCh:
		s.logger.Infof("received an event:%s\n", e.String())
		return e.StreamAddress.Hex(), nil
	}
}

func (s *emitterManager) GetAddressBalance() (*big.Float, error) {
	wei, err := s.ethClient.BalanceAt(context.Background(), s.key.Address, nil)
	if err != nil {
		return nil, err
	}
	vdc, err := convertWeiToVDC(wei)
	if err != nil {
		return nil, err
	}

	return vdc, nil
}
