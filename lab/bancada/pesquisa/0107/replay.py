"""Replay the offline Bateia stages against frozen research inputs."""

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--binary", required=True)
    args = parser.parse_args()
    binary = str(Path(args.binary).resolve())
    root = Path(__file__).resolve().parent
    bundles = json.loads((root / "review-bundles.json").read_text())
    receipts = json.loads((root / "raw-fallback-receipts.json").read_text())
    receipt_by_url = {r["url"]: r for r in receipts}
    checked = 0
    for bundle in bundles:
        assert not bundle["complete"] and bundle.get("failures")
        ref = bundle["repository"]
        assert len(ref["commit"]) == 40 and ref["commit"] != "HEAD"
        for item in bundle["items"]:
            assert item["ref"] == ref
            digest = hashlib.sha256(item["content"].encode()).hexdigest()
            assert item["content_hash"].removeprefix("sha256:") == digest
            if item["path"] != "@repository-metadata":
                url = f"https://raw.githubusercontent.com/{ref['owner']}/{ref['repo']}/{ref['commit']}/{item['path']}"
                assert receipt_by_url[url]["content_sha256"] == digest
                checked += 1
    def path(name):
        return str(root / name)
    stages = [
        ("automatic-plan.json", ["plan", "--input", path("idea.json")]),
        ("catalog-screen.json", ["catalog", "screen", "--snapshot", path("github-discovery.json")]),
        ("assessment.json", ["evidence", "assess", "--idea", path("idea.json"), "--screening", path("assessment-screening.json"), "--evidence", path("review-bundles.json"), "--rules", path("locator-rules.json")]),
        ("evidence-metrics.json", ["evidence", "evaluate", "--assessment", path("assessment.json"), "--evidence", path("review-bundles.json"), "--judgments", path("independent-judgments.json")]),
        ("insight-result.json", ["insights", "analyze", "--input", path("insight-input.json")]),
    ]
    passed = []
    with tempfile.TemporaryDirectory(prefix="bateia-replay-0107-") as directory:
        for original, arguments in stages:
            output = Path(directory) / original
            result = subprocess.run([binary, "idea"] + arguments + ["--output", str(output)], cwd=directory, capture_output=True, timeout=60)
            assert result.returncode == 0, result.stderr.decode()
            assert output.read_bytes() == (root / original).read_bytes(), original
            passed.append(original)
    metrics = json.loads((root / "evidence-metrics.json").read_text())
    assert metrics["evidence_precision"] is None and not metrics["sample_sufficient"]
    insight = json.loads((root / "insight-result.json").read_text())
    assert insight["status"] == "INCONCLUSIVE" and not insight["mechanisms"]
    assert len(insight["extraction_gaps"]) == 24
    print(json.dumps({"offline_stages_byte_identical": passed, "external_source_files_hash_verified": checked, "semantic_accuracy_measured": False}, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
