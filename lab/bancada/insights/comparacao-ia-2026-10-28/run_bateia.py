"""One fresh bounded anonymous run, isolated from dotenv and candidate execution."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import tempfile
import time
from datetime import datetime, timezone

root=Path(__file__).resolve().parent
if any((root/name).exists() for name in ('bateia-result.json','bateia-run.json')):
    raise SystemExit('Existing historical run preserved: use a new experiment directory.')
repo=root.parents[3]
protocol=json.loads((root/'protocol.json').read_text())
freeze=json.loads((root/'direct-freeze.json').read_text())
assert hashlib.sha256((root/freeze['file']).read_bytes()).hexdigest()==freeze['sha256']
binary=repo/'lab/bancada/insights/0140/final/bateia'
env=os.environ.copy()
for key in ('GITHUB_TOKEN','SOURCEGRAPH_TOKEN','BATEIA_GITEE_TOKEN','BATEIA_ATOMGIT_TOKEN','BATEIA_BITBUCKET_WORKSPACES'):
    env.pop(key,None)
command=[str(binary),'idea','insights','discover','--sources','all','--question',protocol['question_en'],'--max-calls','20','--max-per-query','3','--max-references','6','--max-passages','6','--timeout','180s','--attempt-timeout','8s','--output',str(root/'bateia-result.json')]
for familiar in protocol['familiar_approaches']:command+=['--familiar',familiar]
start=datetime.now(timezone.utc).isoformat()
clock=time.monotonic()
with tempfile.TemporaryDirectory(prefix='bateia-comparison-') as directory:
    result=subprocess.run(command,cwd=directory,env=env,capture_output=True,timeout=195)
elapsed=time.monotonic()-clock
(root/'bateia-stdout.log').write_bytes(result.stdout)
(root/'bateia-stderr.log').write_bytes(result.stderr)
receipt=dict(command=command,started_at_utc=start,elapsed_seconds=round(elapsed,3),exit_code=result.returncode,isolated_without_dotenv=True,credentials_used=False,candidate_code_executed=False,binary_sha256=hashlib.sha256(binary.read_bytes()).hexdigest(),direct_arm_sha256=freeze['sha256'],result_present=(root/'bateia-result.json').exists())
if receipt['result_present']:
    receipt['result_sha256']=hashlib.sha256((root/'bateia-result.json').read_bytes()).hexdigest()
(root/'bateia-run.json').write_text(json.dumps(receipt,ensure_ascii=False,indent=2)+'\n')
print(json.dumps(receipt,ensure_ascii=False,indent=2))
