// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
	"kie-pp-cli/internal/cliutil"
	"kie-pp-cli/internal/cliutil/testenv"
	"kie-pp-cli/internal/config"
)

func isolateAuthSetup(t *testing.T) {
	t.Helper()
	restore, err := cliutil.SetHomeOverride("")
	if err != nil {
		t.Fatalf("reset home override: %v", err)
	}
	t.Cleanup(restore)
	testenv.Isolate(t, cliutil.ConfigDir, cliutil.DataDir, cliutil.StateDir, cliutil.CacheDir)
	t.Setenv("KIE_BEARER_AUTH", "")
}

func executeRoot(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	root := RootCmd()
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)
	err := root.Execute()
	return stdout.String(), stderr.String(), err
}

func TestFirstRunRoutesToSafeSetupHint(t *testing.T) {
	isolateAuthSetup(t)

	stdout, stderr, err := executeRoot(t)
	if err != nil {
		t.Fatalf("root command error = %v, stderr = %s", err, stderr)
	}
	if !strings.Contains(stdout, "support continued development") || !strings.Contains(stdout, cliutil.KieAPIKeyURL) {
		t.Fatalf("first-run output did not route to auth setup: %q", stdout)
	}
	if !strings.Contains(stdout, cliutil.KieAffiliateDisclosure) {
		t.Fatalf("first-run referral was not disclosed: %q", stdout)
	}
	if !strings.Contains(stdout, "interactive terminal") {
		t.Fatalf("redirected first-run output must not prompt: %q", stdout)
	}
}

func TestFirstRunMachineModesDoNotPrompt(t *testing.T) {
	isolateAuthSetup(t)

	stdout, stderr, err := executeRoot(t, "--json")
	if err != nil {
		t.Fatalf("json first-run error = %v, stderr = %s", err, stderr)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("json first-run output is invalid: %v; output=%q", err, stdout)
	}
	if result["setup_required"] != true || result["next_step"] != "kie-pp-cli auth setup" {
		t.Fatalf("unexpected first-run JSON: %#v", result)
	}
	if result["get_api_key"] != cliutil.KieAPIKeyURL || result["get_api_key_link_type"] != "affiliate" || result["affiliate_disclosure"] != cliutil.KieAffiliateDisclosure {
		t.Fatalf("first-run JSON referral metadata = %#v", result)
	}

	stdout, stderr, err = executeRoot(t, "--no-input")
	if err != nil {
		t.Fatalf("non-interactive first-run error = %v, stderr = %s", err, stderr)
	}
	if !strings.Contains(stdout, "interactive terminal") || strings.Contains(stdout, "input hidden") {
		t.Fatalf("--no-input must not start the wizard: %q", stdout)
	}
}

func TestDoctorHumanOutputDisclosesAffiliateLinkImmediately(t *testing.T) {
	isolateAuthSetup(t)
	configPath := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(configPath, []byte("base_url = \"\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	stdout, stderr, err := executeRoot(t, "--config", configPath, "doctor")
	if err != nil {
		t.Fatalf("doctor error = %v, stderr = %s", err, stderr)
	}
	urlIndex := strings.Index(stdout, "Get a key at: "+cliutil.KieAPIKeyURL)
	disclosureIndex := strings.Index(stdout, cliutil.KieAffiliateDisclosure)
	if urlIndex < 0 || disclosureIndex < 0 || disclosureIndex <= urlIndex {
		t.Fatalf("doctor must disclose the affiliate relationship immediately after its referral URL: %q", stdout)
	}
	if strings.Count(stdout[urlIndex:disclosureIndex], "\n") != 1 {
		t.Fatalf("doctor disclosure was not immediately after its referral URL: %q", stdout)
	}
}

func TestFirstRunAgentModeAndExplicitSetupDoNotPrompt(t *testing.T) {
	isolateAuthSetup(t)

	stdout, stderr, err := executeRoot(t, "--agent")
	if err != nil {
		t.Fatalf("agent first-run error = %v, stderr = %s", err, stderr)
	}
	if !strings.Contains(stdout, "setup_required") || strings.Contains(stdout, "input hidden") {
		t.Fatalf("agent mode must return a non-blocking setup hint: %q", stdout)
	}

	stdout, stderr, err = executeRoot(t, "auth", "setup", "--no-input")
	if err != nil {
		t.Fatalf("explicit non-interactive setup error = %v, stderr = %s", err, stderr)
	}
	if !strings.Contains(stdout, "support continued development") || !strings.Contains(stdout, cliutil.KieAffiliateDisclosure) || strings.Contains(stdout, "input hidden") {
		t.Fatalf("explicit non-interactive setup must return directions, not prompt: %q", stdout)
	}
}

func TestMediaSetupReturnsDisclosedAffiliateOnboarding(t *testing.T) {
	isolateAuthSetup(t)

	stdout, stderr, err := executeRoot(t, "media", "setup", "--agent")
	if err != nil {
		t.Fatalf("media setup error = %v, stderr = %s", err, stderr)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("media setup output is invalid: %v; output=%q", err, stdout)
	}
	setup, ok := result["results"].(map[string]any)
	if !ok {
		t.Fatalf("media setup agent envelope = %#v", result)
	}
	if setup["auth_configured"] != false || setup["get_api_key"] != cliutil.KieAPIKeyURL {
		t.Fatalf("media setup referral metadata = %#v", setup)
	}
	if setup["get_api_key_link_type"] != "affiliate" || setup["affiliate_disclosure"] != cliutil.KieAffiliateDisclosure {
		t.Fatalf("media setup referral was not explicitly disclosed: %#v", setup)
	}
}

func TestHelpAndVersionDoNotRouteToSetup(t *testing.T) {
	isolateAuthSetup(t)

	for _, args := range [][]string{{"--help"}, {"--version"}} {
		stdout, stderr, err := executeRoot(t, args...)
		if err != nil {
			t.Fatalf("%v error = %v, stderr = %s", args, err, stderr)
		}
		if strings.Contains(stdout, "Get API key:") {
			t.Fatalf("%v unexpectedly started setup: %q", args, stdout)
		}
	}
}

func TestAuthOffersNoTokenArgumentPath(t *testing.T) {
	isolateAuthSetup(t)

	stdout, stderr, err := executeRoot(t, "auth", "--help")
	if err != nil {
		t.Fatalf("auth help error = %v, stderr = %s", err, stderr)
	}
	if strings.Contains(stdout, "set-token") {
		t.Fatalf("auth help still offers a shell-history token path: %q", stdout)
	}

	const fakeKey = "fake-key-must-not-appear"
	stdout, stderr, err = executeRoot(t, "auth", "set-token", fakeKey)
	if err == nil {
		t.Fatal("removed token-argument command unexpectedly succeeded")
	}
	if strings.Contains(stdout, fakeKey) || strings.Contains(stderr, fakeKey) {
		t.Fatal("rejected token argument was written to command output")
	}
}

func TestGuidedAuthSetupSavesWithoutWritingKeyToOutput(t *testing.T) {
	isolateAuthSetup(t)

	var stdout, stderr bytes.Buffer
	cmd := RootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	fakeKey := "fake-key-for-terminal-capture"
	if err := runAuthSetup(cmd, &rootFlags{}, false, func(*cobra.Command) (string, error) {
		return fakeKey, nil
	}); err != nil {
		t.Fatalf("guided setup error = %v, stderr = %s", err, stderr.String())
	}
	if strings.Contains(stdout.String(), fakeKey) || strings.Contains(stderr.String(), fakeKey) {
		t.Fatal("guided setup wrote the API key to command output")
	}
	if !strings.Contains(stdout.String(), cliutil.KieAPIKeyURL) || !strings.Contains(stdout.String(), cliutil.KieAffiliateDisclosure) {
		t.Fatalf("guided setup omitted disclosed referral: %q", stdout.String())
	}
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("load saved config: %v", err)
	}
	if !cfg.CredentialConfigured() {
		t.Fatal("guided setup did not save a usable credential")
	}
	path, err := cliutil.CredentialsFilePath()
	if err != nil {
		t.Fatalf("credentials path: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat credentials: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("credentials mode = %o, want 600", info.Mode().Perm())
	}
}

const authPTYPrompt = "API key (input hidden): "

type authPTYCapture struct {
	mu     sync.Mutex
	output bytes.Buffer
	prompt chan struct{}
}

func (c *authPTYCapture) Write(data []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, err := c.output.Write(data)
	if strings.Contains(c.output.String(), authPTYPrompt) {
		select {
		case c.prompt <- struct{}{}:
		default:
		}
	}
	return n, err
}

func (c *authPTYCapture) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.output.String()
}

func runAuthPTY(t *testing.T, args []string, devNull bool, key string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestAuthPTYHelper$")
	child.Env = append(os.Environ(),
		"KIE_AUTH_PTY_HELPER=1",
		"KIE_AUTH_PTY_ARGS="+strings.Join(args, " "),
		"KIE_AUTH_PTY_DEV_NULL="+map[bool]string{true: "1", false: ""}[devNull],
		"KIE_BEARER_AUTH=",
		"KIE_NO_LEARN=true",
	)
	ptyFile, err := pty.Start(child)
	if errors.Is(err, pty.ErrUnsupported) {
		t.Skip("PTYs are not supported on this platform")
	}
	if err != nil {
		t.Fatalf("start PTY child: %v", err)
	}
	defer ptyFile.Close() //nolint:errcheck // Test cleanup.

	capture := &authPTYCapture{prompt: make(chan struct{}, 1)}
	readDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(capture, ptyFile)
		close(readDone)
	}()

	if key != "" {
		select {
		case <-capture.prompt:
			// The prompt is written immediately before masked input starts.
			// Let the child finish its terminal-mode change before sending input.
			time.Sleep(25 * time.Millisecond)
			if _, err := io.WriteString(ptyFile, key+"\n"); err != nil {
				t.Fatalf("write API key to PTY: %v", err)
			}
		case <-ctx.Done():
			t.Fatal("PTY child did not prompt for the API key")
		}
	}

	if err := child.Wait(); err != nil {
		t.Fatalf("PTY child failed: %v", err)
	}
	if ctx.Err() != nil {
		t.Fatalf("PTY child timed out: %v", ctx.Err())
	}
	_ = ptyFile.Close()
	<-readDone
	return capture.String()
}

// TestAuthPTYHelper runs the command in a fresh process so Cobra uses the PTY
// as its real stdin and stdout. It is only invoked by runAuthPTY.
func TestAuthPTYHelper(t *testing.T) {
	if os.Getenv("KIE_AUTH_PTY_HELPER") != "1" {
		return
	}

	cmd := RootCmd()
	cmd.SetArgs(strings.Fields(os.Getenv("KIE_AUTH_PTY_ARGS")))
	if os.Getenv("KIE_AUTH_PTY_DEV_NULL") == "1" {
		in, err := os.Open(os.DevNull)
		if err != nil {
			t.Fatalf("open %s: %v", os.DevNull, err)
		}
		defer in.Close() //nolint:errcheck // Test cleanup.
		cmd.SetIn(in)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func TestFirstRunPTYMasksInputAndSavesReusableCredential(t *testing.T) {
	isolateAuthSetup(t)

	key := "fake-key-for-pty-capture"
	output := runAuthPTY(t, nil, false, key)
	if !strings.Contains(output, "Kie API key setup") || !strings.Contains(output, authPTYPrompt) {
		t.Fatal("bare first run did not start the guided setup wizard")
	}
	if strings.Contains(output, key) {
		t.Fatal("guided setup echoed the API key")
	}

	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("load saved config: %v", err)
	}
	if !cfg.CredentialConfigured() || cfg.AuthHeader() == "" {
		t.Fatal("guided setup did not save a reusable credential")
	}
}

func TestFirstRunDevNullPTYDoesNotPrompt(t *testing.T) {
	isolateAuthSetup(t)

	for _, args := range [][]string{nil, {"auth", "setup"}} {
		t.Run(strings.Join(append([]string{"bare"}, args...), "_"), func(t *testing.T) {
			output := runAuthPTY(t, args, true, "")
			if !strings.Contains(output, "interactive terminal") || strings.Contains(output, authPTYPrompt) {
				t.Fatal("/dev/null input must return setup directions without prompting")
			}
		})
	}
}
