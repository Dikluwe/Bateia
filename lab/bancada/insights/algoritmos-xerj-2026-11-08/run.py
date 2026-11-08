"""Audit the open-problem arms, freeze baselines, separate setup and research."""
import argparse,fcntl,hashlib,json,os,re,subprocess,time
from pathlib import Path
from urllib.parse import urlparse
ROOT=Path('/tmp/bateia-xerj-algorithms-0155')
p=argparse.ArgumentParser();p.add_argument('arm',choices=['bateia','xerj']);p.add_argument('action',choices=['tool','acquire']);p.add_argument('args',nargs=argparse.REMAINDER);a=p.parse_args()
arm=ROOT/a.arm;logs=arm/'audit';logs.mkdir(exist_ok=True)
args=a.args[1:] if a.args[:1]==['--'] else a.args
baseline=arm/'baseline.json'
if not baseline.exists():raise SystemExit('Write baseline.json before any product call or acquisition')
json.loads(baseline.read_text())
env=os.environ.copy()
for key in list(env):
 if key.startswith('XERJ_') or key in ['GITHUB_TOKEN','GH_TOKEN','SOURCEGRAPH_TOKEN','BATEIA_GITEE_TOKEN','BATEIA_ATOMGIT_TOKEN']:env.pop(key,None)
env.update({'XERJ_DISABLE_FEEDBACK':'true','XERJ_HUB_HOME':str(arm/'hub'),'XERJ_CODE_HOME':str(arm/'code'),'XERJ_DATA_DIR':str(arm/'node-data')})

registry=logs/'eligible-sources.json'
eligible=json.loads(registry.read_text()) if registry.exists() else {}
def canonical_repo(raw):
 u=urlparse(raw)
 if u.scheme!='https' or u.hostname!='github.com' or u.username or u.password:return None
 parts=u.path.strip('/').split('/')
 if len(parts)<2 or not all(re.fullmatch(r'[A-Za-z0-9_.-]+',p or '') for p in parts[:2]):return None
 return ('https://github.com/'+parts[0]+'/'+parts[1].removesuffix('.git')).lower()
def native_urls(value):
 if isinstance(value,dict):
  for key,v in value.items():
   if key not in ['plan','request','query','question','familiar_approaches','args','command']:yield from native_urls(v)
 elif isinstance(value,list):
  for v in value:yield from native_urls(v)
 elif isinstance(value,str):
  for raw in re.findall(r'https://github\.com/[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+',value):
   key=canonical_repo(raw)
   if key:yield key

def check_tool_sources():
 for argument in args:
  for raw in re.findall(r'https://github\.com/[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+',argument):
   if canonical_repo(raw) not in eligible:raise SystemExit('Protocol rejects repository URL not observed in a product return')
 if a.arm=='xerj' and args[:2]==['autoindex','status']:raise SystemExit('Shared journal status excluded from isolated experiment; verify your own node by search')
 if a.arm=='bateia' and args[:3]==['idea','evidence','collect']:
  if '--plan' not in args:raise SystemExit('plan required')
  planfile=Path(args[args.index('--plan')+1]).resolve()
  if not planfile.is_relative_to(arm.resolve()):raise SystemExit('plan must belong to own arm')
  plans=json.loads(planfile.read_text());plans=plans if isinstance(plans,list) else [plans]
  for plan in plans:
   ref=plan.get('repository',{})
   key=canonical_repo('https://github.com/'+ref.get('owner','')+'/'+ref.get('repo',''))
   if key not in eligible:raise SystemExit('Protocol rejects collection from a repository not observed natively')
 if a.arm=='xerj' and args[:1]==['autoindex'] and len(args)>1 and args[1] not in ['--help','-h','map','remove']:
  target=Path(args[1]).resolve()
  if not target.is_relative_to((arm/'acquired').resolve()):raise SystemExit('Only acquired native-discovered sources may be indexed')

if a.action=='tool':check_tool_sources()

origin=None;kind='research'
if a.action=='tool':
 if not args:raise SystemExit('tool arguments required')
 if a.arm=='xerj' and args[:2]==['hub','feedback']:raise SystemExit('feedback disallowed')
 if '--install-timer' in args:raise SystemExit('timers disallowed')
 if '--help' in args or '-h' in args or args[:2] in [['hub','sync'],['hub','status'],['hub','stop']] or args[:1]==['--version'] or args[:1]==['--insecure']:kind='setup'
 command=[str(arm/'bin'/a.arm)]+args
else:
 q=argparse.ArgumentParser();q.add_argument('url');q.add_argument('name');q.add_argument('--origin',required=True,choices=['tool']);q.add_argument('--basis',required=True);q.add_argument('--branch');opts=q.parse_args(args)
 url=urlparse(opts.url)
 if url.scheme!='https' or url.username or url.password or url.query or url.fragment:raise SystemExit('exact public https git URL without credentials required')
 if '/' in opts.name or opts.name in ['.','..']:raise SystemExit('plain destination name required')
 key=canonical_repo(opts.url)
 if key not in eligible:raise SystemExit('Protocol rejects memory acquisition: URL needs a native product witness')
 kind='transport';origin={'kind':opts.origin,'basis':opts.basis,'url':opts.url,'native_witnesses':eligible[key]}
 dest=arm/'acquired'/opts.name;dest.parent.mkdir(exist_ok=True)
 command=['git','clone','--depth','1']+(['--branch',opts.branch,'--single-branch'] if opts.branch else [])+[opts.url,str(dest)]
start=time.time()
with (logs/'.lock').open('a') as lock:
 fcntl.flock(lock,fcntl.LOCK_EX)
 freeze=logs/'baseline-frozen.json'
 if not freeze.exists():freeze.write_bytes(baseline.read_bytes())
 if hashlib.sha256(baseline.read_bytes()).digest()!=hashlib.sha256(freeze.read_bytes()).digest():raise SystemExit('baseline changed after first audited action')
 calls=[json.loads(f.read_text()) for f in logs.glob('[0-9]*.json')]
 if kind=='research' and sum(c.get('kind')=='research' for c in calls)>=60:raise SystemExit('60 research calls exhausted')
 if kind=='transport' and sum(c.get('kind')=='transport' for c in calls)>=6:raise SystemExit('6 acquisitions exhausted')
 if calls and start-min(c['started_unix'] for c in calls)>1800 and kind not in ['setup']:raise SystemExit('30-minute research window exhausted; consolidate and stop your nodes')
 path=logs/f'{len(calls)+1:03d}.json'
 pending={'arm':a.arm,'action':a.action,'kind':kind,'args':args,'command':command,'origin':origin,'started_unix':start,'pending':True,'baseline_sha256':hashlib.sha256(freeze.read_bytes()).hexdigest()}
 path.write_text(json.dumps(pending)+'\n')
try:r=subprocess.run(command,cwd=arm,env=env,capture_output=True,text=True,timeout=1200 if args[:1]==['--insecure'] else 240)
except subprocess.TimeoutExpired as e:
 def txt(x):return x.decode(errors='replace') if isinstance(x,bytes) else x or ''
 class R:pass
 r=R();r.returncode=124;r.stdout=txt(e.stdout);r.stderr=txt(e.stderr)+'\nAudited command timed out.'
record={**pending,'pending':False,'seconds':time.time()-start,'exit_code':r.returncode,'stdout':r.stdout,'stderr':r.stderr}
if a.action=='acquire' and r.returncode==0:record['commit']=subprocess.check_output(['git','-C',str(dest),'rev-parse','HEAD'],text=True).strip()
path.write_text(json.dumps(record,ensure_ascii=False,indent=2)+'\n')

if a.action=='tool' and kind=='research':
 observations=[]
 try:observations.append(('stdout',json.loads(r.stdout)))
 except (ValueError,TypeError):observations.append(('stdout',r.stdout))
 if '--output' in args:
  output=Path(args[args.index('--output')+1])
  if output.exists() and output.stat().st_mtime>=start and output.resolve().is_relative_to(arm.resolve()):
   try:
    snapshot=logs/(path.stem+'.output.snapshot')
    snapshot.write_bytes(output.read_bytes())
    record['product_output_snapshot']=str(snapshot)
    record['product_output_sha256']=hashlib.sha256(snapshot.read_bytes()).hexdigest()
    observations.append((str(snapshot),json.loads(snapshot.read_text())))
   except ValueError:pass
 for witness,value in observations:
  for key in set(native_urls(value)):
   proof={'audit':str(path),'witness':witness,'action_args':args[:4]}
   if proof not in eligible.setdefault(key,[]):eligible[key].append(proof)
 registry.write_text(json.dumps(eligible,indent=2,ensure_ascii=False)+'\n')
 path.write_text(json.dumps(record,ensure_ascii=False,indent=2)+'\n')

print('Audit:',path);print('Exit:',r.returncode,'Seconds:',round(record['seconds'],3));print(r.stdout[:55000]);print(r.stderr[:10000]);raise SystemExit(r.returncode)
