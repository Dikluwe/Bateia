"""Check current gates, declared workflow fixtures and preservation, without network."""
import hashlib
import json
from pathlib import Path
ROOT=Path(__file__).resolve().parents[4]
LAB=Path(__file__).resolve().parent
CHANGED={'01_core/domain/insight_sites.go','01_core/usecase/insight_site_plan.go','01_core/usecase/insight_sites.go','01_core/usecase/insight_site_reading.go','02_shell/cli/app.go','02_shell/cli/idea_insights.go','02_shell/cli/insight_discovery.go','04_wiring/insight_sites.go'}
NEW={'01_core/domain/insight_exploration_memory.go','01_core/domain/insight_exploration_memory_test.go','01_core/usecase/insight_exploration_memory.go','01_core/usecase/insight_followups.go','01_core/usecase/insight_exploration_memory_test.go','01_core/usecase/insight_followups_test.go','02_shell/cli/insight_memory.go','02_shell/cli/insight_memory_test.go','03_infra/github/insight_reference_resolver.go','03_infra/github/insight_reference_resolver_test.go','04_wiring/insight_memory_wiring_test.go'}
def read(p):return json.loads(p.read_text())
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def main():
 baseline=read(LAB/'baseline.json')
 assert {f for f,h in baseline.items() if sha(ROOT/f)!=h}==CHANGED
 current={str(p.relative_to(ROOT)) for d in ['01_core','02_shell','03_infra','04_wiring'] for p in (ROOT/d).rglob('*.go')}
 assert current-baseline.keys()==NEW
 lineage=[]
 for name in sorted(CHANGED|NEW):
  prompt='00_nucleo/prompts/'+name.removesuffix('.go')+'.md'
  assert '@prompt '+prompt in (ROOT/name).read_text()
  lineage.append(dict(source=name,source_sha256=sha(ROOT/name),prompt=prompt,prompt_sha256=sha(ROOT/prompt)))
 gates=read(LAB/'final/verification.json');assert all(v['exit_code']==0 for v in gates.values())
 assert 'audited 387 productive Go files' in (LAB/'final/architecture.log').read_text()
 for step in ['0106','0134','0135','0136','0137','0138']:
  result=read(LAB/('replay-'+step+'.json'));assert result.get('replay_byte_identical',result.get('exploration_replay_equal',False)),step
 cli=read(LAB/'cli-workflow-check.json');assert cli['result']=='PASS' and cli['binary_entrypoint_verified']
 workflow=read(LAB/'workflow-fixture.json')
 assert workflow['mode']=='OFFLINE_WORKFLOW_FIXTURE' and not workflow['network']
 assert workflow['child_documents']=='SYNTHETIC_NOT_PUBLIC_OBSERVATIONS'
 first,second=workflow['first'],workflow['second']
 assert len(first['bundles'])==2 and len(second['bundles'])==0
 assert set(workflow['first_destinations']).isdisjoint(workflow['second_destinations'])
 assert len(workflow['first_destinations'])==len(workflow['second_destinations'])==2
 assert workflow['first_destinations'][0]=='chnm/serendipomatic'
 for report,examined in [(first,4),(second,6)]:
  memory=report['exploration_memory'];assert memory['schema_version']=='insight-exploration-memory-v1'
  assert len(memory['entries'])==14 and sum(e['state']=='EXAMINED' for e in memory['entries'])==examined
  entries={e['id']:e for e in memory['entries']}
  cards=set(c['id'] for c in report['exploration']['cards'])
  for scoped in report.get('site_contexts',[]):cards.update(c['id'] for c in scoped['cards'])
  total_bytes=sum(b['bytes_read'] for b in report['bundles'])
  for follow in report['followups']:
   assert follow['entry_id'] in entries
   if 'bundle' in follow:
    bundle=follow['bundle'];assert follow['status']=='COLLECTED'
    assert entries[follow['entry_id']]['state']=='EXAMINED'
    assert bundle['repository']['commit']=='c'*40
    for item in bundle['items']:
     assert 'Synthetic destination fixture' in item['content']
     assert item['content_hash']=='sha256:'+hashlib.sha256(item['content'].encode()).hexdigest()
    total_bytes+=bundle['bytes_read'];cards.update(c['id'] for c in follow['cards'])
  assert total_bytes<=128<<10 and len(cards)<=6
  for task in report['tasks']:
   assert set(task['memory_entry_refs'])==set(entries)
   assert set(task['reading_card_refs'])==cards
 assert any(e['state']=='DISCARDED' and e['reason'] for e in workflow['feedback_memory']['entries'])
 (LAB/'lineage.json').write_text(json.dumps(lineage,indent=2)+'\n')
 print(json.dumps(dict(status='PASS',contract='0142',gate_exit_codes={k:v['exit_code'] for k,v in gates.items()},productive_files_audited=387,
 changed_existing_go=sorted(CHANGED),new_go=sorted(NEW),preexisting_tests_unchanged=True,
 preserved_previous_evidence_files=sum(f.startswith('lab/') for f in baseline),historical_replays_passed=6,
 offline_rounds=2,first_followups=2,second_followups=2,examined_targets_not_reread=True,feedback_reason_preserved=True,
 github_follow_reader_protocol_fixture_passed=True,cli_discard_and_reopen_passed=True,live_destination_contents_verified=False,
 child_documents='SYNTHETIC',novelty_assessed=False,effect_verified=False,candidate_code_executed=False,
 final_binary_sha256=sha(LAB/'final/bateia')),ensure_ascii=False,indent=2))
if __name__=='__main__':main()
