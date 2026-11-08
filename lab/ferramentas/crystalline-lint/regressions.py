"""Exercise real Go parser behavior, including rejecting productive regressions."""
import pathlib
import subprocess
import sys
import tempfile

binary = str(pathlib.Path(sys.argv[1]).resolve())
cases = [
    ('composition-function', '04_wiring/main.go', 'package main;func Build()any{return nil}', None, False),
    ('composition-test-double', '04_wiring/main_test.go', 'package main;type double struct{};func(double) Decide(){};func TestComposition(){}', None, False),
    ('productive-method', '04_wiring/main.go', 'package main;type value struct{};func(value) Decide(){}', 'V12', False),
    ('private-method', '04_wiring/main.go', 'package main;type value struct{};func(value) decide(){}', 'V12', False),
    ('productive-interface', '04_wiring/main.go', 'package main;type Policy interface{Decide()}', 'V12', False),
    ('private-interface', '04_wiring/main.go', 'package main;type policy interface{decide()}', 'V12', False),
    ('productive-alias', '04_wiring/main.go', 'package main;type Policy = int', 'V12', False),
    ('productive-gravity', '02_shell/cli/app.go', 'package cli;import _ "bateia/03_infra/host"', 'V3', False),
    ('core-gravity', '01_core/domain/types.go', 'package domain;import _ "bateia/02_shell/cli"', 'V3', False),
    ('external-test-composition', '02_shell/cli/host_test.go', 'package cli_test;import _ "bateia/03_infra/host"', None, False),
    ('strict-test-gravity', '02_shell/cli/host_test.go', 'package cli_test;import _ "bateia/03_infra/host"', 'V3', True),
]
for name, relative, source, expected, strict in cases:
    with tempfile.TemporaryDirectory(prefix='bateia-0129-lint-') as temp:
        root = pathlib.Path(temp)
        config = ('check_test_imports = true\n' if strict else '') + '''[project]
root = "."
[languages]
go = { grammar = "tree-sitter-go", enabled = true }
[layers]
L0 = "00_nucleo"
L1 = "01_core"
L2 = "02_shell"
L3 = "03_infra"
L4 = "04_wiring"
lab = "lab"
'''
        (root / 'crystalline.toml').write_text(config)
        layer = 'L' + str(int(relative[:2]))
        def write_file(path, body):
            prompt = '00_nucleo/prompts/' + str(path.with_suffix('.md'))
            (root / prompt).parent.mkdir(parents=True, exist_ok=True)
            (root / prompt).write_text('Regression fixture.\n')
            (root / path).parent.mkdir(parents=True, exist_ok=True)
            (root / path).write_text('/**\n * Crystalline Lineage\n * @prompt ' + prompt + '\n * @layer ' + layer + '\n * @updated 2026-10-28\n */\n' + body)
        path = pathlib.Path(relative)
        write_file(path, source)
        if not path.name.endswith('_test.go'):
            write_file(path.with_name(path.stem + '_test.go'), 'package fixture\n')
        result = subprocess.run([binary, '--fail-on', 'warning', str(root)], cwd=root, capture_output=True, text=True)
        output = result.stdout + result.stderr
        if expected is None:
            assert result.returncode == 0, (name, output)
        else:
            assert result.returncode != 0 and '[' + expected + ']' in output, (name, output)
        print(name + ': PASS')
print(str(len(cases)) + ' parser regression controls passed')

# Exercise the new per-rule switch against a real Go gravity violation.
# Enabling/disabling is inspected in SARIF, not inferred from a green exit code.
import json
for enabled in [False, True]:
    with tempfile.TemporaryDirectory(prefix='bateia-rule-switch-') as temp:
        root = pathlib.Path(temp)
        (root / 'crystalline.toml').write_text('''[project]
root = "."
[languages]
go = { grammar = "tree-sitter-go", enabled = true }
[layers]
L0 = "00_nucleo"
L1 = "01_core"
L2 = "02_shell"
L3 = "03_infra"
L4 = "04_wiring"
lab = "lab"
[rules]
V3 = { enabled = ''' + str(enabled).lower() + ' }\n')
        for relative, body in [
            ('02_shell/cli/app.go', 'package cli;import _ "bateia/03_infra/host"'),
            ('02_shell/cli/app_test.go', 'package cli'),
        ]:
            path = pathlib.Path(relative)
            prompt = '00_nucleo/prompts/' + str(path.with_suffix('.md'))
            (root / prompt).parent.mkdir(parents=True, exist_ok=True)
            (root / prompt).write_text('Per-rule selection regression.\n')
            (root / path).parent.mkdir(parents=True, exist_ok=True)
            (root / path).write_text('/**\n * Crystalline Lineage\n * @prompt ' + prompt + '\n * @layer L2\n * @updated 2026-11-08\n */\n' + body)
        result = subprocess.run([binary, '--checks', 'all', '--format', 'sarif', '--fail-on', 'warning', str(root)], cwd=root, capture_output=True, text=True)
        report = json.loads(result.stdout)['runs'][0]
        scope = report['properties']['analysisScope']
        assert ('V3' in scope['selectedRules']) == enabled, scope
        assert ('V3' in scope['disabledRules']) != enabled, scope
        assert any(item['ruleId'] == 'V3' for item in report['results']) == enabled, report['results']
        assert (result.returncode != 0) == enabled, (result.returncode, result.stderr)
        print('rule-switch-V3-' + ('on' if enabled else 'off') + ': PASS')
print('2 per-rule switch controls passed')
