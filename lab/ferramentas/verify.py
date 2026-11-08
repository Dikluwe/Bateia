"""Canonical local gate: declared toolchain, structural checks, product tests, vet, build."""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[2]
output = Path(os.environ.get('BATEIA_VERIFY_OUTPUT', tempfile.mkdtemp(prefix='bateia-verify-'))).resolve()
output.mkdir(parents=True, exist_ok=True)
env = os.environ.copy()
env.setdefault('GOCACHE', str(Path(tempfile.gettempdir()) / 'bateia-go-cache'))
n = int(env.get('GIT_CONFIG_COUNT', '0'))
env.update({'GIT_CONFIG_COUNT':str(n+1), 'GIT_CONFIG_KEY_'+str(n):'commit.gpgsign', 'GIT_CONFIG_VALUE_'+str(n):'false'})
commands = [
 ('lint', [sys.executable,'lab/ferramentas/crystalline-lint/run.py','--fail-on','warning','.']),
 ('architecture',['go','test','-v','./lab/bancada/arquitetura/0129/boundaries','./lab/bancada/arquitetura/0131/contracts','./lab/bancada/arquitetura/0145']),
 ('test',['go','test','./...']), ('vet',['go','vet','./...']),
 ('build',['go','build','-o',str(output/'bateia'),'./04_wiring']),
 ('diff',['git','diff','--check']),
]
checks = {}
for name, command in commands:
    with (output/(name+'.log')).open('w') as log:
        result = subprocess.run(command, cwd=ROOT, env=env, stdout=log, stderr=subprocess.STDOUT)
    checks[name] = {'command':command,'exit_code':result.returncode}
    print(name+': '+str(result.returncode), flush=True)
(output/'verification.json').write_text(json.dumps(checks, indent=2)+'\n')
print('Evidence: '+str(output))
raise SystemExit(max(item['exit_code'] for item in checks.values()))
