"""Bounded multi-site public run; credentials and candidate execution are excluded."""
import argparse
import json
import os
from pathlib import Path
import subprocess
import tempfile

parser=argparse.ArgumentParser()
parser.add_argument('--binary',required=True)
args=parser.parse_args()
root=Path(__file__).resolve().parent
env=os.environ.copy()
for key in ('GITHUB_TOKEN','SOURCEGRAPH_TOKEN','BATEIA_GITEE_TOKEN','BATEIA_ATOMGIT_TOKEN','BATEIA_BITBUCKET_WORKSPACES'):
    env.pop(key,None)
command=[str(Path(args.binary).resolve()),'idea','insights','discover','--sources','all','--question','backpressure queue cancellation','--familiar','queue','--max-per-query','1','--max-references','4','--max-passages','4','--max-calls','12','--timeout','45s','--attempt-timeout','2s','--output',str(root/'online-result.json')]
with tempfile.TemporaryDirectory(prefix='bateia-sites-public-') as directory:
    result=subprocess.run(command,cwd=directory,env=env,capture_output=True,timeout=55)
(root/'online-stdout.log').write_bytes(result.stdout)
(root/'online-stderr.log').write_bytes(result.stderr)
receipt=dict(command=command,exit_code=result.returncode,no_credentials=True,isolated_cwd_without_dotenv=True,candidate_code_executed=False,result_present=(root/'online-result.json').exists())
if receipt['result_present']:
    report=json.loads((root/'online-result.json').read_text())
    receipt.update(status=report['status'],sources=[cap['id'] for cap in report['plan']['sources']],calls=report['calls'],candidates=len(report['candidates']),bundles=len(report['bundles']),attempts=[dict(source=a['connector_id'],status=a['status'],code=a.get('error_code','')) for a in report['attempts']])
(root/'online-check.json').write_text(json.dumps(receipt,ensure_ascii=False,indent=2)+'\n')
print(json.dumps({k:v for k,v in receipt.items() if k not in ('command','attempts')},ensure_ascii=False))
