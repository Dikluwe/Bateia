"""Fresh public Bateia run, with a real first-round snapshot as the second-round memory."""
import argparse
import fcntl
import hashlib
import json
import os
import shlex
from pathlib import Path
import subprocess
import tempfile
import time
from datetime import datetime, timezone

parser=argparse.ArgumentParser()
parser.add_argument('--round',type=int,choices=[1,2],required=True)
args=parser.parse_args()
root=Path(__file__).resolve().parent
repo=root.parents[3]
lock=(root/'execution.lock').open('a')
fcntl.flock(lock,fcntl.LOCK_EX|fcntl.LOCK_NB)
protocol=json.loads((root/'protocol.json').read_text())
output=root/f'round-{args.round}.json'
receipt_path=root/f'round-{args.round}-receipt.json'
if output.exists() or receipt_path.exists():raise SystemExit('Historical results preserved; use a new directory for another run.')
binary=repo/'lab/bancada/insights/0142/final/bateia'
expected=json.loads((repo/'lab/bancada/insights/0142/verification.json').read_text())['final_binary_sha256']
assert hashlib.sha256(binary.read_bytes()).hexdigest()==expected
command=[str(binary),'idea','insights','discover','--question',protocol['question_en'],'--sources',protocol['sources'],'--max-calls',str(protocol['max_calls']),'--max-per-query',str(protocol['max_per_query']),'--max-references',str(protocol['max_references']),'--max-passages',str(protocol['max_passages']),'--follow-links',str(protocol['follow_links']),'--timeout',str(protocol['total_timeout_seconds'])+'s','--attempt-timeout',str(protocol['attempt_timeout_seconds'])+'s','--output',str(output)]
for familiar in protocol['familiar_approaches']:command+=['--familiar',familiar]
memory_sha=None
if args.round==2:
    memory=root/'round-1.json'
    first=json.loads(memory.read_text())
    assert first['exploration_memory']['schema_version']=='insight-exploration-memory-v1'
    first_receipt=json.loads((root/'round-1-receipt.json').read_text())
    memory_sha=hashlib.sha256(memory.read_bytes()).hexdigest()
    assert memory_sha==first_receipt['result_sha256']
    command+=['--memory',str(memory)]
env=os.environ.copy()
for key in ['GITHUB_TOKEN','SOURCEGRAPH_TOKEN','BATEIA_GITEE_TOKEN','BATEIA_ATOMGIT_TOKEN','BATEIA_BITBUCKET_WORKSPACES']:env.pop(key,None)
token=None
for line in (repo/'.env').read_text().splitlines():
    line=line.strip()
    if line.startswith('export '):line=line[7:].lstrip()
    if line.startswith('GITHUB_TOKEN='):
        parts=shlex.split(line.split('=',1)[1],comments=True)
        if len(parts)==1:token=parts[0]
if not token:raise SystemExit('Configured GitHub credential unavailable')
env['GITHUB_TOKEN']=token
started=datetime.now(timezone.utc).isoformat();tick=time.monotonic()
with tempfile.TemporaryDirectory(prefix='bateia-real-memory-') as cwd:
    result=subprocess.run(command,cwd=cwd,env=env,capture_output=True,timeout=protocol['total_timeout_seconds']+20)
receipt={'round':args.round,'command':command,'started_at_utc':started,'elapsed_seconds':round(time.monotonic()-tick,3),'exit_code':result.returncode,'network_run':True,'synthetic_data':False,'no_tokens':False,'github_credential_preconfigured':True,'credential_logged':False,'isolated_cwd_without_dotenv':True,'candidate_code_executed':False,'binary_sha256':expected,'memory_input_sha256':memory_sha,'result_present':output.exists()}
(root/f'round-{args.round}-stdout.log').write_bytes(result.stdout)
(root/f'round-{args.round}-stderr.log').write_bytes(result.stderr)
if output.exists():
    receipt['result_sha256']=hashlib.sha256(output.read_bytes()).hexdigest()
    report=json.loads(output.read_text())
    receipt.update(status=report['status'],catalog_calls=report['calls'],candidates=len(report['candidates']),primary_documents=len(report['bundles'])+len(report.get('site_contexts',[])),followup_documents=sum('bundle' in f for f in report.get('followups',[])),memory_entries=len(report['exploration_memory']['entries']))
receipt_path.write_text(json.dumps(receipt,ensure_ascii=False,indent=2)+'\n')
print(json.dumps(receipt,ensure_ascii=False,indent=2))
