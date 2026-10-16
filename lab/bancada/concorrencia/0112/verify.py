"""Replay frozen inputs with a supplied observer binary; no candidate execution."""
from collections import Counter
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tempfile

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BINARY = sys.argv[1] if len(sys.argv) == 2 else "/tmp/bateia-observer-0112"

def execute(args):
    first = subprocess.check_output([BINARY, *args], cwd=ROOT)
    second = subprocess.check_output([BINARY, *args], cwd=ROOT)
    assert first == second
    return json.loads(first), first

corpus = json.loads((HERE / "testdata/flow-corpus.json").read_text())
results = []
with tempfile.TemporaryDirectory() as directory:
    for case in corpus["cases"]:
        filename = Path(directory) / "input.json"
        filename.write_text(json.dumps(case["request"]))
        report, _ = execute(["observe", "--input", str(filename)])
        actual = dict(Counter(e["kind"] for e in report.get("relations", [])))
        assert actual == case["expected_relations"], (case["id"], actual)
        results.append({"id": case["id"], "family": case["family"], "passed": True, "actual_relations": actual})
    wrong = dict(corpus["cases"][0]["expected_relations"])
    wrong["SEND_CASE_CONTAINS_GO_CALL"] = 2
    assert results[0]["actual_relations"] != wrong

output = {"analyzer": "go-concurrency-syntax-0112-v4", "reference_origin": corpus["reference_origin"], "cases": len(results), "families": len({r["family"] for r in results}), "passed": len(results), "relation_profile_agreement": 1, "independent_semantic_accuracy_measured": False, "wrong_profile_negative_control_rejected": True, "results": results}
(HERE / "results/flow-evaluation.json").write_text(json.dumps(output, indent=2) + "\n")
input_args = ["observe", "--input", str(HERE.parent / "0111/testdata/conc-input.json")]
report, raw = execute([*input_args, "--environment", str(HERE.parent / "0111/testdata/environment.json")])
(HERE / "results/conc-flow.json").write_bytes(raw)
assert len(report.get("type_diagnostics", [])) == 0
counts = Counter(e["kind"] for e in report["relations"])
assert counts["SEND_CASE_CONTAINS_GO_CALL"] == counts["DEFERRED_RECEIVER_RECEIVE"] == counts["SAME_DECLARED_CHANNEL_FIELD"] == 1
without, raw_without = execute(input_args)
assert not without.get("relations")
(HERE / "results/without-environment.json").write_bytes(raw_without)
for name, path in [("baseline-evaluation", HERE.parent / "0110/testdata/baseline-v2.json"), ("types-evaluation", HERE.parent / "0110/testdata/types-corpus.json")]:
    result, data = execute(["evaluate", "--corpus", str(path)])
    assert result["passed"] == result["cases"]
    (HERE / "results" / (name + ".json")).write_bytes(data)
receipt = {"analyzer": report["analyzer"], "candidate_code_executed": False, "independent_semantic_accuracy_measured": False, "byte_identical_replays": True, "without_environment_relations": 0, "relations": dict(counts), "environment_sha256": report["environment_sha256"], "corpus_sha256": hashlib.sha256((HERE / "testdata/flow-corpus.json").read_bytes()).hexdigest(), "result_sha256": hashlib.sha256(raw).hexdigest(), "binary_sha256": hashlib.sha256(Path(BINARY).read_bytes()).hexdigest()}
(HERE / "results/smoke-receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({"cases": len(results), "passed": len(results), "conc_relations": sum(counts.values())}))
