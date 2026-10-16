"""Characterize declared-label admission; this is not a semantic accuracy benchmark."""

import argparse
import copy
import json
from pathlib import Path
import subprocess


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--binary", required=True)
    parser.add_argument("--output-dir", required=True)
    args = parser.parse_args()
    binary = str(Path(args.binary).resolve())
    root = Path(__file__).resolve().parents[2] / "insights" / "0106"
    original = json.loads((root / "github-insight-input.json").read_text())
    destination = Path(args.output_dir).resolve()
    destination.mkdir(parents=True, exist_ok=True)
    cases = {name: copy.deepcopy(original) for name in (
        "baseline", "reordered", "duplicate_witness", "no_assessments",
        "tampered_hash", "unknown_rule", "declared_status_promoted", "wrong_declared_class",
    )}
    cases["reordered"]["sources"][0]["bundle"]["items"].reverse()
    cases["duplicate_witness"]["sources"][0]["assessments"][0]["capabilities"][0]["sources"].append(
        copy.deepcopy(original["sources"][0]["assessments"][0]["capabilities"][0]["sources"][0])
    )
    cases["no_assessments"]["sources"][0]["assessments"] = []
    cases["tampered_hash"]["sources"][0]["bundle"]["items"][0]["content"] += " altered"
    for capability in cases["unknown_rule"]["sources"][0]["assessments"][0]["capabilities"]:
        for source in capability["sources"]:
            source["rule_id"] = "research-rule-not-in-catalog"
    for name in ("declared_status_promoted", "wrong_declared_class"):
        for capability in cases[name]["sources"][0]["assessments"][0]["capabilities"]:
            capability["status"] = "PROVEN"
            if name == "wrong_declared_class":
                for source in capability["sources"]:
                    source["rule_id"] = "append-only-recovery-test" if source["kind"] == "test" else "append-only-write"
        if name == "wrong_declared_class":
            cases[name]["mechanism_classes"] = ["APPEND_ONLY_LOG_RECOVERY"]
    summaries = []
    reports = {}
    for name, request in cases.items():
        input_path = destination / f"{name}-input.json"
        input_path.write_text(json.dumps(request, ensure_ascii=False, indent=2) + "\n")
        result = subprocess.run([binary, "idea", "insights", "analyze", "--input", str(input_path)], capture_output=True, timeout=30)
        summary = {"case": name, "exit_code": result.returncode, "stderr": result.stderr.decode()}
        if result.returncode == 0:
            report = json.loads(result.stdout)
            reports[name] = report
            (destination / f"{name}-result.json").write_bytes(result.stdout)
            summary.update({
                "report_status": report["status"],
                "mechanisms": [{"class": m["mechanism_class"], "status": m["status"], "id": m["id"]} for m in report["mechanisms"]],
                "hypotheses": [{"status": h["status"], "target_state": h["target_state"]} for h in report["hypotheses"]],
            })
        summaries.append(summary)
    for name in ("reordered", "duplicate_witness"):
        assert reports[name]["mechanisms"] == reports["baseline"]["mechanisms"]
    for name in ("no_assessments", "unknown_rule"):
        assert reports[name]["status"] == "INCONCLUSIVE" and not reports[name]["mechanisms"]
    tampered = next(s for s in summaries if s["case"] == "tampered_hash")
    assert tampered["exit_code"] != 0 and "integrity" in tampered["stderr"].lower()
    assert reports["baseline"]["mechanisms"][0]["status"] == "PARTIAL"
    assert reports["declared_status_promoted"]["mechanisms"][0]["status"] == "PROVEN"
    assert reports["wrong_declared_class"]["mechanisms"][0]["mechanism_class"] == "APPEND_ONLY_LOG_RECOVERY"
    assert reports["wrong_declared_class"]["mechanisms"][0]["status"] == "PROVEN"
    for report in reports.values():
        assert all(h["target_state"] == "UNKNOWN" and h["status"] != "SUPPORTED" for h in report["hypotheses"])
    receipt = {
        "scope": "declared-label admission characterization, one frozen repository; not classifier accuracy",
        "checks_passed": True,
        "independent_semantic_reference": False,
        "cases": summaries,
    }
    (destination / "receipt.json").write_text(json.dumps(receipt, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps({"cases": len(cases), "checks_passed": True, "semantic_accuracy_measured": False}))


if __name__ == "__main__":
    main()
