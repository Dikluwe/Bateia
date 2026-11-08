"""Inspect the current insight projection on the known, frozen 0106 corpus.

This is an authored diagnostic, not an independent semantic classifier or holdout.
It runs Bateia, never sourcegraph/conc, and preserves all historical artifacts.
"""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess


def sha(data):
    return hashlib.sha256(data).hexdigest()


def write(path, data):
    path.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--binary", required=True)
    args = parser.parse_args()
    output = Path(__file__).resolve().parent
    root = output.parents[3]
    historical = output.parent / "0106"
    baseline = json.loads((output / "baseline.json").read_text())
    for name, digest in baseline["go_sources"].items():
        assert sha((root / name).read_bytes()) == digest, name
    binary = Path(args.binary).resolve()
    command = [str(binary), "idea", "insights", "analyze", "--input", str(historical / "github-insight-input.json")]
    run = subprocess.run(command, capture_output=True, timeout=30)
    assert run.returncode == 0, run.stderr.decode()
    assert run.stdout == (historical / "github-insight-result.json").read_bytes()
    (output / "current-insight-result.json").write_bytes(run.stdout)
    report = json.loads(run.stdout)
    review = json.loads((historical / "source-review.json").read_text())
    bundle = json.loads((historical / "github-bundles.json").read_text())[0]
    assert bundle["repository"] == review["repository"]
    items = {item["path"]: item for item in bundle["items"]}
    witnesses = []
    for observation in review["observations"]:
        item = items[observation["path"]]
        assert sha(item["content"].encode()) == item["content_hash"].removeprefix("sha256:")
        start, end = observation["lines"]
        lines = item["content"].splitlines()
        assert 1 <= start <= end <= len(lines)
        witnesses.append({**observation, "content_hash": item["content_hash"], "excerpt": "\n".join(lines[start - 1:end]), "assessment_origin": "MANUAL_REVIEW_NOT_REASSESSED"})
    write(output / "source-witnesses.json", {"repository": review["repository"], "observations": witnesses, "candidate_tests_executed": False})
    target_ranges = {
        "03_infra/runworker/local_queue.go": [(55, 64), (72, 86), (132, 158), (161, 175)],
        "01_core/usecase/mechanism_extraction.go": [(439, 452)],
        "01_core/usecase/insight_investigation.go": [(185, 213)],
    }
    targets = []
    for name, ranges in target_ranges.items():
        data = (root / name).read_bytes()
        lines = data.decode().splitlines()
        targets.append({"path": name, "sha256": sha(data), "ranges": [{"lines": [start, end], "excerpt": "\n".join(lines[start - 1:end])} for start, end in ranges]})
    write(output / "target-witnesses.json", targets)
    assert len(report["mechanisms"]) == 1 and len(report["hypotheses"]) == 1
    mechanism, hypothesis = report["mechanisms"][0], report["hypotheses"][0]
    catalog = json.loads((root / "01_core/usecase/mechanism-rules.json").read_text())
    rule = next(rule for rule in catalog["rules"] if rule["mechanism_class"] == mechanism["mechanism_class"])
    write(output / "projection-observations.json", {
        "question": "O relatório preserva condições e diferenças específicas da fonte para orientar melhoria do Bateia?",
        "regime": "AUTHORED_DIAGNOSTIC_KNOWN_CORPUS",
        "binary_sha256": sha(binary.read_bytes()), "command": command,
        "replay_byte_identical": True, "result_sha256": sha(run.stdout),
        "strategy_equals_catalog_template": mechanism["strategy"] == rule["strategy"],
        "preconditions_equal_catalog_language_limits": mechanism["preconditions"] == sorted(rule["language_limits"]),
        "mechanism_tradeoffs": mechanism["tradeoffs"], "mechanism_failure_model": mechanism["failure_model"],
        "mechanism_status": mechanism["status"], "hypothesis_status": hypothesis["status"], "target_state": hypothesis["target_state"],
        "manual_observations_not_semantically_reassessed": len(witnesses),
        "all_baseline_go_sources_unchanged": True,
        "limits": ["No live source reread or candidate execution.", "Field comparison is not a semantic accuracy benchmark.", "Queue tests refer to Bateia, not to sourcegraph/conc.", "No improvement to target behavior has been implemented or measured."],
    })
    print("Known-corpus replay, source bindings, target snapshots and projection comparison recorded.")


if __name__ == "__main__":
    main()
