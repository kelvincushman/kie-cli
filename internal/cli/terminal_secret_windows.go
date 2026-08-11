// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

//go:build windows

package cli

import (
	"bufio"
	"os"

	"golang.org/x/sys/windows"
)

func isTerminalFile(f *os.File) bool {
	var mode uint32
	return windows.GetConsoleMode(windows.Handle(f.Fd()), &mode) == nil
}

func readMaskedTerminalLine(in *os.File) (string, error) {
	handle := windows.Handle(in.Fd())
	var oldMode uint32
	if err := windows.GetConsoleMode(handle, &oldMode); err != nil {
		return "", err
	}
	if err := windows.SetConsoleMode(handle, oldMode&^windows.ENABLE_ECHO_INPUT); err != nil {
		return "", err
	}
	defer windows.SetConsoleMode(handle, oldMode) //nolint:errcheck // Restore the user's terminal before returning.

	return bufio.NewReader(in).ReadString('\n')
}
