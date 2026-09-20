#!/usr/bin/env python3
"""Public synthetic CLI scenarios. Supplied PROVEN assessments test selection only."""
import json,pathlib,subprocess,hashlib,copy,sys
binary=pathlib.Path(sys.argv[1]).resolve();root=pathlib.Path(sys.argv[2]).resolve();root.mkdir(parents=True,exist_ok=True)
rules=pathlib.Path(__file__).resolve().with_name('rules.json')
def save(name,value):
 p=root/name;p.write_text(json.dumps(value,ensure_ascii=False,indent=2)+'\n');return str(p)
def read(name):return json.loads((root/name).read_text())
commands=[]
def run(args):
 result=subprocess.run([str(binary),*args],capture_output=True,text=True)
 commands.append({'args':args,'exit_code':result.returncode})
 assert result.returncode==0,result.stderr
 return result.stdout
ref={'owner':'local','repo':'target','commit':'b'*40}
idea=save('idea.json',{'problem':'code intelligence','capabilities':['dependency graph','blast radius']})
bundle={'schema_version':'1','repository':ref,'complete':True,'tree_complete':True,'items':[]}
for name,b in [('gaps',bundle),('partial',copy.deepcopy(bundle))]:
 if name=='partial':
  content='package graph\ntype DependencyGraph struct{}'
  b['items']=[{'kind':'file','ref':ref,'path':'graph.go','content':content,'content_hash':hashlib.sha256(content.encode()).hexdigest()}]
 run(['idea','autofagia','--idea',idea,'--bundles',save(name+'-bundles.json',[b]),'--rules',str(rules),'--target','local/target@'+'b'*40,'--output',str(root/(name+'-auto.json'))])
assert read('partial-auto.json')['partial']==['dependency graph']
assert any(o['action']=='ASSESS_INTERNAL_EXTENSION' for o in read('partial-auto.json')['options'])
policy=save('policy.json',{'schema_version':'1','max_projects':3,'max_components':3,'allowed_licenses':['MIT'],'require_immutable_ref':True})
reuse=save('reuse.json',{'schema_version':'1','allowed_licenses':['MIT']})
def candidate(repo,caps):
 r={'owner':'external','repo':repo,'commit':'a'*40}
 def source(kind,path,content):return {'kind':kind,'ref':r,'path':path,'line':1,'snippet':content,'rule_id':'supplied-assessment-fixture','content_hash':hashlib.sha256(content.encode()).hexdigest()}
 code=source('symbol','src/implementation.go','synthetic declared behavior; not an executed verification')
 return {'repository':r,'problem_alignment':{'status':'PROVEN','sources':[code]},'capabilities':[{'capability':c,'status':'PROVEN','sources':[code]} for c in caps],'license':{'status':'PROVEN','identifier':'MIT','sources':[source('license','LICENSE','MIT License')]},'maintenance':{'status':'PROVEN','sources':[source('repository_metadata','@repository-metadata','synthetic frozen activity record')]}}
a=candidate('graph',['dependency graph']);b=candidate('impact',['blast radius']);both=candidate('both',['dependency graph','blast radius'])
blocked=copy.deepcopy(both);blocked['constraints']=[{'constraint':'privacy','evidence':{'status':'CONFLICT'}}]
partial=copy.deepcopy(a);partial['capabilities'][0]['status']='PARTIAL';partial['maintenance']['status']='INCONCLUSIVE'
results={}
for name,assessments in [('complete',[a,b]),('minimum',[a,b,both]),('partial-coverage',[a]),('blocked',[blocked]),('investigate',[partial]),('empty',[])]:
 output=str(root/(name+'-proposal.json'))
 run(['idea','katamari','--autofagia',str(root/'gaps-auto.json'),'--assessments',save(name+'-assessments.json',{'schema_version':'1','assessments':assessments}),'--policy',policy,'--reuse-policy',reuse,'--output',output])
 p=read(name+'-proposal.json')
 results[name]={'status':p['status'],'components':len(p['components']),'feasible':p['feasible'],'coverage_complete':p['coverage_complete'],'options':len(p['options'])}
 assert not p['feasible']
assert results['complete']['status']=='COVERED_AWAITING_VALIDATION'
assert results['minimum']['components']==1
assert results['partial-coverage']['status']=='PARTIAL_COVERAGE'
assert results['blocked']['components']==0
assert results['investigate']['components']==0 and any(o['action']=='INVESTIGATE_EXTERNAL' for o in read('investigate-proposal.json')['options'])
assert results['empty']['status']!='NO_CHANGE_REQUIRED'
save('commands.json',commands)
save('results.json',{'status':'passed','scope':'synthetic CLI integration; supplied external assessments are not behavioral proof','cases':results})
print(json.dumps(results,ensure_ascii=False,indent=2))
