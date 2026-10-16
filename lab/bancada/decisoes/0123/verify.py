"""Verify real diagnostic traces, without converting author agreement into model accuracy."""
from collections import Counter
import hashlib
from itertools import permutations
import json
from pathlib import Path

HERE = Path(__file__).resolve().parent
RESULTS = HERE / 'results'
pilot_path = HERE.parent / '0122/pilot-requests.json'
pilot = json.loads(pilot_path.read_text())
summary = json.loads((RESULTS / 'diagnostics.json').read_text())
rows = [json.loads(line) for line in (RESULTS / 'predictions.jsonl').read_text().splitlines()]
assert summary['predictions_sha256'] == hashlib.sha256((RESULTS / 'predictions.jsonl').read_bytes()).hexdigest()
assert summary['pilot_sha256'] == hashlib.sha256(pilot_path.read_bytes()).hexdigest()
assert summary['independent_semantic_accuracy_measured'] is False
assert summary['training_performed'] is False and summary['threshold_fitted'] is False
assert len(rows) == summary['runs'] == 28
assert len(pilot['cases']) == summary['unique_cases'] == 4
baseline = []
for case in pilot['cases']:
    selected = [r for r in rows if r['case_id'] == case['id']]
    assert len(selected) == 7
    keys = list(case['request']['questions']['declaration_relation']['criteria'])
    assert [r['option_order'] for r in selected[:6]] == [list(p) for p in permutations(keys)]
    assert selected[0]['raw_result']['answers'] == selected[-1]['raw_result']['answers']
    for row in selected:
        answer = row['raw_result']['answers']['declaration_relation']
        assert set(answer['probabilities']) == set(keys)
        assert abs(sum(answer['probabilities'].values()) - 1) < 0.001
        assert row['agrees_with_author_label'] == (answer['choice'] == case['expected'])
        assert not row['raw_result']['usage']['truncated']
        assert not row['raw_result']['usage']['state_tokens_dropped']
    baseline.append(selected[0])
assert sum(r['agrees_with_author_label'] for r in baseline) == summary['baseline_agrees_with_author']
budget = json.loads((RESULTS / 'question-budget.json').read_text())
assert budget['no_instruction_or_option_text_clipping'] is True
assert budget['instruction_tokens'] <= budget['instruction_budget']
assert all(n <= 48 for n in budget['option_tokens'])
checkpoint = json.loads((RESULTS / 'checkpoint.json').read_text())
assert checkpoint['revision'] == summary['checkpoint_revision']
for file in checkpoint['files']:
    path = Path(checkpoint['local_path']) / file['path']
    h = hashlib.sha256()
    with path.open('rb') as source:
        while block := source.read(1 << 20):
            h.update(block)
    assert h.hexdigest() == file['sha256']
    if file['hub_lfs_sha256']:
        assert h.hexdigest() == file['hub_lfs_sha256']
protocol = json.loads((HERE / 'training-protocol.json').read_text())
assert protocol['status'] == 'PLANNED_NOT_TRAINED'
assert protocol['leakage_controls']['final_test_is_not_used_for_prompt_selection_checkpoint_selection_or_temperature_fitting']
counts = Counter((r['expected_author_label'], r['raw_result']['answers']['declaration_relation']['choice']) for r in baseline)
receipt = {'real_inference_verified': True, 'independent_semantic_accuracy_measured': False,
           'baseline_author_agreement': summary['baseline_agrees_with_author'],
           'unique_diagnostics': 4, 'runs': 28,
           'confusion_on_known_author_labels': [{'expected': a, 'actual': b, 'count': n} for (a, b), n in sorted(counts.items())],
           'checkpoint_hashes_verified': True, 'state_and_question_budgets_verified': True,
           'training_performed': False}
(RESULTS / 'diagnostic-verification.json').write_text(json.dumps(receipt, indent=2) + '\n')
print(json.dumps(receipt))
