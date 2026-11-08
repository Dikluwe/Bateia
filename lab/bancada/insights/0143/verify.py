"""Verify the final gate, compatibility, real continuation and preserved evidence."""
import hashlib
import json
from pathlib import Path

ROOT=Path(__file__).resolve().parents[4]
LAB=Path(__file__).resolve().parent
CHANGED={'01_core/domain/insight_sites.go','01_core/usecase/insight_reference_leads.go','01_core/usecase/insight_followups.go','02_shell/cli/app.go','02_shell/cli/idea_insights.go','04_wiring/main.go'}
NEW={'01_core/domain/insight_neighborhood.go','01_core/domain/insight_neighborhood_test.go','01_core/usecase/insight_neighborhood.go','01_core/usecase/insight_neighborhood_frontier.go','01_core/usecase/insight_neighborhood_test.go','01_core/usecase/insight_neighborhood_frontier_test.go','02_shell/cli/insight_neighborhood.go','02_shell/cli/insight_neighborhood_test.go','04_wiring/insight_neighborhood_wiring_test.go'}
def read(path):return json.loads(path.read_text())
def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()
def main():
 baseline=read(LAB/'baseline.json')
 changed={f for f,h in baseline.items() if f.endswith('.go') and sha(ROOT/f)!=h}
 assert changed==CHANGED,changed
 current={str(p.relative_to(ROOT)) for layer in ['01_core','02_shell','03_infra','04_wiring'] for p in (ROOT/layer).rglob('*.go')}
 assert current-baseline.keys()==NEW,current-baseline.keys()
 assert all(sha(ROOT/f)==h for f,h in baseline.items() if f.startswith('lab/bancada/insights/0142/'))
 assert all(sha(ROOT/f)==h for f,h in baseline.items() if f.endswith('_test.go'))
 gates=read(LAB/'final/verification.json');assert all(v['exit_code']==0 for v in gates.values()),gates
 assert 'audited 391 productive Go files' in (LAB/'final/architecture.log').read_text()
 for step in ['0106','0134','0135','0136','0137','0138']:
  result=read(LAB/f'replay-{step}.json');assert result.get('replay_byte_identical',result.get('exploration_replay_equal',False)),step
 prior=read(ROOT/'lab/bancada/insights/0142/workflow-fixture.json');replay=read(LAB/'workflow-0142-replay.json')
 for key in ['first','second','feedback_memory']:assert prior[key]==replay[key],key
 input_path=ROOT/'lab/bancada/insights/pesquisa-real-memoria-autenticada-2026-10-28/round-2.json'
 old=read(input_path);live=read(LAB/'live-step-1.json');end=read(LAB/'live-step-2.json');receipt=read(LAB/'live-receipt.json')
 assert sha(input_path)==receipt['input_sha256']=='e3d5b57d273a5d076383d17d13c03cb637c1f6ef614ced4f9af7f891baf2519f'
 assert sha(LAB/'live-step-1.json')==receipt['output_sha256']
 assert receipt['exit_code']==0 and receipt['network'] and not receipt['synthetic_data']
 binary_sha=sha(LAB/'final/bateia')
 assert binary_sha==sha(LAB/'candidate-bateia')==receipt['binary_sha256']
 assert live['neighborhood']['input_report_id']==old['id'] and end['neighborhood']['input_report_id']==live['id']
 assert live['plan']['request']['memory']==old['exploration_memory'] and end['plan']['request']['memory']==live['exploration_memory']
 for report in [live,end]:
  assert report['calls']==0 and not report['attempts'] and not report['candidates'] and not report['bundles'] and not report.get('site_contexts')
  assert report['plan']['request']['discovery']['question']==old['plan']['request']['discovery']['question']
  assert report['novelty_state']=='NOT_ASSESSED' and report['effect_state']=='UNVERIFIED'
 assert live['exploration_memory']['generation']==old['exploration_memory']['generation']+1
 assert end['exploration_memory']['generation']==live['exploration_memory']['generation']+1
 assert len(live['followups'])==1 and live['followups'][0]['identity']['external_id']=='vitality-vis/rest-api'
 follow=live['followups'][0];assert follow['status']=='COLLECTED'
 item=follow['bundle']['items'][0]
 assert item['content_hash']=='sha256:'+hashlib.sha256(item['content'].encode()).hexdigest()
 for card in follow['cards']:
  ref=card['reference'];assert ref['content_hash']==item['content_hash']
  assert '\n'.join(item['content'].splitlines()[ref['start_line']-1:ref['end_line']])==ref['excerpt']
 assert end['status']=='FRONTIER_EXHAUSTED' and not end.get('followups')
 assert any(e['state']=='PENDING' for e in end['exploration_memory']['entries'])
 lineage=[]
 for name in sorted(CHANGED|NEW):
  prompt='00_nucleo/prompts/'+name.removesuffix('.go')+'.md'
  assert '@prompt '+prompt in (ROOT/name).read_text(),name
  lineage.append(dict(source=name,source_sha256=sha(ROOT/name),prompt=prompt,prompt_sha256=sha(ROOT/prompt)))
 (LAB/'lineage.json').write_text(json.dumps(lineage,indent=2)+'\n')
 result=dict(status='PASS',contract='0143',gate_exit_codes={k:v['exit_code'] for k,v in gates.items()},productive_files_audited=391,changed_existing_go=sorted(CHANGED),new_go=sorted(NEW),preexisting_tests_unchanged=True,preserved_0142_evidence_files=sum(f.startswith('lab/') for f in baseline),historical_replays_passed=7,real_catalog_calls=0,real_new_readme='vitality-vis/rest-api',real_commit=follow['bundle']['repository']['commit'],real_elapsed_seconds=receipt['elapsed_seconds'],focus_exhausted_without_spill=True,question_inherited=True,novelty_assessed=False,effect_verified=False,candidate_code_executed=False,final_binary_sha256=binary_sha)
 (LAB/'verification.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result,indent=2))
if __name__=='__main__':main()
