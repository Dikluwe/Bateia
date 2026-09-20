#!/bin/sh
set -eu

cd "$(dirname "$0")/../.."
sha256sum -c 04_wiring/step0073_gate/frozen-inputs.sha256
sha256sum -c 04_wiring/step0073_gate/artifacts.sha256

verify_baseline() {
	expected="$1"
	path="$2"
	actual="$(git show "03fd5f1:$path" | sha256sum | awk '{print $1}')"
	if [ "$actual" != "$expected" ]; then
		echo "baseline divergente: $path" >&2
		exit 1
	fi
}

verify_baseline a3387e44d328ef0a1e1239226fa16c3ecac2e136a293bab6af4ffe47dd230e8b 01_core/domain/evidence.go
verify_baseline 2659cf5db5fba30e5003064aa3371aca11ee3673d7665db4b035f6d1d9c56050 01_core/usecase/evidence.go

GOCACHE=/tmp/bateia-step0073-gate-cache go test ./04_wiring/step0073_gate -run '^TestStep0073Gate' -count=1
