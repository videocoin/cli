package cmd

import (
	"fmt"
	"time"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/cli/internal/cloud"
	"github.com/VideoCoin/cli/internal/emitter"
	"github.com/VideoCoin/cli/internal/key"
	"github.com/VideoCoin/cli/internal/transmitter"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdStart = &cobra.Command{
	Use:   "start [rtmp-address]",
	Short: "start streaming to VideoCoin testnet",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fflags := cmd.Flags()
		account, err := fflags.GetString("account")
		if account == "" || err != nil {
			logrus.WithError(err).Fatal("account file is missed")
		}

		password, err := fflags.GetString("password")
		if password == "" || err != nil {
			password, err = promptPassword()
			if password == "" || err != nil {
				logrus.WithError(err).Fatal("account password is missed")
			}

			err = fflags.Set("password", password)
			if err != nil {
				logrus.WithError(err).Fatal("failed to set password flag")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		logger := c.Logger
		account, _ := fflags.GetString("account")
		password, _ := fflags.GetString("password")

		sourceRtmpUrl := args[0]
		err := probeConnection(sourceRtmpUrl)
		if err != nil {
			logger.WithError(err).Fatal("failed to probe input rtmp url")
		}

		spinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		spinner.Start()
		defer spinner.Stop()

		ks := key.NewKeyStore()
		key, err := ks.ImportKey(account, password)
		if err != nil {
			logger.WithError(err).Fatal("failed to import account")
		}

		em, err := emitter.NewEmitterManager(
			emitter.EmitterManagerConfig{
				NodeRPCAddr:     c.NodeRPCAddr,
				ContractAddress: c.ContractAddress,
				Key:             key,
				Logger:          logrus.NewEntry(logger.Logger),
			},
		)
		if err != nil {
			logger.WithError(err).Fatal("failed to create a stream manager")
		}

		balance, err := em.GetAddressBalance()
		if err != nil {
			logger.WithError(err).Fatal("failed to get account balance")
		}

		fbalance, _ := balance.Float64()
		if fbalance < float64(c.MinVDCBalance) {
			logger.WithField("balance", fbalance).Fatalf(
				"insufficient account balance, must be minimum %d VDC available", c.MinVDCBalance)
		}

		cm := cloud.NewCloudManager(
			cloud.CloudManagerConfig{
				ManagerAddr: c.ManagerAddr,
				Logger:      logrus.NewEntry(logger.Logger),
			},
		)

		streamID, err := em.RequestStream()
		if err != nil {
			logger.WithError(err).Fatal("failed to request stream")
		}

		logger.Infof("acquired stream id %s", streamID.String())

		destinationRtmpUrl, err := cm.CreateJob(streamID, key.Address.String())
		if err != nil {
			logger.WithError(err).Fatal("failed to create job")
		}

		logger.Infof("acquired destination rtmp url %s", destinationRtmpUrl)

		_, err = cm.AwaitJobStatus(streamID, proto.WorkOrderStatusApproved)
		if err != nil {
			logger.WithError(err).Fatal("failed to get approved job")
		}

		logger.Infof("acquired approved job")

		contractAddress, err := em.CreateStream(streamID)
		if err != nil {
			logger.WithError(err).Fatal("failed to create stream")
		}

		logger.Infof("acquired stream address %s", contractAddress)

		err = cm.UpdateJobContractAddress(streamID, contractAddress)
		if err != nil {
			logger.WithError(err).Fatal("failed to update job")
		}

		tc := transmitter.TransmitterConfig{
			Source:      sourceRtmpUrl,
			Destination: destinationRtmpUrl,
			Logger:      logrus.NewEntry(logger.Logger),
		}
		transmitter := transmitter.NewTransmitter(tc)

		go func() {
			err := transmitter.Start()
			if err != nil {
				logger.WithError(err).Fatal("failed to start transmitter")
			}
		}()

		job, err := cm.AwaitJobStatus(streamID, proto.WorkOrderStatusReady)
		if err != nil {
			logger.WithError(err).Fatal("failed to get ready job")
		}

		spinner.Stop()
		fmt.Printf(
			"Your stream is going to be available shortly. Use next URL to access it: %s\n", job.OutputURL)

		fmt.Println("Stop streaming with cmd+c.")

		done := exitSignal()
		<-done
		transmitter.Stop()
	},
}
