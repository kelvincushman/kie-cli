import unittest

from patch_generated_cli import (
    patch_chat_stream,
    patch_first_run_doctor,
    patch_first_run_mcp,
    patch_grok_status,
    patch_market_stdin,
)


class GeneratedCLIPatchTests(unittest.TestCase):
    def test_stream_patch_changes_default_and_adds_guard(self):
        source = '''
func generated() {
\tdata, statusCode, err := c.PostWithParams(ctx, path, params, body)
\tcmd.Flags().BoolVar(&bodyStream, "stream", true, "Stream output")
}
'''
        patched = patch_chat_stream(source, "example.go")
        self.assertIn('BoolVar(&bodyStream, "stream", false,', patched)
        self.assertIn("if bodyStream && !flags.dryRun", patched)
        self.assertLess(patched.index("if bodyStream"), patched.index("data, statusCode"))

    def test_stream_patch_is_idempotent(self):
        source = '''
func generated() {
\tif bodyStream && !flags.dryRun {
\t\treturn fmt.Errorf("blocked")
\t}
\tdata, statusCode, err := c.PostWithParams(ctx, path, params, body)
\tcmd.Flags().BoolVar(&bodyStream, "stream", false, "Stream output")
}
'''
        self.assertEqual(patch_chat_stream(source, "example.go"), source)

    def test_grok_discards_only_unused_status(self):
        source = "data, statusCode, err := c.PostWithParamsAndHeaders(ctx, path, params, body, headers)\n"
        self.assertIn("data, _, err :=", patch_grok_status(source))

    def test_grok_keeps_status_when_it_is_consumed(self):
        source = '''
data, statusCode, err := c.PostWithParamsAndHeaders(ctx, path, params, body, headers)
fmt.Println(statusCode)
'''
        self.assertEqual(patch_grok_status(source), source)

    def test_market_stdin_uses_cobra_input(self):
        source = "stdinData, err := io.ReadAll(os.Stdin)\n"
        self.assertIn("io.ReadAll(cmd.InOrStdin())", patch_market_stdin(source))

    def test_doctor_restores_guided_auth_setup(self):
        source = '''
report["auth_hint"] = "Set it with: kie-pp-cli auth set-token <token> or export KIE_BEARER_AUTH=\\"your-token-here\\""
report["auth_key_url"] = "https://kie.ai/api-key"
if keyURL, ok := report["auth_key_url"]; ok {
\tfmt.Fprintf(w, "  Get a key at: %v\\n", keyURL)
}
warning = "; run auth set-token or auth logout to consolidate and remove legacy secrets"
warning = "; run auth set-token or auth logout to consolidate"
'''
        patched = patch_first_run_doctor(source)
        self.assertIn("kie-pp-cli auth setup in an interactive terminal", patched)
        self.assertIn("cliutil.KieAPIKeyURL", patched)
        self.assertIn('report["auth_key_url_type"] = "affiliate"', patched)
        self.assertIn('report["auth_affiliate_disclosure"] = cliutil.KieAffiliateDisclosure', patched)
        url_renderer = 'fmt.Fprintf(w, "  Get a key at: %v\\n", keyURL)'
        disclosure_renderer = 'if disclosure, ok := report["auth_affiliate_disclosure"]; ok {'
        self.assertEqual(patched.count(disclosure_renderer), 1)
        self.assertLess(patched.index(url_renderer), patched.index(disclosure_renderer))
        self.assertNotIn("auth set-token", patched)

    def test_doctor_repairs_each_partial_affiliate_field_independently(self):
        source = '''
report["auth_hint"] = "Run kie-pp-cli auth setup in an interactive terminal"
report["auth_key_url"] = cliutil.KieAPIKeyURL
report["auth_affiliate_disclosure"] = cliutil.KieAffiliateDisclosure
if keyURL, ok := report["auth_key_url"]; ok {
\tfmt.Fprintf(w, "  Get a key at: %v\\n", keyURL)
}
warning = "; run auth setup or auth logout to consolidate"
'''
        patched = patch_first_run_doctor(source)
        self.assertIn('report["auth_key_url_type"] = "affiliate"', patched)
        self.assertEqual(patched.count('report["auth_affiliate_disclosure"] = cliutil.KieAffiliateDisclosure'), 1)
        self.assertIn('if disclosure, ok := report["auth_affiliate_disclosure"]; ok {', patched)

    def test_mcp_restores_guided_auth_setup(self):
        source = '''
hint := "\\n      Set it with: kie-pp-cli auth set-token <token> or export KIE_BEARER_AUTH=\\\"your-token-here\\\"" +
    "\\n      Get a key at: https://kie.ai/api-key" +
    "\\n      Run doctor"
context := map[string]any{
    "key_url": "https://kie.ai/api-key",
}
'''
        patched = patch_first_run_mcp(source)
        self.assertIn("kie-pp-cli auth setup", patched)
        self.assertIn("cliutil.KieAPIKeyURL", patched)
        self.assertIn('"key_url_type": "affiliate"', patched)
        self.assertIn('"affiliate_disclosure": cliutil.KieAffiliateDisclosure,', patched)
        self.assertIn('"\\n      " + cliutil.KieAffiliateDisclosure +', patched)
        self.assertNotIn("auth set-token", patched)

    def test_mcp_repairs_partial_error_and_metadata_disclosures(self):
        source = '''
hintA := "\\n      Run 'kie-pp-cli auth setup'" +
    "\\n      Get a key at: " + cliutil.KieAPIKeyURL +
    "\\n      " + cliutil.KieAffiliateDisclosure +
    "\\n      Run doctor"
hintB := "\\n      Run 'kie-pp-cli auth setup'" +
    "\\n      Get a key at: " + cliutil.KieAPIKeyURL +
    "\\n      Run doctor"
context := map[string]any{
    "key_url": cliutil.KieAPIKeyURL,
    "affiliate_disclosure": cliutil.KieAffiliateDisclosure,
}
'''
        patched = patch_first_run_mcp(source)
        self.assertEqual(patched.count('"\\n      Get a key at: " + cliutil.KieAPIKeyURL +'), 2)
        self.assertEqual(patched.count('"\\n      " + cliutil.KieAffiliateDisclosure +'), 2)
        self.assertIn('"key_url_type": "affiliate",', patched)
        self.assertEqual(patched.count('"affiliate_disclosure": cliutil.KieAffiliateDisclosure,'), 1)


if __name__ == "__main__":
    unittest.main()
