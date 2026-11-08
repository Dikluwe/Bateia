"""Structural documentation checks, with transparent semantic limits."""
from pathlib import Path
import hashlib
import json
import re
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[4]
OUT = Path(__file__).resolve().parent
baseline = json.loads((OUT / 'baseline.json').read_text())
prompts = [Path(p.name.removesuffix('.txt').replace('__', '/')) for p in sorted((OUT / 'before').glob('*.txt'))]
additions = [Path(p) for p in [
    '00_nucleo/contratos/README.md',
    '00_nucleo/contratos/exploracao-de-vizinhancas-v1.md',
    '00_nucleo/contratos/exploracao-de-vizinhancas-aceitacao-v1.md',
    '00_nucleo/documentacao/escrever-prompts.md',
    '00_nucleo/passos/0146-consolidar-prompts-de-exploracao.md',
]]
indexes = [Path(p) for p in ['00_nucleo/arquitetura/README.md', '00_nucleo/documentacao/README.md', '00_nucleo/passos/README.md']]
allowed = set(map(str, prompts + additions + indexes))
errors = []
def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()
def slug(s):
    return re.sub(r'[^\w\- ]', '', s.lower().replace('`', '')).replace(' ', '-')
def link_errors(path, text):
    failures = []
    for raw in re.findall(r'(?<!!)\[[^\]]*\]\(([^)]+)\)', text):
        raw = raw.strip('<>')
        url = urlsplit(raw)
        if url.scheme or url.netloc:
            continue
        target = (path.parent / unquote(url.path)).resolve() if url.path else path
        if not target.exists():
            failures.append('missing target ' + raw)
        elif url.fragment and target.suffix == '.md':
            anchors = {slug(m.group(1)) for m in re.finditer(r'^#{1,6} (.+)$', target.read_text(), re.M)}
            if unquote(url.fragment) not in anchors:
                failures.append('missing anchor ' + raw)
    return failures
matrix_path = ROOT / additions[2]
matrix = matrix_path.read_text()
cases = set(re.findall(r'^\| ([A-Z]+-\d+) \|', matrix, re.M))
assert len(cases) == 24, len(cases)
def prompt_errors(path, text):
    required = ['## Interface e contrato', '## Comportamento esperado', '## Limites e dependências', '## Origem e preservação'] if path.stem.endswith('_test') else ['## Interface e responsabilidade','## Comportamento vigente','## Limites e dependências','## Aceitação','## Origem e preservação']
    failures = ['missing section ' + h for h in required if h not in text]
    used = set(re.findall(r'\[([A-Z]+-\d+)\]', text))
    if not used:
        failures.append('no acceptance criteria')
    failures += ['unknown criterion ' + c for c in used - cases]
    if '/contratos/exploracao-de-vizinhancas-v1.md#' not in text:
        failures.append('no scoped contract dependency')
    return failures
consumers = {case: [] for case in cases}
for rel in prompts:
    if rel.stem.endswith('_test'):
        for case in set(re.findall(r'\[([A-Z]+-\d+)\]', (ROOT / rel).read_text())):
            if case in consumers:
                consumers[case].append(str(rel))
for case, paths in consumers.items():
    if not paths:
        errors.append('criterion has no consumer prompt: ' + case)
for rel in prompts + additions + indexes:
    path = ROOT / rel
    text = path.read_text()
    errors += [str(rel) + ': ' + e for e in link_errors(path, text)]
    if rel in prompts:
        errors += [str(rel) + ': ' + e for e in prompt_errors(path, text)]
for rel in prompts:
    snapshot = OUT / 'before' / (str(rel).replace('/', '__') + '.txt')
    if sha(snapshot) != baseline[str(rel)]:
        errors.append('snapshot differs from initial file: ' + str(rel))
changed = []
for rel, digest in baseline.items():
    path = ROOT / rel
    if not path.exists() or sha(path) != digest:
        changed.append(rel)
        if rel not in allowed:
            errors.append('unexpected existing-file change: ' + rel)
# Detect new Go as well as changed/missing existing Go within product scopes.
current_go = {str(p.relative_to(ROOT)) for scope in ['01_core','02_shell','03_infra','04_wiring'] for p in (ROOT / scope).rglob('*.go')}
baseline_go = {p for p in baseline if p.endswith('.go') and p.split('/')[0] in {'01_core','02_shell','03_infra','04_wiring'}}
if current_go != baseline_go:
    errors.append('productive/test Go inventory changed')
probe = ROOT / prompts[0]
negative_controls = {
    'broken_link': bool(link_errors(probe, '[broken](missing-0146.md)')),
    'broken_anchor': bool(link_errors(probe, '[broken](../../../contratos/exploracao-de-vizinhancas-v1.md#missing-0146)')),
    'missing_section': bool(prompt_errors(probe, probe.read_text().replace('## Interface e responsabilidade', '## Other'))),
    'unknown_criterion': bool(prompt_errors(probe, probe.read_text() + '\n[UNKNOWN-99](x.md)')),
}
if not all(negative_controls.values()):
    errors.append('negative structural control failed')
result = dict(prompt_count=len(prompts), acceptance_case_count=len(cases), checked_documents=len(prompts + additions + indexes), changed_existing_files=changed, unchanged_go_count=len(baseline_go), negative_controls=negative_controls, errors=errors, semantic_completeness_proven=False, known_implementation_gaps=['MEM-04'], acceptance_consumers=consumers)
(OUT / 'documentation-check.json').write_text(json.dumps(result, indent=2) + '\n')
print(json.dumps(result, indent=2))
raise SystemExit(bool(errors))
