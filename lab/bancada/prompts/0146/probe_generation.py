"""Synthetic offline counterexample, never a real exploration history."""
from pathlib import Path
import hashlib
import json
import subprocess

ROOT = Path(__file__).resolve().parents[4]
OUT = Path(__file__).resolve().parent / 'generation-probe'
OUT.mkdir(exist_ok=True)
BINARY = ROOT / 'lab/bancada/arquitetura/0145/final/bateia'
def canonical(obj):
    return json.dumps(obj, ensure_ascii=False, separators=(',', ':')).encode()
def digest(obj):
    return hashlib.sha256(canonical(obj)).hexdigest()
connector = 'github-repository-catalog'
external = 'synthetic/generation-boundary'
identity = dict(connector_id=connector, external_id=external, canonical_url='https://github.com/' + external)
entry = dict(id='exploration_entry_' + digest(connector + '\x00' + external), identity=identity, state='PENDING', first_generation=1, last_question='queue', read_attempts=0)
memory = dict(schema_version='insight-exploration-memory-v1', id='', generation=999999, entries=[entry])
memory['id'] = 'exploration_memory_' + digest(memory)
source = OUT / 'input.json'
source.write_bytes(canonical(memory) + b'\n')
def feedback(src, dst):
    return subprocess.run([str(BINARY), 'idea', 'insights', 'memory', '--input', str(src), '--output', str(dst), '--target', identity['canonical_url'], '--state', 'pending', '--reason', 'Synthetic boundary probe'], capture_output=True, text=True)
first = feedback(source, OUT / 'output.json')
second = feedback(OUT / 'output.json', OUT / 'rejected.json')
result = dict(synthetic=True, binary=str(BINARY.relative_to(ROOT)), binary_sha256=hashlib.sha256(BINARY.read_bytes()).hexdigest(), input_generation=999999, first_exit=first.returncode, first_stderr=first.stderr, second_exit=second.returncode, second_stderr=second.stderr, rejected_file_exists=(OUT / 'rejected.json').exists())
if first.returncode == 0:
    result['output_generation'] = json.loads((OUT / 'output.json').read_text())['generation']
(OUT / 'result.json').write_text(json.dumps(result, indent=2) + '\n')
print(json.dumps(result, indent=2))
assert first.returncode == 0 and result['output_generation'] == 1000000
assert second.returncode != 0 and not result['rejected_file_exists']
