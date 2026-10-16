"""Freeze public upstream files at explicitly selected commits, without execution."""
from concurrent.futures import ThreadPoolExecutor
import hashlib
import json
from pathlib import Path
import urllib.request

HERE = Path(__file__).resolve().parent
REPOS = [('Jeffail/tunny', 'a274c3ce48a6dc4f7b5fb5b8eaefbbd9e23574b9'),
         ('alitto/pond', '7b7c16ad4903c9d7014beac3ed3d2533b410f170')]
def get(url):
    request = urllib.request.Request(url, headers={'User-Agent': 'Bateia-contract-evaluation'})
    with urllib.request.urlopen(request, timeout=30) as response:
        data = response.read(1024 * 1024 + 1)
        assert len(data) <= 1024 * 1024
        return data

for repo, commit in REPOS:
    tree_url = f'https://api.github.com/repos/{repo}/git/trees/{commit}?recursive=1'
    tree = json.loads(get(tree_url))
    assert not tree.get('truncated')
    paths = sorted(item['path'] for item in tree['tree'] if item['type'] == 'blob'
                   and not item['path'].startswith('examples/')
                   and (item['path'].endswith('.go') or item['path'] in ('README.md', 'go.mod', 'LICENSE')))
    def collect(path):
        url = f'https://raw.githubusercontent.com/{repo}/{commit}/{path}'
        data = get(url)
        return {'path': path, 'url': url, 'sha256': hashlib.sha256(data).hexdigest(),
                'content': data.decode('utf-8')}
    with ThreadPoolExecutor(max_workers=4) as executor:
        files = list(executor.map(collect, paths))
    snapshot = {'repository': repo, 'commit': commit, 'date': '2026-10-07',
                'selection': 'GO_AND_CONTRACT_FILES_EXCLUDING_EXAMPLES',
                'candidate_executed': False, 'files': files}
    with (HERE / 'testdata' / (repo.split('/')[-1] + '-snapshot.json')).open('x') as output:
        json.dump(snapshot, output, indent=2, ensure_ascii=False)
        output.write('\n')
    print(json.dumps({'repository': repo, 'commit': commit, 'files': len(files)}))
