#!/usr/bin/env bash
set -u

repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$repo_root" || exit 2

sha256sum --check 04_wiring/step0071_gate/frozen-inputs.sha256 || exit 3

GOCACHE=/tmp/bateia-step0071-gocache go test ./04_wiring/step0071_gate -run '^TestStep0071Gate' -count=1
