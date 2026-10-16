"""Execute unchanged upstream tests from a hash-checked frozen bundle."""
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import subprocess
import tempfile

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BUNDLE = ROOT / 'lab/bancada/insights/0106/github-bundles.json'
COMMIT = '5f936abd7ae87036af1f75c95fb9d0daaf00116b'
RESULTS = HERE / 'results'
RESULTS.mkdir(exist_ok=True)
bundle = json.loads(BUNDLE.read_text())[0]
work = Path(tempfile.mkdtemp(prefix='bateia-conc-0119-', dir='/tmp'))
sources = {}
for item in bundle['items']:
    path = PurePosixPath(item['path'])
    if item['path'] == '@repository-metadata':
        continue
    assert not path.is_absolute() and '..' not in path.parts
    assert item['ref']['commit'] == COMMIT
    data = item['content'].encode()
    digest = hashlib.sha256(data).hexdigest()
    assert digest == item['content_hash']
    target = work / path
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_bytes(data)
    sources[item['path']] = digest

env = dict(os.environ, GOWORK='off', GOTOOLCHAIN='local')
command = ['go', 'test', '-mod=mod', '-race', '-count=3', '-timeout=90s',
           '-json', '-run', '^(TestPool|ExamplePool)$', './pool']
completed = subprocess.run(command, cwd=work, env=env, capture_output=True,
                           timeout=180)
(RESULTS / 'upstream-tests.jsonl').write_bytes(completed.stdout)
(RESULTS / 'upstream-tests.stderr').write_bytes(completed.stderr)
events = [json.loads(line) for line in completed.stdout.splitlines()]
passed = {}
for event in events:
    if event['Action'] == 'pass' and 'Test' in event:
        name = event['Test']
        passed[name] = passed.get(name, 0) + 1
for path, digest in sources.items():
    assert hashlib.sha256((work / path).read_bytes()).hexdigest() == digest, path
version = subprocess.check_output(['go', 'version'], cwd=work, env=env).decode().strip()
modules = subprocess.check_output(['go', 'list', '-mod=mod', '-m', '-json', 'all'],
                                  cwd=work, env=env)
(RESULTS / 'modules.jsonstream').write_bytes(modules)
go_sum = (work / 'go.sum').read_bytes() if (work / 'go.sum').exists() else b''
(RESULTS / 'go.sum').write_bytes(go_sum)
receipt = {
    'schema_version': 'upstream-reference-execution-v1',
    'reference_origin': 'UPSTREAM_TESTS_UNMODIFIED',
    'candidate_code_executed': True,
    'independent_semantic_accuracy_measured': False,
    'complete_repository_checkout_certified': False,
    'repository_commit': COMMIT,
    'bundle_sha256': hashlib.sha256(BUNDLE.read_bytes()).hexdigest(),
    'temporary_module': str(work), 'toolchain': version,
    'command': command, 'exit_code': completed.returncode,
    'source_hashes_before_and_after': sources,
    'pass_counts': passed,
    'modules_sha256': hashlib.sha256(modules).hexdigest(),
    'go_sum_sha256': hashlib.sha256(go_sum).hexdigest(),
    'stdout_sha256': hashlib.sha256(completed.stdout).hexdigest(),
    'stderr_sha256': hashlib.sha256(completed.stderr).hexdigest(),
    'limits': ['Selected upstream tests only; no concurrent configuration test.',
               'Passing race detector is restricted to these executions.',
               'No independent semantic classification corpus.']
}
(RESULTS / 'execution.json').write_text(json.dumps(receipt, indent=2) + '\n')
assert completed.returncode == 0, completed.stderr.decode()
assert passed.get('ExamplePool') == 1  # Go repeats tests, but runs examples once.
for name in ('TestPool',
             'TestPool/panics_on_configuration_after_init/before_wait',
             'TestPool/panics_on_configuration_after_init/after_wait',
             'TestPool/is_reusable'):
    assert passed.get(name) == 3, (name, passed.get(name))
print(json.dumps({'exit_code': completed.returncode,
                  'passed_test_names': len(passed), 'test_repetitions': 3,
                  'example_repetitions': 1,
                  'candidate_code_executed': True}))
