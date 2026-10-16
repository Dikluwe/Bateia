"""Exploratory correlated context controls; not a model quality evaluation."""
import hashlib,json,time
from pathlib import Path
import torch,laya
P=Path(__file__).resolve().parent
controls=json.loads((P/'controls.json').read_text())
checkpoint=json.loads((P.parent/'0123/results/checkpoint.json').read_text())
torch.set_num_threads(8);torch.set_num_interop_threads(1);torch.use_deterministic_algorithms(True)
agent=laya.load(checkpoint['local_path'],device='cpu')
rows=[]
for case in controls['cases']:
    start=time.perf_counter();result=agent.predict(case['state'],controls['question'])
    assert not result['usage']['truncated'] and result['usage']['state_tokens_dropped']==0
    rows.append({**case,'expected_author_only':controls['expected_author_only'][case['id']],'raw_result':result,'elapsed_seconds':time.perf_counter()-start})
(P/'results/context-controls.json').write_text(json.dumps({'checkpoint_revision':checkpoint['revision'],'question_sha256':hashlib.sha256(json.dumps(controls['question'],sort_keys=True).encode()).hexdigest(),'training_performed':False,'independent_accuracy_measured':False,'cases':rows},indent=2)+'\n')
for r in rows:
    a=r['raw_result']['answers']['declaration_relation']; print(r['id'],a['choice'],a.get('answer_confidence'))
