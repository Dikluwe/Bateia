#!/usr/bin/env python3
"""Open diagnostic, not an independent benchmark. Run against extracted dossier directory."""
import json,pathlib,sys,hashlib
root=pathlib.Path(sys.argv[1])
def read(name):return json.loads((root/name).read_text())
cases={}
for name in ['out-events','out-negated','out-notes','out-self','out-comparison']:
    a=read(name+'/assessment-report.json');d=read(name+'/reuse-decision.json')
    cases[name]={'decision':d['decision'],'assessments':[{'repository':x['repository'],
        'capabilities':[{'capability':c['capability'],'status':c['status'],'explanation':c['explanation']} for c in x['capabilities']],
        'constraints':x.get('constraints'),'maintenance':x['maintenance']['status'],
        'license':x['license']} for x in a['assessments']]}
assert cases['out-events']['assessments'][0]['capabilities'][0]['status']=='PROVEN'
assert cases['out-negated']['assessments'][0]['constraints'][0]['evidence']['status']=='CONFLICT'
assert cases['out-notes']['assessments'][0]['capabilities'][0]['explanation']=='nenhuma regra versionada cobre esta capacidade'
self_assessment=cases['out-self']['assessments'][0]
assert self_assessment['license']['sources'][0]['path']=='01_core/parser/license_checker_test.go'
self_items=[i['path'] for b in read('out-self/evidence/bundles.json') for i in b['items']]
missing=[p for p in ['01_core/usecase/recommendation.go','01_core/usecase/idea_run.go','01_core/usecase/evidence.go'] if p not in self_items]
assert len(missing)==3
screen=read('out-comparison/screening-report.json')['decisions']
assert len(screen)==2 and all(d['signals']['description_alignment']==1 and d['signals']['capability_coverage']==1 for d in screen)
replay={}
for name in ['reuse-decision.json','assessment-report.json','alternatives.json','expansion-plan.json']:
    a=(root/'out-events'/name).read_bytes();b=(root/'replay-events'/name).read_bytes()
    assert a==b;replay[name]=hashlib.sha256(a).hexdigest()
controls={}
for name in ['empty','failed']:
    service=read(name+'-service.json')['Recommendation']
    mcp=json.loads(read(name+'-mcp.json')['result']['content'][0]['text'])
    requests=read(name+'-requests.json');assert len(requests)>0
    assert service['decision']['decision']==mcp['decision']['decision']=='INCONCLUSIVE'
    assert bool(service['failures'])==bool(mcp['failures'])==(name=='failed')
    controls[name]={'decision':'INCONCLUSIVE','provider_requests':len(requests),'failures':mcp['failures']}
print(json.dumps({'diagnostic':'completed_with_quality_findings','cases':cases,'self_collected_paths':self_items,
 'self_missing_core_files':missing,'comparison_scores':[d['ranking_score'] for d in screen],
 'replay_identical_sha256':replay,'controlled_service_and_mcp':controls},indent=2,ensure_ascii=False))
