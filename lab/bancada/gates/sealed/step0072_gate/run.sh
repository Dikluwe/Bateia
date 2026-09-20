#!/bin/sh
set -eu

cd "$(dirname "$0")/../.."
sha256sum -c 04_wiring/step0072_gate/frozen-inputs.sha256
sha256sum -c 04_wiring/step0072_gate/artifacts.sha256

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
verify_baseline 8292e9a0b1425aabdad38fbf6a6c4b6964d146116e4b8ba5c1f3b982aee50ad4 01_core/domain/interfaces.go
verify_baseline 2659cf5db5fba30e5003064aa3371aca11ee3673d7665db4b035f6d1d9c56050 01_core/usecase/evidence.go
verify_baseline 29b8176ab30bc07ec9ea818e7a2a6608fa0409b8d0255fa5b12b6268581bb012 03_infra/local/client.go

GOCACHE=/tmp/bateia-step0072-gate-cache go test ./04_wiring/step0072_gate -run '^TestStep0072Gate' -count=1
