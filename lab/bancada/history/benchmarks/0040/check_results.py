#!/usr/bin/env python3
"""Public regression checks, not an independent or live-discovery evaluation."""
import json, pathlib, hashlib, sys
root=pathlib.Path(sys.argv[1])
def read(p):return json.loads((root/p).read_text())
def assessment(case):return read('out-'+case+'/assessment-report.json')['assessments'][0]
events=assessment('events');negated=assessment('negated');self_a=assessment('self')
assert events['capabilities'][0]['status']=='PARTIAL'
assert events['capabilities'][0]['confidence']<=.65
assert not events['reusable_parts']
assert self_a['license']['status']=='INCONCLUSIVE'
c=negated['constraints'][0]['evidence']
assert c['status']=='INCONCLUSIVE' and c['sources'][0]['line']==3
screen=read('out-comparison/screening-report.json')['decisions']
assert len(screen)==2 and all(d['signals']['description_alignment']==0 and d['signals']['capability_coverage']==0 for d in screen)
bundles=read('out-self/evidence/bundles.json')
skipped=[t['target']['path'] for b in bundles for t in b['target_results'] if t['status']=='SKIPPED_BUDGET']
assert skipped and not all(b['complete'] for b in bundles)
assert '01_core/usecase/evidence.go' in skipped
assert 'manutenção não verificada' in (root/'out-events/decision.md').read_text()
replay={}
for name in ['reuse-decision.json','assessment-report.json','alternatives.json','expansion-plan.json']:
 a=(root/'out-events'/name).read_bytes();b=(root/'replay-events'/name).read_bytes()
 assert a==b;replay[name]=hashlib.sha256(a).hexdigest()
controls={}
for name in ['empty','failed']:
 service=read(name+'-service.json')['Recommendation']
 mcp=json.loads(read(name+'-mcp.json')['result']['content'][0]['text'])
 assert service['decision']['decision']==mcp['decision']['decision']=='INCONCLUSIVE'
 assert bool(service['failures'])==bool(mcp['failures'])==(name=='failed')
 assert len(read(name+'-requests.json'))>0
 controls[name]='INCONCLUSIVE'
assert all(c['exit_code']==0 for c in read('commands.json'))
print(json.dumps({'status':'regressions_passed','independent_evaluation':False,'events_capability':'PARTIAL','events_reusable_parts':len(events['reusable_parts'] or []),'self_license':'INCONCLUSIVE','negated_non_goal':{'status':c['status'],'line':c['sources'][0]['line']},'screening':[{'score':d['ranking_score'],'description_alignment':d['signals']['description_alignment'],'capability_coverage':d['signals']['capability_coverage']} for d in screen],'self_skipped_count':len(skipped),'self_core_evidence_file':'SKIPPED_BUDGET','replay_sha256':replay,'service_mcp':controls,'notes_capability':assessment('notes')['capabilities'][0]},ensure_ascii=False,indent=2))
