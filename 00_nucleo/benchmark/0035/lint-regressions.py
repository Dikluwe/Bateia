#!/usr/bin/env python3
"""CLI regressions for the 0035 linter patch. Uses only temporary fixtures."""
import json
import pathlib
import subprocess
import sys
import tempfile

binary = pathlib.Path(sys.argv[1]).resolve()

def check(path, code, harness=False):
    with tempfile.TemporaryDirectory(prefix="bateia-0035-lint-case-") as tmp:
        root = pathlib.Path(tmp)
        for layer in ["01_core", "02_shell", "03_infra", "04_wiring"]:
            (root / layer).mkdir()
        config = '''[layers]
L1 = "01_core"
L2 = "02_shell"
L3 = "03_infra"
L4 = "04_wiring"
[wiring_exceptions]
allow_adapter_structs = false
'''
        if harness:
            config += '[wiring_exceptions.harness_files]\n"04_wiring/gate.go" = "verification fixture"\n'
        (root / "crystalline.toml").write_text(config)
        source = root / path
        source.parent.mkdir(parents=True, exist_ok=True)
        source.write_text(code)
        result = subprocess.run([str(binary), str(root), "--checks", "v1,v3,v12", "--format", "sarif"], capture_output=True, text=True)
        assert result.returncode in (0, 1), result.stderr
        return json.loads(result.stdout)["runs"][0]["results"]

def count(findings, rule):
    return sum(f["ruleId"] == rule for f in findings)

# All fixtures lack lineage deliberately. V12 exemption must not hide V1.
findings = check("04_wiring/main_test.go", "package main; type fake struct{}; func (fake) run() {}")
assert count(findings, "V12") == 0
assert count(findings, "V1") == 1
findings = check("04_wiring/gate.go", "package main; type port interface{ Run() }", harness=True)
assert count(findings, "V12") == 0
assert count(findings, "V1") == 1
findings = check("04_wiring/sub/gate.go", "package main; type port interface{ Run() }", harness=True)
assert count(findings, "V12") == 1
findings = check("04_wiring/production.go", "package main; type hidden interface{ Run() }; type alias = int; type service struct{}; func (service) run() {}")
assert count(findings, "V12") == 4
findings = check("04_wiring/main.go", "package main; func NewServer() {}; func helper() {}")
assert count(findings, "V12") == 0
# Go tests continue to report forbidden imports (L2 -> L3).
findings = check("02_shell/app_test.go", 'package shell; import "fixture/03_infra/db"; func TestX() { db.Run() }')
assert count(findings, "V3") == 1
assert count(findings, "V1") == 1
print("PASS: six CLI fixtures; V12 production types/methods, exact harness scope, V1 and V3 retained")
