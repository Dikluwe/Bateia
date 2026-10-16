"""Run a frozen observer against references frozen before any response."""
from collections import Counter
import hashlib
import json
from pathlib import Path
import subprocess

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
DATA = HERE / 'testdata'
RESULTS = HERE / 'results'
def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()
freeze = json.loads((DATA / 'observer-freeze.json').read_text())
assert digest(Path(freeze['binary_path'])) == freeze['binary_sha256']
actual_sources = {str(p.relative_to(ROOT)): digest(p) for p in sorted((ROOT / 'lab/bancada/concorrencia/0109').rglob('*.go'))}
assert actual_sources == freeze['source_hashes']
reference = json.loads((DATA / 'reference.json').read_text())
assert reference['independent_semantic_review'] is False
assert not json.loads((DATA / 'protocol.json').read_text())['adapter_is_classifier']
for project, snapshot_hash in reference['snapshots_sha256'].items():
    assert digest(DATA / (project + '-snapshot.json')) == snapshot_hash
for q in reference['questions']:
    snapshot = json.loads((DATA / (q['project'] + '-snapshot.json')).read_text())
    items = {f['path']: f for f in snapshot['files']}
    for ref in q['evidence']:
        item = items[ref['path']]
        assert hashlib.sha256(item['content'].encode()).hexdigest() == ref['source_sha256'] == item['sha256']
        segment = '\n'.join(item['content'].splitlines()[ref['first_line']-1:ref['last_line']])
        assert hashlib.sha256(segment.encode()).hexdigest() == ref['segment_sha256']
cases = json.loads((DATA / 'cases.json').read_text())['cases']
rows = []
profiles = []
for case in cases:
    args = [freeze['binary_path'], 'observe', '--input', str(DATA / case['input'])]
    if case['environment']:
        args += ['--environment', str(DATA / case['environment'])]
    first = subprocess.check_output(args, cwd=ROOT)
    second = subprocess.check_output(args, cwd=ROOT)
    assert first == second
    report = json.loads(first)
    assert report['analyzer'] == freeze['analyzer']
    request = json.loads((DATA / case['input']).read_text())
    files = {f['path']: f for f in request['files']}
    def verify_locations(value):
        if isinstance(value, dict):
            if 'content_sha256' in value:
                item = files[value['path']]
                assert value['content_sha256'] == item['sha256']
                raw = item['content'].encode()
                if 'offset' in value:
                    assert 0 <= value['offset'] <= value['end_offset'] <= len(raw)
                    assert value['line'] == raw[:value['offset']].count(b'\n') + 1
            for child in value.values():
                verify_locations(child)
        elif isinstance(value, list):
            for child in value:
                verify_locations(child)
    verify_locations(report)
    (RESULTS / (case['id'] + '.json')).write_bytes(first)
    guards = report.get('configuration_guards', [])
    promotions = sum(g['atomicity'] != 'UNRESOLVED' for g in guards)
    diagnostics = len(report.get('type_diagnostics', []))
    gaps = dict(Counter(g['code'] for g in report['gaps']))
    profiles.append({'id': case['id'], 'project': case['project'], 'profile': case['profile'],
                     'type_diagnostics': diagnostics, 'gap_counts': gaps,
                     'configuration_guards': len(guards),
                     'atomicity_promotions': promotions,
                     'relation_counts': dict(Counter(r['kind'] for r in report.get('relations', []))),
                     'result_sha256': hashlib.sha256(first).hexdigest(),
                     'input_sha256': digest(DATA / case['input']),
                     'environment_sha256': digest(DATA / case['environment']) if case['environment'] else None})
    for q in reference['questions']:
        if q['id'] not in case['question_ids']:
            continue
        row = {'case': case['id'], 'question': q['id'], 'axis': q['axis'],
               'reference_label': q['reference_label']}
        if q['axis'] == 'STATIC_RELATION':
            matches = [r for r in report.get('relations', []) if r['kind'] == 'STATIC_CALL_TARGET'
                       and r['from']['path'] == q['call_path'] and r['from']['line'] == q['call_line']
                       and r['to']['path'] == q['call_path'] and r['to']['line'] == q['target_line']]
            row.update(outcome='REFERENCE_LOCATED' if matches else 'NO_SUPPORTED_RELATION',
                       evidence=matches)
        elif q['axis'] == 'USAGE_CONTRACT':
            row.update(outcome='CONTRACT_VERDICT_NOT_IMPLEMENTED',
                       semantic_correctness='NOT_SCORED')
        else:
            row.update(outcome='NO_ATOMICITY_PROOF_EMITTED' if promotions == 0 else 'UNSUPPORTED_PROMOTION',
                       semantic_correctness='COMPATIBLE_ABSTENTION_NOT_PROOF_OF_ACCURACY')
        rows.append(row)
complete_ids = {c['id'] for c in cases if c['profile'] == 'complete'}
calibration_args = [freeze['binary_path'], 'observe', '--input',
                    str(ROOT / 'lab/bancada/concorrencia/0111/testdata/conc-input.json'),
                    '--environment', str(ROOT / 'lab/bancada/concorrencia/0111/testdata/environment.json')]
calibration_raw = subprocess.check_output(calibration_args, cwd=ROOT)
assert calibration_raw == subprocess.check_output(calibration_args, cwd=ROOT)
calibration = json.loads(calibration_raw)
assert len(calibration.get('configuration_guards', [])) == 1
assert calibration['configuration_guards'][0]['atomicity'] == 'UNRESOLVED'
(RESULTS / 'conc-calibration.json').write_bytes(calibration_raw)
primary = [r for r in rows if r['case'] in complete_ids]
structural = [r for r in primary if r['axis'] == 'STATIC_RELATION']
contract = [r for r in primary if r['axis'] == 'USAGE_CONTRACT']
evaluation = {'schema_version': 'external-configuration-evaluation-v1',
              'analyzer': freeze['analyzer'], 'reference_origin': reference['reference_origin'],
              'independent_semantic_accuracy_measured': False,
              'candidate_code_executed': False, 'representative_sample': False,
              'projects': 2, 'questions': 6, 'profiles': 6,
              'profile_question_rows': len(rows), 'variants_are_not_independent_samples': True,
              'byte_identical_replays': True,
              'reference_sha256': digest(DATA / 'reference.json'),
              'protocol_sha256': digest(DATA / 'protocol.json'),
              'observer_freeze_sha256': digest(DATA / 'observer-freeze.json'),
              'calibration_excluded_from_external_metrics': {
                  'project': 'sourcegraph/conc', 'guards': 1, 'atomicity': 'UNRESOLVED',
                  'result_sha256': hashlib.sha256(calibration_raw).hexdigest()},
              'metrics_complete_profiles_only': {
                  'static_references_located': sum(r['outcome'] == 'REFERENCE_LOCATED' for r in structural),
                  'static_reference_questions': len(structural),
                  'contract_verdicts_emitted': 0, 'contract_questions': len(contract),
                  'unproved_atomicity_questions_without_promotion': sum(r['outcome'] == 'NO_ATOMICITY_PROOF_EMITTED' for r in primary if r['axis'] == 'UNIVERSAL_ATOMICITY'),
                  'unsupported_atomicity_promotions': sum(p['atomicity_promotions'] for p in profiles if p['id'] in complete_ids),
                  'semantic_precision_recall': None},
              'profiles_results': profiles, 'question_results': rows}
(RESULTS / 'evaluation.json').write_text(json.dumps(evaluation, indent=2) + '\n')
print(json.dumps(evaluation['metrics_complete_profiles_only']))
