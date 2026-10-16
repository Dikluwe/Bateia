"""Verify collected bytes and scope accounting; no semantic or model evaluation."""
from collections import Counter
import hashlib
import json
from pathlib import Path

HERE = Path(__file__).resolve().parent
bundles_path = HERE / 'results/github-bundles.json'
bundles = json.loads(bundles_path.read_text())
plans = json.loads((HERE / 'plans.json').read_text())
assert len(bundles) == len(plans) == 3
summary = []
for bundle, plan in zip(bundles, plans):
    assert bundle['repository'] == plan['repository']
    for item in bundle['items']:
        assert item['ref'] == bundle['repository']
        assert hashlib.sha256(item['content'].encode()).hexdigest() == item['content_hash']
    summary.append({'repository': bundle['repository'], 'items': len(bundle['items']),
                    'complete': bundle['complete'], 'bytes_read': bundle.get('bytes_read'),
                    'failures': bundle.get('failures', []),
                    'target_status_counts': dict(Counter(t['status'] for t in bundle.get('target_results', [])))})
fallback_path = HERE / 'results/decision-index-fallback.json'
fallback = json.loads(fallback_path.read_text())
assert fallback['repository'] == bundles[2]['repository']
assert hashlib.sha256(fallback['content'].encode()).hexdigest() == fallback['sha256']
assert fallback['origin'] == 'DIRECT_PUBLIC_HTTP_FALLBACK_AFTER_BATEIA_GITHUB_403'
assert any(f['stage'] == 'tree' and '403' in f['message'] for f in bundles[2]['failures'])
assert all(not b['complete'] for b in bundles)
assert {'README.md', 'laya/serve.py', 'pyproject.toml'} <= {i['path'] for i in bundles[0]['items']}
pilot = json.loads((HERE / 'pilot-requests.json').read_text())
root = HERE.parents[3]
conc = json.loads((root / 'lab/bancada/insights/0106/github-bundles.json').read_text())[0]
snapshots = [json.loads((root / ('lab/bancada/concorrencia/0120/testdata/' + name + '-snapshot.json')).read_text()) for name in ('tunny', 'pond')]
assert pilot['inference_executed'] is False and pilot['accuracy_measured'] is False
assert len(pilot['cases']) == 4
for case in pilot['cases']:
    assert set(case['request']['questions']['declaration_relation']['criteria']) == {'supports', 'contradicts', 'undetermined'}
    assert 'min_confidence' not in case['request']
    ref = case['source']
    if ref['repository'] == 'sourcegraph/conc':
        assert ref['commit'] == conc['repository']['commit']
        item = next(i for i in conc['items'] if i['path'] == ref['path'])
    else:
        snapshot = next(s for s in snapshots if s['repository'] == ref['repository'])
        assert ref['commit'] == snapshot['commit']
        item = next(i for i in snapshot['files'] if i['path'] == ref['path'])
    assert hashlib.sha256(item['content'].encode()).hexdigest() == ref['content_sha256']
    span = '\n'.join(item['content'].splitlines()[ref['first_line'] - 1:ref['last_line']])
    assert hashlib.sha256(span.encode()).hexdigest() == ref['segment_sha256']
    assert case['request']['state'].endswith('\nSOURCE: ' + span)
receipt = {'collector': 'BATEIA_EVIDENCE_COLLECT', 'candidate_code_executed': False,
           'model_inference_executed': False, 'semantic_accuracy_measured': False,
           'bundles': summary,
           'bundle_file_sha256': hashlib.sha256(bundles_path.read_bytes()).hexdigest(),
           'fallback_sha256': hashlib.sha256(fallback_path.read_bytes()).hexdigest(),
           'pilot_requests_sha256': hashlib.sha256((HERE / 'pilot-requests.json').read_bytes()).hexdigest()}
(HERE / 'results/collection-review.json').write_text(json.dumps(receipt, indent=2) + '\n')
print(json.dumps({'bundles_verified': 3, 'items_verified': sum(len(b['items']) for b in bundles),
                  'fallback_verified': True, 'model_inference_executed': False}))
