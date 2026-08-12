// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

//go:build !darwin && !linux && !windows

package cli

import (
	"fmt"
	"os"
)

func isTerminalFile(_ *os.File) bool {
	return false
}

func readMaskedTerminalLine(_ *os.File) (string, error) {
	return "", fmt.Errorf("masked terminal input is not supported on this platform")
}
