#!/usr/bin/env bash
set -euo pipefail
artifact_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "$artifact_dir/../../.." && pwd)
linter_source=${BATEIA_LINTER_SOURCE:-/repos/Antigravity/tekt-linter}
expected_commit=49f46885fd03f753e9f1cf37e271d6baf30ed725
if [ "$(git -C "$linter_source" rev-parse HEAD)" != "$expected_commit" ]; then
  echo "Linter source must be at $expected_commit" >&2
  exit 2
fi
if [ -n "$(git -C "$linter_source" status --porcelain)" ]; then
  echo "Linter source must have a clean worktree" >&2
  exit 2
fi
linter_build=$(mktemp -d /tmp/bateia-0035-linter-reproduce.XXXXXX)
# Source remains read-only; the complete patch is applied to a temporary copy.
cp -a "$linter_source/." "$linter_build/"
git -C "$linter_build" apply --check "$artifact_dir/crystalline-lint.patch"
git -C "$linter_build" apply "$artifact_dir/crystalline-lint.patch"
linter_target=${CARGO_TARGET_DIR:-$linter_build/target}
cargo test --offline --manifest-path "$linter_build/Cargo.toml"
cargo build --offline --manifest-path "$linter_build/Cargo.toml"
python3 "$artifact_dir/lint-regressions.py" "$linter_target/debug/crystalline-lint"
"$linter_target/debug/crystalline-lint" "$repo_root"
printf 'Validated binary: %s\n' "$linter_target/debug/crystalline-lint"
