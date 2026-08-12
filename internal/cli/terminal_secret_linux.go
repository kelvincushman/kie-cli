// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

//go:build linux

package cli

import (
	"bufio"
	"os"

	"golang.org/x/sys/unix"
)

func isTerminalFile(f *os.File) bool {
	_, err := unix.IoctlGetTermios(int(f.Fd()), unix.TCGETS)
	return err == nil
}

func readMaskedTerminalLine(in *os.File) (string, error) {
	oldState, err := unix.IoctlGetTermios(int(in.Fd()), unix.TCGETS)
	if err != nil {
		return "", err
	}
	newState := *oldState
	newState.Lflag &^= unix.ECHO
	if err := unix.IoctlSetTermios(int(in.Fd()), unix.TCSETS, &newState); err != nil {
		return "", err
	}
	defer unix.IoctlSetTermios(int(in.Fd()), unix.TCSETS, oldState) //nolint:errcheck // Restore the user's terminal before returning.

	return bufio.NewReader(in).ReadString('\n')
}
