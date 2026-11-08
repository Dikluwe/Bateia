"""Continue a broader real focus with one attempt; suppress its known .git alias."""
import fcntl, hashlib, json, os, shlex, subprocess, tempfile, time
from pathlib import Path
root=Path(__file__).resolve().parent
repo=root.parents[3]
lock=(root/'live.lock').open('a');fcntl.flock(lock,fcntl.LOCK_EX|fcntl.LOCK_NB)
input_path=repo/'lab/bancada/insights/pesquisa-real-memoria-autenticada-2026-10-28/round-2.json'
output=root/'live-step-1.json'
if output.exists() or (root/'live-receipt.json').exists():raise SystemExit('Preserve previous live results')
binary=root/'candidate-bateia'
command=[str(binary),'idea','insights','expand','--input',str(input_path),'--focus','https://github.com/vitality-vis/frontend','--follow-links','1','--max-passages','8','--timeout','60s','--attempt-timeout','12s','--output',str(output)]
env=os.environ.copy()
for key in ['GITHUB_TOKEN','SOURCEGRAPH_TOKEN','BATEIA_GITEE_TOKEN','BATEIA_ATOMGIT_TOKEN','BATEIA_BITBUCKET_WORKSPACES']:env.pop(key,None)
for line in (repo/'.env').read_text().splitlines():
 line=line.strip()
 if line.startswith('export '):line=line[7:].lstrip()
 if line.startswith('GITHUB_TOKEN='):
  parts=shlex.split(line.split('=',1)[1],comments=True)
  if len(parts)==1:env['GITHUB_TOKEN']=parts[0]
if not env.get('GITHUB_TOKEN'):raise SystemExit('Configured GitHub credential unavailable')
input_hash=hashlib.sha256(input_path.read_bytes()).hexdigest()
start=time.monotonic()
with tempfile.TemporaryDirectory(prefix='bateia-neighborhood-live-') as cwd:
 result=subprocess.run(command,cwd=cwd,env=env,capture_output=True,timeout=80)
(root/'live-stdout.log').write_bytes(result.stdout);(root/'live-stderr.log').write_bytes(result.stderr)
receipt=dict(command=command,exit_code=result.returncode,elapsed_seconds=round(time.monotonic()-start,3),input_sha256=input_hash,binary_sha256=hashlib.sha256(binary.read_bytes()).hexdigest(),network=True,synthetic_data=False,candidate_code_executed=False)
assert input_hash==hashlib.sha256(input_path.read_bytes()).hexdigest()
if output.exists():
 report=json.loads(output.read_text())
 receipt.update(output_sha256=hashlib.sha256(output.read_bytes()).hexdigest(),catalog_calls=report['calls'],status=report['status'],reads=[dict(repository=f['identity']['external_id'],status=f['status']) for f in report.get('followups',[])])
(root/'live-receipt.json').write_text(json.dumps(receipt,indent=2)+'\n');print(json.dumps(receipt,indent=2))
