"""Verify trace, scope and deterministic mechanism controls, not semantic accuracy."""
import hashlib,json
from pathlib import Path
P=Path(__file__).resolve().parent
b=json.loads((P/'results/bundles.json').read_text())[0]
assert b['repository']['commit']=='3cf26cbcb18725dbc2d127bb8bb2c4c43243ae63'
for item in b['items']:
 assert item['content_hash'].removeprefix('sha256:')==hashlib.sha256(item['content'].encode()).hexdigest()
for t in json.loads((P/'plans.json').read_text())[0]['targets']:
 assert any(i['path']==t['path'] for i in b['items'])
ins=json.loads((P/'results/insights.json').read_text())
assert ins['status']=='INCONCLUSIVE' and not ins['mechanisms'] and not ins['hypotheses']
assert not json.loads((P/'insights-input.json').read_text())['sources'][0]['assessments']
c=json.loads((P/'results/context-controls.json').read_text())
assert len(c['cases'])==6 and not c['training_performed']
assert c['cases'][0]['raw_result']['answers']==c['cases'][-1]['raw_result']['answers']
for row in c['cases']:
 assert not row['raw_result']['usage']['truncated']
m=json.loads((P/'results/mechanism-controls.json').read_text())
assert m['split_calibration']['overlapping_families']==['one_repository']
assert m['question_budget']['full_instruction_tokens']>m['question_budget']['used_instruction_tokens']
assert not m['question_budget']['state_stats']['truncated']
# Confirm imported common.py matches the actual newly collected source.
installed=Path('/tmp/bateia-laya-0123/venv/lib/python3.11/site-packages/laya/common.py')
collected=next(i for i in b['items'] if i['path']=='laya/common.py')
assert installed.read_bytes()==collected['content'].encode()
print('PASS: collected hashes, targets, explicit empty assessments, diagnostic trace and mechanism controls')
