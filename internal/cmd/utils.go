package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nareix/joy4/av/avutil"
	"golang.org/x/crypto/ssh/terminal"
)

func promptPassword() (string, error) {
	fmt.Println("Enter account password:")
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func probeConnection(source string) error {
	if source == "" {
		return fmt.Errorf("source stream is not set")
	}

	conn, err := avutil.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source connection: %s", err.Error())
	}
	defer conn.Close()

	_, err = conn.Streams()
	if err != nil {
		return fmt.Errorf("failed to acquire source connection streams: %s", err.Error())
	}

	return nil
}

func exitSignal() chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	return done
}
