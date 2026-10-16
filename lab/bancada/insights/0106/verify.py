"""Reproduce step 0106 without executing code from the collected repository."""

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile
import urllib.request


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--binary", required=True)
    parser.add_argument("--online", action="store_true")
    args = parser.parse_args()
    root = Path(__file__).resolve().parent
    binary = str(Path(args.binary).resolve())
    request = json.loads((root / "github-insight-input.json").read_text())
    bundle = request["sources"][0]["bundle"]
    assert bundle == json.loads((root / "github-bundles.json").read_text())[0]
    assert bundle["complete"] and not bundle.get("failures")
    ref = bundle["repository"]
    assert ref["owner"] == "sourcegraph" and ref["repo"] == "conc"
    assert ref["commit"] == "5f936abd7ae87036af1f75c95fb9d0daaf00116b"
    items = {item["path"]: item for item in bundle["items"]}
    assert len(items) == len(bundle["items"])
    for item in items.values():
        assert item["ref"] == ref
        assert item["content_hash"].removeprefix("sha256:") == hashlib.sha256(
            item["content"].encode()
        ).hexdigest()
    for witness in request["sources"][0]["assessments"][0]["capabilities"][0]["sources"]:
        item = items[witness["path"]]
        assert witness["ref"] == ref
        assert witness["content_hash"] == item["content_hash"]
        assert witness["snippet"] == item["content"].splitlines()[witness["line"] - 1]

    online_verified = []
    if args.online:
        for path in ("pool/pool.go", "pool/pool_test.go"):
            url = f"https://raw.githubusercontent.com/{ref['owner']}/{ref['repo']}/{ref['commit']}/{path}"
            req = urllib.request.Request(url, headers={"User-Agent": "Bateia-step-0106"})
            with urllib.request.urlopen(req, timeout=30) as response:
                assert response.status == 200
                data = response.read((1 << 20) + 1)
                assert len(data) <= 1 << 20
            assert data == items[path]["content"].encode()
            online_verified.append({"path": path, "url": url, "sha256": hashlib.sha256(data).hexdigest()})

    original = (root / "github-insight-result.json").read_bytes()
    command = [binary, "idea", "insights", "analyze", "--input"]
    replay = subprocess.run(command + [str(root / "github-insight-input.json")], capture_output=True, timeout=30)
    assert replay.returncode == 0, replay.stderr.decode()
    assert replay.stdout == original
    report = json.loads(original)
    assert report["status"] == "PARTIAL"
    assert report["mechanisms"][0]["status"] == "PARTIAL"
    assert report["hypotheses"][0]["status"] == "INCONCLUSIVE"
    assert report["hypotheses"][0]["target_state"] == "UNKNOWN"
    with tempfile.TemporaryDirectory(prefix="bateia-check-0106-") as directory:
        scratch = Path(directory)
        altered = json.loads(json.dumps(request))
        altered["sources"][0]["bundle"]["items"][0]["content"] += " altered"
        input_path = scratch / "input.json"
        output_path = scratch / "protected.json"
        input_path.write_text(json.dumps(altered))
        output_path.write_bytes(original)
        rejection = subprocess.run(command + [str(input_path), "--output", str(output_path)], capture_output=True, timeout=30)
        assert rejection.returncode != 0
        assert "integrity" in rejection.stderr.decode().lower() or "integridade" in rejection.stderr.decode().lower()
        assert output_path.read_bytes() == original
        request["sources"][0]["assessments"] = []
        input_path.write_text(json.dumps(request))
        control = subprocess.run(command + [str(input_path)], capture_output=True, timeout=30)
        assert control.returncode == 0, control.stderr.decode()
        assert json.loads(control.stdout)["status"] == "INCONCLUSIVE"
        assert not json.loads(control.stdout)["mechanisms"]
    print(json.dumps({
        "schema_version": "1", "repository": ref, "items_verified": len(items),
        "online_witness_files_verified": online_verified,
        "replay_byte_identical": True, "tampered_item_rejected": True,
        "previous_output_preserved": True, "no_assessments_control": "INCONCLUSIVE",
        "result_sha256": hashlib.sha256(original).hexdigest(),
        "report_status": report["status"], "mechanism_status": "PARTIAL",
        "hypothesis_status": "INCONCLUSIVE", "target_state": "UNKNOWN",
        "candidate_tests_executed": False,
    }, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
