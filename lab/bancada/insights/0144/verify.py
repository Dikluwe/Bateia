"""Verify gates, appended tests, immutable history and real alias-aware continuation."""
import hashlib
import json
from pathlib import Path
ROOT=Path(__file__).resolve().parents[4]
LAB=Path(__file__).resolve().parent
CHANGED={'01_core/domain/insight_neighborhood.go','01_core/domain/insight_neighborhood_test.go','01_core/usecase/insight_neighborhood.go','01_core/usecase/insight_neighborhood_frontier.go','01_core/usecase/insight_followups.go','02_shell/cli/insight_neighborhood.go','02_shell/cli/insight_neighborhood_test.go'}
NEW={'01_core/usecase/insight_neighborhood_selection.go','01_core/usecase/insight_neighborhood_selection_test.go'}
APPENDED={'01_core/domain/insight_neighborhood_test.go':'domain-test.go.txt','02_shell/cli/insight_neighborhood_test.go':'cli-test.go.txt'}
def read(p):return json.loads(p.read_text())
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def main():
 baseline=read(LAB/'baseline.json')
 changed={f for f,h in baseline.items() if f.endswith('.go') and sha(ROOT/f)!=h};assert changed==CHANGED,changed
 current={str(p.relative_to(ROOT)) for layer in ['01_core','02_shell','03_infra','04_wiring'] for p in (ROOT/layer).rglob('*.go')}
 assert current-baseline.keys()==NEW,current-baseline.keys()
 for name,old in APPENDED.items():
  previous=(LAB/'before'/old).read_bytes();assert hashlib.sha256(previous).hexdigest()==baseline[name]
  assert (ROOT/name).read_bytes().startswith(previous),name
 assert all(sha(ROOT/f)==h for f,h in baseline.items() if f.endswith('_test.go') and f not in APPENDED)
 preserved=[f for f in baseline if f.startswith('lab/')];assert all(sha(ROOT/f)==baseline[f] for f in preserved)
 gates=read(LAB/'final/verification.json');assert all(v['exit_code']==0 for v in gates.values()),gates
 assert 'audited 392 productive Go files' in (LAB/'final/architecture.log').read_text()
 for step in ['0106','0134','0135','0136','0137','0138']:
  result=read(LAB/f'replay-{step}.json');assert result.get('replay_byte_identical',result.get('exploration_replay_equal',False)),step
 previous=read(ROOT/'lab/bancada/insights/0142/workflow-fixture.json');replay=read(LAB/'workflow-0142-replay.json')
 for k in ['first','second','feedback_memory']:assert previous[k]==replay[k]
 compatibility=read(LAB/'replay-0143-input.json');old=read(ROOT/'lab/bancada/insights/0143/live-step-1.json')
 assert compatibility['neighborhood']['input_report_id']==old['id'] and compatibility['neighborhood']['focus']==old['neighborhood']['focus']
 assert compatibility['calls']==0 and not compatibility.get('followups') and compatibility['status']=='FRONTIER_EXHAUSTED'
 input_path=ROOT/'lab/bancada/insights/pesquisa-real-memoria-autenticada-2026-10-28/round-2.json'
 original=read(input_path);first=read(LAB/'live-step-1.json');second=read(LAB/'live-step-2.json');general=read(LAB/'live-general-preview.json');receipt=read(LAB/'live-receipt.json')
 assert sha(input_path)==receipt['input_sha256']=='e3d5b57d273a5d076383d17d13c03cb637c1f6ef614ced4f9af7f891baf2519f'
 assert receipt['exit_code']==0 and receipt['network'] and not receipt['synthetic_data']
 assert receipt['output_sha256']==sha(LAB/'live-step-1.json')
 binary_sha=sha(LAB/'final/bateia');assert binary_sha==sha(LAB/'candidate-bateia')==receipt['binary_sha256']
 assert first['neighborhood']['input_report_id']==original['id'] and first['plan']['request']['memory']==original['exploration_memory']
 assert second['neighborhood']['input_report_id']==first['id'] and second['plan']['request']['memory']==first['exploration_memory']
 assert first['neighborhood']['focus']==second['neighborhood']['focus'] and first['neighborhood']['focus']['external_id']=='vitality-vis/frontend'
 for report in [first,second,general]:
  assert report['calls']==0 and not report['attempts'] and not report['candidates'] and not report['bundles']
  assert report['plan']['request']['discovery']['question']==original['plan']['request']['discovery']['question']
  assert report['novelty_state']=='NOT_ASSESSED' and report['effect_state']=='UNVERIFIED'
 assert len(first['followups'])==1 and first['followups'][0]['identity']['external_id']=='vitality-vis/rest-api'
 follow=first['followups'][0];assert follow['status']=='COLLECTED'
 item=follow['bundle']['items'][0];assert item['content_hash']=='sha256:'+hashlib.sha256(item['content'].encode()).hexdigest()
 for card in follow['cards']:
  ref=card['reference'];assert ref['content_hash']==item['content_hash']
  assert '\n'.join(item['content'].splitlines()[ref['start_line']-1:ref['end_line']])==ref['excerpt']
 neighbors=first['neighborhood']['neighbors'];assert [(n['identity']['external_id'],n['selection_status'],n['distance']) for n in neighbors]==[('vitality-vis/rest-api','FRONTIER',1),('vitality-vis/frontend.git','ALIAS_ALREADY_EXAMINED',0)]
 raw_alias=next(e for e in original['exploration_memory']['entries'] if e['identity']['external_id']=='vitality-vis/frontend.git')
 for report in [first,second,general]:assert next(e for e in report['exploration_memory']['entries'] if e['id']==raw_alias['id'])==raw_alias
 assert second['status']=='FRONTIER_EXHAUSTED' and not second.get('followups')
 assert any(e['state']=='PENDING' for e in second['exploration_memory']['entries'])
 assert general['neighborhood'].get('focus') is None and general['status']=='FRONTIER_PENDING'
 assert len(general['neighborhood']['frontier_entry_ids'])==11 and all('distance' not in n for n in general['neighborhood']['neighbors'])
 lineage=[]
 for name in sorted(CHANGED|NEW):
  prompt='00_nucleo/prompts/'+name.removesuffix('.go')+'.md';assert '@prompt '+prompt in (ROOT/name).read_text()
  lineage.append(dict(source=name,source_sha256=sha(ROOT/name),prompt=prompt,prompt_sha256=sha(ROOT/prompt)))
 (LAB/'lineage.json').write_text(json.dumps(lineage,indent=2)+'\n')
 result=dict(status='PASS',contract='0144',gate_exit_codes={k:v['exit_code'] for k,v in gates.items()},productive_files_audited=392,changed_go=sorted(CHANGED),new_go=sorted(NEW),preexisting_test_bodies_preserved=True,preserved_previous_evidence_files=len(preserved),historical_controls_passed=8,real_catalog_calls=0,real_read='vitality-vis/rest-api',read_budget=1,known_alias_suppressed=True,suppressed_alias_memory_unchanged=True,inherited_focus_verified=True,all_trails_preview_pending=11,real_elapsed_seconds=receipt['elapsed_seconds'],real_commit=follow['bundle']['repository']['commit'],candidate_code_executed=False,novelty_assessed=False,effect_verified=False,final_binary_sha256=binary_sha)
 (LAB/'verification.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result,indent=2))
if __name__=='__main__':main()
