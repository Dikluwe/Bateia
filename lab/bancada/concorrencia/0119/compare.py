"""Compare selected behavioral claims with unchanged upstream execution events."""
import hashlib
import json
from pathlib import Path

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
RESULTS = HERE / 'results'
def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

execution = json.loads((RESULTS / 'execution.json').read_text())
claims_path = HERE.parent / '0118/claims.json'
review = json.loads(claims_path.read_text())
bundle_path = ROOT / 'lab/bancada/insights/0106/github-bundles.json'
assert execution['bundle_sha256'] == digest(bundle_path)
assert execution['repository_commit'] == review['repository_commit']
assert execution['exit_code'] == 0
assert execution['reference_origin'] == 'UPSTREAM_TESTS_UNMODIFIED'
assert execution['candidate_code_executed'] is True
assert execution['independent_semantic_accuracy_measured'] is False
for key, path in [('stdout_sha256', 'upstream-tests.jsonl'),
                  ('stderr_sha256', 'upstream-tests.stderr'),
                  ('modules_sha256', 'modules.jsonstream'), ('go_sum_sha256', 'go.sum')]:
    assert execution[key] == digest(RESULTS / path)
bundle = json.loads(bundle_path.read_text())[0]
expected_sources = {}
for item in bundle['items']:
    if item['path'] == '@repository-metadata':
        continue
    assert item['ref']['commit'] == execution['repository_commit']
    expected_sources[item['path']] = hashlib.sha256(item['content'].encode()).hexdigest()
    assert expected_sources[item['path']] == item['content_hash']
assert execution['source_hashes_before_and_after'] == expected_sources
events = [json.loads(line) for line in (RESULTS / 'upstream-tests.jsonl').read_text().splitlines()]
assert not any(event['Action'] == 'fail' for event in events)
assert any(event['Action'] == 'pass' and 'Test' not in event for event in events)
counts = {}
for event in events:
    if event['Action'] == 'pass' and 'Test' in event:
        counts[event['Test']] = counts.get(event['Test'], 0) + 1
assert counts == execution['pass_counts']
references = {
    'CONFIGURE_BEFORE_GO': ('ExamplePool', 1, 'CONFIGURE_BEFORE_GO_SUCCEEDS_IN_EXAMPLE'),
    'CONFIGURE_AFTER_GO': ('TestPool/panics_on_configuration_after_init/before_wait', 3, 'EXPECTED_PANIC_OBSERVED'),
    'CONFIGURE_AFTER_WAIT': ('TestPool/panics_on_configuration_after_init/after_wait', 3, 'EXPECTED_PANIC_OBSERVED_CLAIM_CONTRADICTED'),
    'REUSE_AFTER_WAIT': ('TestPool/is_reusable', 3, 'TWO_BATCHES_COMPLETE_IN_UPSTREAM_TEST')
}
comparison = []
for claim in review['claims']:
    entry = {'claim_id': claim['id'], 'prior_judgment': claim['judgment']}
    if claim['id'] in references:
        name, repetitions, outcome = references[claim['id']]
        assert counts.get(name) == repetitions
        entry.update(upstream_test=name, passes=repetitions,
                     status='CONSISTENT_WITH_PRIOR_REVIEW_IN_EXECUTED_CASE',
                     observed_behavior=outcome)
    else:
        entry.update(status='NOT_EVALUATED_BY_SELECTED_TESTS')
    comparison.append(entry)
report = {'schema_version': 'bounded-behavioral-comparison-v1',
          'reference_origin': execution['reference_origin'],
          'candidate_code_executed': True,
          'independent_semantic_accuracy_measured': False,
          'evaluated_behavioral_claims': 4, 'unassessed_claims': 4,
          'selection': 'PURPOSEFUL_EXISTING_CLAIMS_NOT_REPRESENTATIVE_CORPUS',
          'claims_sha256': digest(claims_path),
          'execution_sha256': digest(RESULTS / 'execution.json'),
          'comparisons': comparison}
(RESULTS / 'comparison.json').write_text(json.dumps(report, indent=2) + '\n')
print(json.dumps({'behavioral_claims_consistent': 4, 'unassessed': 4,
                  'independent_semantic_accuracy_measured': False}))
