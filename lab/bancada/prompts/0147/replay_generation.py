"""Replay the preserved synthetic counterexample with the fixed executable."""
from pathlib import Path
import hashlib
import json
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[4]
HERE = Path(__file__).resolve().parent
BASE = HERE / 'generation-replay'
BASE.mkdir(exist_ok=True)
OUT = Path(tempfile.mkdtemp(prefix='run-', dir=BASE))
BINARY = HERE / 'final/bateia'
OLD_INPUT = ROOT / 'lab/bancada/prompts/0146/generation-probe/input.json'
original = OLD_INPUT.read_bytes()
memory = json.loads(original)
target = memory['entries'][0]['identity']['canonical_url']
def feedback(source, output):
    completed = subprocess.run([str(BINARY), 'idea', 'insights', 'memory', '--input', str(source), '--output', str(output), '--target', target, '--state', 'pending', '--reason', 'Synthetic generation replay'], capture_output=True, text=True)
    return dict(exit_code=completed.returncode, stderr=completed.stderr, stdout=completed.stdout, output_exists=output.exists())
exhausted = feedback(OLD_INPUT, OUT / 'invalid.json')
assert exhausted['exit_code'] != 0 and not exhausted['output_exists'] and not exhausted['stdout']
assert 'alteração de memória rejeitada' in exhausted['stderr']
memory['generation'], memory['id'] = 999998, ''
canonical = json.dumps(memory, ensure_ascii=False, separators=(',', ':')).encode()
memory['id'] = 'exploration_memory_' + hashlib.sha256(canonical).hexdigest()
control = OUT / 'last-valid-input.json'
control.write_text(json.dumps(memory, ensure_ascii=False, separators=(',', ':')) + '\n')
last = feedback(control, OUT / 'last-valid-output.json')
assert last['exit_code'] == 0
produced = json.loads((OUT / 'last-valid-output.json').read_bytes())
assert produced['generation'] == 999999
closed = feedback(OUT / 'last-valid-output.json', OUT / 'must-not-exist.json')
assert closed['exit_code'] != 0 and not closed['output_exists'] and not closed['stdout']
# The diagnostic distinguishes a valid snapshot whose advance is refused from an invalid memory.
assert 'alteração de memória rejeitada' in closed['stderr']
assert OLD_INPUT.read_bytes() == original
result = dict(synthetic=True, binary_sha256=hashlib.sha256(BINARY.read_bytes()).hexdigest(), preserved_input=str(OLD_INPUT.relative_to(ROOT)), preserved_input_sha256=hashlib.sha256(original).hexdigest(), exhausted=exhausted, last_valid=last, output_generation=produced['generation'], next_advance_rejected=closed)
(OUT / 'result.json').write_text(json.dumps(result, indent=2) + '\n')
print(json.dumps(dict(evidence=str(OUT.relative_to(ROOT)), exhausted_exit=exhausted['exit_code'], exhausted_published=exhausted['output_exists'], last_valid_exit=last['exit_code'], last_valid_generation=produced['generation'], next_advance_exit=closed['exit_code']), indent=2))
