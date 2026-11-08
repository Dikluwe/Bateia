import os,sys,json,subprocess,time
from pathlib import Path
if len(sys.argv)!=2: raise SystemExit('usage: search.py "query"')
query=sys.argv[1]
env=os.environ.copy();env.update({'XERJ_DISABLE_FEEDBACK':'true','XERJ_HUB_OFFLINE':'1'})
command=['/tmp/xerj-official-rc101/unpacked/xerj-1.0.0-rc.101-x86_64-unknown-linux-gnu/xerj','search','--url','http://127.0.0.1:19390','--prefix','agenttest','-k','8','--full','2200','--json',query]
r=subprocess.run(command,env=env,cwd='/tmp',capture_output=True,text=True)
with Path('/tmp/xerj-agent-session-20261108/queries.jsonl').open('a') as f:f.write(json.dumps({'query':query,'command':command,'exit_code':r.returncode,'stdout':r.stdout,'stderr':r.stderr})+'\n')
print(r.stdout);print(r.stderr,file=sys.stderr);raise SystemExit(r.returncode)
