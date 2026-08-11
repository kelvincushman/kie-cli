// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"kie-pp-cli/internal/cli"
	"kie-pp-cli/internal/mediamcp"
)

const defaultHTTPAddr = "127.0.0.1:7780"

var version = "0.0.0-dev"

func main() {
	_ = os.Setenv("KIE_LEARN_SURFACE", "mcp")
	if err := cli.BindMCPServerProfile(); err != nil {
		fail("MCP client-profile bind failed: %v", err)
	}

	transport := flag.String("transport", defaultTransport(), "MCP transport: stdio | http")
	addr := flag.String("addr", defaultHTTPAddr, "loopback bind address for the stateless http transport")
	flag.Parse()

	server := mediamcp.NewServer(version, nil)
	switch strings.ToLower(strings.TrimSpace(*transport)) {
	case "stdio":
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			fail("MCP server error: %v", err)
		}
	case "http":
		if err := requireLoopbackAddress(*addr); err != nil {
			fail("unsafe --addr: %v", err)
		}
		httpServer := &http.Server{
			Addr:              *addr,
			Handler:           mediamcp.NewHTTPHandler(server),
			ReadHeaderTimeout: 10 * time.Second,
		}
		fmt.Fprintf(os.Stderr, "kie-media-mcp serving stateless MCP 2026-07-28 over HTTP at %s\n", *addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fail("MCP server error: %v", err)
		}
	default:
		fail("unknown --transport %q (supported: stdio, http)", *transport)
	}
}

func defaultTransport() string {
	if value := strings.TrimSpace(os.Getenv("KIE_MEDIA_MCP_TRANSPORT")); value != "" {
		return value
	}
	return "stdio"
}

func requireLoopbackAddress(addr string) error {
	host, _, err := net.SplitHostPort(strings.TrimSpace(addr))
	if err != nil {
		return fmt.Errorf("expected host:port: %w", err)
	}
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("kie-media-mcp is local-only; use 127.0.0.1, [::1], or localhost")
	}
	return nil
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
