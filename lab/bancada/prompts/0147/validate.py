"""Check this patch's scope and documentary links without rewriting history."""
from pathlib import Path
import hashlib
import json
import re
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[4]
OUT = Path(__file__).resolve().parent
baseline = json.loads((OUT / 'baseline.json').read_text())
go_names = [
    '01_core/usecase/insight_exploration_memory.go',
    '01_core/usecase/insight_exploration_memory_test.go',
    '01_core/usecase/insight_site_plan.go',
    '01_core/usecase/insight_site_plan_test.go',
    '01_core/usecase/insight_neighborhood_plan.go',
    '01_core/usecase/insight_neighborhood_plan_test.go',
    '01_core/usecase/insight_followup_selection_test.go',
    '01_core/usecase/insight_neighborhood_sources_test.go',
    '02_shell/cli/insight_memory_test.go',
]
prompt_names = ['00_nucleo/prompts/' + name.removesuffix('.go') + '.md' for name in go_names]
prompt_names += ['00_nucleo/prompts/02_shell/cli/insight_memory.md']
other_docs = [
    '00_nucleo/contratos/README.md',
    '00_nucleo/contratos/exploracao-de-vizinhancas-v1.md',
    '00_nucleo/contratos/exploracao-de-vizinhancas-aceitacao-v1.md',
    '00_nucleo/arquitetura/README.md',
    '00_nucleo/passos/README.md',
]
new_step = '00_nucleo/passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md'
allowed = set(go_names + prompt_names + other_docs)
errors, changed = [], []
def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()
for rel, old in baseline.items():
    file = ROOT / rel
    if not file.is_file() or sha(file) != old:
        changed.append(rel)
        if rel not in allowed:
            errors.append('unexpected change: ' + rel)
for name in go_names:
    prompt = '00_nucleo/prompts/' + name.removesuffix('.go') + '.md'
    raw = (ROOT / name).read_text()
    if name in changed and (prompt not in changed or '@prompt ' + prompt not in raw):
        errors.append('changed Go lacks updated owner prompt: ' + name)
for file in (OUT / 'before').glob('*.txt'):
    original = file.name.removesuffix('.txt').replace('__', '/')
    if original not in baseline or sha(file) != baseline[original]:
        errors.append('before snapshot differs from baseline: ' + original)
current_go = {str(p.relative_to(ROOT)) for scope in ['01_core','02_shell','03_infra','04_wiring'] for p in (ROOT / scope).rglob('*.go')}
old_go = {rel for rel in baseline if rel.endswith('.go') and rel.split('/')[0] in {'01_core','02_shell','03_infra','04_wiring'}}
if current_go != old_go:
    errors.append('unexpected Go inventory change')
def slug(text):
    return re.sub(r'[^\w\- ]', '', text.lower().replace('`', '')).replace(' ', '-')
for rel in prompt_names + other_docs + [new_step]:
    file = ROOT / rel
    for raw in re.findall(r'(?<!!)\[[^\]]*\]\(([^)]+)\)', file.read_text()):
        url = urlsplit(raw.strip('<>'))
        if url.scheme or url.netloc:
            continue
        target = (file.parent / unquote(url.path)).resolve() if url.path else file
        if not target.exists():
            errors.append(rel + ': missing target ' + raw)
        elif url.fragment and target.suffix == '.md':
            anchors = {slug(m.group(1)) for m in re.finditer(r'^#{1,6} (.+)$', target.read_text(), re.M)}
            if unquote(url.fragment) not in anchors:
                errors.append(rel + ': missing anchor ' + raw)
red = (OUT / 'generation-before.log').read_text()
if 'feedback returned an unusable next generation' not in red or 'exhausted generation reached discovery effects' not in red:
    errors.append('missing regression-before-fix evidence')
result = dict(changed_existing_files=changed, changed_go_count=len(set(changed) & set(go_names)), preserved_go_count=len(old_go - set(go_names)), prompt_count=len(prompt_names), historical_files_and_rules_unchanged=not any(p not in allowed for p in changed), errors=errors, semantic_completeness_proven=False)
(OUT / 'scope-and-docs.json').write_text(json.dumps(result, indent=2) + '\n')
print(json.dumps(result, indent=2))
raise SystemExit(bool(errors))
