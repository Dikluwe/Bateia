"""Replay known frozen sources; this is development verification, not a holdout."""
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BINARY = '/tmp/bateia-declarations-0121'
DATA = HERE / 'testdata'
RESULTS = HERE / 'results'
def load(path):
    return json.loads(path.read_text())
def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()
freeze = load(HERE.parent / '0120/testdata/observer-freeze.json')
assert {str(p.relative_to(ROOT)): digest(p) for p in sorted((ROOT / 'lab/bancada/concorrencia/0109').rglob('*.go'))} == freeze['source_hashes']
def run(request):
    with tempfile.TemporaryDirectory() as directory:
        path = Path(directory) / 'input.json'
        path.write_text(json.dumps(request))
        args = [BINARY, '--input', str(path)]
        first = subprocess.check_output(args, cwd=ROOT)
        assert first == subprocess.check_output(args, cwd=ROOT)
        return load_result(first), first
def load_result(raw):
    return json.loads(raw)
def request_for(project, profile='complete'):
    return {'source': load(HERE.parent / ('0120/testdata/' + project + '-' + profile + '-input.json')),
            'environment': load(HERE.parent / ('0120/testdata/' + project + '-environment.json'))}
results = []
for project, target, spans in [('tunny', '*Pool.SetSize', [(92, 101)]),
                               ('pond', '*pool.Resize', [(488, 520)])]:
    request = request_for(project)
    initial, _ = run(request)
    owner = next(o for o in initial['api_owners'] if o['name'] == target)
    snapshot = load(HERE.parent / ('0120/testdata/' + project + '-snapshot.json'))
    document = next(f for f in snapshot['files'] if f['path'] == 'README.md')
    request['documents'] = [{k: document[k] for k in ('path', 'content', 'sha256')}]
    request['bindings'] = [{'document_path': document['path'], 'first_line': first,
                            'last_line': last, 'owner': owner['location']} for first, last in spans]
    if project == 'pond':
        interface = next(o for o in initial['api_owners'] if o['name'] == 'BasePool.Resize')
        request['bindings'].append({'document_path': document['path'], 'first_line': 517,
                                    'last_line': 520, 'owner': interface['location']})
    report, raw = run(request)
    assert report['observer_report'] == load(HERE.parent / ('0120/results/' + project + '-complete.json'))
    linked = [d for d in report['declarations'] if d['association'] == 'CALLER_PROVIDED_BINDING']
    assert len(linked) == len(request['bindings'])
    assert any(d['owner']['name'] == target and d['evidence'] for d in linked)
    if project == 'pond':
        api_doc = next(d for d in report['declarations'] if d['owner']['name'] == 'BasePool.Resize' and d['association'] == 'AST_ATTACHED_GO_DOC')
        assert 'greater than 0' in api_doc['text'] and not api_doc['evidence']
        external = next(d for d in linked if d['owner']['name'] == 'BasePool.Resize')
        assert 'greater than or equal to 0' in external['text'] and not external['evidence']
    (DATA / (project + '-request.json')).write_text(json.dumps(request, indent=2) + '\n')
    (RESULTS / (project + '.json')).write_bytes(raw)
    results.append({'project': project, 'owners': len(report['api_owners']),
                    'ast_documents': sum(d['association'] == 'AST_ATTACHED_GO_DOC' for d in report['declarations']),
                    'provided_bindings': len(linked), 'evidence_links': sum(len(d['evidence']) for d in report['declarations']),
                    'observer_projection_preserved': True, 'result_sha256': hashlib.sha256(raw).hexdigest()})
    # Type failure does not remove syntactically attached API documentation.
    fragment, fragment_raw = run(request_for(project, 'source_fragment'))
    assert fragment['observer_report']['type_diagnostics']
    assert any(d['owner']['name'] == (target if project == 'tunny' else 'BasePool.Resize') for d in fragment['declarations'])
    (RESULTS / (project + '-fragment.json')).write_bytes(fragment_raw)

conc_request = {'source': load(HERE.parent / '0111/testdata/conc-input.json'),
                'environment': load(HERE.parent / '0111/testdata/environment.json')}
conc, raw = run(conc_request)
assert conc['observer_report'] == load(HERE.parent / '0116/results/conc-guard.json')
config = next(d for d in conc['declarations'] if d['owner']['name'] == '*Pool.WithMaxGoroutines')
assert {'CONFIGURATION_WRITE', 'CONFIGURATION_GUARD_CALL'} <= {e['kind'] for e in config['evidence']}
assert conc['observer_report']['configuration_guards'][0]['atomicity'] == 'UNRESOLVED'
(DATA / 'conc-request.json').write_text(json.dumps(conc_request, indent=2) + '\n')
(RESULTS / 'conc.json').write_bytes(raw)
results.append({'project': 'sourcegraph/conc', 'owners': len(conc['api_owners']),
                'ast_documents': len(conc['declarations']), 'provided_bindings': 0,
                'evidence_links': sum(len(d['evidence']) for d in conc['declarations']),
                'observer_projection_preserved': True, 'result_sha256': hashlib.sha256(raw).hexdigest()})
receipt = {'schema_version': 'api-declarations-replay-v1',
           'reference_origin': 'KNOWN_DEVELOPMENT_SOURCES_NOT_HOLDOUT',
           'candidate_code_executed': False, 'semantic_accuracy_measured': False,
           'byte_identical_replays': True, 'observer_v8_unchanged': True,
           'external_associations_automatically_discovered': False,
           'semantic_verdicts_implemented': False,
           'binary_sha256': digest(Path(BINARY)), 'results': results}
(RESULTS / 'replay-receipt.json').write_text(json.dumps(receipt, indent=2) + '\n')
print(json.dumps(results))
