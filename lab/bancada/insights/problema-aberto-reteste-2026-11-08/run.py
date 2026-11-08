"""Audit the open-problem arms, freeze baselines, separate setup and research."""
import argparse,fcntl,hashlib,json,os,subprocess,time
from pathlib import Path
from urllib.parse import urlparse
ROOT=Path('/tmp/bateia-xerj-open-0153')
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
origin=None;kind='research'
if a.action=='tool':
 if not args:raise SystemExit('tool arguments required')
 if a.arm=='xerj' and args[:2]==['hub','feedback']:raise SystemExit('feedback disallowed')
 if '--install-timer' in args:raise SystemExit('timers disallowed')
 if '--help' in args or '-h' in args or args[:2] in [['hub','sync'],['hub','status'],['hub','stop']] or args[:1]==['--version'] or args[:1]==['--insecure']:kind='setup'
 command=[str(arm/'bin'/a.arm)]+args
else:
 q=argparse.ArgumentParser();q.add_argument('url');q.add_argument('name');q.add_argument('--origin',required=True,choices=['tool','model']);q.add_argument('--basis',required=True);q.add_argument('--branch');opts=q.parse_args(args)
 url=urlparse(opts.url)
 if url.scheme!='https' or url.username or url.password or url.query or url.fragment:raise SystemExit('exact public https git URL without credentials required')
 if '/' in opts.name or opts.name in ['.','..']:raise SystemExit('plain destination name required')
 kind='transport';origin={'kind':opts.origin,'basis':opts.basis,'url':opts.url}
 dest=arm/'acquired'/opts.name;dest.parent.mkdir(exist_ok=True)
 command=['git','clone','--depth','1']+(['--branch',opts.branch,'--single-branch'] if opts.branch else [])+[opts.url,str(dest)]
start=time.time()
with (logs/'.lock').open('a') as lock:
 fcntl.flock(lock,fcntl.LOCK_EX)
 freeze=logs/'baseline-frozen.json'
 if not freeze.exists():freeze.write_bytes(baseline.read_bytes())
 if hashlib.sha256(baseline.read_bytes()).digest()!=hashlib.sha256(freeze.read_bytes()).digest():raise SystemExit('baseline changed after first audited action')
 calls=[json.loads(f.read_text()) for f in logs.glob('[0-9]*.json')]
 if kind=='research' and sum(c.get('kind')=='research' for c in calls)>=40:raise SystemExit('40 research calls exhausted')
 if kind=='transport' and sum(c.get('kind')=='transport' for c in calls)>=6:raise SystemExit('6 acquisitions exhausted')
 if calls and start-min(c['started_unix'] for c in calls)>1200 and kind not in ['setup']:raise SystemExit('20-minute research window exhausted; consolidate and stop your nodes')
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
print('Audit:',path);print('Exit:',r.returncode,'Seconds:',round(record['seconds'],3));print(r.stdout[:55000]);print(r.stderr[:10000]);raise SystemExit(r.returncode)
