// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"kie-pp-cli/internal/cli"
	"kie-pp-cli/internal/cliutil"
	"kie-pp-cli/internal/media"
	"kie-pp-cli/internal/mediamcp"
)

func TestGrillMeCLIAndMCPReturnEquivalentDirectorState(t *testing.T) {
	restoreHome, err := cliutil.SetHomeOverride("")
	if err != nil {
		t.Fatal(err)
	}
	defer restoreHome()

	request := "Create a 5 second vertical silent text-to-video TikTok video ad in a cinematic documentary style"
	cliTurn := runCLIGrill(t, filepath.Join(t.TempDir(), "cli-home"), request)
	mcpTurn := runMCPGrill(t, filepath.Join(t.TempDir(), "mcp-media"), request)

	if got, want := canonicalTurn(cliTurn), canonicalTurn(mcpTurn); !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.MarshalIndent(got, "", "  ")
		wantJSON, _ := json.MarshalIndent(want, "", "  ")
		t.Fatalf("CLI and MCP director state differ\nCLI:\n%s\nMCP:\n%s", gotJSON, wantJSON)
	}
}

func runCLIGrill(t *testing.T, home, request string) media.Turn {
	t.Helper()
	cmd := cli.RootCmd()
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--home", home, "--agent", "grill-me", request})
	if _, err := cmd.ExecuteC(); err != nil {
		t.Fatalf("CLI grill failed: %v\nstderr: %s\nstdout: %s", err, stderr.String(), stdout.String())
	}
	var envelope struct {
		Results media.Turn `json:"results"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatalf("decode CLI grill: %v\n%s", err, stdout.String())
	}
	return envelope.Results
}

func runMCPGrill(t *testing.T, storeRoot, request string) media.Turn {
	t.Helper()
	store := media.NewStore(storeRoot)
	server := mediamcp.NewServer("test", &mediamcp.Dependencies{
		Store: func() (*media.Store, error) { return store, nil },
	})
	httpServer := httptest.NewServer(mediamcp.NewHTTPHandler(server))
	defer httpServer.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "parity-test", Version: "test"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "media_grill_start",
		Arguments: map[string]any{"request": request},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("MCP grill failed: %#v", result.Content)
	}
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var turn media.Turn
	if err := json.Unmarshal(encoded, &turn); err != nil {
		t.Fatal(err)
	}
	return turn
}

func canonicalTurn(turn media.Turn) media.Turn {
	turn.ResumeCommand = ""
	if turn.Brief == nil {
		return turn
	}
	copy := *turn.Brief
	copy.ID = ""
	copy.CreatedAt = time.Time{}
	copy.UpdatedAt = time.Time{}
	turn.Brief = &copy
	return turn
}
