// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package main

import "testing"

func TestRequireLoopbackAddress(t *testing.T) {
	for _, addr := range []string{"127.0.0.1:7780", "localhost:7780", "[::1]:7780"} {
		if err := requireLoopbackAddress(addr); err != nil {
			t.Errorf("%s rejected: %v", addr, err)
		}
	}
	for _, addr := range []string{":7780", "0.0.0.0:7780", "192.0.2.1:7780", "bad"} {
		if err := requireLoopbackAddress(addr); err == nil {
			t.Errorf("unsafe address %s accepted", addr)
		}
	}
}
