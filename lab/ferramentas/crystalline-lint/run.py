"""Build the pinned development tool and run its unmodified architecture rules."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import tarfile
import tempfile

HERE = Path(__file__).resolve().parent

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def ensure_tool():
    manifest = json.loads((HERE / 'manifest.json').read_text())
    archive_path = HERE / manifest.get('archive_file', 'source.tar.gz')
    assert digest(archive_path) == manifest['archive_sha256'], 'linter source digest mismatch'
    assert digest(HERE / 'go-parser.patch') == manifest['patch_sha256'], 'linter patch digest mismatch'
    identity = manifest['archive_sha256'][:16] + '-' + manifest['patch_sha256'][:16]
    cache = Path(os.environ.get('BATEIA_TOOL_CACHE', str(Path(tempfile.gettempdir()) / 'bateia-tool-cache')))
    root = cache / identity
    binary = root / 'target/debug/crystalline-lint'
    receipt = root / 'validated.json'
    if binary.is_file() and receipt.is_file():
        data = json.loads(receipt.read_text())
        if data.get('binary_sha256') == digest(binary):
            return binary
    root.mkdir(parents=True, exist_ok=True, mode=0o700)
    source = root / 'source'
    if not source.exists():
        source.mkdir()
        with tarfile.open(archive_path) as archive:
            archive.extractall(source, filter='data')
        subprocess.run(['git', 'apply', '--check', str(HERE / 'go-parser.patch')], cwd=source, check=True)
        subprocess.run(['git', 'apply', str(HERE / 'go-parser.patch')], cwd=source, check=True)
    env = os.environ.copy()
    env['CARGO_TARGET_DIR'] = str(root / 'target')
    for command in [['cargo', 'test', '--locked', '--offline', '--lib'], ['cargo', 'build', '--locked', '--offline']]:
        subprocess.run(command, cwd=source, env=env, check=True, stdout=sys.stderr)
    subprocess.run([sys.executable, str(HERE / 'regressions.py'), str(binary)], check=True, stdout=sys.stderr)
    receipt.write_text(json.dumps({'revision':manifest['revision'], 'archive_sha256':manifest['archive_sha256'], 'patch_sha256':manifest['patch_sha256'], 'binary_sha256':digest(binary)}, indent=2) + '\n')
    return binary

if __name__ == '__main__':
    tool = ensure_tool()
    subprocess.run([sys.executable, str(HERE / 'regressions.py'), str(tool)], check=True)
    env = os.environ.copy()
    env['PATH'] = str(tool.parent) + os.pathsep + env['PATH']
    raise SystemExit(subprocess.run(['crystalline-lint', *(sys.argv[1:] or ['--fail-on','warning','.'])], env=env).returncode)
