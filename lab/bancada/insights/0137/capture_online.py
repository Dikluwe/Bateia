"""Bounded public trial of the actual CLI, without credentials or candidate execution."""
import argparse
import json
import os
from pathlib import Path
import subprocess
import tempfile

parser = argparse.ArgumentParser()
parser.add_argument('--binary', required=True)
args = parser.parse_args()
root = Path(__file__).resolve().parent
binary = str(Path(args.binary).resolve())
env = os.environ.copy()
env.pop('GITHUB_TOKEN', None)
command = [binary, 'idea', 'insights', 'discover', '--question', 'backpressure queue cancellation',
           '--familiar', 'queue', '--familiar', 'thread pool', '--max-per-query', '1',
           '--max-references', '2', '--max-passages', '4', '--timeout', '45s',
           '--attempt-timeout', '6s', '--output', str(root/'online-result.json')]
with tempfile.TemporaryDirectory(prefix='bateia-insight-discovery-public-') as directory:
    result = subprocess.run(command, cwd=directory, env=env, capture_output=True, timeout=55)
(root/'online-stdout.log').write_bytes(result.stdout)
(root/'online-stderr.log').write_bytes(result.stderr)
receipt = dict(command=command, exit_code=result.returncode, public_no_github_token=True,
               isolated_cwd_without_dotenv=True, candidate_code_executed=False,
               result_present=(root/'online-result.json').exists())
if receipt['result_present']:
    report = json.loads((root/'online-result.json').read_text())
    receipt.update(status=report['status'], attempts=len(report['attempts']),
                   attempt_statuses=[a['status'] for a in report['attempts']],
                   candidates=len(report['candidates']), readmes=[r['status'] for r in report['readmes']],
                   cards=len(report['exploration']['cards']))
(root/'online-check.json').write_text(json.dumps(receipt, ensure_ascii=False, indent=2)+'\n')
print(json.dumps({k:v for k,v in receipt.items() if k!='command'}, ensure_ascii=False))
