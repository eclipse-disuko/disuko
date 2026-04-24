# SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
#
# SPDX-License-Identifier: Apache-2.0

#!/usr/bin/env python3
"""
Generates and reformats backend/docs/swagger.yaml from Swaggo annotations.

Steps:
  1. Install swag if not on PATH
  2. Run `swag init` to regenerate swagger.yaml
  3. Reorder top-level YAML keys
  4. Inline wrapped single-quoted description scalars -> double-quoted
  5. Convert single-quoted $ref scalars -> double-quoted
  6. Persist the result with a single trailing newline

Execute this script via `python3 backend/docs/generate-swagger-spec.py`
"""

import re
import shutil
import subprocess
import sys
from pathlib import Path

BACKEND_DIR = Path(__file__).resolve().parents[1]
SWAGGER_FILE = BACKEND_DIR / "docs" / "swagger.yaml"

DESIRED_KEY_ORDER = [
    "swagger",
    "host",
    "basePath",
    "info",
    "schemes",
    "securityDefinitions",
    "paths",
    "definitions",
    "externalDocs",
]


def ensure_swag() -> str:
    if path := shutil.which("swag"):
        return path
    print("swag not found on PATH – installing via go install …", file=sys.stderr)
    subprocess.run(
        ["go", "install", "github.com/swaggo/swag/cmd/swag@latest"],
        check=True,
    )
    path = shutil.which("swag")
    if not path:
        # Try GOPATH/bin explicitly
        gopath = subprocess.check_output(["go", "env", "GOPATH"], text=True).strip()
        candidate = Path(gopath) / "bin" / "swag"
        if candidate.exists():
            return str(candidate)
        sys.exit("ERROR: swag not found after installation. Add $(go env GOPATH)/bin to $PATH.")
    return path


def run_swag(swag_bin: str) -> None:
    print("Running swag init …", file=sys.stderr)
    subprocess.run(
        [
            swag_bin,
            "init",
            "--ot", "yaml",
            "-g", "server/routes.go",
            "--output", "docs",
        ],
        cwd=BACKEND_DIR,
        check=True,
    )
    print(f"Generated {SWAGGER_FILE}", file=sys.stderr)


def reorder_top_level_keys(lines: list[str]) -> list[str]:
    top_re = re.compile(r'^([a-zA-Z_][a-zA-Z0-9_]*)\s*:')

    blocks: dict[str, list[str]] = {}
    block_order: list[str] = []
    current_key: str | None = None
    current_lines: list[str] = []

    for line in lines:
        is_top = (
            not line.startswith(' ')
            and not line.startswith('-')
            and top_re.match(line)
        )
        if is_top:
            if current_key is not None:
                blocks[current_key] = current_lines
            current_key = is_top.group(1)
            current_lines = [line]
            if current_key not in block_order:
                block_order.append(current_key)
        elif current_key is not None:
            current_lines.append(line)

    if current_key is not None:
        blocks[current_key] = current_lines

    ordered: list[str] = []
    for key in DESIRED_KEY_ORDER:
        if key in blocks:
            ordered.extend(blocks[key])
    for key in block_order:
        if key not in DESIRED_KEY_ORDER and key in blocks:
            ordered.extend(blocks[key])

    return ordered


def transform_quotes(lines: list[str]) -> list[str]:
    """Inline wrapped single-quoted descriptions and convert $ref quotes."""
    result: list[str] = []
    i = 0
    while i < len(lines):
        line = lines[i]

        # Single-quoted description, fully on one line: description: 'foo'
        m_single = re.match(r"^(\s*(?:- )?description: )'(.*)'$", line)
        # Open single-quoted description (no closing quote on this line)
        m_open = re.match(r"^(\s*(?:- )?description: )'(.*)$", line)

        if m_single:
            content = m_single.group(2).replace("''", "'").replace('"', '\\"')
            result.append(f'{m_single.group(1)}"{content}"')
            i += 1
        elif m_open and not line.rstrip().endswith("'"):
            prefix = m_open.group(1)
            accumulated = m_open.group(2)
            i += 1
            while i < len(lines):
                cont = lines[i].strip()
                if cont.endswith("'"):
                    accumulated += ' ' + cont[:-1]
                    i += 1
                    break
                accumulated += ' ' + cont
                i += 1
            content = accumulated.replace("''", "'").replace('"', '\\"')
            result.append(f'{prefix}"{content}"')
        else:
            # Convert single-quoted $ref -> double-quoted
            line = re.sub(r"\$ref: '(#[^']+)'", r'$ref: "\1"', line)
            result.append(line)
            i += 1

    return result


def main() -> None:
    swag_bin = ensure_swag()
    run_swag(swag_bin)

    with open(SWAGGER_FILE, "r") as f:
        lines = f.read().splitlines()

    lines = reorder_top_level_keys(lines)
    lines = transform_quotes(lines)

    with open(SWAGGER_FILE, "w") as f:
        f.write('\n'.join(lines) + '\n')

    print(f"Reformatted {SWAGGER_FILE} ({len(lines)} lines)", file=sys.stderr)


if __name__ == "__main__":
    main()
