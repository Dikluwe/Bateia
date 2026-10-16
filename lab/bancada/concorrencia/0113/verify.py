"""Verify frozen callback profiles and replay with the supplied binary; no candidate execution."""
from collections import Counter
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tempfile

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BINARY = sys.argv[1] if len(sys.argv) == 2 else "/tmp/bateia-observer-0113"

def run(args):
    a = subprocess.check_output([BINARY, *args], cwd=ROOT)
    b = subprocess.check_output([BINARY, *args], cwd=ROOT)
    assert a == b
    return json.loads(a), a

corpus = json.loads((HERE / "testdata/callback-corpus.json").read_text())
results = []
with tempfile.TemporaryDirectory() as d:
    filename = Path(d) / "input.json"
    for c in corpus["cases"]:
        filename.write_text(json.dumps(c["request"]))
        r, _ = run(["observe", "--input", str(filename)])
        actual = dict(Counter(e["kind"] for e in r.get("callback_relations", [])))
        assert actual == c["expected_callbacks"], (c["id"], actual)
        results.append({"id": c["id"], "family": c["family"], "passed": True, "actual_callbacks": actual})
    wrong = dict(corpus["cases"][0]["expected_callbacks"])
    wrong["CALLBACK_PARAMETER_INVOCATION"] = 2
    assert results[0]["actual_callbacks"] != wrong
    historical = json.loads((HERE.parent / "0112/testdata/flow-corpus.json").read_text())
    for c in historical["cases"]:
        filename.write_text(json.dumps(c["request"]))
        r, _ = run(["observe", "--input", str(filename)])
        assert dict(Counter(e["kind"] for e in r.get("relations", []))) == c["expected_relations"]
report = {"analyzer": "go-concurrency-syntax-0113-v5", "reference_origin": corpus["reference_origin"], "cases": len(results), "families": len({x["family"] for x in results}), "passed": len(results), "callback_profile_agreement": 1, "independent_semantic_accuracy_measured": False, "wrong_profile_negative_control_rejected": True, "legacy_flow_cases_passed": len(historical["cases"]), "results": results}
(HERE / "results/callback-evaluation.json").write_text(json.dumps(report, indent=2) + "\n")
args = ["observe", "--input", str(HERE.parent / "0111/testdata/conc-input.json"), "--environment", str(HERE.parent / "0111/testdata/environment.json")]
r, raw = run(args)
(HERE / "results/conc-callbacks.json").write_bytes(raw)
old = json.loads((HERE.parent / "0112/results/conc-flow.json").read_text())
assert r["relations"] == old["relations"]
bindings = [e for e in r["callback_relations"] if e["kind"] == "CALLBACK_ARGUMENT_BINDING" and e["to"]["path"] == "waitgroup.go"]
forward = next(e for e in r["callback_relations"] if e["kind"] == "CALLBACK_PARAMETER_FORWARDING" and e["from"]["path"] == "waitgroup.go")
invoke = next(e for e in r["callback_relations"] if e["kind"] == "CALLBACK_PARAMETER_INVOCATION" and e["from"]["path"] == "panics/panics.go")
assert len(bindings) == 2 and all(e["to"] == forward["parameter_declaration"] for e in bindings)
assert forward["to"] == invoke["from"]
without, data = run(args[:3])
assert not without.get("callback_relations")
(HERE / "results/without-environment.json").write_bytes(data)
for name, path in [("baseline-evaluation", HERE.parent / "0110/testdata/baseline-v2.json"), ("types-evaluation", HERE.parent / "0110/testdata/types-corpus.json")]:
    evaluation, data = run(["evaluate", "--corpus", str(path)])
    assert evaluation["passed"] == evaluation["cases"]
    (HERE / "results" / (name + ".json")).write_bytes(data)
receipt = {"analyzer": r["analyzer"], "candidate_code_executed": False, "independent_semantic_accuracy_measured": False, "byte_identical_replays": True, "structural_callback_chain_linked": True, "callback_execution_verified": False, "legacy_relations_identical": True, "callbacks": dict(Counter(e["kind"] for e in r["callback_relations"])), "environment_sha256": r["environment_sha256"], "result_sha256": hashlib.sha256(raw).hexdigest(), "corpus_sha256": hashlib.sha256((HERE / "testdata/callback-corpus.json").read_bytes()).hexdigest(), "binary_sha256": hashlib.sha256(Path(BINARY).read_bytes()).hexdigest()}
(HERE / "results/smoke-receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({"cases": len(results), "passed": len(results), "legacy_flow_cases": len(historical["cases"]), "conc_callbacks": len(r["callback_relations"])}))
