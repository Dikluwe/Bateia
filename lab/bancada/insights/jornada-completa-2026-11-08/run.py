"""Audit product calls and exact-URL transport in the two cold-start arms."""
import argparse,hashlib,json,os,subprocess,time
from pathlib import Path
from urllib.parse import urlparse
ROOT=Path('/tmp/bateia-xerj-cold-0150')
p=argparse.ArgumentParser();p.add_argument('arm',choices=['bateia','xerj']);p.add_argument('action',choices=['tool','acquire']);p.add_argument('args',nargs=argparse.REMAINDER);a=p.parse_args()
arm=ROOT/a.arm;logs=arm/'audit';logs.mkdir(exist_ok=True)
args=a.args[1:] if a.args[:1]==['--'] else a.args
calls=list(logs.glob('*.json'))
if a.action=='tool' and sum(json.loads(f.read_text())['action']=='tool' for f in calls)>=24:raise SystemExit('24 product calls exhausted')
env=os.environ.copy()
for key in list(env):
 if key.startswith('XERJ_') or key in ['GITHUB_TOKEN','GH_TOKEN','SOURCEGRAPH_TOKEN','BATEIA_GITEE_TOKEN','BATEIA_ATOMGIT_TOKEN']:env.pop(key,None)
env.update({'XERJ_DISABLE_FEEDBACK':'true','XERJ_HUB_HOME':str(arm/'hub'),'XERJ_CODE_HOME':str(arm/'code'),'XERJ_DATA_DIR':str(arm/'node-data')})
origin=None
if a.action=='tool':
 if not args:raise SystemExit('tool arguments required')
 if a.arm=='xerj' and args[0]=='hub' and len(args)>1 and args[1]=='feedback':raise SystemExit('feedback disallowed')
 if '--install-timer' in args:raise SystemExit('timers disallowed')
 command=[str(arm/'bin'/a.arm)]+args
else:
 q=argparse.ArgumentParser();q.add_argument('url');q.add_argument('name');q.add_argument('--origin',required=True,choices=['target','tool','model']);q.add_argument('--basis',required=True);q.add_argument('--branch');opts=q.parse_args(args)
 url=urlparse(opts.url)
 if url.scheme!='https' or url.username or url.password or url.query or url.fragment:raise SystemExit('exact public https git URL without credentials required')
 if '/' in opts.name or opts.name in ['.','..']:raise SystemExit('plain destination name required')
 origin={'kind':opts.origin,'basis':opts.basis,'url':opts.url}
 dest=arm/'acquired'/opts.name;dest.parent.mkdir(exist_ok=True)
 command=['git','clone','--depth','1']+(['--branch',opts.branch,'--single-branch'] if opts.branch else [])+[opts.url,str(dest)]
start=time.time()
# Reserve an audit entry before waiting, so a long-lived node cannot overwrite
# the record of a concurrent query in the same arm.
index=1
while True:
 path=logs/f'{index:03d}.json'
 try:
  with path.open('x') as f:f.write(json.dumps({'arm':a.arm,'action':a.action,'args':args,'command':command,'origin':origin,'started_unix':start,'pending':True}))
  break
 except FileExistsError:index+=1
try:r=subprocess.run(command,cwd=arm,env=env,capture_output=True,text=True,timeout=1200 if args and args[0]=='--insecure' else 240)
except subprocess.TimeoutExpired as e:
 def txt(x):return x.decode(errors='replace') if isinstance(x,bytes) else x or ''
 class R:pass
 r=R();r.returncode=124;r.stdout=txt(e.stdout);r.stderr=txt(e.stderr)+'\nTransport/product call exceeded 240 seconds.'
record={'arm':a.arm,'action':a.action,'args':args,'command':command,'origin':origin,'started_unix':start,'seconds':time.time()-start,'exit_code':r.returncode,'stdout':r.stdout,'stderr':r.stderr}
path.write_text(json.dumps(record,ensure_ascii=False,indent=2)+'\n')
print('Audit:',path);print('Exit:',r.returncode,'Seconds:',round(record['seconds'],3));print(r.stdout[:50000]);print(r.stderr[:12000])
raise SystemExit(r.returncode)
