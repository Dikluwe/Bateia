"""Freeze existing collector bytes and local standard-library export data; never run conc."""
import base64
import hashlib
import json
import os
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parents[4]
OUT = Path(__file__).resolve().parent
COMMIT = "5f936abd7ae87036af1f75c95fb9d0daaf00116b"
bundle_path = ROOT / "lab/bancada/insights/0106/github-bundles.json"
bundle = json.loads(bundle_path.read_text())[0]
files = []
source_receipts = []
for item in bundle["items"]:
    name = item["path"]
    if name.endswith(".go") and not name.endswith("_test.go") and (name.startswith("pool/") or name.startswith("panics/") or name == "waitgroup.go"):
        assert item["ref"]["commit"] == COMMIT
        content = item["content"]
        digest = hashlib.sha256(content.encode()).hexdigest()
        assert digest == item["content_hash"]
        assert "//go:build" not in content and "// +build" not in content
        files.append({"path": name, "content": content, "sha256": digest})
        source_receipts.append({"path": name, "sha256": digest})
files.sort(key=lambda f: f["path"])
process_env = os.environ.copy()
process_env.update({"GOOS": "linux", "GOARCH": "amd64", "CGO_ENABLED": "0", "GOWORK": "off", "GOTOOLCHAIN": "local", "GOPROXY": "off"})
command = ["go", "list", "-deps", "-export", "-json", "context", "sync", "errors", "fmt", "sync/atomic", "runtime", "runtime/debug", "sort"]
raw = subprocess.check_output(command, env=process_env, cwd=ROOT)
decoder = json.JSONDecoder()
text = raw.decode()
packages = []
while text.strip():
    item, end = decoder.raw_decode(text.lstrip())
    packages.append(item)
    text = text.lstrip()[end:]
exports = []
for p in packages:
    assert p.get("Standard") and not p.get("Error")
    if p["ImportPath"] == "unsafe":
        continue
    archive = Path(p["Export"]).read_bytes()
    assert archive.startswith(b"!<arch>\n")
    offset = 8
    data = None
    while offset < len(archive):
        header = archive[offset:offset + 60]
        size = int(header[48:58])
        if header[:16].strip().rstrip(b"/") == b"__.PKGDEF":
            data = archive[offset + 60:offset + 60 + size]
            break
        offset += 60 + size + size % 2
    assert data is not None
    header = f"{'__.PKGDEF':<16}{0:<12}{0:<6}{0:<6}{644:<8}{len(data):<10}`\n".encode()
    assert len(header) == 60
    data = b"!<arch>\n" + header + data + (b"\n" if len(data) % 2 else b"")
    exports.append({"import_path": p["ImportPath"], "origin": "LOCAL_GO_STANDARD_EXPORT", "data": base64.b64encode(data).decode(), "sha256": hashlib.sha256(data).hexdigest()})
version = subprocess.check_output(["go", "version"], text=True).strip()
environment = {"schema_version": "1", "toolchain": version, "goos": "linux", "goarch": "amd64", "packages": [{"import_path": "github.com/sourcegraph/conc", "directory": "."}, {"import_path": "github.com/sourcegraph/conc/panics", "directory": "panics"}, {"import_path": "github.com/sourcegraph/conc/pool", "directory": "pool"}], "exports": sorted(exports, key=lambda e: e["import_path"])}
request = {"schema_version": "1", "origin": "sourcegraph/conc@" + COMMIT + "; collector 0106 non-test pool closure", "build": {"go_version": version, "goos": "linux", "goarch": "amd64"}, "files": files}
receipt = {"schema_version": "1", "source_bundle_sha256": hashlib.sha256(bundle_path.read_bytes()).hexdigest(), "commit": COMMIT, "selection": "COLLECTED_NON_TEST_GO_POOL_PANICS_ROOT_NO_BUILD_TAGS", "complete_repository_build_selection_verified": False, "test_files_selected": False, "standard_export_command": command, "toolchain": version, "goos": "linux", "goarch": "amd64", "cgo_enabled": False, "candidate_code_executed": False, "standard_export_bytes": sum(len(base64.b64decode(e["data"])) for e in exports), "sources": source_receipts, "standard_exports": [{"import_path": e["import_path"], "sha256": e["sha256"]} for e in environment["exports"]]}
for name, value in [("conc-input.json", request), ("environment.json", environment), ("preparation-receipt.json", receipt)]:
    destination = OUT / "testdata" / name
    # Preserve old receipts: preparation refuses to replace an existing freeze.
    with destination.open("x") as stream:
        json.dump(value, stream, indent=2)
        stream.write("\n")
print(json.dumps({"source_files": len(files), "standard_exports": len(exports), "export_bytes": receipt["standard_export_bytes"]}))
