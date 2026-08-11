#!/usr/bin/env python3
"""Reapply documented endpoint compatibility patches after a pristine reprint."""

from __future__ import annotations

import argparse
import re
from pathlib import Path


CHAT_STREAM_FILES = (
    "internal/cli/gemini-2-5-flash_chat-completions.go",
    "internal/cli/gemini-2-5-pro_chat-completions.go",
    "internal/cli/gemini_3-6-flash.go",
    "internal/cli/gemini_3-flash-v1betamodels.go",
    "internal/cli/gemini-3-pro_chat-completions.go",
    "internal/cli/gemini_3-5-flash.go",
    "internal/cli/promoted_claude.go",
    "internal/cli/promoted_gemini-3-5-flash-openai.go",
    "internal/cli/promoted_gemini-3-1-pro.go",
    "internal/cli/promoted_gemini-3-6-flash-openai.go",
    "internal/cli/promoted_gemini-3-flash.go",
)

STREAM_GUARD = (
    'if bodyStream && !flags.dryRun {\n'
    '\treturn fmt.Errorf("--stream is not supported by this command: the API returns server-sent events '
    'that this CLI cannot decode into JSON/table output; omit --stream (the default) for a normal JSON response")\n'
    "}"
)


class PatchError(RuntimeError):
    """Raised when generated code no longer matches the documented patch shape."""


def patch_chat_stream(text: str, relative_path: str) -> str:
    text, replacements = re.subn(
        r'(cmd\.Flags\(\)\.BoolVar\(&bodyStream,\s*"stream",\s*)true(,)',
        r"\g<1>false\2",
        text,
        count=1,
    )
    if replacements == 0 and not re.search(
        r'cmd\.Flags\(\)\.BoolVar\(&bodyStream,\s*"stream",\s*false,', text
    ):
        raise PatchError(f"{relative_path}: could not locate the --stream flag")

    if "if bodyStream && !flags.dryRun" not in text:
        request = re.search(
            r"(?m)^(?P<indent>\s*)data,\s*(?:statusCode|_),\s*err := c\.Post",
            text,
        )
        if not request:
            raise PatchError(f"{relative_path}: could not locate the request call for the stream guard")
        indent = request.group("indent")
        guard = "\n".join(indent + line for line in STREAM_GUARD.splitlines()) + "\n"
        text = text[: request.start()] + guard + text[request.start() :]
    return text


def patch_grok_status(text: str) -> str:
    if text.count("statusCode") == 1:
        text, replacements = re.subn(
            r"data, statusCode, err := c\.PostWithParamsAndHeaders",
            "data, _, err := c.PostWithParamsAndHeaders",
            text,
            count=1,
        )
        if replacements != 1:
            raise PatchError("internal/cli/grok_responses.go: unused statusCode assignment changed shape")
    return text


def patch_market_stdin(text: str) -> str:
    text = text.replace("io.ReadAll(os.Stdin)", "io.ReadAll(cmd.InOrStdin())")
    if "io.ReadAll(cmd.InOrStdin())" not in text:
        raise PatchError("internal/cli/kie-ai-jobs_market-create-task.go: stdin reader changed shape")
    return text


def patch_first_run_doctor(text: str) -> str:
    replacements = {
        'report["auth_hint"] = "Set it with: kie-pp-cli auth set-token <token> or export KIE_BEARER_AUTH=\\"your-token-here\\""':
            'report["auth_hint"] = "Run kie-pp-cli auth setup in an interactive terminal, or set KIE_BEARER_AUTH through your environment\'s secret store"',
        'report["auth_key_url"] = "https://kie.ai/api-key"':
            'report["auth_key_url"] = cliutil.KieAPIKeyURL',
        "run auth set-token or auth logout": "run auth setup or auth logout",
    }
    for old, new in replacements.items():
        text = text.replace(old, new)
    required = (
        "Run kie-pp-cli auth setup in an interactive terminal",
        "cliutil.KieAPIKeyURL",
        "run auth setup or auth logout to consolidate",
    )
    missing = [marker for marker in required if marker not in text]
    if missing:
        raise PatchError("internal/cli/doctor.go: first-run guidance patch no longer matches: " + ", ".join(missing))
    return text


def patch_first_run_mcp(text: str) -> str:
    text = text.replace(
        'Set it with: kie-pp-cli auth set-token <token> or export KIE_BEARER_AUTH=\\"your-token-here\\"',
        "Run 'kie-pp-cli auth setup' in an interactive terminal, or set KIE_BEARER_AUTH through your environment's secret store.",
    )
    text = text.replace('"\\n      Get a key at: https://kie.ai/api-key"', '"\\n      Get a key at: " + cliutil.KieAPIKeyURL')
    text = text.replace('"key_url": "https://kie.ai/api-key"', '"key_url": cliutil.KieAPIKeyURL')
    if "auth set-token" in text:
        raise PatchError("internal/mcp/tools.go: stale direct-token guidance survived")
    required = ("kie-pp-cli auth setup", "cliutil.KieAPIKeyURL")
    missing = [marker for marker in required if marker not in text]
    if missing:
        raise PatchError("internal/mcp/tools.go: first-run guidance patch no longer matches: " + ", ".join(missing))
    return text


def validate_preserved_extensions(root: Path) -> None:
    required = {
        "internal/cli/auth.go": ("newAuthSetupCmd", "runAuthSetup"),
        "internal/cli/root.go": ("runAuthSetup(cmd, flags, false, nil)", "novelCommandHooks"),
        "internal/cli/helpers.go": ("kie-pp-cli auth setup",),
        "internal/client/client.go": ("authPlaceholderCredentialErrorWithSetup",),
        "internal/config/config.go": ("CredentialConfigured",),
    }
    for relative, markers in required.items():
        path = root / relative
        if not path.is_file():
            raise PatchError(f"missing preserved extension file: {relative}")
        text = path.read_text()
        missing = [marker for marker in markers if marker not in text]
        if missing:
            raise PatchError(f"{relative}: local extension markers disappeared: {', '.join(missing)}")

    doctor = root / "internal/cli/doctor.go"
    original = doctor.read_text()
    patched = patch_first_run_doctor(original)
    if patched != original:
        doctor.write_text(patched)

    mcp_tools = root / "internal/mcp/tools.go"
    if not mcp_tools.is_file():
        raise PatchError("missing preserved extension file: internal/mcp/tools.go")
    original = mcp_tools.read_text()
    patched = patch_first_run_mcp(original)
    if patched != original:
        mcp_tools.write_text(patched)


def patch_tree(root: Path, preserved: bool = False) -> None:
    if preserved:
        validate_preserved_extensions(root)
    for relative in CHAT_STREAM_FILES:
        path = root / relative
        if not path.is_file():
            raise PatchError(f"missing generated endpoint: {relative}")
        original = path.read_text()
        patched = patch_chat_stream(original, relative)
        if patched != original:
            path.write_text(patched)

    grok = root / "internal/cli/grok_responses.go"
    if not grok.is_file():
        raise PatchError("missing generated endpoint: internal/cli/grok_responses.go")
    original = grok.read_text()
    patched = patch_grok_status(original)
    if patched != original:
        grok.write_text(patched)

    market = root / "internal/cli/kie-ai-jobs_market-create-task.go"
    if not market.is_file():
        raise PatchError("missing generated endpoint: internal/cli/kie-ai-jobs_market-create-task.go")
    original = market.read_text()
    patched = patch_market_stdin(original)
    if patched != original:
        market.write_text(patched)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, required=True, help="Root of the freshly printed CLI")
    parser.add_argument(
        "--preserved",
        action="store_true",
        help="Also verify and restore the local first-run extension in an AST-preserved tree",
    )
    args = parser.parse_args()
    patch_tree(args.root.resolve(), preserved=args.preserved)
    print(f"reapplied generated endpoint patches in {args.root.resolve()}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
