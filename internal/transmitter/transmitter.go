package transmitter

import (
	"fmt"
	"io"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"

	"github.com/sirupsen/logrus"
)

func init() {
	format.RegisterAll()
}

type TransmitterConfig struct {
	Source      string
	Destination string
	Logger      *logrus.Entry
}

type Transmitter struct {
	source      string
	destination string
	logger      *logrus.Entry

	srcConn av.DemuxCloser
	dstConn av.DemuxCloser
}

func NewTransmitter(c TransmitterConfig) *Transmitter {
	return &Transmitter{
		source:      c.Source,
		destination: c.Destination,
		logger:      c.Logger.WithField("component", "transmitter"),
	}
}

func (t *Transmitter) Start() error {
	srcConn, err := avutil.Open(t.source)
	if err != nil {
		return fmt.Errorf("failed to open source connection: %s", err.Error())
	}
	t.srcConn = srcConn
	defer srcConn.Close()

	dstConn, err := rtmp.Dial(t.destination)
	if err != nil {
		return fmt.Errorf("failed to dial destination connection: %s", err.Error())
	}
	t.dstConn = dstConn
	defer dstConn.Close()

	streams, err := srcConn.Streams()
	if err != nil {
		return fmt.Errorf("failed to dial source connection streams: %s", err.Error())
	}

	if err := dstConn.WriteHeader(streams); err != nil {
		return fmt.Errorf("failed to write header: %s", err.Error())
	}

	for {
		var pkt av.Packet
		if pkt, err = srcConn.ReadPacket(); err != nil {
			if err == io.EOF {
				return nil
			}

			return fmt.Errorf("read source packet failed with error: %s", err)
		}

		err := dstConn.WritePacket(pkt)
		if err != nil {
			return fmt.Errorf("write destination packet failed with error: %s", err)
		}
	}
}

func (t *Transmitter) Stop() {
	if t.srcConn != nil {
		t.srcConn.Close()
	}
	if t.dstConn != nil {
		t.dstConn.Close()
	}
}
