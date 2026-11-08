"""Run the old and reorganized CLI offline against the same frozen reports."""
import hashlib,json,os,subprocess,tempfile
from pathlib import Path
ROOT=Path(__file__).resolve().parents[4]
LAB=Path(__file__).resolve().parent
old=ROOT/'lab/bancada/insights/0144/final/bateia';new=LAB/'candidate-bateia'
assert hashlib.sha256(old.read_bytes()).hexdigest()==json.loads((ROOT/'lab/bancada/insights/0144/verification.json').read_text())['final_binary_sha256']
out=LAB/'comparisons';out.mkdir(exist_ok=False)
env=os.environ.copy()
for key in ['GITHUB_TOKEN','SOURCEGRAPH_TOKEN','BATEIA_GITEE_TOKEN','BATEIA_ATOMGIT_TOKEN','BATEIA_BITBUCKET_WORKSPACES']:env.pop(key,None)
r2=ROOT/'lab/bancada/insights/pesquisa-real-memoria-autenticada-2026-10-28/round-2.json'
first=ROOT/'lab/bancada/insights/0144/live-step-1.json';second=ROOT/'lab/bancada/insights/0144/live-step-2.json'
with tempfile.TemporaryDirectory(prefix='bateia-tekt-preservation-') as cwd:
 feedback=out/'feedback.json'
 command=[str(old),'idea','insights','memory','--input',str(r2),'--target','https://github.com/vitality-vis/rest-api','--state','discarded','--reason','Offline preservation case','--output',str(feedback)]
 prepared=subprocess.run(command,cwd=cwd,env=env,capture_output=True,timeout=20);assert prepared.returncode==0,prepared.stderr
 cases={
  'focus':['--input',str(r2),'--focus','https://github.com/vitality-vis/frontend'],
  'inherited-focus':['--input',str(first)],
  'all-trails':['--input',str(second),'--all-trails'],
  'alias-focus':['--input',str(r2),'--focus','https://github.com/vitality-vis/frontend.git'],
  'other-trail':['--input',str(r2),'--focus','https://github.com/MuseCompNet/muse-tech-central'],
  'feedback':['--input',str(r2),'--focus','https://github.com/vitality-vis/frontend','--memory',str(feedback)],
  'overrides':['--input',str(r2),'--all-trails','--max-passages','1','--timeout','30s','--attempt-timeout','5s'],
 }
 checks=[]
 for name,args in cases.items():
  statuses=[];raw=[]
  for tag,binary in [('old',old),('new',new)]:
   target=out/f'{name}.{tag}.json';cmd=[str(binary),'idea','insights','expand',*args,'--follow-links','0','--output',str(target)]
   result=subprocess.run(cmd,cwd=cwd,env=env,capture_output=True,timeout=20)
   (out/f'{name}.{tag}.stderr.log').write_bytes(result.stderr)
   assert result.returncode==0,(name,tag,result.stderr)
   data=target.read_bytes();report=json.loads(data)
   assert report['calls']==0 and all('bundle' not in f for f in report.get('followups',[]))
   statuses.append(report['status']);raw.append(data)
  assert raw[0]==raw[1],name
  checks.append(dict(case=name,byte_identical=True,status=statuses[0],result_sha256=hashlib.sha256(raw[0]).hexdigest()))
 for name,args in {
  'unknown-focus':['--input',str(r2),'--focus','https://github.com/unknown/repository'],
  'conflicting-scope':['--input',str(r2),'--focus','https://github.com/vitality-vis/frontend','--all-trails'],
 }.items():
  messages=[]
  for tag,binary in [('old',old),('new',new)]:
   target=out/f'{name}.{tag}.json';target.write_bytes(b'preserve')
   result=subprocess.run([str(binary),'idea','insights','expand',*args,'--output',str(target)],cwd=cwd,env=env,capture_output=True,timeout=20)
   assert result.returncode!=0 and target.read_bytes()==b'preserve'
   message=result.stderr.decode().split('Erro na expansão de insights: ',1)[-1]
   messages.append(message)
  assert messages[0]==messages[1]
  checks.append(dict(case=name,rejection_identical=True,previous_output_preserved=True))
(LAB/'comparison.json').write_text(json.dumps(dict(network=False,document_reads=0,old_binary_sha256=hashlib.sha256(old.read_bytes()).hexdigest(),new_binary_sha256=hashlib.sha256(new.read_bytes()).hexdigest(),cases=checks),ensure_ascii=False,indent=2)+'\n')
print(json.dumps(checks,ensure_ascii=False,indent=2))
