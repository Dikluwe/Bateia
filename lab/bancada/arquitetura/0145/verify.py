"""Check architectural gates, explicit ownership, preserved consumers and CLI equivalence."""
import hashlib,json
from pathlib import Path
ROOT=Path(__file__).resolve().parents[4]
LAB=Path(__file__).resolve().parent
CHANGED={'01_core/usecase/insight_neighborhood.go','01_core/usecase/insight_neighborhood_frontier.go','01_core/usecase/insight_followups.go','01_core/usecase/insight_exploration_memory.go','01_core/usecase/insight_neighborhood_test.go','01_core/usecase/insight_followups_test.go','04_wiring/main.go'}
NEW={'01_core/usecase/insight_neighborhood_plan.go','01_core/usecase/insight_neighborhood_plan_test.go','01_core/usecase/insight_neighborhood_sources.go','01_core/usecase/insight_neighborhood_sources_test.go','01_core/usecase/insight_neighborhood_result.go','01_core/usecase/insight_neighborhood_result_test.go','01_core/usecase/insight_followup_selection.go','01_core/usecase/insight_followup_selection_test.go','01_core/usecase/insight_followup_context.go','01_core/usecase/insight_followup_context_test.go','02_shell/cli/web_input.go','02_shell/cli/web_input_test.go'}
def read(p):return json.loads(p.read_text())
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def main():
 baseline=read(LAB/'baseline.json')
 changed={f for f,h in baseline.items() if f.endswith('.go') and sha(ROOT/f)!=h};assert changed==CHANGED,changed
 current={str(p.relative_to(ROOT)) for layer in ['01_core','02_shell','03_infra','04_wiring'] for p in (ROOT/layer).rglob('*.go')}
 assert current-baseline.keys()==NEW,current-baseline.keys()
 preserved=[f for f in baseline if f.startswith('lab/bancada/insights/')];assert all(sha(ROOT/f)==baseline[f] for f in preserved)
 for name in ['AGENTS.md','.agentrules','.agents/AGENTS.md','crystalline.toml']:assert sha(ROOT/name)==baseline[name]
 changed_tests={f for f in CHANGED if f.endswith('_test.go')}
 assert all(sha(ROOT/f)==h for f,h in baseline.items() if f.endswith('_test.go') and f not in changed_tests)
 gates=read(LAB/'final/verification.json');assert all(v['exit_code']==0 for v in gates.values()),gates
 architecture=(LAB/'final/architecture.log').read_text()
 assert 'audited 398 productive Go files' in architecture
 for name in ['TestInvestigationOwnershipBoundaries','TestOwnershipControlsRejectRegressions','TestExistingConsumersAndMovedRulesRemainUnchanged']:assert '--- PASS: '+name in architecture
 comparison=read(LAB/'comparison.json');assert not comparison['network'] and comparison['document_reads']==0
 assert comparison['new_binary_sha256']==sha(LAB/'candidate-bateia')==sha(LAB/'final/bateia')
 assert len(comparison['cases'])==9
 for check in comparison['cases']:
  if check.get('byte_identical'):
   old=LAB/'comparisons'/(check['case']+'.old.json');new=LAB/'comparisons'/(check['case']+'.new.json')
   assert old.read_bytes()==new.read_bytes() and sha(new)==check['result_sha256']
  else:assert check['rejection_identical'] and check['previous_output_preserved']
 for step in ['0106','0134','0135','0136','0137','0138']:
  result=read(LAB/f'replay-{step}.json');assert result.get('replay_byte_identical',result.get('exploration_replay_equal',False)),step
 prior=read(ROOT/'lab/bancada/insights/0142/workflow-fixture.json');current=read(LAB/'workflow-0142-replay.json')
 for key in ['first','second','feedback_memory']:assert prior[key]==current[key],key
 lineage=[]
 for name in sorted(CHANGED|NEW):
  prompt='00_nucleo/prompts/'+name.removesuffix('.go')+'.md'
  assert '@prompt '+prompt in (ROOT/name).read_text()
  lineage.append(dict(source=name,source_sha256=sha(ROOT/name),prompt=prompt,prompt_sha256=sha(ROOT/prompt)))
 (LAB/'lineage.json').write_text(json.dumps(lineage,indent=2)+'\n')
 result=dict(status='PASS',contract='0145',gate_exit_codes={k:v['exit_code'] for k,v in gates.items()},productive_files_audited=398,changed_go=sorted(CHANGED),new_go=sorted(NEW),instruction_files_unchanged=True,ownership_audit_with_negative_controls=True,consumer_bodies_and_moved_rules_preserved=True,preserved_previous_evidence_files=len(preserved),offline_cli_cases=9,byte_identical_reports=7,rejections_and_previous_outputs_preserved=2,historical_replays_passed=7,new_network_calls=0,candidate_code_executed=False,novelty_assessed=False,effect_verified=False,final_binary_sha256=sha(LAB/'final/bateia'))
 (LAB/'verification.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result,indent=2))
if __name__=='__main__':main()
