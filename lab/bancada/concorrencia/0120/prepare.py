"""Freeze references before observer responses; interpretations are authored."""
import hashlib
import json
from pathlib import Path

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
def digest(data):
    return hashlib.sha256(data).hexdigest()
def save(name, value):
    with (HERE / 'testdata' / name).open('x') as output:
        json.dump(value, output, indent=2, ensure_ascii=False)
        output.write('\n')
snapshots = {name: json.loads((HERE / 'testdata' / (name + '-snapshot.json')).read_text())
             for name in ('tunny', 'pond')}
def ref(project, path, first, last, role):
    item = next(f for f in snapshots[project]['files'] if f['path'] == path)
    assert digest(item['content'].encode()) == item['sha256']
    lines = item['content'].splitlines()
    assert 1 <= first <= last <= len(lines)
    return {'path': path, 'first_line': first, 'last_line': last,
            'source_sha256': item['sha256'],
            'segment_sha256': digest('\n'.join(lines[first-1:last]).encode()),
            'url': item['url'], 'role': role}
questions = [
    {'id': 'tunny_dynamic', 'project': 'tunny', 'axis': 'USAGE_CONTRACT',
     'question': 'SetSize pode ser chamado enquanto outras goroutines processam?',
     'reference_label': 'DECLARED_SUPPORTED_WITH_STOPPED_POOL_EXCEPTION',
     'evidence': [ref('tunny', 'README.md', 92, 101, 'UPSTREAM_DECLARATION'),
                  ref('tunny', 'tunny.go', 263, 268, 'UPSTREAM_DECLARATION')],
     'limit': 'Contrato declarado, sem prova de execução ou garantia universal.'},
    {'id': 'tunny_initialization', 'project': 'tunny', 'axis': 'STATIC_RELATION',
     'question': 'New chama SetSize para a configuração inicial?',
     'reference_label': 'DIRECT_CALL_IN_SUPPLIED_SOURCE',
     'call_path': 'tunny.go', 'call_line': 125, 'target_line': 266,
     'evidence': [ref('tunny', 'tunny.go', 116, 128, 'STATIC_SOURCE'),
                  ref('tunny', 'tunny.go', 263, 268, 'STATIC_SOURCE')],
     'limit': 'Chamada estrutural não certifica inicialização concluída.'},
    {'id': 'tunny_atomicity', 'project': 'tunny', 'axis': 'UNIVERSAL_ATOMICITY',
     'question': 'O recorte comprova atomicidade de todas as interações de configuração?',
     'reference_label': 'NOT_PROVEN_BY_BOUNDED_SOURCE_REVIEW',
     'evidence': [ref('tunny', 'tunny.go', 263, 293, 'STATIC_SOURCE')],
     'limit': 'Mutex e contrato publicado não substituem prova de todas as interações.'},
    {'id': 'pond_dynamic', 'project': 'pond', 'axis': 'USAGE_CONTRACT',
     'question': 'Resize é apresentado como operação de ajuste durante o uso?',
     'reference_label': 'DECLARED_DYNAMIC_RECONFIGURATION',
     'evidence': [ref('pond', 'README.md', 488, 520, 'UPSTREAM_DECLARATION')],
     'limit': 'Exemplo/declaração não demonstram todas as intercalacões.'},
    {'id': 'pond_zero_declarations', 'project': 'pond', 'axis': 'USAGE_CONTRACT',
     'question': 'As declarações de contrato concordam sobre Resize(0)?',
     'reference_label': 'CONFLICTING_DECLARATIONS_ZERO_BOUNDARY',
     'evidence': [ref('pond', 'README.md', 517, 520, 'UPSTREAM_DECLARATION'),
                  ref('pond', 'pool.go', 89, 92, 'UPSTREAM_DECLARATION'),
                  ref('pond', 'pool.go', 174, 181, 'STATIC_SOURCE')],
     'limit': 'README admite zero, interface diz maior que zero; implementação trata zero. Não é defeito de execução demonstrado.'},
    {'id': 'pond_atomicity', 'project': 'pond', 'axis': 'UNIVERSAL_ATOMICITY',
     'question': 'O recorte comprova atomicidade de todas as interações de configuração?',
     'reference_label': 'NOT_PROVEN_BY_BOUNDED_SOURCE_REVIEW',
     'evidence': [ref('pond', 'pool.go', 174, 205, 'STATIC_SOURCE')],
     'limit': 'Lock observado e execução posterior de workers não fornecem prova universal.'}
]
reference = {'schema_version': 'external-configuration-reference-v1',
             'phase': 'REFERENCE_FROZEN_BEFORE_OBSERVER_RESPONSES',
             'reference_origin': 'UPSTREAM_DECLARATIONS_WITH_AUTHOR_INTERPRETATION',
             'independent_semantic_review': False,
             'candidate_code_executed': False, 'projects': 2,
             'representative_sample': False, 'questions': questions,
             'snapshots_sha256': {name: digest((HERE / 'testdata' / (name + '-snapshot.json')).read_bytes())
                                 for name in snapshots}}
save('reference.json', reference)
standard = json.loads((ROOT / 'lab/bancada/concorrencia/0111/testdata/environment.json').read_text())
cases = []
for name, snapshot in snapshots.items():
    files = [{'path': item['path'], 'content': item['content'], 'sha256': item['sha256']}
             for item in snapshot['files'] if item['path'].endswith('.go')
             and not item['path'].endswith('_test.go')
             and not item['path'].startswith('internal/assert/')]
    assert len(files) <= 32
    assert all('//go:build' not in f['content'] and '// +build' not in f['content'] for f in files)
    module = next(f['content'].splitlines()[0].split()[1] for f in snapshot['files'] if f['path'] == 'go.mod')
    directories = sorted({str(Path(f['path']).parent) for f in files})
    environment = dict(standard, packages=[{'import_path': module + ('/' + d if d != '.' else ''),
                                         'directory': d} for d in directories])
    save(name + '-environment.json', environment)
    for profile in ('complete', 'no_exports', 'source_fragment'):
        selected = files if profile != 'source_fragment' else [next(f for f in files if f['path'] == ('tunny.go' if name == 'tunny' else 'pool.go'))]
        case_id = name + '-' + profile
        request = {'schema_version': '1', 'origin': snapshot['repository'] + '@' + snapshot['commit'] + '; ' + profile,
                   'build': {'go_version': standard['toolchain'], 'goos': 'linux', 'goarch': 'amd64'},
                   'files': selected}
        save(case_id + '-input.json', request)
        cases.append({'id': case_id, 'project': name, 'profile': profile,
                      'input': case_id + '-input.json',
                      'environment': name + '-environment.json' if profile != 'no_exports' else None,
                      'question_ids': [q['id'] for q in questions if q['project'] == name]})
save('cases.json', {'cases': cases, 'independent_project_clusters': 2,
                    'variants_are_not_independent_samples': True})
save('protocol.json', {'schema_version': 'configuration-evaluation-protocol-v1',
                      'adapter_is_classifier': False,
                      'no_guard_means_no_supported_fact_not_no_protection': True,
                      'contract_verdict_emitted_by_v8': False,
                      'static_relation_rule': 'STATIC_CALL_TARGET site.path/line and target.path/line exactly match frozen reference.',
                      'contract_rule': 'NOT_EMITTED: no automated contract verdict exists; count unanswered, not semantic error or correct abstention.',
                      'atomicity_rule': 'No guarantee may be promoted from configuration_guards; UNRESOLVED stays unresolved.',
                      'metrics': ['static_reference_localization', 'contract_answer_coverage',
                                  'unsupported_atomicity_promotions', 'incomplete_input_gaps'],
                      'semantic_precision_recall': 'NOT_MEASURED_NO_INDEPENDENT_SEMANTIC_ORACLE'})
print(json.dumps({'projects': 2, 'questions': len(questions), 'input_profiles': len(cases)}))
