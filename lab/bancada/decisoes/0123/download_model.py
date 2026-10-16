"""Download only a pinned English checkpoint into the isolated /tmp environment."""
import hashlib
import json
from pathlib import Path
from huggingface_hub import HfApi, snapshot_download

HERE = Path(__file__).resolve().parent
metadata = json.loads((HERE / 'results/hub-metadata.json').read_text())
destination = Path('/tmp/bateia-laya-0123/english-checkpoint')
info = HfApi().model_info(metadata['id'], revision=metadata['sha'], files_metadata=True, token=False)
assert info.sha == metadata['sha']
snapshot_download(metadata['id'], revision=metadata['sha'], token=False,
                  allow_patterns=['model.safetensors', 'rl_agent_config.json',
                                  'encoder/config.json', 'tokenizer/*'],
                  local_dir=destination, max_workers=2)
files = []
for item in info.siblings:
    path = destination / item.rfilename
    if not path.is_file():
        continue
    h = hashlib.sha256()
    with path.open('rb') as source:
        while block := source.read(1 << 20):
            h.update(block)
    expected = item.lfs.sha256 if item.lfs else None
    if expected:
        assert h.hexdigest() == expected
    files.append({'path': item.rfilename, 'bytes': path.stat().st_size,
                  'sha256': h.hexdigest(), 'hub_lfs_sha256': expected})
assert any(f['path'] == 'model.safetensors' for f in files)
(HERE / 'results/checkpoint.json').write_text(json.dumps({
    'repository': metadata['id'], 'revision': metadata['sha'],
    'checkpoint': 'ENGLISH_ROOT', 'local_path': str(destination),
    'files': files}, indent=2) + '\n')
print(json.dumps({'files': len(files), 'bytes': sum(f['bytes'] for f in files),
                  'revision': metadata['sha']}))
