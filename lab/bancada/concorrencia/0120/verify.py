"""Check evaluation accounting; negative controls do not measure semantic accuracy."""
from copy import deepcopy
import hashlib
import json
from pathlib import Path

HERE = Path(__file__).resolve().parent
DATA = HERE / 'testdata'
RESULTS = HERE / 'results'
def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()
cases = json.loads((DATA / 'cases.json').read_text())['cases']
reference = json.loads((DATA / 'reference.json').read_text())
questions = {q['id']: q for q in reference['questions']}
def validate(evaluation):
    assert evaluation['projects'] == len({c['project'] for c in cases}) == 2
    assert evaluation['questions'] == len(questions) == 6
    assert evaluation['profiles'] == len(cases) == 6
    assert evaluation['profile_question_rows'] == len(evaluation['question_results']) == 18
    assert evaluation['independent_semantic_accuracy_measured'] is False
    assert evaluation['variants_are_not_independent_samples'] is True
    assert evaluation['reference_sha256'] == digest(DATA / 'reference.json')
    assert evaluation['protocol_sha256'] == digest(DATA / 'protocol.json')
    assert evaluation['observer_freeze_sha256'] == digest(DATA / 'observer-freeze.json')
    profiles = {p['id']: p for p in evaluation['profiles_results']}
    rows = {(r['case'], r['question']): r for r in evaluation['question_results']}
    assert len(rows) == 18 and set(profiles) == {c['id'] for c in cases}
    located = 0
    incomplete = 0
    for case in cases:
        profile = profiles[case['id']]
        report_path = RESULTS / (case['id'] + '.json')
        assert profile['result_sha256'] == digest(report_path)
        assert profile['input_sha256'] == digest(DATA / case['input'])
        if case['environment']:
            assert profile['environment_sha256'] == digest(DATA / case['environment'])
        report = json.loads(report_path.read_text())
        assert report['status'] == 'OBSERVED_SYNTAX'  # Syntax availability is not a contract verdict.
        assert not any(g['atomicity'] != 'UNRESOLVED' for g in report.get('configuration_guards', []))
        diagnostics = len(report.get('type_diagnostics', []))
        assert diagnostics == profile['type_diagnostics']
        if case['profile'] == 'complete':
            assert diagnostics == 0
        else:
            assert diagnostics > 0
            assert any(g['code'] == 'SUPPLIED_PACKAGE_TYPECHECK_FAILED' for g in report['gaps'])
            incomplete += 1
        for question_id in case['question_ids']:
            question = questions[question_id]
            row = rows[(case['id'], question_id)]
            assert row['axis'] == question['axis'] and row['reference_label'] == question['reference_label']
            if question['axis'] == 'STATIC_RELATION':
                matches = [r for r in report.get('relations', []) if r['kind'] == 'STATIC_CALL_TARGET'
                           and r['from']['path'] == question['call_path'] and r['from']['line'] == question['call_line']
                           and r['to']['path'] == question['call_path'] and r['to']['line'] == question['target_line']]
                assert row['evidence'] == matches
                assert row['outcome'] == ('REFERENCE_LOCATED' if matches else 'NO_SUPPORTED_RELATION')
                if case['profile'] == 'complete':
                    located += bool(matches)
            elif question['axis'] == 'USAGE_CONTRACT':
                assert row['outcome'] == 'CONTRACT_VERDICT_NOT_IMPLEMENTED'
                assert row['semantic_correctness'] == 'NOT_SCORED'
            else:
                assert row['outcome'] == 'NO_ATOMICITY_PROOF_EMITTED'
                assert row['semantic_correctness'] == 'COMPATIBLE_ABSTENTION_NOT_PROOF_OF_ACCURACY'
    assert evaluation['metrics_complete_profiles_only'] == {
        'static_references_located': located, 'static_reference_questions': 1,
        'contract_verdicts_emitted': 0, 'contract_questions': 3,
        'unproved_atomicity_questions_without_promotion': 2,
        'unsupported_atomicity_promotions': 0, 'semantic_precision_recall': None}
    assert incomplete == 4
    calibration = evaluation['calibration_excluded_from_external_metrics']
    assert calibration['result_sha256'] == digest(RESULTS / 'conc-calibration.json')
    actual = json.loads((RESULTS / 'conc-calibration.json').read_text())
    assert len(actual['configuration_guards']) == calibration['guards'] == 1
    assert actual['configuration_guards'][0]['atomicity'] == calibration['atomicity'] == 'UNRESOLVED'

evaluation = json.loads((RESULTS / 'evaluation.json').read_text())
validate(evaluation)
controls = []
for name, mutate in [
    ('count_profiles_as_independent_projects', lambda r: r.update(projects=6)),
    ('promote_unanswered_to_contract_coverage', lambda r: r['metrics_complete_profiles_only'].update(contract_verdicts_emitted=3)),
    ('claim_independent_semantic_accuracy', lambda r: r.update(independent_semantic_accuracy_measured=True)),
    ('count_uncertainty_as_semantic_success', lambda r: next(x for x in r['question_results'] if x['axis'] == 'UNIVERSAL_ATOMICITY').update(semantic_correctness='CORRECT')),
    ('fabricate_atomicity_promotion_count', lambda r: r['metrics_complete_profiles_only'].update(unsupported_atomicity_promotions=1))]:
    changed = deepcopy(evaluation)
    mutate(changed)
    try:
        validate(changed)
    except AssertionError:
        controls.append({'control': name, 'rejected': True})
    else:
        raise AssertionError(name)
(RESULTS / 'integrity-controls.json').write_text(json.dumps({
    'negative_controls': controls,
    'semantic_accuracy_measured': False,
    'evaluation_sha256': digest(RESULTS / 'evaluation.json')}, indent=2) + '\n')
print(json.dumps({'evaluation_integrity': 'PASS', 'negative_controls': len(controls)}))
