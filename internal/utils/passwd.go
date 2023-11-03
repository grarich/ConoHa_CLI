/*
Copyright © 2023 grarich <grarich@grawlily.com>
*/
package utils

import (
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/signal"
	"syscall"
)

func ReadPassword() ([]byte, error) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	defer signal.Stop(signalChan)

	cs, err := terminal.GetState(syscall.Stdin)
	if err != nil {
		return nil, err
	}

	go func() {
		<-signalChan
		if err := terminal.Restore(syscall.Stdin, cs); err != nil {
			return
		}
		os.Exit(1)
	}()

	return terminal.ReadPassword(syscall.Stdin)
}
