#!/usr/bin/env bash
set -euo pipefail
artifact_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "$artifact_dir/../../../.." && pwd)
linter_source=${BATEIA_LINTER_SOURCE:-/repos/Antigravity/tekt-linter}
revision=e53de6fd2bff9f3a87a3e8a076f5e470ce4bb7a0
linter_build=$(mktemp -d /tmp/bateia-0129-linter-reproduce.XXXXXX)
# Immutable source, isolated build: no external checkout or global installation is modified.
git -C "$linter_source" archive "$revision" | tar -x -C "$linter_build"
git -C "$linter_build" apply --check "$artifact_dir/crystalline-go-parser.patch"
git -C "$linter_build" apply "$artifact_dir/crystalline-go-parser.patch"
cargo test --offline --lib --manifest-path "$linter_build/Cargo.toml"
cargo build --offline --manifest-path "$linter_build/Cargo.toml"
linter_target=${CARGO_TARGET_DIR:-$linter_build/target}
python3 "$artifact_dir/lint-regressions.py" "$linter_target/debug/crystalline-lint"
# Run the required command with the validated parser on this process's PATH.
(cd "$repo_root" && PATH="$linter_target/debug:$PATH" crystalline-lint --fail-on warning .)
printf 'Validated binary: %s\n' "$linter_target/debug/crystalline-lint"
