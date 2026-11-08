import json,os,signal,subprocess,sys,time
from pathlib import Path
root=Path('/tmp/bateia-xerj-demo-0156/xerj');state=root/'node-lifecycle.json';binary=root/'bin/xerj';cmd=[str(binary),'--insecure','--port','19701','--data-dir',str(root/'node-data'),'--config',str(root/'node.toml')]
if sys.argv[1]=='start':
 if state.exists():raise SystemExit('Node already recorded')
 env=os.environ.copy()
 for key in list(env):
  if key.startswith('XERJ_') or key in ['GITHUB_TOKEN','GH_TOKEN','SOURCEGRAPH_TOKEN']:env.pop(key,None)
 env.update({'XERJ_DISABLE_FEEDBACK':'true','XERJ_HUB_HOME':str(root/'hub'),'XERJ_CODE_HOME':str(root/'code'),'XERJ_DATA_DIR':str(root/'node-data')})
 with (root/'node.log').open('wb') as log:p=subprocess.Popen(cmd,cwd=root,env=env,stdin=subprocess.DEVNULL,stdout=log,stderr=log,start_new_session=True)
 record={'pid':p.pid,'command':cmd,'port':19701,'started_unix':time.time(),'log':str(root/'node.log')};state.write_text(json.dumps(record,indent=2)+'\n');print(json.dumps(record))
elif sys.argv[1]=='stop':
 record=json.loads(state.read_text());pid=record['pid'];proc=Path('/proc')/str(pid)
 if proc.exists():
  actual=(proc/'cmdline').read_bytes().split(b'\0')
  if str(binary).encode() not in actual or str(root/'node-data').encode() not in actual:raise SystemExit('PID no longer belongs to this node; no signal sent')
  os.kill(pid,signal.SIGINT);record['signal']='SIGINT'
  for _ in range(30):
   if not proc.exists() or (proc/'stat').read_text().split()[2]=='Z':break
   time.sleep(.1)
 record['stopped_unix']=time.time();record['process_exited']=not proc.exists() or (proc/'stat').read_text().split()[2]=='Z';state.write_text(json.dumps(record,indent=2)+'\n');print(json.dumps(record))
